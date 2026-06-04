package container

import (
	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

func (fc *FacadeContainer) providerAuthFacade() (any, error) {
	authHelperInst, err := container.GetInstance[auth.Helper](InstanceAuthHelper)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}
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

	return facade.NewAuthAuditFacade(
		authHelperInst,
		authAuditUCInst,
		authListByPeriodUCInst,
		authListByUsernameUCInst,
	), nil
}

func (fc *FacadeContainer) providerDataFacade() (any, error) {
	authHelperInst, err := container.GetInstance[auth.Helper](InstanceAuthHelper)
	if err != nil {
		return nil, errs.NewContainerError(fc.GetName(), "provider: retrieve instance failed", err)
	}
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

	return facade.NewDataAuditFacade(
		authHelperInst,
		dataAuditUCInst,
		dataListByPeriodUCInst,
		dataListByInstanceUCInst,
	), nil
}
