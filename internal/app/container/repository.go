package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

const (
	InstanceAuthAuditRepo string = "authAuditRepo"
	InstanceDataAuditRepo string = "dataAuditRepo"
)

type RepositoryContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*RepositoryContainer)(nil)
var _ container.LazyContainer = (*RepositoryContainer)(nil)

func NewRepositoryContainer(orchestrator container.Orchestrator) *RepositoryContainer {
	return &RepositoryContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(RepositoryContainerName, orchestrator),
	}
}

func (rc *RepositoryContainer) Init(ctx context.Context) error {
	err := errors.Join(
		rc.RegisterProvider(InstanceAuthAuditRepo, rc.providerAuthAuditRepo),
		rc.RegisterProvider(InstanceDataAuditRepo, rc.providerDataAuditRepo),
	)
	if err != nil {
		return errs.NewContainerError(rc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
