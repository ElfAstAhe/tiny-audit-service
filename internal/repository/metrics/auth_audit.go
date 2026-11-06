package metrics

import (
	"context"
	"time"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/repository/metrics"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type AuthAuditRepository struct {
	*metrics.CRUD[*domain.AuthAudit, string]
	next domain.AuthAuditRepository
}

var _ libdom.CRUDRepository[*domain.AuthAudit, string] = (*AuthAuditRepository)(nil)
var _ domain.AuthAuditRepository = (*AuthAuditRepository)(nil)

func NewAuthAuditRepository(next domain.AuthAuditRepository) *AuthAuditRepository {
	return &AuthAuditRepository{
		next: next,
		CRUD: metrics.NewCRUD[*domain.AuthAudit, string]("AuthAuditRepository", next),
	}
}

func (aam *AuthAuditRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) (res []*domain.AuthAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(aam.GetRepositoryName(), "ListByPeriod", err, start)
	}(time.Now())

	return aam.next.ListByPeriod(ctx, from, till, limit, offset)
}

func (aam *AuthAuditRepository) ListByUsername(ctx context.Context, username string, offset, limit int) (res []*domain.AuthAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(aam.GetRepositoryName(), "ListByUsername", err, start)
	}(time.Now())

	return aam.next.ListByUsername(ctx, username, offset, limit)
}
