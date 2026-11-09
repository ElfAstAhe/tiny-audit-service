package domain

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

// AuthAuditRepository определяет контракт для выполнения базовых CRUD-операций
// и специализированных выборок над сущностями AuthAudit в хранилище данных.
type AuthAuditRepository interface {
	domain.CRUDRepository[*AuthAudit, string]

	// ListByPeriod возвращает срез записей аудита аутентификации за указанный временной интервал с поддержкой пагинации.
	ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*AuthAudit, error)

	// ListByUsername возвращает срез записей аудита аутентификации для конкретного пользователя с поддержкой пагинации.
	ListByUsername(ctx context.Context, username string, offset, limit int) ([]*AuthAudit, error)
}
