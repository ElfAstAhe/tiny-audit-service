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

type AuthListByUsernameTraceInteractor struct {
	*telemetry.BaseTelemetry
	spanName string
	next     usecase.AuthListByUsernameUseCase
}

var _ usecase.AuthListByUsernameUseCase = (*AuthListByUsernameTraceInteractor)(nil)

func NewAuthListByUsernameUseCase(ucName string, next usecase.AuthListByUsernameUseCase) *AuthListByUsernameTraceInteractor {
	return &AuthListByUsernameTraceInteractor{
		spanName:      fmt.Sprintf("%s.List", ucName),
		next:          next,
		BaseTelemetry: telemetry.NewBaseTelemetry(ucName),
	}
}

func (lut *AuthListByUsernameTraceInteractor) List(ctx context.Context, username string, limit, offset int) ([]*domain.AuthAudit, error) {
	ctx, span := lut.StartSpan(ctx, lut.spanName)
	defer span.End()

	span.SetAttributes(
		attribute.String("username", username),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := lut.next.List(ctx, username, limit, offset)
	if err != nil {
		span.AddEvent("List_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
