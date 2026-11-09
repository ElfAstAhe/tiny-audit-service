package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/migration"
)

// Константы, определяющие строковые идентификаторы инстансов базы данных, мигратора и менеджера транзакций в DI-контейнере.
const (
	// InstanceDB определяет имя инстанса подключения к базе данных.
	InstanceDB string = "DB"

	// InstanceDBMigrator определяет имя инстанса утилиты выполнения миграций БД.
	InstanceDBMigrator string = "db-migrator"

	// InstanceTM определяет имя инстанса менеджера транзакций (Transaction Manager).
	InstanceTM string = "transaction-manager"
)

// PgContainer представляет собой ленивый (lazy) DI-контейнер для управления
// жизненным циклом пула соединений PostgreSQL, пингования БД и автоматического наката миграций.
type PgContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*PgContainer)(nil)
var _ container.LazyContainer = (*PgContainer)(nil)

// NewPgContainer создает и инициализирует новый экземпляр контейнера базы данных PgContainer.
func NewPgContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *PgContainer {
	return &PgContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(DBContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init выполняет инициализацию контейнера: регистрирует провайдеры, пингует базу данных для проверки
// связи, инициализирует мигратор и накатывает новые SQL-скрипты изменений (Migrate Up).
func (pc *PgContainer) Init(initCtx context.Context) error {
	// add providers
	err := errors.Join(
		pc.RegisterProvider(InstanceDB, pc.providerDB),
		pc.RegisterProvider(InstanceDBMigrator, pc.providerDBMigrator),
		pc.RegisterProvider(InstanceTM, pc.providerTM),
	)
	if err != nil {
		return errs.NewContainerError(pc.GetName(), "container init: register providers failed", err)
	}
	// init db instance
	dbInst, err := container.GetInstance[db.DB](InstanceDB)
	if err != nil {
		return errs.NewContainerError(pc.GetName(), "container init: init db failed", err)
	}
	// check db connection
	err = dbInst.Ping(initCtx)
	if err != nil {
		return errs.NewContainerError(pc.GetName(), "container init: check db failed", err)
	}
	// data migration
	migrator, err := container.GetInstance[migration.Migrator](InstanceDBMigrator)
	if err != nil {
		return errs.NewContainerError(pc.GetName(), "container init: init migrator failed", err)
	}
	// migrate up
	err = migrator.Up(initCtx)
	if err != nil {
		return errs.NewContainerError(pc.GetName(), "container init: up migrator failed", err)
	}

	return nil
}
