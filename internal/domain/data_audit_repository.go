package domain

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

type DataAuditRepository interface {
	domain.CRUDRepository[*DataAudit, string]
	ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*DataAudit, error)
	ListByInstance(ctx context.Context, typeName string, instanceID string, limit, offset int) ([]*DataAudit, error)
}
