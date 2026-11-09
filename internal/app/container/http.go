package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие строковые идентификаторы инстансов HTTP-роутера и раннера сервера в DI-контейнере.
const (
	// InstanceHTTPRouter определяет имя инстанса прикладного HTTP-роутера.
	InstanceHTTPRouter string = "HTTPRouter"

	// InstanceHTTPRunner определяет имя инстанса запускаемого HTTP-сервера (Runner).
	InstanceHTTPRunner string = "HTTPRunner"
)

// HTTPContainer представляет собой ленивый (lazy) DI-контейнер для управления
// жизненным циклом, конфигурацией роутеров и запуском основного веб-сервера приложения.
type HTTPContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*HTTPContainer)(nil)
var _ container.LazyContainer = (*HTTPContainer)(nil)

// NewHTTPContainer создает и инициализирует новый экземпляр ленивого контейнера HTTPContainer.
func NewHTTPContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *HTTPContainer {
	return &HTTPContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(HTTPContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init выполняет первичную инициализацию контейнера HTTPContainer, регистрируя провайдеры для роутера и HTTP-раннера.
func (hc *HTTPContainer) Init(initCtx context.Context) error {
	err := errors.Join(
		hc.RegisterProvider(InstanceHTTPRouter, hc.providerChiRouter),
		hc.RegisterProvider(InstanceHTTPRunner, hc.providerHTTPRunner),
	)
	if err != nil {
		return errs.NewContainerError(hc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
