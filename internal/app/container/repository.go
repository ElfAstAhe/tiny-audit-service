package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

const (
	InstanceAuthAuditRepo        string = "authAuditRepo"
	InstanceAuthAuditMetricsRepo string = "authAuditMetricsRepo"
	InstanceAuthAuditTraceRepo   string = "authAuditTraceRepo"
	InstanceDataAuditRepo        string = "dataAuditRepo"
	InstanceDataAuditMetricsRepo string = "dataAuditMetricsRepo"
	InstanceDataAuditTraceRepo   string = "dataAuditTraceRepo"
)

type RepositoryContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*RepositoryContainer)(nil)
var _ container.LazyContainer = (*RepositoryContainer)(nil)

func NewRepositoryContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *RepositoryContainer {
	return &RepositoryContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(RepositoryContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

func (rc *RepositoryContainer) Init(ctx context.Context) error {
	err := errors.Join(
		rc.RegisterProvider(InstanceAuthAuditRepo, rc.providerAuthAuditRepo),
		rc.RegisterProvider(InstanceAuthAuditMetricsRepo, rc.providerAuthAuditMetricsRepo),
		rc.RegisterProvider(InstanceAuthAuditTraceRepo, rc.providerAuthAuditTraceRepo),
		rc.RegisterProvider(InstanceDataAuditRepo, rc.providerDataAuditRepo),
		rc.RegisterProvider(InstanceDataAuditMetricsRepo, rc.providerDataAuditMetricsRepo),
		rc.RegisterProvider(InstanceDataAuditTraceRepo, rc.providerDataAuditTraceRepo),
	)
	if err != nil {
		return errs.NewContainerError(rc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
