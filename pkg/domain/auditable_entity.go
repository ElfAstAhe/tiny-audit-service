package domain

import (
	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
)

type AuditableEntity[ID comparable] interface {
	libdom.Entity[ID]
	Auditable
}
