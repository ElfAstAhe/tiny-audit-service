package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие строковые идентификаторы инстансов бизнес-сценариев (Use Cases) и Unit of Work в DI-контейнере.
const (
	// InstanceUnitOfWork определяет имя инстанса паттерна Unit of Work для управления транзакциями.
	InstanceUnitOfWork string = "unit-of-work"

	// InstanceAuthAuditUC определяет имя инстанса UseCase регистрации событий аудита аутентификации.
	InstanceAuthAuditUC string = "AuthAuditUC"

	// InstanceAuthListByPeriodUC определяет имя инстанса UseCase получения списка логов аутентификации за период.
	InstanceAuthListByPeriodUC string = "AuthListByPeriodUC"

	// InstanceAuthListByUsernameUC определяет имя инстанса UseCase получения логов аутентификации конкретного пользователя.
	InstanceAuthListByUsernameUC string = "AuthListByUsernameUC"

	// InstanceDataAuditUC определяет имя инстанса UseCase регистрации событий аудита изменений данных.
	InstanceDataAuditUC string = "DataAuditUC"

	// InstanceDataListByPeriodUC определяет имя инстанса UseCase получения списка логов изменений данных за период.
	InstanceDataListByPeriodUC string = "DataListByPeriodUC"

	// InstanceDataListByInstanceUC определяет имя инстанса UseCase получения логов изменений конкретного инстанса данных.
	InstanceDataListByInstanceUC string = "DataListByInstanceUC"

	// InstanceAuthAuditTailGetUC определяет имя инстанса UseCase извлечения устаревших "хвостов" таблиц аудита аутентификации.
	InstanceAuthAuditTailGetUC string = "AuthAuditTailGetUC"

	// InstanceAuthAuditTailCutUC определяет имя инстанса UseCase удаления устаревших "хвостов" таблиц аудита аутентификации.
	InstanceAuthAuditTailCutUC string = "AuthAuditTailCutUC"

	// InstanceDataAuditTailGetUC определяет имя инстанса UseCase извлечения устаревших "хвостов" таблиц аудита данных.
	InstanceDataAuditTailGetUC string = "DataAuditTailGetUC"

	// InstanceDataAuditTailCutUC определяет имя инстанса UseCase удаления устаревших "хвостов" таблиц аудита данных.
	InstanceDataAuditTailCutUC string = "DataAuditTailCutUC"
)

// UseCaseContainer представляет собой ленивый (lazy) DI-контейнер для конфигурации,
// регистрации и управления жизненным циклом прикладных сценариев и интеракторов бизнес-логики (BLL).
type UseCaseContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*UseCaseContainer)(nil)
var _ container.LazyContainer = (*UseCaseContainer)(nil)

// NewUseCaseContainer создает и инициализирует новый экземпляр ленивого контейнера сценариев UseCaseContainer.
func NewUseCaseContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *UseCaseContainer {
	return &UseCaseContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(UseCaseContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init выполняет первичную инициализацию контейнера UseCaseContainer, регистрируя провайдеры для Unit of Work и всех UseCase.
func (ucc *UseCaseContainer) Init(ctx context.Context) error {
	err := errors.Join(
		ucc.RegisterProvider(InstanceUnitOfWork, ucc.providerUnitOfWork),
		ucc.RegisterProvider(InstanceAuthAuditUC, ucc.providerAuthAuditUC),
		ucc.RegisterProvider(InstanceAuthListByPeriodUC, ucc.providerAuthListByPeriodUC),
		ucc.RegisterProvider(InstanceAuthListByUsernameUC, ucc.providerAuthListByUsernameUC),
		ucc.RegisterProvider(InstanceDataAuditUC, ucc.providerDataAuditUC),
		ucc.RegisterProvider(InstanceDataListByPeriodUC, ucc.providerDataListByPeriodUC),
		ucc.RegisterProvider(InstanceDataListByInstanceUC, ucc.providerDataListByInstanceUC),
		ucc.RegisterProvider(InstanceAuthAuditTailGetUC, ucc.providerAuthAuditTailGetUC),
		ucc.RegisterProvider(InstanceAuthAuditTailCutUC, ucc.providerAuthAuditTailCutUC),
		ucc.RegisterProvider(InstanceDataAuditTailGetUC, ucc.providerDataAuditTailGetUC),
		ucc.RegisterProvider(InstanceDataAuditTailCutUC, ucc.providerDataAuditTailCutUC),
	)
	if err != nil {
		return errs.NewContainerError(ucc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
