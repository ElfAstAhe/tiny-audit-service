package client

import (
	"strings"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
)

// Константы со значениями по умолчанию для конфигурации пула воркеров клиента аудита.
const (
	DefaultWorkerCount     int  = 2
	DefaultDataCapacity    int  = 32
	DefaultCompleteProcess bool = true
)

// BaseAuditClientOption определяет сигнатуру функции для конфигурирования опций клиента аудита.
type BaseAuditClientOption[D any] func(options *BaseAuditClientOptions[D])

// BaseAuditClientOptions содержит конфигурационные параметры для базового клиента аудита.
// Использует generic-тип D для определения структуры передаваемых данных аудита.
type BaseAuditClientOptions[D any] struct {
	Pool          *worker.BasePoolOptions[D]
	TokenProvider auth.TokenProvider
	AuditAction   AuditAction[D]
}

// NewBaseAuditClientOptions создает и возвращает новый экземпляр BaseAuditClientOptions
// с заполненными значениями по умолчанию.
func NewBaseAuditClientOptions[D any]() *BaseAuditClientOptions[D] {
	res := &BaseAuditClientOptions[D]{
		Pool: worker.NewBasePoolOptions[D](),
	}
	res.Pool.WorkerCount = DefaultWorkerCount
	res.Pool.DataCapacity = DefaultDataCapacity
	res.Pool.CompleteProcess = DefaultCompleteProcess

	return res
}

// Validate проверяет корректность заполнения всех обязательных полей конфигурации.
// Возвращает ошибку, если хотя бы одно из ключевых полей не инициализировано или имеет невалидное значение.
//
//goland:noinspection DuplicatedCode
func (aco *BaseAuditClientOptions[D]) Validate() error {
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
	if utils.IsNil(aco.AuditAction) {
		return errs.NewTlCommonError("Validate", "audit action is required", nil)
	}

	return nil
}

// WithName задает уникальное имя для клиента аудита.
func WithName[D any](name string) BaseAuditClientOption[D] {
	return func(options *BaseAuditClientOptions[D]) {
		options.Pool.Name = name
	}
}

// WithPoolWorkerCount задает количество параллельных воркеров для обработки данных пула.
func WithPoolWorkerCount[D any](workerCount int) BaseAuditClientOption[D] {
	return func(options *BaseAuditClientOptions[D]) {
		options.Pool.WorkerCount = workerCount
	}
}

// WithPoolDataCapacity устанавливает емкость буферизированного канала (очереди) для данных пула.
func WithPoolDataCapacity[D any](dataCapacity int) BaseAuditClientOption[D] {
	return func(options *BaseAuditClientOptions[D]) {
		options.Pool.DataCapacity = dataCapacity
	}
}

// WithPoolCompleteProcess управляет флагом завершения обработки всех оставшихся в очереди данных пула при остановке.
func WithPoolCompleteProcess[D any](flag bool) BaseAuditClientOption[D] {
	return func(options *BaseAuditClientOptions[D]) {
		options.Pool.CompleteProcess = flag
	}
}

// WithPoolStopTimeout задает максимальное время ожидания корректного завершения работы пула (Graceful Shutdown).
func WithPoolStopTimeout[D any](timeout time.Duration) BaseAuditClientOption[D] {
	return func(options *BaseAuditClientOptions[D]) {
		options.Pool.StopTimeout = timeout
	}
}

// WithTokenProvider устанавливает провайдер токенов для аутентификации запросов аудита.
func WithTokenProvider[D any](provider auth.TokenProvider) BaseAuditClientOption[D] {
	return func(options *BaseAuditClientOptions[D]) {
		options.TokenProvider = provider
	}
}

// WithAuditAction задает целевое действие выполнения (обработчик) для поступающих данных.
func WithAuditAction[D any](action AuditAction[D]) BaseAuditClientOption[D] {
	return func(options *BaseAuditClientOptions[D]) {
		options.AuditAction = action
	}
}

// WithLogger конфигурирует логгер для записи системных событий клиента.
func WithLogger[D any](log logger.Logger) BaseAuditClientOption[D] {
	return func(options *BaseAuditClientOptions[D]) {
		options.Pool.Logger = log
	}
}
