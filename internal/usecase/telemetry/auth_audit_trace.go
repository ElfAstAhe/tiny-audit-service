package telemetry

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/infra/telemetry"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"go.opentelemetry.io/otel/codes"
)

type AuthAuditTraceInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.AuthAuditUseCase
}

var _ usecase.AuthAuditUseCase = (*AuthAuditTraceInteractor)(nil)

func NewAuthAuditTraceUseCase(
	ucName string,
	next usecase.AuthAuditUseCase,
) *AuthAuditTraceInteractor {
	return &AuthAuditTraceInteractor{
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
		spanName:      fmt.Sprintf("%s.Audit", ucName),
		next:          next,
	}
}

func (aat *AuthAuditTraceInteractor) Audit(ctx context.Context, data *domain.AuthAudit) error {
	ctx, span := aat.StartSpan(ctx, aat.spanName)
	defer span.End()

	err := aat.next.Audit(ctx, data)
	if err != nil {
		span.AddEvent("Audit_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
