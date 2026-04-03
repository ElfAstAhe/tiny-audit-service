package trace

import (
	"context"
	"fmt"
	"time"

	libdomain "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type AuthAuditTraceRepository struct {
	*repository.BaseCRUDTraceRepository[*domain.AuthAudit, string]
	next domain.AuthAuditRepository
}

var _ libdomain.CRUDRepository[*domain.AuthAudit, string] = (*AuthAuditTraceRepository)(nil)
var _ domain.AuthAuditRepository = (*AuthAuditTraceRepository)(nil)

func NewAuthAuditTraceRepository(next domain.AuthAuditRepository) *AuthAuditTraceRepository {
	return &AuthAuditTraceRepository{
		next:                    next,
		BaseCRUDTraceRepository: repository.NewBaseCRUDTraceRepository[*domain.AuthAudit, string]("AuthAuditRepository", next),
	}
}

func (aat *AuthAuditTraceRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.AuthAudit, error) {
	ctx, span := aat.StartSpan(ctx, fmt.Sprintf("%s.ListByPeriod", aat.BaseCRUDTraceRepository.GetRepositoryName()))
	defer span.End()

	span.SetAttributes(
		attribute.String("param.from", from.Format(time.DateTime)),
		attribute.String("param.till", till.Format(time.DateTime)),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := aat.next.ListByPeriod(ctx, from, till, limit, offset)
	if err != nil {
		span.AddEvent("ListByPeriod_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}

func (aat *AuthAuditTraceRepository) ListByUsername(ctx context.Context, username string, limit, offset int) ([]*domain.AuthAudit, error) {
	ctx, span := aat.StartSpan(ctx, fmt.Sprintf("%s.ListByUsername", aat.BaseCRUDTraceRepository.GetRepositoryName()))
	defer span.End()

	span.SetAttributes(
		attribute.String("param.username", username),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := aat.next.ListByUsername(ctx, username, limit, offset)
	if err != nil {
		span.AddEvent("ListByUsername_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
