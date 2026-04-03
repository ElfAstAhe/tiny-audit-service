package app

import (
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/metrics"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/trace"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase/telemetry"
)

//goland:noinspection DuplicatedCode
func (app *App) initDependencies() error {
	var err error
	// transaction manager
	app.tm = db.NewTxManager(app.db)
	var (
		authAuditRepo domain.AuthAuditRepository
		dataAuditRepo domain.DataAuditRepository

		authAuditUC          usecase.AuthAuditUseCase
		authListByPeriodUC   usecase.AuthListByPeriodUseCase
		authListByUsernameUC usecase.AuthListByUsernameUseCase

		dataAuditUC          usecase.DataAuditUseCase
		dataListByPeriodUC   usecase.DataListByPeriodUseCase
		dataListByInstanceUC usecase.DataListByInstanceUseCase
	)
	// repositories
	{
		authAuditRepo, err = postgres.NewAuthAuditPgRepository(app.db, app.db)
		if err != nil {
			return err
		}
		authAuditRepo = trace.NewAuthAuditTraceRepository(metrics.NewAuthAuditMetricsRepository(authAuditRepo))

		dataAuditRepo, err = postgres.NewDataAuditPgRepository(app.db, app.db)
		if err != nil {
			return err
		}
		dataAuditRepo = trace.NewDataAuditTraceRepository(metrics.NewDataAuditMetricsRepository(dataAuditRepo))
	}
	// use cases
	{
		authAuditUC = telemetry.NewAuthAuditUseCase("AuthAuditUseCase", usecase.NewAuthAuditUseCase(app.tm, authAuditRepo))
		authListByPeriodUC = telemetry.NewAuthListByPeriodUseCase("AuthListByPeriodUseCase", usecase.NewAuthListByPeriodUseCase(authAuditRepo))
		authListByUsernameUC = telemetry.NewAuthListByUsernameUseCase("AuthListByUsernameUseCase", usecase.NewAuthListByUsernameUseCase(authAuditRepo))

		dataAuditUC = telemetry.NewDataAuditUseCase("DataAuditUseCase", usecase.NewDataAuditUseCase(app.tm, dataAuditRepo))
		dataListByPeriodUC = telemetry.NewDataListByPeriodUseCase("DataListByPeriodUseCase", usecase.NewDataListByPeriodUseCase(dataAuditRepo))
		dataListByInstanceUC = telemetry.NewDataListByInstanceUseCase("", usecase.NewDataListByInstanceUseCase(dataAuditRepo))
	}
	// facade
	{
		app.authFacade = facade.NewAuthAuditFacade(
			app.authHelper,
			authAuditUC,
			authListByPeriodUC,
			authListByUsernameUC,
		)
		app.dataFacade = facade.NewDataAuditFacade(
			app.authHelper,
			dataAuditUC,
			dataListByPeriodUC,
			dataListByInstanceUC,
		)
	}

	return nil
}
