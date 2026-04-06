package app

import (
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/helper"
	"github.com/ElfAstAhe/go-service-template/pkg/infra/telemetry"
	libmigr "github.com/ElfAstAhe/go-service-template/pkg/migration/goose"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
	_ "github.com/ElfAstAhe/tiny-audit-service/migrations/tiny-audit-service"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/hellofresh/health-go/v5"
	healthPg "github.com/hellofresh/health-go/v5/checks/pgx5"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func (app *App) initHelpers() error {
	var err error
	// hash cipher
	app.hashCipher = utils.NewSHA256Hash()
	// data cipher
	var dataCipherKey []byte
	dataCipherKey, err = app.hashCipher.Encrypt([]byte(app.config.App.CipherKey))
	if err != nil {
		return errs.NewCommonError("build correct data cipher key", err)
	}
	app.dataCipher, err = utils.NewAesGcmCipher(dataCipherKey)
	if err != nil {
		return errs.NewCommonError("create data cipher", err)
	}
	// data cipher helper
	app.cipherHelper = helper.NewCipherHelper(app.dataCipher)
	// keys helper
	app.keysHelper = helper.NewRSAKeysHelper(helper.RSAKey2048)
	// jwt helper
	app.jwtHelper = helper.NewJWTHelper(
		app.config.App.TokenIssuer,
		jwt.GetSigningMethod(app.config.Auth.JWTSigningMethod),
		app.config.Auth.JWTSecret,
		app.config.Auth.AccessTokenTTL,
		func() string {
			template := "tas-%v"
			rnd, err := uuid.NewRandom()
			if err != nil {
				return fmt.Sprintf(template, time.Now().Nanosecond())
			}

			return fmt.Sprintf(template, rnd.String())
		},
	)
	// jwt http helper
	app.jwtHTTPHelper = helper.NewJWTHTTPHelper(app.jwtHelper)
	// jwt gRPC helper
	app.jwtGRPCHelper = helper.NewJWTGRPCHelper(app.jwtHelper)
	// auth helper
	app.authHelper = auth.NewHelper(
		auth.DefaultHeaderName,
		auth.DefaultCookieName,
		auth.DefaultMetadataName,
		app.jwtHelper,
		app.jwtHTTPHelper,
		app.jwtGRPCHelper,
	)

	return nil
}

func (app *App) initDB() error {
	var err error
	app.db, err = postgres.NewPgDB(app.config.DB)
	if err != nil {
		return errs.NewCommonError("create db client", err)
	}

	return nil
}

func (app *App) migrateDB() error {
	migrator, err := libmigr.NewGooseDBMigrator(app.ctx, app.db, app.logger)
	if err != nil {
		return errs.NewCommonError("create migrator", err)
	}
	if err = migrator.Initialize(); err != nil {
		return errs.NewCommonError("init migrator", err)
	}
	if err = migrator.Up(); err != nil {
		return errs.NewCommonError("migrator up", err)
	}

	return nil
}

func (app *App) initTelemetry() error {
	// Вызываем нашу настройку
	shutdown, err := telemetry.SetupOTel(app.ctx, app.config.Telemetry)
	if err != nil {
		return errs.NewCommonError("failed to setup telemetry", err)
	}

	// Сохраняем shutdown в App, чтобы вызвать его в конце main
	app.telemetryShutdown = shutdown

	return nil
}

func (app *App) initMetrics() error {
	// Регистрация стандартных метрик Go (Memory, Goroutines, GC, Stack)
	// Они автоматически полетят в prometheus.DefaultRegisterer
	//if err := prometheus.Register(collectors.NewGoCollector()); err != nil {
	//    return errs.NewCommonError("failed to register go run-time collector", err)
	//}

	// Регистрация метрик процесса (CPU, Open FDs, Threads)
	//if err := prometheus.Register(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{})); err != nil {
	//	return errs.NewCommonError("failed to register process collector", err)
	//}

	if app.db != nil {
		if err := prometheus.Register(collectors.NewDBStatsCollector(app.db.GetDB(), config.AppName)); err != nil {
			return errs.NewCommonError("failed to register db stats", err)
		}
	}

	return nil
}

func (app *App) initStartupServices() error {
	// here start any startup services (workers, observers, etc.)
	// auth audit tail cutter
	if err := app.authAuditTailCutter.Start(); err != nil {
		return errs.NewCommonError("startup auth audit tail cutter", err)
	}
	// data audit tail cutter
	if err := app.dataAuditTailCutter.Start(); err != nil {
		return errs.NewCommonError("startup data audit tail cutter", err)
	}

	return nil
}

func (app *App) initHealth() error {
	healthChecker, err := health.New(health.WithComponent(health.Component{
		Name:    config.AppName,
		Version: config.AppVersion,
	}))
	if err != nil {
		return errs.NewCommonError("failed create health checker", err)
	}

	// Регистрируем Postgres. Либа сама будет делать Ping
	err = healthChecker.Register(health.Config{
		Name:      app.db.GetDriver(),
		Timeout:   app.config.DB.ConnTimeout,
		SkipOnErr: false,
		Check: healthPg.New(healthPg.Config{
			DSN: app.config.DB.DSN,
		}),
	})
	if err != nil {
		return errs.NewCommonError("failed to register pg health checker", err)
	}

	app.health = healthChecker

	return nil
}
