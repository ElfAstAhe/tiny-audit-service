package container

import (
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	libuc "github.com/ElfAstAhe/go-service-template/pkg/usecase"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase/telemetry"
)

func (ucc *UseCaseContainer) providerUnitOfWork() (any, error) {
	tmInst, err := container.GetInstance[db.TransactionManager](InstanceTM)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return libuc.NewUnitOfWork(tmInst, nil), nil
}

func (ucc *UseCaseContainer) providerAuthAuditUC() (any, error) {
	uwInst, err := container.GetInstance[libdom.UnitOfWork](InstanceTM)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	authAuditRepoInst, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return telemetry.NewAuthAuditUseCase("AuthAuditUseCase", usecase.NewAuthAuditUseCase(uwInst, authAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerAuthListByPeriodUC() (any, error) {
	authAuditRepoInst, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return telemetry.NewAuthListByPeriodUseCase("AuthListByPeriodUseCase", usecase.NewAuthListByPeriodUseCase(authAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerAuthListByUsernameUC() (any, error) {
	authAuditRepoInst, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return telemetry.NewAuthListByUsernameUseCase("AuthListByUsernameUseCase", usecase.NewAuthListByUsernameUseCase(authAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerDataAuditUC() (any, error) {
	uwInst, err := container.GetInstance[libdom.UnitOfWork](InstanceTM)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	dataAuditRepoInst, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return telemetry.NewDataAuditUseCase("DataAuditUseCase", usecase.NewDataAuditUseCase(uwInst, dataAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerDataListByPeriodUC() (any, error) {
	dataAuditRepoInst, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return telemetry.NewDataListByPeriodUseCase("DataListByPeriodUseCase", usecase.NewDataListByPeriodUseCase(dataAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerDataListByInstanceUC() (any, error) {
	dataAuditRepoInst, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return telemetry.NewDataListByInstanceUseCase("DataListByInstanceUseCase", usecase.NewDataListByInstanceUseCase(dataAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerAuthAuditTailGetUC() (any, error) {
	tailRepoInst, err := container.GetInstance[domain.TailRepository[string]](InstanceAuthAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return usecase.NewTailGetUseCase[string](tailRepoInst), nil
}

func (ucc *UseCaseContainer) providerAuthAuditTailCutUC() (any, error) {
	uwInst, err := container.GetInstance[libdom.UnitOfWork](InstanceTM)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	tailRepoInst, err := container.GetInstance[domain.TailRepository[string]](InstanceAuthAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return usecase.NewTailCutUseCase[string](uwInst, tailRepoInst), nil
}

func (ucc *UseCaseContainer) providerDataAuditTailGetUC() (any, error) {
	tailRepoInst, err := container.GetInstance[domain.TailRepository[string]](InstanceDataAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return usecase.NewTailGetUseCase[string](tailRepoInst), nil
}

func (ucc *UseCaseContainer) providerDataAuditTailCutUC() (any, error) {
	uwInst, err := container.GetInstance[libdom.UnitOfWork](InstanceTM)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	tailRepoInst, err := container.GetInstance[domain.TailRepository[string]](InstanceDataAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return usecase.NewTailCutUseCase[string](uwInst, tailRepoInst), nil
}
