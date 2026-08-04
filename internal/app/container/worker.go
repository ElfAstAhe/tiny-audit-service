package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

const (
	InstanceAuthAuditTailCutter       string = "authAuditTailCutter"
	InstanceDataAuditTailCutter       string = "dataAuditTailCutter"
	InstanceLoginAttemptsAMQPListener string = "login-attempts-amqp-listener"
)

type WorkerContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*WorkerContainer)(nil)
var _ container.LazyContainer = (*WorkerContainer)(nil)

func NewWorkerContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *WorkerContainer {
	return &WorkerContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(WorkerContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

func (wc *WorkerContainer) Init(initCtx context.Context) error {
	err := errors.Join(
		wc.RegisterProvider(InstanceAuthAuditTailCutter, wc.providerAuthAuditTailCutter),
		wc.RegisterProvider(InstanceDataAuditTailCutter, wc.providerDataAuditTailCutter),
		wc.RegisterProvider(InstanceLoginAttemptsAMQPListener, wc.providerLoginAttemptsAMQPListener),
	)
	if err != nil {
		return errs.NewContainerError(wc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
