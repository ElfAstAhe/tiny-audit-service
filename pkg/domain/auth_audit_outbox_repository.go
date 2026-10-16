package domain

import "github.com/ElfAstAhe/go-service-template/pkg/domain"

type AuthAuditOutboxRepository interface {
	*domain.CRUDRepository[*AuthAuditOutbox, string]
}
