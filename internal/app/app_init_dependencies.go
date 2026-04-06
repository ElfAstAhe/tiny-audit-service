package app

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/metrics"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/trace"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase/telemetry"
)

//goland:noinspection DuplicatedCode
func (app *App) initDependencies() error {
	var err error
	// transaction manager
	app.tm = db.NewTxManager(app.db)
	var (
		authAuditRepository *postgres.AuthAuditPgRepository
		dataAuditRepository *postgres.DataAuditPgRepository

		authAuditRepo domain.AuthAuditRepository
		dataAuditRepo domain.DataAuditRepository

		authAuditUC          usecase.AuthAuditUseCase
		authListByPeriodUC   usecase.AuthListByPeriodUseCase
		authListByUsernameUC usecase.AuthListByUsernameUseCase

		dataAuditUC          usecase.DataAuditUseCase
		dataListByPeriodUC   usecase.DataListByPeriodUseCase
		dataListByInstanceUC usecase.DataListByInstanceUseCase

		authAuditTailGetUC usecase.TailGetUseCase[string]
		authAuditTailCutUC usecase.TailCutUseCase[string]
		dataAuditTailGetUC usecase.TailGetUseCase[string]
		dataAuditTailCutUC usecase.TailCutUseCase[string]
	)
	// repositories
	{
		authAuditRepository, err = postgres.NewAuthAuditPgRepository(app.db, app.db)
		if err != nil {
			return err
		}
		authAuditRepo = trace.NewAuthAuditTraceRepository(metrics.NewAuthAuditMetricsRepository(authAuditRepository))

		dataAuditRepository, err = postgres.NewDataAuditPgRepository(app.db, app.db)
		if err != nil {
			return err
		}
		dataAuditRepo = trace.NewDataAuditTraceRepository(metrics.NewDataAuditMetricsRepository(dataAuditRepository))
	}
	// use cases
	{
		authAuditUC = telemetry.NewAuthAuditUseCase("AuthAuditUseCase", usecase.NewAuthAuditUseCase(app.tm, authAuditRepo))
		authListByPeriodUC = telemetry.NewAuthListByPeriodUseCase("AuthListByPeriodUseCase", usecase.NewAuthListByPeriodUseCase(authAuditRepo))
		authListByUsernameUC = telemetry.NewAuthListByUsernameUseCase("AuthListByUsernameUseCase", usecase.NewAuthListByUsernameUseCase(authAuditRepo))

		dataAuditUC = telemetry.NewDataAuditUseCase("DataAuditUseCase", usecase.NewDataAuditUseCase(app.tm, dataAuditRepo))
		dataListByPeriodUC = telemetry.NewDataListByPeriodUseCase("DataListByPeriodUseCase", usecase.NewDataListByPeriodUseCase(dataAuditRepo))
		dataListByInstanceUC = telemetry.NewDataListByInstanceUseCase("", usecase.NewDataListByInstanceUseCase(dataAuditRepo))

		authAuditTailGetUC = usecase.NewTailGetUseCase[string](authAuditRepository)
		authAuditTailCutUC = usecase.NewTailCutUseCase[string](authAuditRepository)

		dataAuditTailGetUC = usecase.NewTailGetUseCase[string](dataAuditRepository)
		dataAuditTailCutUC = usecase.NewTailCutUseCase[string](dataAuditRepository)
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
	// workers, observers, etc
	{
		app.authAuditTailCutter = worker.NewTailCutter(
			app.ctx,
			"auth",
			worker.NewTailCutterConfig(
				3*time.Second,
				app.config.App.AuthTailJobRepeatDuration,
				2,
				128,
				app.config.App.AuthTailDuration,
				app.config.App.AuthTailCut,
			),
			authAuditTailGetUC,
			authAuditTailCutUC,
			app.logger,
		)
		app.dataAuditTailCutter = worker.NewTailCutter(
			app.ctx,
			"data",
			worker.NewTailCutterConfig(
				3*time.Second,
				app.config.App.DataTailJobRepeatDuration,
				2,
				128,
				app.config.App.DataTailDuration,
				app.config.App.DataTailCut,
			),
			dataAuditTailGetUC,
			dataAuditTailCutUC,
			app.logger,
		)
	}

	return nil
}
