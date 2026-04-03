package dto

type AuditInstanceDTO struct {
	TypeName   string `json:"type_name,omitempty"`
	InstanceID string `json:"instance_id,omitempty"`
	Limit      int    `json:"limit,omitempty"`
	Offset     int    `json:"offset,omitempty"`
} // @name AuditInstanceDTO
