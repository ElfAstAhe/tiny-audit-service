package dto

import (
	"time"
)

// AuthAuditDTO представляет собой объект передачи данных (DTO) для сериализации
// и десериализации записей аудита событий аутентификации пользователей на внешнем API-слое.
type AuthAuditDTO struct {
	ID           string    `json:"id,omitempty"`
	Source       string    `json:"source,omitempty"`
	EventDate    time.Time `json:"event_date,omitzero"`
	Event        string    `json:"event,omitempty"`
	Status       string    `json:"status,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	TraceID      string    `json:"trace_id,omitempty"`
	Username     string    `json:"username,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitzero"`
	UpdatedAt    time.Time `json:"updated_at,omitzero"`
} // @name AuthAuditDTO
