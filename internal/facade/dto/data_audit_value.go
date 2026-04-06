package dto

type DataAuditValueDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Before      string `json:"before"`
	After       string `json:"after"`
} // @name DataAuditValueDTO
