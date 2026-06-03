package container

import (
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/metrics"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/trace"
)

//goland:noinspection DuplicatedCode
func (rc *RepositoryContainer) providerAuthAuditRepo() (any, error) {
	dbInst, err := container.GetInstance[*postgres.PgDB](InstanceDB)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), "provider: retrieve instance failed", err)
	}
	res, err := postgres.NewAuthAuditPgRepository(dbInst, dbInst)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceAuthAuditRepo), err)
	}

	return trace.NewAuthAuditTraceRepository(metrics.NewAuthAuditMetricsRepository(res)), nil
}

//goland:noinspection DuplicatedCode
func (rc *RepositoryContainer) providerDataAuditRepo() (any, error) {
	dbInst, err := container.GetInstance[*postgres.PgDB](InstanceDB)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), "provider: retrieve instance failed", err)
	}
	res, err := postgres.NewDataAuditPgRepository(dbInst, dbInst)
	if err != nil {
		return nil, errs.NewContainerError(rc.GetName(), fmt.Sprintf("provider: create %s instance failed", InstanceDataAuditRepo), err)
	}

	return trace.NewDataAuditTraceRepository(metrics.NewDataAuditMetricsRepository(res)), nil
}
