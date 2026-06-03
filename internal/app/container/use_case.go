package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

const (
	InstanceTM                   string = "TransactionManager"
	InstanceAuthAuditUC          string = "AuthAuditUC"
	InstanceAuthListByPeriodUC   string = "AuthListByPeriodUC"
	InstanceAuthListByUsernameUC string = "AuthListByUsernameUC"
	InstanceDataAuditUC          string = "DataAuditUC"
	InstanceDataListByPeriodUC   string = "DataListByPeriodUC"
	InstanceDataListByInstanceUC string = "DataListByInstanceUC"
	InstanceAuthAuditTailGetUC   string = "AuthAuditTailGetUC"
	InstanceAuthAuditTailCutUC   string = "AuthAuditTailCutUC"
	InstanceDataAuditTailGetUC   string = "DataAuditTailGetUC"
	InstanceDataAuditTailCutUC   string = "DataAuditTailCutUC"
)

type UseCaseContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*UseCaseContainer)(nil)
var _ container.LazyContainer = (*UseCaseContainer)(nil)

func NewUseCaseContainer(orchestrator container.Orchestrator) *UseCaseContainer {
	return &UseCaseContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(UseCaseContainerName, orchestrator),
	}
}

func (ucc *UseCaseContainer) Init(ctx context.Context) error {
	err := errors.Join(
		ucc.RegisterProvider(InstanceTM, ucc.providerTM),
		ucc.RegisterProvider(InstanceAuthAuditUC, ucc.providerAuthAuditUC),
		ucc.RegisterProvider(InstanceAuthListByPeriodUC, ucc.providerAuthListByPeriodUC),
		ucc.RegisterProvider(InstanceAuthListByUsernameUC, ucc.providerAuthListByUsernameUC),
		ucc.RegisterProvider(InstanceDataAuditUC, ucc.providerDataAuditUC),
		ucc.RegisterProvider(InstanceDataListByPeriodUC, ucc.providerDataListByPeriodUC),
		ucc.RegisterProvider(InstanceDataListByInstanceUC, ucc.providerDataListByInstanceUC),
		ucc.RegisterProvider(InstanceAuthAuditTailGetUC, ucc.providerAuthAuditTailGetUC),
		ucc.RegisterProvider(InstanceAuthAuditTailCutUC, ucc.providerAuthAuditTailCutUC),
		ucc.RegisterProvider(InstanceDataAuditTailGetUC, ucc.providerDataAuditTailGetUC),
		ucc.RegisterProvider(InstanceDataAuditTailCutUC, ucc.providerDataAuditTailCutUC),
	)
	if err != nil {
		return errs.NewContainerError(ucc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
