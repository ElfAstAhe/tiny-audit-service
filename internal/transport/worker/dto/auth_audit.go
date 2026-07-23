package dto

import (
	"github.com/Azure/go-amqp"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

type AuthAuditDTO struct {
	*dto.AuthAuditDTO
	Message *libamqp.Message[*amqp.MessageHeader]
}
