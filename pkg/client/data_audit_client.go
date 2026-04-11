package client

import (
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

type DataAuditClient interface {
	AuditClient[*dto.DataAuditDTO]
}
