package trace

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

// AuthListByPeriodInteractor представляет собой декоратор UseCase-слоя, предназначенный
// для распределенной трассировки и логирования спанов при выборке записей аудита аутентификации за период.
type AuthListByPeriodInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.AuthListByPeriodUseCase
}

var _ usecase.AuthListByPeriodUseCase = (*AuthListByPeriodInteractor)(nil)

// NewAuthListByPeriodUseCase создает и инициализирует новый экземпляр трассирующего декоратора AuthListByPeriodInteractor.
func NewAuthListByPeriodUseCase(ucName string, next usecase.AuthListByPeriodUseCase) *AuthListByPeriodInteractor {
	return &AuthListByPeriodInteractor{
		spanName:      fmt.Sprintf("%s.List", ucName),
		next:          next,
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
	}
}

// List выполняет декорируемый вызов UseCase для получения списка логов аутентификации за указанный временной период,
// обогащая текущий спан параметрами временных меток, лимитов, а также фиксируя возникающие ошибки и их статус.
func (lpt *AuthListByPeriodInteractor) List(ctx context.Context, from, till time.Time, limit int, offset int) ([]*domain.AuthAudit, error) {
	ucCtx, span := lpt.StartSpan(ctx, lpt.spanName)
	defer span.End()

	span.SetAttributes(
		attribute.String("param.from", from.Format(time.DateTime)),
		attribute.String("param.till", till.Format(time.DateTime)),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := lpt.next.List(ucCtx, from, till, limit, offset)
	if err != nil {
		span.AddEvent("List_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
