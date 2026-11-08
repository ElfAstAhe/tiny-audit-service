package utils

import (
	"context"
)

// AuditRequestIDContextKey определяет изолированный тип ключа контекста для безопасного хранения Request ID.
type AuditRequestIDContextKey struct{}

// AuditTraceIDContextKey определяет изолированный тип ключа контекста для безопасного хранения Trace ID распределенной трассировки.
type AuditTraceIDContextKey struct{}

var AuditReqID = AuditRequestIDContextKey{}
var AuditTrcID = AuditTraceIDContextKey{}

// RequestIDFromContext выполняет безопасное извлечение и приведение типа для Request ID из контекста выполнения.
func RequestIDFromContext(ctx context.Context) string {
	val, ok := ctx.Value(AuditReqID).(string)
	if !ok {
		return UnknownValue
	}

	return val
}

// TraceIDFromContext выполняет безопасное извлечение и приведение типа для Trace ID распределенной трассировки из контекста выполнения.
func TraceIDFromContext(ctx context.Context) string {
	val, ok := ctx.Value(AuditTrcID).(string)
	if !ok {
		return UnknownValue
	}

	return val
}

// WithRequestID обогащает контекст выполнения уникальным идентификатором входящего запроса Request ID.
func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, AuditReqID, requestID)
}

// WithTraceID обогащает контекст выполнения уникальным идентификатором распределенной трассировки Trace ID.
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, AuditTrcID, traceID)
}
