package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

const (
	InstanceAuthFacade string = "AuthFacade"
	InstanceDataFacade string = "DataFacade"
)

type FacadeContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*FacadeContainer)(nil)
var _ container.LazyContainer = (*FacadeContainer)(nil)

func NewFacadeContainer(orchestrator container.Orchestrator) *FacadeContainer {
	return &FacadeContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(FacadeContainerName, orchestrator),
	}
}

func (fc *FacadeContainer) Init(ctx context.Context) error {
	err := errors.Join(
		fc.RegisterProvider(InstanceAuthFacade, fc.providerAuthFacade),
		fc.RegisterProvider(InstanceDataFacade, fc.providerDataFacade),
	)
	if err != nil {
		return errs.NewContainerError(fc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
