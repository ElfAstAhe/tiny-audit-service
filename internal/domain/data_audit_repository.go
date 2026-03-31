package domain

import (
	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

type DataAuditRepository interface {
	domain.CRUDRepository[*DataAudit, string]
}
