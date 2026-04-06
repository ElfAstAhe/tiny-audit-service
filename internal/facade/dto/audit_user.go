package dto

type AuditUserDTO struct {
	Username string `json:"username,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
} // @name AuditUserDTO
