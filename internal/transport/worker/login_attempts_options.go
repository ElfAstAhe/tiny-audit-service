package worker

import (
	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
)

type LoginAttemptsOptions struct {
	DispatcherOpts *worker.BaseSchedulerDispatcherConfig
	ReceiveOpts    *amqp.ReceiveOptions
}

func NewLoginAttemptsOptions(
	dispatcherOpts *worker.BaseSchedulerDispatcherConfig,
	ReceiveOpts *amqp.ReceiveOptions,
) *LoginAttemptsOptions {
	return &LoginAttemptsOptions{
		DispatcherOpts: dispatcherOpts,
		ReceiveOpts:    ReceiveOpts,
	}
}

func (lao LoginAttemptsOptions) Validate() error {
	if utils.IsNil(lao.DispatcherOpts) {
		return errs.NewCommonError("dispatcher options not provided (nil)", nil)
	}

	return nil
}
