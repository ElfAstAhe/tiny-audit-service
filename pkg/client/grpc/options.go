package grpc

import (
	"crypto/x509"
	"strings"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
)

// Константы со значениями по умолчанию для конфигурации gRPC-клиента.
const (
	DefaultTarget                       string        = "localhost:51052"
	DefaultConnectTimeout               time.Duration = 15 * time.Second
	DefaultIdleTimeout                  time.Duration = 15 * time.Second
	DefaultKeepAliveTime                time.Duration = 15 * time.Second
	DefaultKeepAliveTimeout             time.Duration = 15 * time.Second
	DefaultKeepAlivePermitWithoutStream bool          = true
	DefaultWorkerCount                  int           = 2
	DefaultDataCapacity                 int           = 10000
	DefaultCompleteProcessing           bool          = true
	DefaultStopTimeout                  time.Duration = 10 * time.Second
)

// Option определяет сигнатуру функции для настройки параметров gRPC-клиента.
type Option[D any] func(options *Options[D])

// Options содержит низкоуровневые параметры конфигурации gRPC-соединения,
// включая настройки безопасности и механизмы KeepAlive (KA).
type Options[D any] struct {
	*client.BaseAuditClientOptions[D]
	Target                       string
	CertPool                     *x509.CertPool
	CertServerNameOverride       string
	ConnTimeout                  time.Duration
	IdleTimeout                  time.Duration
	KeepAliveTime                time.Duration
	KeepAliveTimeout             time.Duration
	KeepAlivePermitWithoutStream bool
}

// NewOptions создает и возвращает новый экземпляр Options
// с заполненными значениями параметров KeepAlive и таймаутов по умолчанию.
func NewOptions[D any]() *Options[D] {
	res := &Options[D]{
		BaseAuditClientOptions:       client.NewBaseAuditClientOptions[D](),
		Target:                       DefaultTarget,
		ConnTimeout:                  DefaultConnectTimeout,
		IdleTimeout:                  DefaultIdleTimeout,
		KeepAliveTime:                DefaultKeepAliveTime,
		KeepAliveTimeout:             DefaultKeepAliveTimeout,
		KeepAlivePermitWithoutStream: DefaultKeepAlivePermitWithoutStream,
	}
	res.Pool.WorkerCount = DefaultWorkerCount
	res.Pool.DataCapacity = DefaultDataCapacity
	res.Pool.CompleteProcess = DefaultCompleteProcessing
	res.Pool.StopTimeout = DefaultStopTimeout

	return res
}

// Validate проверяет корректность заполнения числовых параметров и таймаутов gRPC-клиента.
// Возвращает ошибку, если временные интервалы имеют невалидные (отрицательные или нулевые) значения.
//
//goland:noinspection DuplicatedCode
func (co *Options[D]) Validate() error {
	if utils.IsNil(co.BaseAuditClientOptions) || utils.IsNil(co.Pool) {
		return errs.NewTlCommonError("Validate", "base audit client options and pool configuration are required", nil)
	}
	if strings.TrimSpace(co.Pool.Name) == "" {
		return errs.NewTlCommonError("Validate", "name is required", nil)
	}
	if co.Pool.WorkerCount <= 0 {
		return errs.NewTlCommonError("Validate", "worker count is required", nil)
	}
	if co.Pool.DataCapacity <= 0 {
		return errs.NewTlCommonError("Validate", "data capacity is required", nil)
	}
	if co.Pool.StopTimeout <= 0 {
		return errs.NewTlCommonError("Validate", "stop timeout is required", nil)
	}
	if utils.IsNil(co.Pool.Logger) {
		return errs.NewTlCommonError("Validate", "logger is required", nil)
	}
	if utils.IsNil(co.TokenProvider) {
		return errs.NewTlCommonError("Validate", "token provider is required", nil)
	}
	if utils.IsNil(co.AuditAction) {
		return errs.NewTlCommonError("Validate", "audit action is required", nil)
	}
	if strings.TrimSpace(co.Target) == "" {
		return errs.NewTlCommonError("Validate", "target is required", nil)
	}
	if co.ConnTimeout <= 0 {
		return errs.NewTlCommonError("Validate", "connect timeout is required", nil)
	}
	if co.IdleTimeout <= 0 {
		return errs.NewTlCommonError("Validate", "idle timeout is required", nil)
	}
	if co.KeepAliveTime <= 0 {
		return errs.NewTlCommonError("Validate", "keep alive time is required", nil)
	}
	if co.KeepAliveTimeout <= 0 {
		return errs.NewTlCommonError("Validate", "keep alive timeout is required", nil)
	}

	return nil
}

// WithName задает уникальное имя для пула воркеров gRPC-клиента.
func WithName[D any](name string) Option[D] {
	return func(options *Options[D]) {
		options.Pool.Name = name
	}
}

// WithTarget настраивает хост и порт сервера
func WithTarget[D any](target string) Option[D] {
	return func(options *Options[D]) {
		options.Target = target
	}
}

// WithCertPool управляет использованием защищенного TLS-соединения для gRPC-клиента.
func WithCertPool[D any](pool *x509.CertPool) Option[D] {
	return func(options *Options[D]) {
		options.CertPool = pool
	}
}

// WithIdleTimeout задаёт idle timeout соединения
func WithIdleTimeout[D any](timeout time.Duration) Option[D] {
	return func(options *Options[D]) {
		options.IdleTimeout = timeout
	}
}

// WithConnectTimeout задает максимальное время ожидания при установке gRPC-соединения.
func WithConnectTimeout[D any](timeout time.Duration) Option[D] {
	return func(options *Options[D]) {
		options.ConnTimeout = timeout
	}
}

// WithKeepAliveTime определяет интервал отправки пингов KeepAlive для проверки активности соединения.
func WithKeepAliveTime[D any](duration time.Duration) Option[D] {
	return func(options *Options[D]) {
		options.KeepAliveTime = duration
	}
}

// WithKeepAliveTimeout задает время ожидания ответа на пинг KeepAlive перед закрытием соединения как неактивного.
func WithKeepAliveTimeout[D any](timeout time.Duration) Option[D] {
	return func(options *Options[D]) {
		options.KeepAliveTimeout = timeout
	}
}

// WithKeepAlivePermitWithoutStream разрешает отправку KeepAlive пингов даже при отсутствии активных RPC-стримов.
func WithKeepAlivePermitWithoutStream[D any](flag bool) Option[D] {
	return func(options *Options[D]) {
		options.KeepAlivePermitWithoutStream = flag
	}
}

// WithWorkerCount задает количество параллельных воркеров отправки gRPC-сообщений.
func WithWorkerCount[D any](count int) Option[D] {
	return func(options *Options[D]) {
		options.Pool.WorkerCount = count
	}
}

// WithDataCapacity устанавливает максимальную емкость буфера очереди отправки пула.
func WithDataCapacity[D any](capacity int) Option[D] {
	return func(options *Options[D]) {
		options.Pool.DataCapacity = capacity
	}
}

// WithCompleteProcess определяет, нужно ли досылать оставшиеся в буфере логи при остановке пула воркеров.
func WithCompleteProcess[D any](flag bool) Option[D] {
	return func(options *Options[D]) {
		options.Pool.CompleteProcess = flag
	}
}

// WithStopTimeout задает максимальное время ожидания Graceful Shutdown для пула воркеров gRPC-клиента.
func WithStopTimeout[D any](timeout time.Duration) Option[D] {
	return func(options *Options[D]) {
		options.Pool.StopTimeout = timeout
	}
}

// WithTokenProvider конфигурирует провайдер авторизационных токенов для gRPC-запросов.
func WithTokenProvider[D any](provider auth.TokenProvider) Option[D] {
	return func(options *Options[D]) {
		options.TokenProvider = provider
	}
}

// withAuditAction задает целевое действие выполнения (обработчик) для поступающих данных.
func withAuditAction[D any](action client.AuditAction[D]) Option[D] {
	return func(options *Options[D]) {
		options.AuditAction = action
	}
}

// WithLogger настраивает системный логгер для пула воркеров gRPC-клиента.
func WithLogger[D any](log logger.Logger) Option[D] {
	return func(options *Options[D]) {
		options.Pool.Logger = log
	}
}
