package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие уникальные имена регистрируемых инстансов воркеров и слушателей в DI-контейнере.
const (
	InstanceAuthAuditTailCutter   string = "authAuditTailCutter"
	InstanceDataAuditTailCutter   string = "dataAuditTailCutter"
	InstanceLoginAttemptsListener string = "login-attempts-listener"
)

// WorkerContainer представляет собой ленивый (lazy) DI-контейнер для управления
// жизненным циклом фоновых воркеров, слушателей очередей и планировщиков очистки хвостов таблиц.
type WorkerContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*WorkerContainer)(nil)
var _ container.LazyContainer = (*WorkerContainer)(nil)

// NewWorkerContainer создает и инициализирует новый экземпляр контейнера фоновых воркеров WorkerContainer.
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

// Init выполняет первичную инициализацию контейнера, регистрируя провайдеры для всех фоновых задач.
func (wc *WorkerContainer) Init(initCtx context.Context) error {
	err := errors.Join(
		wc.RegisterProvider(InstanceAuthAuditTailCutter, wc.providerAuthAuditTailCutter),
		wc.RegisterProvider(InstanceDataAuditTailCutter, wc.providerDataAuditTailCutter),
		wc.RegisterProvider(InstanceLoginAttemptsListener, wc.providerLoginAttemptsListener),
	)
	if err != nil {
		return errs.NewContainerError(wc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
