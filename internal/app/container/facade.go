package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие уникальные строковые идентификаторы инстансов фасадов и их RBAC-декораторов в DI-контейнере.
const (
	InstanceAuthFacade     string = "auth-facade"
	InstanceDataFacade     string = "data-facade"
	InstanceRBACAuthFacade string = "rbac-auth-facade"
	InstanceRBACDataFacade string = "rbac-data-facade"
)

// FacadeContainer представляет собой ленивый (lazy) DI-контейнер для конфигурации,
// регистрации и управления жизненным циклом прикладных фасадов приложения и слоев авторизации.
type FacadeContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*FacadeContainer)(nil)
var _ container.LazyContainer = (*FacadeContainer)(nil)

// NewFacadeContainer создает и инициализирует новый экземпляр контейнера прикладных фасадов FacadeContainer.
func NewFacadeContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *FacadeContainer {
	return &FacadeContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(FacadeContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init выполняет первичную инициализацию контейнера, регистрируя провайдеры для базовых фасадов и их RBAC-оберток.
func (fc *FacadeContainer) Init(ctx context.Context) error {
	err := errors.Join(
		fc.RegisterProvider(InstanceAuthFacade, fc.providerAuthFacade),
		fc.RegisterProvider(InstanceDataFacade, fc.providerDataFacade),
		fc.RegisterProvider(InstanceRBACAuthFacade, fc.providerRBACAuthFacade),
		fc.RegisterProvider(InstanceRBACDataFacade, fc.providerRBACDataFacade),
	)
	if err != nil {
		return errs.NewContainerError(fc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
