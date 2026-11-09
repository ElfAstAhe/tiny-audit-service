package dto

import (
	"github.com/ElfAstAhe/go-service-template/pkg/transport/broker"
)

// LoginAttemptWorkerJob представляет собой структуру задачи для воркера,
// объединяющую десериализованные данные события аутентификации и исходное сообщение из брокера.
type LoginAttemptWorkerJob struct {
	Data    *LoginAttemptEventDTO
	Message broker.Message
}
