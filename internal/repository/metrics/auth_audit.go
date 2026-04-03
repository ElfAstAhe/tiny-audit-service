package metrics

import (
	"context"
	"time"

	libdomain "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/infra/metrics"
	"github.com/ElfAstAhe/go-service-template/pkg/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type AuthAuditMetricsRepository struct {
	*repository.BaseCRUDMetricsRepository[*domain.AuthAudit, string]
	next domain.AuthAuditRepository
}

var _ libdomain.CRUDRepository[*domain.AuthAudit, string] = (*AuthAuditMetricsRepository)(nil)
var _ domain.AuthAuditRepository = (*AuthAuditMetricsRepository)(nil)

func NewAuthAuditMetricsRepository(next domain.AuthAuditRepository) *AuthAuditMetricsRepository {
	return &AuthAuditMetricsRepository{
		next:                      next,
		BaseCRUDMetricsRepository: repository.NewBaseCRUDMetricsRepository[*domain.AuthAudit, string]("AuthAuditRepository", next),
	}
}

func (aam *AuthAuditMetricsRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) (res []*domain.AuthAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(aam.BaseCRUDMetricsRepository.GetRepositoryName(), "ListByPeriod", err, start)
	}(time.Now())

	return aam.next.ListByPeriod(ctx, from, till, limit, offset)
}

func (aam *AuthAuditMetricsRepository) ListByUsername(ctx context.Context, username string, offset, limit int) (res []*domain.AuthAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(aam.BaseCRUDMetricsRepository.GetRepositoryName(), "ListByUsername", err, start)
	}(time.Now())

	return aam.next.ListByUsername(ctx, username, limit, offset)
}
