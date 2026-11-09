package dto

// DataAuditValueDTO представляет собой объект передачи данных (DTO) для сериализации
// изменений конкретного поля сущности, фиксируя его наименование, описание и состояния до и после мутации.
type DataAuditValueDTO struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Before      string `json:"before"`
	After       string `json:"after"`
} // @name DataAuditValueDTO
