package metrics

import (
	"context"
	"time"

	libdomain "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/infra/metrics"
	"github.com/ElfAstAhe/go-service-template/pkg/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type DataAuditMetricsRepository struct {
	*repository.BaseCRUDMetricsRepository[*domain.DataAudit, string]
	next domain.DataAuditRepository
}

var _ libdomain.CRUDRepository[*domain.DataAudit, string] = (*DataAuditMetricsRepository)(nil)
var _ domain.DataAuditRepository = (*DataAuditMetricsRepository)(nil)

func NewDataAuditMetricsRepository(next domain.DataAuditRepository) *DataAuditMetricsRepository {
	return &DataAuditMetricsRepository{
		next:                      next,
		BaseCRUDMetricsRepository: repository.NewBaseCRUDMetricsRepository[*domain.DataAudit, string]("DataAuditRepository", next),
	}
}

func (dam *DataAuditMetricsRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) (res []*domain.DataAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(dam.BaseCRUDMetricsRepository.GetRepositoryName(), "ListByPeriod", err, start)
	}(time.Now())

	return dam.next.ListByPeriod(ctx, from, till, limit, offset)
}

func (dam *DataAuditMetricsRepository) ListByInstance(ctx context.Context, typeName string, instanceID string, limit, offset int) (res []*domain.DataAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(dam.BaseCRUDMetricsRepository.GetRepositoryName(), "ListByInstance", err, start)
	}(time.Now())

	return dam.next.ListByInstance(ctx, typeName, instanceID, limit, offset)
}
