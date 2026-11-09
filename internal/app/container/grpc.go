package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие строковые идентификаторы инстансов gRPC-сервисов и раннера сервера в DI-контейнере.
const (
	// InstanceAuthAuditGRPCService определяет имя инстанса gRPC-сервиса аудита аутентификации.
	InstanceAuthAuditGRPCService string = "authAuditGRPCService"

	// InstanceDataAuditGRPCService определяет имя инстанса gRPC-сервиса аудита изменений данных.
	InstanceDataAuditGRPCService string = "dataAuditGRPCService"

	// InstanceGRPCRunner определяет имя инстанса запускаемого gRPC-сервера (Runner).
	InstanceGRPCRunner string = "grpcRunner"
)

// GRPCContainer представляет собой ленивый (lazy) DI-контейнер для управления
// жизненным циклом, регистрацией gRPC-сервисов приложения и запуском gRPC-сервера.
type GRPCContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*GRPCContainer)(nil)
var _ container.LazyContainer = (*GRPCContainer)(nil)

// NewGRPCContainer создает и инициализирует новый экземпляр ленивого контейнера GRPCContainer.
func NewGRPCContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *GRPCContainer {
	return &GRPCContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(GRPCContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init выполняет первичную инициализацию контейнера GRPCContainer, регистрируя провайдеры для gRPC-сервисов и раннера.
func (gc *GRPCContainer) Init(ctx context.Context) error {
	err := errors.Join(
		gc.RegisterProvider(InstanceAuthAuditGRPCService, gc.providerAuthAuditGRPCService),
		gc.RegisterProvider(InstanceDataAuditGRPCService, gc.providerDataAuditGRPCService),
		gc.RegisterProvider(InstanceGRPCRunner, gc.providerGRPCRunner),
	)
	if err != nil {
		return errs.NewContainerError(gc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
