package dto

import (
	"time"
)

type DataAuditValueDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Before      string `json:"before"`
	After       string `json:"after"`
}

type DataAuditDTO struct {
	Source           string               `json:"source,omitempty"`
	EventDate        time.Time            `json:"event_date,omitempty"`
	Event            string               `json:"event,omitempty"`
	Status           string               `json:"status,omitempty"`
	RequestID        string               `json:"request_id,omitempty"`
	TraceID          string               `json:"trace_id,omitempty"`
	Username         string               `json:"username,omitempty"`
	InternalTypeName string               `json:"internal_type_name,omitempty"`
	TypeName         string               `json:"type_name,omitempty"`
	TypeDescription  string               `json:"type_description,omitempty"`
	InstanceID       string               `json:"instance_id,omitempty"`
	InstanceName     string               `json:"instance_name,omitempty"`
	Values           []*DataAuditValueDTO `json:"values,omitempty"`
}
