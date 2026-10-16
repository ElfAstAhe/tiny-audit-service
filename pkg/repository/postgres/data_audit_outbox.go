package postgres

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/db"
	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/repository"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/domain"
)

const (
	sourceLabelListPending string = "Pending"
)

type DataAuditOutboxRepository struct {
	*repository.BaseCRUDRepository[*domain.DataAuditOutbox, string]
}

var _ libdom.CRUDRepository[*domain.DataAuditOutbox, string] = (*DataAuditOutboxRepository)(nil)
var _ domain.DataAuditOutboxRepository = (*DataAuditOutboxRepository)(nil)
var _ domain.OutboxRepository[*domain.DataAuditOutbox, string] = (*DataAuditOutboxRepository)(nil)

func NewDataAuditOutboxRepository(
	querier db.Querier,
	decipher db.ErrorDecipher,
) (*DataAuditOutboxRepository, error) {
	// instance
	res := &DataAuditOutboxRepository{}
	// sql
	// callbacks
	// embed
	// setup

	return res, nil
}

func (dao *DataAuditOutboxRepository) ListPending(ctx context.Context, limit int) ([]*domain.DataAuditOutbox, error) {

}
