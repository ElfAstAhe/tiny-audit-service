package container

import (
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/metrics"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/trace"
)

//goland:noinspection DuplicatedCode
func (rc *RepositoryContainer) providerAuthAuditRepo() (any, error) {
	dbInst, err := container.GetInstance[db.DB](InstanceDB)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), "provider: retrieve instance failed", err)
	}
	res, err := repository.NewAuthAudit(dbInst.GetDriver(), dbInst, dbInst)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceAuthAuditRepo), err)
	}

	return res, nil
}

func (rc *RepositoryContainer) providerAuthAuditMetricsRepo() (any, error) {
	repo, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), "provider: retrieve instance failed", err)
	}

	return metrics.NewAuthAuditRepository(repo), nil
}

func (rc *RepositoryContainer) providerAuthAuditTraceRepo() (any, error) {
	repo, err := container.GetInstance[domain.AuthAuditRepository](InstanceAuthAuditMetricsRepo)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), "provider: retrieve instance failed", err)
	}

	return trace.NewAuthAuditRepository(repo), nil
}

//goland:noinspection DuplicatedCode
func (rc *RepositoryContainer) providerDataAuditRepo() (any, error) {
	dbInst, err := container.GetInstance[db.DB](InstanceDB)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), "provider: retrieve instance failed", err)
	}
	res, err := repository.NewDataAudit(dbInst.GetDriver(), dbInst, dbInst)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceDataAuditRepo), err)
	}

	return res, nil
}

func (rc *RepositoryContainer) providerDataAuditMetricsRepo() (any, error) {
	repo, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditRepo)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), "provider: retrieve instance failed", err)
	}

	return metrics.NewDataAuditRepository(repo), nil
}

func (rc *RepositoryContainer) providerDataAuditTraceRepo() (any, error) {
	repo, err := container.GetInstance[domain.DataAuditRepository](InstanceDataAuditMetricsRepo)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), "provider: retrieve instance failed", err)
	}

	return trace.NewDataAuditRepository(repo), nil
}
