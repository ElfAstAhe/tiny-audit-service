package dto

import (
	"time"
)

type DataAuditDTO struct {
	ID              string               `json:"id,omitempty"`
	Source          string               `json:"source,omitempty"`
	EventDate       time.Time            `json:"event_date,omitempty"`
	Event           string               `json:"event,omitempty"`
	Status          string               `json:"status,omitempty"`
	RequestID       string               `json:"request_id,omitempty"`
	TraceID         string               `json:"trace_id,omitempty"`
	Username        string               `json:"username,omitempty"`
	TypeName        string               `json:"type_name,omitempty"`
	TypeDescription string               `json:"type_description,omitempty"`
	InstanceID      string               `json:"instance_id,omitempty"`
	InstanceName    string               `json:"instance_name,omitempty"`
	Values          []*DataAuditValueDTO `json:"values,omitempty"`
	CreatedAt       time.Time            `json:"created_at,omitempty"`
	UpdatedAt       time.Time            `json:"updated_at,omitempty"`
} // @name DataAuditDTO
