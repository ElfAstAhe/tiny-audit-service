package worker

import (
	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

type LoginAttempts struct {
	*worker.BaseSchedulerDispatcher[*dto.AuthAuditDTO]
	conf        *config.LoginAttemptsConfig
	receiver    libamqp.Receiver[*amqp.ReceiveOptions, *amqp.MessageHeader]
	authAuditUC usecase.AuthAuditUseCase
}

var _ worker.Scheduler = (*LoginAttempts)(nil)
var _ worker.CommonWorker = (*LoginAttempts)(nil)
var _ container.Runner = (*LoginAttempts)(nil)

func NewLoginAttempts(
	opts ...LoginAttemptOption,
) (*LoginAttempts, error) {
	res := &LoginAttempts{
		BaseSchedulerDispatcher: worker.NewBaseSchedulerDispatcher[*dto.AuthAuditDTO{}](),
	}

	for _, option := range opts {
		option(res)
	}

	if err := res.validate(); err != nil {
		return nil, errs.NewCommonError("login attempts scheduled dispatcher options validation failed", err)
	}

	return res, nil
}

func (la *LoginAttempts) validate() error {
	if la.GetName() == "" {
		return errs.NewCommonError("name empty", nil)
	}
	if utils.IsNil(la.receiver) {
		return errs.NewCommonError("amqp receiver is nil", nil)
	}
	if utils.IsNil(la.authAuditUC) {
		return errs.NewCommonError("auth audit use case is nil", nil)
	}

	return nil
}
