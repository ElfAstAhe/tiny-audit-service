package container

import (
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase/telemetry"
)

func (ucc *UseCaseContainer) providerTM() (any, error) {
	dbInst, err := container.GetInstance[db.DB](InstanceDB)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return db.NewTxManager(dbInst), nil
}

func (ucc *UseCaseContainer) providerAuthAuditUC() (any, error) {
	tmInst, err := container.GetInstance[db.TransactionManager](InstanceTM)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	authAuditRepoInst, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return telemetry.NewAuthAuditUseCase("AuthAuditUseCase", usecase.NewAuthAuditUseCase(tmInst, authAuditRepoInst)), nil
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
	tmInst, err := container.GetInstance[db.TransactionManager](InstanceTM)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	dataAuditRepoInst, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return telemetry.NewDataAuditUseCase("DataAuditUseCase", usecase.NewDataAuditUseCase(tmInst, dataAuditRepoInst)), nil
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
	tailRepoInst, err := container.GetInstance[domain.TailRepository[string]](InstanceAuthAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return usecase.NewTailCutUseCase[string](tailRepoInst), nil
}

func (ucc *UseCaseContainer) providerDataAuditTailGetUC() (any, error) {
	tailRepoInst, err := container.GetInstance[domain.TailRepository[string]](InstanceDataAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return usecase.NewTailGetUseCase[string](tailRepoInst), nil
}

func (ucc *UseCaseContainer) providerDataAuditTailCutUC() (any, error) {
	tailRepoInst, err := container.GetInstance[domain.TailRepository[string]](InstanceDataAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return usecase.NewTailCutUseCase[string](tailRepoInst), nil
}
