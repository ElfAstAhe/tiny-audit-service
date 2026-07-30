package dto

import (
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

type AuthAuditDTO struct {
	*dto.AuthAuditDTO
	Message libamqp.Message
}

func NewAuthAudit(data *dto.AuthAuditDTO, msg libamqp.Message) *AuthAuditDTO {
	return &AuthAuditDTO{
		AuthAuditDTO: data,
		Message:      msg,
	}
}
