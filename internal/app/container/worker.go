package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

type WorkerContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*WorkerContainer)(nil)
var _ container.LazyContainer = (*WorkerContainer)(nil)

func NewWorkerContainer(orchestrator container.Orchestrator) *HTTPContainer {
	return &HTTPContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(HTTPContainerName, orchestrator),
	}
}

func (wc *WorkerContainer) Init(initCtx context.Context) error {
	err := errors.Join(
	//
	)
	if err != nil {
		return errs.NewContainerError(wc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
