package utils

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/domain"
)

type AuditRequestIDContextKey struct{}
type AuditTraceIDContextKey struct{}

var AuditReqID = AuditRequestIDContextKey{}
var AuditTrcID = AuditTraceIDContextKey{}

// HasChanges test two instances for changes
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

func GetAuditStatus(err error) string {
	if err == nil {
		return dto.AuditStatusSuccess
	}

	return dto.AuditStatusFail
}

func UsernameFromContext(ctx context.Context) string {
	subj := auth.FromContext(ctx)
	if subj == nil {
		return "unknown"
	}

	return subj.Name
}

func RequestIDFromContext(ctx context.Context) string {
	val, ok := ctx.Value(AuditReqID).(string)
	if !ok {
		return "unknown"
	}

	return val
}

func TraceIDFromContext(ctx context.Context) string {
	val, ok := ctx.Value(AuditTrcID).(string)
	if !ok {
		return "unknown"
	}

	return val
}

func WithRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, AuditReqID, requestID)
}

func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, AuditTrcID, traceID)
}
