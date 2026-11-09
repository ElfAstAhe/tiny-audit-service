package domain

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

// DataAuditRepository определяет контракт для выполнения базовых CRUD-операций
// и специализированных архивных выборок над сущностями DataAudit в хранилище данных.
type DataAuditRepository interface {
	domain.CRUDRepository[*DataAudit, string]

	// ListByPeriod возвращает срез записей аудита данных за указанный временной интервал с поддержкой пагинации.
	ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*DataAudit, error)

	// ListByInstance возвращает срез записей аудита для конкретного экземпляра сущности по её типу и уникальному ID.
	ListByInstance(ctx context.Context, typeName string, instanceID string, limit, offset int) ([]*DataAudit, error)
}
