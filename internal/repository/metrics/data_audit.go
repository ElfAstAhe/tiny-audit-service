package metrics

import (
	"context"
	"time"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/repository/metrics"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type DataAuditMetricsRepository struct {
	*metrics.BaseCRUDMetricsRepository[*domain.DataAudit, string]
	next domain.DataAuditRepository
}

var _ libdom.CRUDRepository[*domain.DataAudit, string] = (*DataAuditMetricsRepository)(nil)
var _ domain.DataAuditRepository = (*DataAuditMetricsRepository)(nil)

func NewDataAuditMetricsRepository(next domain.DataAuditRepository) *DataAuditMetricsRepository {
	return &DataAuditMetricsRepository{
		next:                      next,
		BaseCRUDMetricsRepository: metrics.NewBaseCRUDMetricsRepository[*domain.DataAudit, string]("DataAuditRepository", next),
	}
}

func (dam *DataAuditMetricsRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) (res []*domain.DataAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(dam.GetRepositoryName(), "ListByPeriod", err, start)
	}(time.Now())

	return dam.next.ListByPeriod(ctx, from, till, limit, offset)
}

func (dam *DataAuditMetricsRepository) ListByInstance(ctx context.Context, typeName string, instanceID string, limit, offset int) (res []*domain.DataAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(dam.GetRepositoryName(), "ListByInstance", err, start)
	}(time.Now())

	return dam.next.ListByInstance(ctx, typeName, instanceID, limit, offset)
}
