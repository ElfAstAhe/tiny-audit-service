package dto

import (
	"time"
)

type AuthAuditDTO struct {
	Source       string    `json:"source,omitempty"`
	EventDate    time.Time `json:"event_date,omitempty"`
	Event        string    `json:"event,omitempty"`
	Status       string    `json:"status,omitempty"`
	RequestID    string    `json:"request_id,omitempty"`
	TraceID      string    `json:"trace_id,omitempty"`
	Username     string    `json:"username,omitempty"`
	AccessToken  string    `json:"access_token,omitempty"`
	RefreshToken string    `json:"refresh_token,omitempty"`
}
