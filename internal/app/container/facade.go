package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

const (
	InstanceAuthFacade     string = "auth-facade"
	InstanceDataFacade     string = "data-facade"
	InstanceRBACAuthFacade string = "rbac-auth-facade"
	InstanceRBACDataFacade string = "rbac-data-facade"
)

type FacadeContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*FacadeContainer)(nil)
var _ container.LazyContainer = (*FacadeContainer)(nil)

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
