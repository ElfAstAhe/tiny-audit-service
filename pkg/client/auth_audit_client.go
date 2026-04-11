package client

import (
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

type AuthAuditClient interface {
	AuditClient[*dto.AuthAuditDTO]
}
