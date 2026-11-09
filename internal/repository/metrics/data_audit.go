package metrics

import (
	"context"
	"time"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/repository/metrics"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

// DataAuditMetricsRepository представляет собой декоратор репозитория, предназначенный
// для автоматического сбора бизнес-метрик и замера времени выполнения операций с данными аудита.
type DataAuditMetricsRepository struct {
	*metrics.CRUD[*domain.DataAudit, string]
	next domain.DataAuditRepository
}

var _ libdom.CRUDRepository[*domain.DataAudit, string] = (*DataAuditMetricsRepository)(nil)
var _ domain.DataAuditRepository = (*DataAuditMetricsRepository)(nil)

// NewDataAuditRepository создает новый экземпляр декоратора DataAuditMetricsRepository
// для автоматического замера метрик выполнения операций с базой данных.
func NewDataAuditRepository(next domain.DataAuditRepository) *DataAuditMetricsRepository {
	return &DataAuditMetricsRepository{
		next: next,
		CRUD: metrics.NewCRUD[*domain.DataAudit, string]("DataAuditRepository", next),
	}
}

// ListByPeriod возвращает список записей аудита данных за указанный промежуток времени,
// автоматически замеряя время выполнения операции и фиксируя метрики ошибок.
func (dam *DataAuditMetricsRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) (res []*domain.DataAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(dam.GetRepositoryName(), "ListByPeriod", err, start)
	}(time.Now())

	return dam.next.ListByPeriod(ctx, from, till, limit, offset)
}

// ListByInstance возвращает список записей аудита для конкретного экземпляра сущности по её типу и ID,
// автоматически замеряя время выполнения операции и фиксируя метрики ошибок.
func (dam *DataAuditMetricsRepository) ListByInstance(ctx context.Context, typeName string, instanceID string, limit, offset int) (res []*domain.DataAudit, err error) {
	defer func(start time.Time) {
		metrics.ObserveRepositoryOp(dam.GetRepositoryName(), "ListByInstance", err, start)
	}(time.Now())

	return dam.next.ListByInstance(ctx, typeName, instanceID, limit, offset)
}
