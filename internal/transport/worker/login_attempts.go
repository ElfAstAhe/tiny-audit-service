package worker

import (
	"context"
	"errors"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker/mapper"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

type LoginAttempts struct {
	*worker.BaseSchedulerDispatcher[*dto.LoginAttemptWorkerJob]
	opts               *LoginAttemptsOptions
	receiver           libamqp.Receiver[*amqp.ReceiveOptions]
	authAuditUC        usecase.AuthAuditUseCase
	batchSize          int
	batchReadTimeout   time.Duration
	acknowledgeTimeout time.Duration
}

var _ worker.Scheduler = (*LoginAttempts)(nil)
var _ worker.CommonWorker = (*LoginAttempts)(nil)
var _ container.Runner = (*LoginAttempts)(nil)

func NewLoginAttempts(
	opts ...LoginAttemptsOption,
) (*LoginAttempts, error) {
	localOpts := NewLoginAttemptsOptions()

	for _, opt := range opts {
		opt(localOpts)
	}

	if err := localOpts.Validate(); err != nil {
		return nil, errs.NewCommonError("login attempts scheduled dispatcher options validation failed", err)
	}

	res := &LoginAttempts{
		authAuditUC:        localOpts.AuthAuditUC,
		receiver:           localOpts.Receiver,
		opts:               localOpts,
		batchSize:          localOpts.BatchSize,
		batchReadTimeout:   localOpts.BatchReadTimeout,
		acknowledgeTimeout: localOpts.AcknowledgeTimeout,
	}

	res.BaseSchedulerDispatcher = worker.NewBaseSchedulerDispatcher[*dto.LoginAttemptWorkerJob](
		localOpts.Name,
		localOpts.DispatcherOpts,
		res.dataProvider,
		res.storeAuthAudit,
		localOpts.Logger,
	)

	return res, nil
}

func (la *LoginAttempts) dataProvider(ctx context.Context, eventTime time.Time) ([]*dto.LoginAttemptWorkerJob, error) {
	la.GetLogger().Debugf("login attempts receiver %s time event %s data provider start", la.GetName(), eventTime.Format(time.DateTime))
	defer la.GetLogger().Debugf("login attempts receiver %s time event %s data provider finish", la.GetName(), eventTime.Format(time.DateTime))

	resData := make([]*dto.LoginAttemptWorkerJob, 0)

	// timed context
	brokerCtx, brokerCancel := context.WithTimeout(ctx, la.batchReadTimeout)
	defer brokerCancel()

	for {
		if len(resData) >= la.batchSize {
			break
		}

		select {
		case <-brokerCtx.Done():
			if brokerCtx.Err() != nil {
				la.GetLogger().Debugf("receiver context done with err: %v", brokerCtx.Err())

				// cancellation, return nothing and not errors
				if errors.Is(brokerCtx.Err(), context.Canceled) {
					return []*dto.LoginAttemptWorkerJob{}, nil
				}
				// deadline exceeded, return all what we got by read timeout
				if errors.Is(brokerCtx.Err(), context.DeadlineExceeded) {
					return resData, nil
				}
			}

			return resData, nil
		default:
			// receive message
			message, err := la.receiver.Receive(brokerCtx, la.opts.ReceiveOpts)
			// got error ?
			if err != nil {
				// handled errors
				// context canceled
				if errors.Is(err, context.Canceled) {
					return []*dto.LoginAttemptWorkerJob{}, nil
				}
				// read timeout
				if errors.Is(err, context.DeadlineExceeded) {
					return resData, nil
				}

				// unhandled error
				return nil, errs.NewCommonError("receiver receive failed", err)
			}
			// convert message to dto
			job, err := mapper.MapMessageToLoginAttemptWorkerJob(message)
			if err != nil {
				// reject message (send it into DLQ)
				errReject := la.receiver.Reject(brokerCtx, message, err)
				// got any error, exit loop
				if errReject != nil {
					return nil, errs.NewCommonError("receiver reject failed", errReject)
				}

				continue
			}
			// add data
			resData = append(resData, job)
		}
	}

	return resData, nil
}

func (la *LoginAttempts) storeAuthAudit(ctx context.Context, workerIndex int, data *dto.LoginAttemptWorkerJob) error {
	la.GetLogger().Debugf("login attempts receiver %s worker %v store auth audit start", la.GetName(), workerIndex)
	defer la.GetLogger().Debugf("login attempts receiver %s worker %v store auth audit finish", la.GetName(), workerIndex)

	err := la.authAuditUC.Audit(ctx, mapper.MapLoginAttemptEventDTOToAuthAuditModel(data.Data))

	brokerCtx, brokerCancel := context.WithTimeout(ctx, la.acknowledgeTimeout)
	defer brokerCancel()

	if err != nil {
		la.GetLogger().Debugf("login attempts receiver %s worker %v store auth audit failed %v", la.GetName(), workerIndex, err)
		// unique violation, pass it, data already exists in ms storage
		if _, ok := errors.AsType[*errs.BllUniqueError](err); ok {
			if acceptErr := la.receiver.Accept(brokerCtx, data.Message); acceptErr != nil {
				return errs.NewCommonError("receiver accept failed", acceptErr)
			}

			return nil
		}

		// any errors
		if releaseErr := la.receiver.Release(brokerCtx, data.Message); releaseErr != nil {
			return errs.NewCommonError("receiver release failed", releaseErr)
		}

		return errs.NewCommonError("store auth audit failed", err)
	}

	if acceptErr := la.receiver.Accept(brokerCtx, data.Message); acceptErr != nil {
		return errs.NewCommonError("receiver accept failed", acceptErr)
	}

	return nil
}
