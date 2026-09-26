package dto

import (
	"github.com/ElfAstAhe/go-service-template/pkg/transport/broker"
)

type LoginAttemptWorkerJob struct {
	Data    *LoginAttemptEventDTO
	Message broker.Message
}
