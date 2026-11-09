package domain

import (
	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
)

// AuditableEntity маркерный интерфейс, объединяет в себе domain.Entity и Auditable
type AuditableEntity[ID comparable] interface {
	libdom.Entity[ID]
	Auditable
}
