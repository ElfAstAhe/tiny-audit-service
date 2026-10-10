package grpc

import (
	"strings"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
)

// Константы со значениями по умолчанию для конфигурации gRPC-клиента аудита.
const (
	DefaultTarget             string        = "localhost:51052"
	DefaultWorkerCount        int           = 2
	DefaultDataCapacity       int           = 10000
	DefaultCompleteProcessing bool          = true
	DefaultStopTimeout        time.Duration = 10 * time.Second
)

// AuditClientOption определяет сигнатуру функции для конфигурирования параметров gRPC-клиента.
type AuditClientOption[D any] func(options *AuditClientOptions[D])

// AuditClientOptions содержит конфигурационные параметры для gRPC-клиента аудита.
// Включает в себя базовые опции через эмбеддинг BaseAuditClientOptions.
type AuditClientOptions[D any] struct {
	*client.BaseAuditClientOptions[D]
	Target string
	GRPC   *RawClientOptions
}

// NewAuditClientOptions создает и возвращает новый экземпляр AuditClientOptions
// с переопределенными базовыми параметрами пула и gRPC-настройками по умолчанию.
func NewAuditClientOptions[D any]() *AuditClientOptions[D] {
	res := &AuditClientOptions[D]{
		Target:                 DefaultTarget,
		GRPC:                   NewRawClientOptions(),
		BaseAuditClientOptions: client.NewBaseAuditClientOptions[D](),
	}
	res.BaseAuditClientOptions.Pool.WorkerCount = DefaultWorkerCount
	res.BaseAuditClientOptions.Pool.DataCapacity = DefaultDataCapacity
	res.BaseAuditClientOptions.Pool.CompleteProcess = DefaultCompleteProcessing
	res.BaseAuditClientOptions.Pool.StopTimeout = DefaultStopTimeout

	return res
}

// Validate проверяет корректность заполнения всех обязательных fields конфигурации.
// Возвращает ошибку, если хотя бы одно из ключевых полей не инициализировано или имеет невалидное значение.
//
//goland:noinspection DuplicatedCode
func (aco *AuditClientOptions[D]) Validate() error {
	if aco.BaseAuditClientOptions == nil || aco.Pool == nil {
		return errs.NewTlCommonError("Validate", "base audit client options and pool configuration are required", nil)
	}
	if strings.TrimSpace(aco.Pool.Name) == "" {
		return errs.NewTlCommonError("Validate", "name is required", nil)
	}
	if aco.Pool.WorkerCount <= 0 {
		return errs.NewTlCommonError("Validate", "worker count is required", nil)
	}
	if aco.Pool.DataCapacity <= 0 {
		return errs.NewTlCommonError("Validate", "data capacity is required", nil)
	}
	if aco.Pool.StopTimeout <= 0 {
		return errs.NewTlCommonError("Validate", "stop timeout is required", nil)
	}
	if utils.IsNil(aco.Pool.Logger) {
		return errs.NewTlCommonError("Validate", "logger is required", nil)
	}
	if utils.IsNil(aco.TokenProvider) {
		return errs.NewTlCommonError("Validate", "token provider is required", nil)
	}
	if strings.TrimSpace(aco.Target) == "" {
		// Исправлено: текст ошибки теперь указывает на target, а не base url
		return errs.NewTlCommonError("Validate", "target address is required", nil)
	}
	// Исправлено: защита от nil pointer паники
	if aco.GRPC == nil {
		return errs.NewTlCommonError("Validate", "grpc raw client options are required", nil)
	}
	if err := aco.GRPC.Validate(); err != nil {
		return errs.NewTlCommonError("Validate", "grpc client options validation failed", err)
	}

	return nil
}

// WithName задает уникальное имя для пула воркеров gRPC-клиента.
func WithName[D any](name string) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.Name = name
	}
}

// WithTarget конфигурирует целевой сетевой адрес удаленного gRPC-сервиса аудита.
func WithTarget[D any](target string) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Target = target
	}
}

// WithWorkerCount задает количество параллельных воркеров отправки gRPC-сообщений.
func WithWorkerCount[D any](count int) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.WorkerCount = count
	}
}

// WithDataCapacity устанавливает максимальную емкость буфера очереди отправки пула.
func WithDataCapacity[D any](capacity int) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.DataCapacity = capacity
	}
}

// WithCompleteProcess определяет, нужно ли досылать оставшиеся в буфере логи при остановке пула воркеров.
func WithCompleteProcess[D any](flag bool) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.CompleteProcess = flag
	}
}

// WithStopTimeout задает максимальное время ожидания Graceful Shutdown для пула воркеров gRPC-клиента.
func WithStopTimeout[D any](timeout time.Duration) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.StopTimeout = timeout
	}
}

// WithTokenProvider конфигурирует провайдер авторизационных токенов для gRPC-запросов.
func WithTokenProvider[D any](provider auth.TokenProvider) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.TokenProvider = provider
	}
}

// WithLogger настраивает системный логгер для пула воркеров gRPC-клиента.
func WithLogger[D any](log logger.Logger) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.Logger = log
	}
}

// WithSecure управляет использованием защищенного TLS-соединения для gRPC-клиента.
func WithSecure[D any](flag bool) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.GRPC.Secure = flag
	}
}

// WithConnectionTimeout задает максимальное время ожидания при установке gRPC-соединения.
func WithConnectionTimeout[D any](timeout time.Duration) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.GRPC.ConnTimeout = timeout
	}
}

// WithKeepAliveTime определяет интервал отправки пингов KeepAlive для проверки активности соединения.
func WithKeepAliveTime[D any](duration time.Duration) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.GRPC.KATime = duration
	}
}

// WithKeepAliveTimeout задает время ожидания ответа на пинг KeepAlive перед закрытием соединения как неактивного.
func WithKeepAliveTimeout[D any](timeout time.Duration) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.GRPC.KATimeout = timeout
	}
}

// WithKeepAlivePermitWithoutStream разрешает отправку KeepAlive пингов даже при отсутствии активных RPC-стримов.
func WithKeepAlivePermitWithoutStream[D any](flag bool) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.GRPC.KAPermitWithoutStream = flag
	}
}
