package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/tiny-audit-service/internal/config"
	"github.com/hellofresh/health-go/v5"
	healthPg "github.com/hellofresh/health-go/v5/checks/pgx5"
)

// InstanceHealthStatus определяет строковый идентификатор инстанса службы проверки здоровья (Health) в DI-контейнере.
const (
	InstanceHealthStatus string = "health-status"
)

// ServiceContainer представляет собой ленивый (lazy) DI-контейнер для конфигурации
// и инициализации внутренних инфраструктурных служб, таких как проверки работоспособности (Health Checks).
type ServiceContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*ServiceContainer)(nil)
var _ container.LazyContainer = (*ServiceContainer)(nil)

// NewServiceContainer создает и инициализирует новый экземпляр ленивого контейнера служб ServiceContainer.
func NewServiceContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *ServiceContainer {
	return &ServiceContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(ServiceContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init выполняет первичную инициализацию контейнера ServiceContainer, регистрируя провайдеры
// и настраивая автоматические проверки доступности базы данных (PostgreSQL) в системе Health Check.
func (sc *ServiceContainer) Init(initCtx context.Context) error {
	// add providers
	err := errors.Join(
		sc.RegisterProvider(InstanceHealthStatus, sc.providerHealthStatus),
	)
	if err != nil {
		return errs.NewContainerError(sc.GetName(), "container init: register providers failed", err)
	}
	confInst, err := container.GetInstance[*config.Config](InstanceConfig)
	if err != nil {
		return errs.NewContainerError(sc.GetName(), "container init: retrieve instance failed", err)
	}
	// init health checks
	dbInst, err := container.GetInstance[db.DB](InstanceDB)
	if err != nil {
		return errs.NewContainerError(sc.GetName(), "container init: retrieve instance failed", err)
	}
	healthStatusInst, err := container.GetInstance[*health.Health](InstanceHealthStatus)
	if err != nil {
		return errs.NewContainerError(sc.GetName(), "container init: retrieve instance failed", err)
	}
	// Регистрируем Postgres. Либа сама будет делать Ping
	err = healthStatusInst.Register(health.Config{
		Name:      dbInst.GetDriver(),
		Timeout:   confInst.DB.ConnTimeout,
		SkipOnErr: false,
		Check: healthPg.New(healthPg.Config{
			DSN: confInst.DB.DSN,
		}),
	})
	if err != nil {
		return errs.NewContainerError(sc.GetName(), "container init: register health status check failed", err)
	}

	return nil
}
