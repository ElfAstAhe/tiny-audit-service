package facade

import (
	"context"
	"time"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

type DataAuditFacade interface {
	Audit(ctx context.Context, data *dto.DataAuditDTO) error
	ListByPeriod(ctx context.Context, from, till time.Time, limit int) ([]*dto.DataAuditDTO, error)
	ListByInstance(ctx context.Context, typeName string, instanceID string, limit int) ([]*dto.DataAuditDTO, error)
}
