package worker

import (
	"strings"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/broker"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/transport/worker/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

const (
	defaultBatchSize        int           = 50
	defaultBatchReadTimeout time.Duration = 2 * time.Second
)

// LoginAttemptsOption определяет сигнатуру функции для конфигурирования структуры LoginAttemptsOptions.
type LoginAttemptsOption func(*LoginAttemptsOptions)

// LoginAttemptsOptions группирует настройки периодического диспетчера задач, брокера сообщений,
// UseCase-слоя и таймаутов, необходимых для асинхронного разбора событий попыток аутентификации.
type LoginAttemptsOptions struct {
	Name               string
	DispatcherOpts     *worker.PeriodicDispatcherOptions[*dto.LoginAttemptWorkerJob]
	Receiver           broker.Receiver
	AuthAuditUC        usecase.AuthAuditUseCase
	BatchSize          int
	BatchReadTimeout   time.Duration
	AcknowledgeTimeout time.Duration
	Logger             logger.Logger
}

// NewLoginAttemptsOptions представляет собой фабричный метод для создания опций LoginAttemptsOptions со значениями по умолчанию.
func NewLoginAttemptsOptions() *LoginAttemptsOptions {
	return &LoginAttemptsOptions{
		BatchSize:        defaultBatchSize,
		BatchReadTimeout: defaultBatchReadTimeout,
		DispatcherOpts:   worker.NewPeriodicDispatcherOptions[*dto.LoginAttemptWorkerJob](),
	}
}

// Validate проверяет полноту, корректность заполнения и валидность интервалов всех зависимостей в LoginAttemptsOptions.
func (lao LoginAttemptsOptions) Validate() error {
	if strings.TrimSpace(lao.Name) == "" {
		return errs.NewCommonError("name not applied (empty)", nil)
	}
	if utils.IsNil(lao.Logger) {
		return errs.NewCommonError("logger not provided (nil)", nil)
	}
	if utils.IsNil(lao.DispatcherOpts) {
		return errs.NewCommonError("dispatcher options not provided (nil)", nil)
	}
	if utils.IsNil(lao.Receiver) {
		return errs.NewCommonError("receiver not provided (nil)", nil)
	}
	if utils.IsNil(lao.AuthAuditUC) {
		return errs.NewCommonError("auth audit use case not provided (nil)", nil)
	}
	if lao.BatchSize <= 0 {
		return errs.NewCommonError("batch size must be greater than 0", nil)
	}
	if lao.BatchReadTimeout <= 0 {
		return errs.NewCommonError("batch read timeout must be greater than 0", nil)
	}
	if lao.AcknowledgeTimeout <= 0 {
		return errs.NewCommonError("acknowledge timeout must be greater than 0", nil)
	}

	return nil
}

// WithLAOName возвращает опцию для настройки уникального строкового имени воркера.
func WithLAOName(name string) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.Name = name
	}
}

// WithLAODispatcherOpts возвращает опцию для инициализации и конфигурации вложенных параметров PeriodicDispatcherOptions.
func WithLAODispatcherOpts(options ...worker.PeriodicDispatcherOption[*dto.LoginAttemptWorkerJob]) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.DispatcherOpts = worker.NewPeriodicDispatcherOptions[*dto.LoginAttemptWorkerJob]()
		for _, option := range options {
			option(opts.DispatcherOpts)
		}
	}
}

// WithLAOReceiver возвращает опцию для привязки интерфейса получения сообщений из брокера (Receiver).
func WithLAOReceiver(receiver broker.Receiver) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.Receiver = receiver
	}
}

// WithLAOAuthAuditUseCase возвращает опцию для привязки UseCase-интерактора регистрации логов аутентификации.
func WithLAOAuthAuditUseCase(authAuditUC usecase.AuthAuditUseCase) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.AuthAuditUC = authAuditUC
	}
}

// WithLAOBatchSize возвращает опцию для изменения размера пачки (Batch Size) единовременно вычитываемых сообщений.
func WithLAOBatchSize(batchSize int) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.BatchSize = batchSize
	}
}

// WithLAOBatchReadTimeout возвращает опцию для настройки максимального таймаута ожидания накопления пачки сообщений.
func WithLAOBatchReadTimeout(timeout time.Duration) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.BatchReadTimeout = timeout
	}
}

// WithLAOAcknowledgeTimeout возвращает опцию для изменения времени ожидания подтверждения (Acknowledge) факта обработки сообщений.
func WithLAOAcknowledgeTimeout(timeout time.Duration) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.AcknowledgeTimeout = timeout
	}
}

// WithLAOLogger возвращает опцию для привязки изолированного системного логгера к настройкам воркера.
func WithLAOLogger(logger logger.Logger) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.Logger = logger
	}
}
