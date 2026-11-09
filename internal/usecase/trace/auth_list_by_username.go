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

// AuthListByUsernameInteractor представляет собой декоратор UseCase-слоя, предназначенный
// для распределенной трассировки и логирования спанов при выборке записей аудита аутентификации по имени пользователя.
type AuthListByUsernameInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.AuthListByUsernameUseCase
}

var _ usecase.AuthListByUsernameUseCase = (*AuthListByUsernameInteractor)(nil)

// NewAuthListByUsernameUseCase создает и инициализирует новый экземпляр трассирующего декоратора AuthListByUsernameInteractor.
func NewAuthListByUsernameUseCase(ucName string, next usecase.AuthListByUsernameUseCase) *AuthListByUsernameInteractor {
	return &AuthListByUsernameInteractor{
		spanName:      fmt.Sprintf("%s.List", ucName),
		next:          next,
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
	}
}

// List выполняет декорируемый вызов UseCase для получения списка логов аутентификации конкретного пользователя,
// обогащая текущий спан именем пользователя, параметрами пагинации, а также фиксируя возникающие ошибки и их статус.
func (lut *AuthListByUsernameInteractor) List(ctx context.Context, username string, limit, offset int) ([]*domain.AuthAudit, error) {
	ucCtx, span := lut.StartSpan(ctx, lut.spanName)
	defer span.End()

	span.SetAttributes(
		attribute.String("param.username", username),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := lut.next.List(ucCtx, username, limit, offset)
	if err != nil {
		span.AddEvent("List_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
