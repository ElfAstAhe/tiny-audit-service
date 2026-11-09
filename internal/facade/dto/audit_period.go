package dto

import (
	"time"
)

// AuditPeriodDTO представляет собой объект передачи данных (DTO) для формирования
// поисковых запросов и фильтрации записей аудита за определенный временной интервал с поддержкой пагинации.
type AuditPeriodDTO struct {
	From   time.Time `json:"from,omitzero"`
	Till   time.Time `json:"till,omitzero"`
	Limit  int       `json:"limit,omitempty"`
	Offset int       `json:"offset,omitempty"`
} // @name AuditPeriodDTO
