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
	*worker.BaseSchedulerDispatcher[*dto.AuthAuditDTO]
	opts             *LoginAttemptsOptions
	receiver         libamqp.Receiver[*amqp.ReceiveOptions]
	authAuditUC      usecase.AuthAuditUseCase
	batchSize        int
	batchReadTimeout time.Duration
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
		authAuditUC:      localOpts.AuthAuditUC,
		receiver:         localOpts.Receiver,
		opts:             localOpts,
		batchSize:        localOpts.BatchSize,
		batchReadTimeout: localOpts.BatchReadTimeout,
	}

	res.BaseSchedulerDispatcher = worker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		localOpts.Name,
		localOpts.DispatcherOpts,
		res.dataProvider,
		res.storeAuthAudit,
		localOpts.Logger,
	)

	return res, nil
}

func (la *LoginAttempts) dataProvider(ctx context.Context, eventTime time.Time) ([]*dto.AuthAuditDTO, error) {
	la.GetLogger().Debugf("login attempts receiver %s time event %s data provider start", la.GetName(), eventTime.Format(time.DateTime))
	defer la.GetLogger().Debugf("login attempts receiver %s time event %s data provider finish", la.GetName(), eventTime.Format(time.DateTime))

	resData := make([]*dto.AuthAuditDTO, 0)

	// timed context
	receiveCtx, receiveCancel := context.WithTimeout(ctx, 1*time.Second)
	defer receiveCancel()

	for {
		if len(resData) > la.batchSize {
			break
		}

		select {
		case <-receiveCtx.Done():
			if ctx.Err() != nil {
				la.GetLogger().Warn("parent context canceled, release all received messages")
				la.releaseAll(resData)

				return nil, ctx.Err()
			}

			return resData, nil
		default:
			message, err := la.receiver.Receive(receiveCtx, la.opts.ReceiveOpts)
			if err != nil {
				if ctx.Err() != nil || errors.Is(err, context.Canceled) {
					la.releaseAll(resData)

					return nil, err
				}
				if errors.Is(err, context.DeadlineExceeded) {
					continue
				}

				continue
			}
			data, err := mapper.ToAuthAuditDTO(message)
			if err != nil {

			}

			resData = append(resData)

		}
	}

	return resData, nil
}

func (la *LoginAttempts) storeAuthAudit(ctx context.Context, workerIndex int, data *dto.AuthAuditDTO) error {

}
