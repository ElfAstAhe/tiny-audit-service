package telemetry

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/infra/telemetry"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type DataListByInstanceTraceInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.DataListByInstanceUseCase
}

var _ usecase.DataListByInstanceUseCase = (*DataListByInstanceTraceInteractor)(nil)

func NewDataListByInstanceUseCase(ucName string, next usecase.DataListByInstanceUseCase) *DataListByInstanceTraceInteractor {
	return &DataListByInstanceTraceInteractor{
		spanName:      fmt.Sprintf("%s.List", ucName),
		next:          next,
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
	}
}

func (lit *DataListByInstanceTraceInteractor) List(ctx context.Context, typeName, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	ctx, span := lit.StartSpan(ctx, lit.spanName)
	defer span.End()

	span.SetAttributes(
		attribute.String("param.typeName", typeName),
		attribute.String("param.instanceID", instanceID),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := lit.next.List(ctx, typeName, instanceID, limit, offset)
	if err != nil {
		span.AddEvent("List_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
