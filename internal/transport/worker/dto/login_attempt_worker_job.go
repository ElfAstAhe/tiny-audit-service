package dto

import (
	"github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
)

type LoginAttemptWorkerJob struct {
	Data    *LoginAttemptEventDTO
	Message amqp.Message
}
