package container

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие строковые идентификаторы инстансов ядра приложения, конфигурации и логгера в DI-контейнере.
const (
	// InstanceApplication определяет имя инстанса самого приложения в контейнере.
	InstanceApplication string = "application"

	// InstanceApplicationReady определяет имя инстанса функции или флага готовности приложения к работе.
	InstanceApplicationReady string = "application-ready"

	// InstanceApplicationHealth определяет имя инстанса проверки общего статуса здоровья приложения.
	InstanceApplicationHealth string = "application-health"

	// InstanceConfig определяет имя инстанса корневого объекта конфигурации приложения.
	InstanceConfig string = "config"

	// InstanceLogger определяет имя инстанса системного логгера.
	InstanceLogger string = "logger"
)

// AppContainer представляет собой контейнер объектов ядра приложения (App Container),
// управляющий регистрацией конфигурации, логгеров и базовых флагов жизненного цикла.
type AppContainer struct {
	*container.BaseContainer
}

var _ container.Container = (*AppContainer)(nil)

// NewAppContainer создает и инициализирует новый экземпляр базового контейнера AppContainer.
func NewAppContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *AppContainer {
	return &AppContainer{
		BaseContainer: container.NewBaseContainer(
			container.WithName(AppContainerName),
			container.WithOrchestrator(orchestrator),
			container.WithLogger(log),
		),
	}
}

// Init выполняет первичную инициализацию контейнера AppContainer.
func (ac *AppContainer) Init(ctx context.Context) error {
	return nil
}

// Close выполняет корректное завершение работы контейнера и освобождение занятых им ресурсов.
func (ac *AppContainer) Close(closeCtx context.Context) error {
	return nil
}
