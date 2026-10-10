package grpc

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// Константы со значениями по умолчанию для конфигурации gRPC-клиента.
const (
	DefaultSecure                       bool          = false
	DefaultConnectionTimeout            time.Duration = 15 * time.Second
	DefaultKeepAliveTime                time.Duration = 15 * time.Second
	DefaultKeepAliveTimeout             time.Duration = 15 * time.Second
	DefaultKeepAlivePermitWithoutStream bool          = true
)

// RawClientOption определяет сигнатуру функции для настройки параметров gRPC-клиента.
type RawClientOption func(options *RawClientOptions)

// RawClientOptions содержит низкоуровневые параметры конфигурации gRPC-соединения,
// включая настройки безопасности и механизмы KeepAlive (KA).
type RawClientOptions struct {
	Secure                bool
	ConnTimeout           time.Duration
	KATime                time.Duration
	KATimeout             time.Duration
	KAPermitWithoutStream bool
}

// NewRawClientOptions создает и возвращает новый экземпляр RawClientOptions
// с заполненными значениями параметров KeepAlive и таймаутов по умолчанию.
func NewRawClientOptions() *RawClientOptions {
	return &RawClientOptions{
		Secure:                DefaultSecure,
		ConnTimeout:           DefaultConnectionTimeout,
		KATime:                DefaultKeepAliveTime,
		KATimeout:             DefaultKeepAliveTimeout,
		KAPermitWithoutStream: DefaultKeepAlivePermitWithoutStream,
	}
}

// Validate проверяет корректность заполнения числовых параметров и таймаутов gRPC-клиента.
// Возвращает ошибку, если временные интервалы имеют невалидные (отрицательные или нулевые) значения.
func (rac *RawClientOptions) Validate() error {
	if rac.ConnTimeout <= 0 {
		return errs.NewCommonError("connection timeout is required", nil)
	}
	if rac.KATime <= 0 {
		return errs.NewCommonError("ka time is required", nil)
	}
	if rac.KATimeout <= 0 {
		return errs.NewCommonError("ka timeout is required", nil)
	}

	return nil
}

// WithRawSecure управляет использованием защищенного TLS-соединения для gRPC-клиента.
func WithRawSecure(flag bool) RawClientOption {
	return func(options *RawClientOptions) {
		options.Secure = flag
	}
}

// WithRawConnectionTimeout задает максимальное время ожидания при установке gRPC-соединения.
func WithRawConnectionTimeout(timeout time.Duration) RawClientOption {
	return func(options *RawClientOptions) {
		options.ConnTimeout = timeout
	}
}

// WithRawKeepAliveTime определяет интервал отправки пингов KeepAlive для проверки активности соединения.
func WithRawKeepAliveTime(duration time.Duration) RawClientOption {
	return func(options *RawClientOptions) {
		options.KATime = duration
	}
}

// WithRawKeepAliveTimeout задает время ожидания ответа на пинг KeepAlive перед закрытием соединения как неактивного.
func WithRawKeepAliveTimeout(timeout time.Duration) RawClientOption {
	return func(options *RawClientOptions) {
		options.KATimeout = timeout
	}
}

// WithRawKeepAlivePermitWithoutStream разрешает отправку KeepAlive пингов даже при отсутствии активных RPC-стримов.
func WithRawKeepAlivePermitWithoutStream(flag bool) RawClientOption {
	return func(options *RawClientOptions) {
		options.KAPermitWithoutStream = flag
	}
}
