package utils

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/domain"
)

// AuditRequestIDContextKey определяет изолированный тип ключа контекста для безопасного хранения Request ID.
type AuditRequestIDContextKey struct{}

// AuditTraceIDContextKey определяет изолированный тип ключа контекста для безопасного хранения Trace ID распределенной трассировки.
type AuditTraceIDContextKey struct{}

var AuditReqID = AuditRequestIDContextKey{}
var AuditTrcID = AuditTraceIDContextKey{}

// HasChanges выполняет сверку двух экземпляров сущностей на наличие изменений, сопоставляя их атомарные хэш-коды HashCode.
func HasChanges(before, after domain.Auditable) bool {
	if utils.IsNil(before) && !utils.IsNil(after) {
		return true
	}
	if !utils.IsNil(before) && utils.IsNil(after) {
		return true
	}
	if utils.IsNil(before) && utils.IsNil(after) {
		return false
	}

	return before.HashCode() != after.HashCode()
}

// BuildSingleDataAuditValues компилирует односторонний срез изменений полей (только состояние «До» или только «После») для лога аудита.
func BuildSingleDataAuditValues(instance domain.Auditable, before bool) []*dto.DataAuditValueDTO {
	if utils.IsNil(instance) {
		return nil
	}

	auditMap := instance.ToAuditMap()

	res := make([]*dto.DataAuditValueDTO, 0, len(auditMap))
	for field, auditField := range instance.ToAuditMap() {
		dataValue := &dto.DataAuditValueDTO{
			Name:        field,
			Description: auditField.Description,
		}
		if before {
			dataValue.Before = auditField.Value
		} else {
			dataValue.After = auditField.Value
		}

		res = append(res, dataValue)
	}

	return res
}

// BuildBothDataAuditValues вычисляет дифференциал изменений полей между состояниями «До» и «После», фиксируя только мутировавшие атрибуты.
func BuildBothDataAuditValues(before, after domain.Auditable) []*dto.DataAuditValueDTO {
	if utils.IsNil(before) || utils.IsNil(after) {
		return nil
	}
	beforeFields := before.ToAuditMap()
	afterFields := after.ToAuditMap()

	res := make([]*dto.DataAuditValueDTO, 0, len(before.ToAuditMap()))

	for field, auditField := range beforeFields {
		// add only changes, equal fields ignored
		if beforeFields[field].Value != afterFields[field].Value {
			dataValue := &dto.DataAuditValueDTO{
				Name:        field,
				Description: auditField.Description,
				Before:      beforeFields[field].Value,
				After:       afterFields[field].Value,
			}
			res = append(res, dataValue)
		}
	}

	return res
}

// GetAuditStatus транслирует прикладную ошибку операции в детерминированный строковый маркер статуса (success/fail).
func GetAuditStatus(err error) string {
	if err == nil {
		return dto.AuditStatusSuccess
	}

	return dto.AuditStatusFail
}

// UsernameFromContext извлекает профиль авторизованного субъекта из контекста запроса и возвращает его имя, либо маркер unknown.
func UsernameFromContext(ctx context.Context) string {
	subj := auth.FromContext(ctx)
	if subj == nil {
		return "unknown"
	}

	return subj.Name
}

// RequestIDFromContext выполняет безопасное извлечение и приведение типа для Request ID из контекста выполнения.
func RequestIDFromContext(ctx context.Context) string {
	val, ok := ctx.Value(AuditReqID).(string)
	if !ok {
		return "unknown"
	}

	return val
}

// TraceIDFromContext выполняет безопасное извлечение и приведение типа для Trace ID распределенной трассировки из контекста выполнения.
func TraceIDFromContext(ctx context.Context) string {
	val, ok := ctx.Value(AuditTrcID).(string)
	if !ok {
		return "unknown"
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
