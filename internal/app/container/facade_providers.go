package container

import (
	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/rbac"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

func (fc *FacadeContainer) providerAuthFacade() (any, error) {
	authAuditUCInst, err := container.GetInstance[usecase.AuthAuditUseCase](InstanceAuthAuditUC)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}
	authListByPeriodUCInst, err := container.GetInstance[usecase.AuthListByPeriodUseCase](InstanceAuthListByPeriodUC)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}
	authListByUsernameUCInst, err := container.GetInstance[usecase.AuthListByUsernameUseCase](InstanceAuthListByUsernameUC)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}

	return facade.NewAuthAudit(
		authAuditUCInst,
		authListByPeriodUCInst,
		authListByUsernameUCInst,
	), nil
}

func (fc *FacadeContainer) providerDataFacade() (any, error) {
	dataAuditUCInst, err := container.GetInstance[usecase.DataAuditUseCase](InstanceDataAuditUC)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}
	dataListByPeriodUCInst, err := container.GetInstance[usecase.DataListByPeriodUseCase](InstanceDataListByPeriodUC)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}
	dataListByInstanceUCInst, err := container.GetInstance[usecase.DataListByInstanceUseCase](InstanceDataListByInstanceUC)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}

	return facade.NewDataAudit(
		dataAuditUCInst,
		dataListByPeriodUCInst,
		dataListByInstanceUCInst,
	), nil
}

func (fc *FacadeContainer) providerRBACAuthFacade() (any, error) {
	authHelperInst, err := container.GetInstance[auth.Helper](InstanceAuthHelper)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}
	authAuditFacadeInst, err := container.GetInstance[facade.AuthAudit](InstanceAuthFacade)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}

	return rbac.NewAuthAudit(authHelperInst, authAuditFacadeInst), nil
}

func (fc *FacadeContainer) providerRBACDataFacade() (any, error) {
	authHelperInst, err := container.GetInstance[auth.Helper](InstanceAuthHelper)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}
	dataAuditFacadeInst, err := container.GetInstance[facade.DataAudit](InstanceDataFacade)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}

	return rbac.NewDataAudit(authHelperInst, dataAuditFacadeInst), nil
}
