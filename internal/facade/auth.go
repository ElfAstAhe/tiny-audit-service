package facade

import (
	"context"
	"time"

	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

type AuthAuditFacade interface {
	Audit(ctx context.Context, data *dto.AuthAuditDTO) error
	ListByPeriod(ctx context.Context, from, till time.Time, limit int, offset int) ([]*dto.AuthAuditDTO, error)
	ListByUsername(ctx context.Context, username string, limit int, offset int) ([]*dto.AuthAuditDTO, error)
}
