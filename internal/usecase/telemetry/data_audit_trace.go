package telemetry

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/infra/telemetry"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"go.opentelemetry.io/otel/codes"
)

type DataAuditTraceInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.DataAuditUseCase
}

var _ usecase.DataAuditUseCase = (*DataAuditTraceInteractor)(nil)

func NewDataAuditTraceUseCase(
	ucName string,
	next usecase.DataAuditUseCase,
) *DataAuditTraceInteractor {
	return &DataAuditTraceInteractor{
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
		spanName:      fmt.Sprintf("%s.Audit", ucName),
		next:          next,
	}
}

func (dat *DataAuditTraceInteractor) Audit(ctx context.Context, data *domain.DataAudit) error {
	ctx, span := dat.StartSpan(ctx, dat.spanName)
	defer span.End()

	err := dat.next.Audit(ctx, data)
	if err != nil {
		span.AddEvent("Audit_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
