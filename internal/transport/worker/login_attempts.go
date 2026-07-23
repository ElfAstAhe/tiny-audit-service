package worker

import (
	"context"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

type LoginAttempts struct {
	*worker.BaseSchedulerDispatcher[*dto.AuthAuditDTO]
	opts        *LoginAttemptsOptions
	receiver    libamqp.Receiver[*amqp.ReceiveOptions, *amqp.MessageHeader]
	authAuditUC usecase.AuthAuditUseCase
}

var _ worker.Scheduler = (*LoginAttempts)(nil)
var _ worker.CommonWorker = (*LoginAttempts)(nil)
var _ container.Runner = (*LoginAttempts)(nil)

func NewLoginAttempts(
	name string,
	receiver libamqp.Receiver[*amqp.ReceiveOptions, *amqp.MessageHeader],
	authAuditUC usecase.AuthAuditUseCase,
	opts *LoginAttemptsOptions,
	log logger.Logger,
) (*LoginAttempts, error) {
	if err := opts.Validate(); err != nil {
		return nil, errs.NewCommonError("login attempts scheduled dispatcher options validation failed", err)
	}

	res := &LoginAttempts{
		authAuditUC: authAuditUC,
		receiver:    receiver,
		opts:        opts,
	}

	res.BaseSchedulerDispatcher = worker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO](
		name,
		opts.DispatcherOpts,
		res.dataProvider,
		res.storeAuthAudit,
		log,
	)

	return res, nil
}

func (la *LoginAttempts) dataProvider(ctx context.Context, eventTime time.Time) ([]*dto.AuthAuditDTO, error) {
	la.GetLogger().Debugf("login attempts receiver %s time event %s data provider start", la.GetName(), eventTime.Format(time.DateTime))
	defer la.GetLogger().Debugf("login attempts receiver %s time event %s data provider finish", la.GetName(), eventTime.Format(time.DateTime))

	messages := make([]*libamqp.Message[*amqp.MessageHeader], 0)

	receiveCtx, receiveCancel := context.WithTimeout(ctx, 1*time.Second)
	defer receiveCancel()

	for {
		message, err := la.receiver.Receive(receiveCtx, la.opts.ReceiveOpts)
		if err != nil {

		}
	}
}

func (la *LoginAttempts) storeAuthAudit(ctx context.Context, workerIndex int, data *dto.AuthAuditDTO) error {

}
