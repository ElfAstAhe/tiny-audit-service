package container

import (
	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	libextdom "github.com/ElfAstAhe/go-service-template/pkg/ext/domain"
	libextuc "github.com/ElfAstAhe/go-service-template/pkg/ext/usecase"
	libuc "github.com/ElfAstAhe/go-service-template/pkg/usecase"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase/trace"
)

func (ucc *UseCaseContainer) providerUnitOfWork() (any, error) {
	tmInst, err := container.GetInstance[db.TransactionManager](InstanceTM)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return libuc.NewUnitOfWork(tmInst, nil), nil
}

func (ucc *UseCaseContainer) providerAuthAuditUC() (any, error) {
	uwInst, err := container.GetInstance[libdom.UnitOfWork](InstanceUnitOfWork)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	authAuditRepoInst, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return trace.NewAuthAuditUseCase("AuthAuditUseCase", usecase.NewAuthAudit(uwInst, authAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerAuthListByPeriodUC() (any, error) {
	authAuditRepoInst, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return trace.NewAuthListByPeriodUseCase("AuthListByPeriodUseCase", usecase.NewAuthListByPeriod(authAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerAuthListByUsernameUC() (any, error) {
	authAuditRepoInst, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return trace.NewAuthListByUsernameUseCase("AuthListByUsernameUseCase", usecase.NewAuthListByUsername(authAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerDataAuditUC() (any, error) {
	uwInst, err := container.GetInstance[libdom.UnitOfWork](InstanceUnitOfWork)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	dataAuditRepoInst, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return trace.NewDataAuditUseCase("DataAuditUseCase", usecase.NewDataAudit(uwInst, dataAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerDataListByPeriodUC() (any, error) {
	dataAuditRepoInst, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return trace.NewDataListByPeriodUseCase("DataListByPeriodUseCase", usecase.NewDataListByPeriod(dataAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerDataListByInstanceUC() (any, error) {
	dataAuditRepoInst, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditTraceRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return trace.NewDataListByInstanceUseCase("DataListByInstanceUseCase", usecase.NewDataListByInstance(dataAuditRepoInst)), nil
}

func (ucc *UseCaseContainer) providerAuthAuditTailGetUC() (any, error) {
	tailRepoInst, err := container.GetInstance[libextdom.TailRepository[string]](InstanceAuthAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return libextuc.NewTailGetUseCase[string](tailRepoInst), nil
}

//goland:noinspection DuplicatedCode
func (ucc *UseCaseContainer) providerAuthAuditTailCutUC() (any, error) {
	uwInst, err := container.GetInstance[libdom.UnitOfWork](InstanceUnitOfWork)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	tailRepoInst, err := container.GetInstance[libextdom.TailRepository[string]](InstanceAuthAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return libextuc.NewTailCutUseCase[string](uwInst, tailRepoInst), nil
}

func (ucc *UseCaseContainer) providerDataAuditTailGetUC() (any, error) {
	tailRepoInst, err := container.GetInstance[libextdom.TailRepository[string]](InstanceDataAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return libextuc.NewTailGetUseCase[string](tailRepoInst), nil
}

//goland:noinspection DuplicatedCode
func (ucc *UseCaseContainer) providerDataAuditTailCutUC() (any, error) {
	uwInst, err := container.GetInstance[libdom.UnitOfWork](InstanceUnitOfWork)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}
	tailRepoInst, err := container.GetInstance[libextdom.TailRepository[string]](InstanceDataAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(ucc.GetName(), "provider: retrieve instance failed", err)
	}

	return libextuc.NewTailCutUseCase[string](uwInst, tailRepoInst), nil
}
