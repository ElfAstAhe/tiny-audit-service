package client

import (
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

// AuthAuditClient определяет специализированный контракт для клиента,
// осуществляющего отправку объектов аудит-логов аутентификации dto.AuthAuditDTO.
type AuthAuditClient interface {
	AuditClient[*dto.AuthAuditDTO]
}
