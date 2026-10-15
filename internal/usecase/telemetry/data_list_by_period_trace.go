package telemetry

import (
	"context"
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/infra/telemetry"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// DataListByPeriodTraceInteractor implements the usecase.DataListByPeriodUseCase interface,
// acting as a non-invasive distributed tracing decorator component.
// It manages OpenTelemetry span lifecycles and records parameter attributes for period-based audit data retrieval operations.
type DataListByPeriodTraceInteractor struct {
	*telemetry.BaseTelemetry                                 // Generic framework-level tracing base handle infrastructure
	spanName                 string                          // Pre-calculated target telemetry instrumentation name for execution boundaries
	next                     usecase.DataListByPeriodUseCase // Downstream execution layer or concrete application usecase interactor implementation path
}

// Compile-time interface compliance verification
var _ usecase.DataListByPeriodUseCase = (*DataListByPeriodTraceInteractor)(nil)

// NewDataListByPeriodUseCase acts as a factory constructor mounting required tracing infrastructure over a domain usecase interactor.
func NewDataListByPeriodUseCase(ucName string, next usecase.DataListByPeriodUseCase) *DataListByPeriodTraceInteractor {
	return &DataListByPeriodTraceInteractor{
		spanName:      fmt.Sprintf("%s.List", ucName),
		next:          next,
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
	}
}

// List intercepts period-based lookup invocations, capturing context duration boundaries and injecting transaction state metadata into OpenTelemetry spans.
func (pti *DataListByPeriodTraceInteractor) List(ctx context.Context, from, till time.Time, limit int, offset int) ([]*domain.DataAudit, error) {
	ctx, span := pti.StartSpan(ctx, pti.spanName)
	defer span.End()

	span.SetAttributes(
		attribute.Int64("param.from", from.Unix()),
		attribute.Int64("param.till", till.Unix()),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := pti.next.List(ctx, from, till, limit, offset)
	if err != nil {
		span.AddEvent("List_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
