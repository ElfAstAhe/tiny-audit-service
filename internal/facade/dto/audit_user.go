package dto

// AuditUserDTO представляет собой объект передачи данных (DTO) для формирования
// поискового запроса записей аудита по имени конкретного пользователя с поддержкой пагинации.
type AuditUserDTO struct {
	Username string `json:"username,omitempty"`
	Limit    int    `json:"limit,omitempty"`
	Offset   int    `json:"offset,omitempty"`
} // @name AuditUserDTO
