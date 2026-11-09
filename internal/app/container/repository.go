package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие строковые идентификаторы инстансов репозиториев, их метрик и трассирующих оберток в DI-контейнере.
const (
	// InstanceAuthAuditRepo определяет имя инстанса базового репозитория аудита аутентификации.
	InstanceAuthAuditRepo string = "authAuditRepo"

	// InstanceAuthAuditMetricsRepo определяет имя инстанса репозитория аудита аутентификации, обогащенного сбором метрик.
	InstanceAuthAuditMetricsRepo string = "authAuditMetricsRepo"

	// InstanceAuthAuditTraceRepo определяет имя инстанса репозитория аудита аутентификации, обогащенного распределенной трассировкой.
	InstanceAuthAuditTraceRepo string = "authAuditTraceRepo"

	// InstanceDataAuditRepo определяет имя инстанса базового репозитория аудита изменений данных.
	InstanceDataAuditRepo string = "dataAuditRepo"

	// InstanceDataAuditMetricsRepo определяет имя инстанса репозитория аудита изменений данных, обогащенного сбором метрик.
	InstanceDataAuditMetricsRepo string = "dataAuditMetricsRepo"

	// InstanceDataAuditTraceRepo определяет имя инстанса репозитория аудита изменений данных, обогащенного распределенной трассировкой.
	InstanceDataAuditTraceRepo string = "dataAuditTraceRepo"
)

// RepositoryContainer представляет собой ленивый (lazy) DI-контейнер для управления
// жизненным циклом и сквозной регистрацией репозиториев (базовых, метрик и трассировок) слоя доступа к данным (DAL).
type RepositoryContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*RepositoryContainer)(nil)
var _ container.LazyContainer = (*RepositoryContainer)(nil)

// NewRepositoryContainer создает и инициализирует новый экземпляр ленивого контейнера репозиториев RepositoryContainer.
func NewRepositoryContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *RepositoryContainer {
	return &RepositoryContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(RepositoryContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init выполняет первичную инициализацию контейнера RepositoryContainer, регистрируя провайдеры для всех типов репозиториев.
func (rc *RepositoryContainer) Init(ctx context.Context) error {
	err := errors.Join(
		rc.RegisterProvider(InstanceAuthAuditRepo, rc.providerAuthAuditRepo),
		rc.RegisterProvider(InstanceAuthAuditMetricsRepo, rc.providerAuthAuditMetricsRepo),
		rc.RegisterProvider(InstanceAuthAuditTraceRepo, rc.providerAuthAuditTraceRepo),
		rc.RegisterProvider(InstanceDataAuditRepo, rc.providerDataAuditRepo),
		rc.RegisterProvider(InstanceDataAuditMetricsRepo, rc.providerDataAuditMetricsRepo),
		rc.RegisterProvider(InstanceDataAuditTraceRepo, rc.providerDataAuditTraceRepo),
	)
	if err != nil {
		return errs.NewContainerError(rc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
