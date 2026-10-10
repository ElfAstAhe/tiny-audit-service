package rest

import (
	"net/url"
	"strings"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
)

// Константы со значениями по умолчанию для конфигурации REST-клиента аудита.
const (
	DefaultBaseURL            string        = "http://localhost:8080/"
	DefaultTimeout            time.Duration = 5 * time.Second
	DefaultWorkerCount        int           = 2
	DefaultDataCapacity       int           = 10000
	DefaultCompleteProcessing bool          = true
	DefaultStopTimeout        time.Duration = 10 * time.Second
)

// AuditClientOption определяет сигнатуру функции для конфигурирования параметров REST-клиента.
type AuditClientOption[D any] func(options *AuditClientOptions[D])

// AuditClientOptions содержит конфигурационные параметры для REST-клиента аудита.
// Включает в себя базовые опции через эмбеддинг BaseAuditClientOptions.
type AuditClientOptions[D any] struct {
	*client.BaseAuditClientOptions[D]
	BaseURL     string
	ReadTimeout time.Duration
}

// NewAuditClientOptions создает и возвращает новый экземпляр AuditClientOptions
// с переопределенными базовыми параметрами пула и REST-настройками по умолчанию.
func NewAuditClientOptions[D any]() *AuditClientOptions[D] {
	res := &AuditClientOptions[D]{
		BaseURL:                DefaultBaseURL,
		ReadTimeout:            DefaultTimeout,
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
	if aco.ReadTimeout <= 0 {
		return errs.NewTlCommonError("Validate", "read timeout is required", nil)
	}
	if utils.IsNil(aco.TokenProvider) {
		return errs.NewTlCommonError("Validate", "token provider is required", nil)
	}
	if strings.TrimSpace(aco.BaseURL) == "" {
		return errs.NewTlCommonError("Validate", "base url is required", nil)
	}
	if _, err := url.Parse(aco.BaseURL); err != nil {
		return errs.NewTlCommonError("Validate", "base url is invalid", err)
	}

	return nil
}

// WithName задает уникальное имя для пула воркеров REST-клиента.
func WithName[D any](name string) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.Name = name
	}
}

// WithBaseURL конфигурирует базовый URL адрес удаленного сервиса аудита.
func WithBaseURL[D any](url string) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.BaseURL = url
	}
}

// WithReadTimeout задает таймаут на чтение/выполнение HTTP-запросов.
func WithReadTimeout[D any](timeout time.Duration) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.ReadTimeout = timeout
	}
}

// WithWorkerCount задает количество параллельных воркеров отправки REST-сообщений.
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

// WithCompleteProcess определяет, нужно ли досылать оставшиеся в буфере логи при остановке пула.
func WithCompleteProcess[D any](flag bool) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.CompleteProcess = flag
	}
}

// WithStopTimeout задает максимальное время ожидания Graceful Shutdown для пула воркеров.
func WithStopTimeout[D any](timeout time.Duration) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.StopTimeout = timeout
	}
}

// WithTokenProvider конфигурирует провайдер авторизационных токенов для HTTP-запросов.
func WithTokenProvider[D any](provider auth.TokenProvider) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.TokenProvider = provider
	}
}

// WithLogger настраивает системный логгер для пула воркеров REST-клиента.
func WithLogger[D any](log logger.Logger) AuditClientOption[D] {
	return func(options *AuditClientOptions[D]) {
		options.Pool.Logger = log
	}
}
