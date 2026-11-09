package dto

import (
	"time"
)

// LoginAttemptEventDTO представляет собой объект передачи данных (DTO) для десериализации
// поступающих из брокера сообщений событий о попытках входа пользователей в систему.
type LoginAttemptEventDTO struct {
	NodeName  string    `json:"node_name"`
	EventDate time.Time `json:"event_date"`
	RequestID string    `json:"request_id,omitempty"`
	TraceID   string    `json:"trace_id,omitempty"`
	IP        string    `json:"ip,omitempty"`
	Username  string    `json:"username"`
	Success   bool      `json:"success"`
	Error     string    `json:"error"`
}
