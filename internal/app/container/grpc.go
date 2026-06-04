package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

const (
	InstanceAuthAuditGRPCService string = "authAuditGRPCService"
	InstanceDataAuditGRPCService string = "dataAuditGRPCService"
	InstanceGRPCRunner           string = "grpcRunner"
)

type GRPCContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*GRPCContainer)(nil)
var _ container.LazyContainer = (*GRPCContainer)(nil)

func NewGRPCContainer(orchestrator container.Orchestrator) *GRPCContainer {
	return &GRPCContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(GRPCContainerName, orchestrator),
	}
}

func (gc *GRPCContainer) Init(ctx context.Context) error {
	err := errors.Join(
		gc.RegisterProvider(InstanceAuthAuditGRPCService, gc.providerAuthAuditGRPCService),
		gc.RegisterProvider(InstanceDataAuditGRPCService, gc.providerDataAuditGRPCService),
		gc.RegisterProvider(InstanceGRPCRunner, gc.providerGRPCRunner),
	)
	if err != nil {
		return errs.NewContainerError(gc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
