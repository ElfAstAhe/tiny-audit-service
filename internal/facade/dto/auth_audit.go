package dto

import (
	"time"
)

type AuthAuditDTO struct {
	ID           string    `json:"id,omitempty"`
	Source       string    `json:"source,omitempty"`
	EventDate    time.Time `json:"event_date,omitempty"`
	Event        string    `json:"event,omitempty"`
	Status       string    `json:"status,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	Username     string    `json:"username,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
	CreatedAt    time.Time `json:"created_at,omitempty"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
} // @name AuthAuditDTO
