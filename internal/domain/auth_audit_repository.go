package domain

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

type AuthAuditRepository interface {
	domain.CRUDRepository[*AuthAudit, string]
	ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*AuthAudit, error)
	ListByUsername(ctx context.Context, username string, offset, limit int) ([]*AuthAudit, error)
}
