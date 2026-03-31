package domain

import (
	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

type AuthAuditRepository interface {
	domain.CRUDRepository[*AuthAudit, string]
}
