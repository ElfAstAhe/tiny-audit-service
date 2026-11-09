package trace

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/infra/telemetry"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// DataListByInstanceInteractor представляет собой декоратор UseCase-слоя, предназначенный
// для распределенной трассировки и логирования спанов при выборке аудита данных по конкретному экземпляру сущности.
type DataListByInstanceInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.DataListByInstanceUseCase
}

var _ usecase.DataListByInstanceUseCase = (*DataListByInstanceInteractor)(nil)

// NewDataListByInstanceUseCase создает и инициализирует новый экземпляр трассирующего декоратора DataListByInstanceInteractor.
func NewDataListByInstanceUseCase(ucName string, next usecase.DataListByInstanceUseCase) *DataListByInstanceInteractor {
	return &DataListByInstanceInteractor{
		spanName:      fmt.Sprintf("%s.List", ucName),
		next:          next,
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
	}
}

// List выполняет декорируемый вызов UseCase для получения записей аудита конкретного инстанса,
// обогащая текущий спан параметрами типа, идентификатора, лимитов, а также фиксируя возникающие ошибки.
func (lit *DataListByInstanceInteractor) List(ctx context.Context, typeName, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	ucCtx, span := lit.StartSpan(ctx, lit.spanName)
	defer span.End()

	span.SetAttributes(
		attribute.String("param.typeName", typeName),
		attribute.String("param.instanceID", instanceID),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := lit.next.List(ucCtx, typeName, instanceID, limit, offset)
	if err != nil {
		span.AddEvent("List_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
