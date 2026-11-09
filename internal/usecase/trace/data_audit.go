package trace

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/infra/telemetry"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"go.opentelemetry.io/otel/codes"
)

// DataAuditInteractor представляет собой декоратор UseCase-слоя, предназначенный
// для распределенной трассировки и логирования спанов при выполнении операций аудита изменений данных.
type DataAuditInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.DataAuditUseCase
}

var _ usecase.DataAuditUseCase = (*DataAuditInteractor)(nil)

// NewDataAuditUseCase создает и инициализирует новый экземпляр трассирующего декоратора DataAuditInteractor.
func NewDataAuditUseCase(
	ucName string,
	next usecase.DataAuditUseCase,
) *DataAuditInteractor {
	return &DataAuditInteractor{
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
		spanName:      fmt.Sprintf("%s.Audit", ucName),
		next:          next,
	}
}

// Audit выполняет декорируемый вызов UseCase для регистрации записи аудита данных,
// обогащая текущий контекст и фиксируя в спане возникающие ошибки и их статус.
func (dat *DataAuditInteractor) Audit(ctx context.Context, data *domain.DataAudit) error {
	ucCtx, span := dat.StartSpan(ctx, dat.spanName)
	defer span.End()

	err := dat.next.Audit(ucCtx, data)
	if err != nil {
		span.AddEvent("Audit_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return err
}
