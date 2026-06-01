package container

import (
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/migration/goose"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
	_ "github.com/ElfAstAhe/tiny-audit-service/migrations/tiny-audit-service"
)

func (pc *PgContainer) providerDB(name string) (any, error) {
	appCnt, err := pc.GetOrchestrator().GetContainer(AppContainerName)
	if err != nil {
		return nil, errs.NewContainerError(pc.GetName(), "provider: retrieve container failed", err)
	}
	confInst, err := container.GetInstance[*config.Config](appCnt, InstanceConfig)
	if err != nil {
		return nil, errs.NewContainerError(pc.GetName(), "provider: retrieve instance failed", err)
	}
	res, err := postgres.NewPgDB(confInst.DB)
	if err != nil {
		return nil, errs.NewContainerError(pc.GetName(), fmt.Sprintf("provider: create %s instance failed", name), err)
	}

	return res, nil
}

func (pc *PgContainer) providerDBMigrator(name string) (any, error) {
	appCnt, err := pc.GetOrchestrator().GetContainer(AppContainerName)
	if err != nil {
		return nil, errs.NewContainerError(pc.GetName(), "provider: retrieve container failed", err)
	}
	logInst, err := container.GetInstance[logger.Logger](appCnt, InstanceLogger)
	if err != nil {
		return nil, errs.NewContainerError(pc.GetName(), "provider: retrieve instance failed", err)
	}
	dbInst, err := container.GetInstance[*postgres.PgDB](pc, InstanceDB)
	if err != nil {
		return nil, errs.NewContainerError(pc.GetName(), "provider: retrieve instance failed", err)
	}
	res, err := goose.NewDBMigrator(dbInst, logInst)
	if err != nil {
		return nil, errs.NewContainerError(pc.GetName(), fmt.Sprintf("provider: create %s instance failed", name), err)
	}
	err = res.Initialize()
	if err != nil {
		return nil, errs.NewContainerError(pc.GetName(), fmt.Sprintf("provider: initialize %s instance failed", name), err)
	}

	return res, nil
}
