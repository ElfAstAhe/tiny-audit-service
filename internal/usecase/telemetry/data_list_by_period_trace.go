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

type DataListByPeriodTraceInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.DataListByPeriodUseCase
}

var _ usecase.DataListByPeriodUseCase = (*DataListByPeriodTraceInteractor)(nil)

func NewDataListByPeriodUseCase(ucName string, next usecase.DataListByPeriodUseCase) *DataListByPeriodTraceInteractor {
	return &DataListByPeriodTraceInteractor{
		spanName:      fmt.Sprintf("%s.List", ucName),
		next:          next,
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
	}
}

func (pti *DataListByPeriodTraceInteractor) List(ctx context.Context, from, till time.Time, limit int, offset int) ([]*domain.DataAudit, error) {
	ctx, span := pti.StartSpan(ctx, pti.spanName)
	defer span.End()

	span.SetAttributes(
		attribute.String("param.from", from.Format(time.DateTime)),
		attribute.String("param.till", till.Format(time.DateTime)),
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
