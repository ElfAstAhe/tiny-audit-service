package container

import (
	"context"
	"errors"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
)

// Константы, определяющие строковые идентификаторы инстансов криптографических утилит, хелперов и JWT-инструментов в DI-контейнере.
const (
	// InstanceHashCipher определяет имя инстанса утилиты хеширования SHA256.
	InstanceHashCipher string = "sha256-cipher"

	// InstanceDataCipher определяет имя инстанса утилиты симметричного шифрования AES-GCM.
	InstanceDataCipher string = "aes-gcm-cipher"

	// InstanceDataCipherHelper определяет имя инстанса хелпера для шифрования данных.
	InstanceDataCipherHelper string = "data-cipher-helper"

	// InstanceKeysHelper определяет имя инстанса хелпера для управления асимметричными ключами RSA-2048.
	InstanceKeysHelper string = "rsa-2048-keys-helper"

	// InstanceJWTHelper определяет имя инстанса базового хелпера для работы с JWT-токенами.
	InstanceJWTHelper string = "jwt-helper"

	// InstanceJWTHTTPHelper определяет имя инстанса хелпера работы с JWT-токенами в HTTP-запросах.
	InstanceJWTHTTPHelper string = "jwt-http-helper"

	// InstanceJWTGRPCHelper определяет имя инстанса хелпера работы с JWT-токенами в gRPC-контексте.
	InstanceJWTGRPCHelper string = "jwt-grpc-helper"

	// InstanceAuthHelper определяет имя инстанса хелпера проверки авторизации субъектов.
	InstanceAuthHelper string = "auth-helper"
)

// ToolsContainer представляет собой ленивый (lazy) DI-контейнер для управления
// жизненным циклом криптографических инструментов, шифровщиков, генераторов ключей и JWT-хелперов.
type ToolsContainer struct {
	*container.BaseLazyContainer
}

var _ container.Container = (*ToolsContainer)(nil)
var _ container.LazyContainer = (*ToolsContainer)(nil)

// NewToolsContainer создает и инициализирует новый экземпляр ленивого контейнера инструментов ToolsContainer.
func NewToolsContainer(
	orchestrator container.Orchestrator,
	log logger.Logger,
) *ToolsContainer {
	return &ToolsContainer{
		BaseLazyContainer: container.NewBaseLazyContainer(
			container.WithLazyName(ToolsContainerName),
			container.WithLazyOrchestrator(orchestrator),
			container.WithLazyLogger(log),
		),
	}
}

// Init выполняет первичную инициализацию контейнера ToolsContainer, регистрируя провайдеры для всех криптографических хелперов и шифров.
func (tc *ToolsContainer) Init(ctx context.Context) error {
	err := errors.Join(
		tc.RegisterProvider(InstanceHashCipher, tc.providerHashCipher),
		tc.RegisterProvider(InstanceDataCipher, tc.providerDataCipher),
		tc.RegisterProvider(InstanceDataCipherHelper, tc.providerDataCipherHelper),
		tc.RegisterProvider(InstanceKeysHelper, tc.providerKeysHelper),
		tc.RegisterProvider(InstanceJWTHelper, tc.providerJWTHelper),
		tc.RegisterProvider(InstanceJWTHTTPHelper, tc.providerJWTHTTPHelper),
		tc.RegisterProvider(InstanceJWTGRPCHelper, tc.providerJWTGRPCHelper),
		tc.RegisterProvider(InstanceAuthHelper, tc.providerAuthHelper),
	)
	if err != nil {
		return errs.NewContainerError(tc.GetName(), "container init: register providers failed", err)
	}

	return nil
}
