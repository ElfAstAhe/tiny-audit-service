package dto

import (
	"time"
)

type AuditPeriodDTO struct {
	From   time.Time `json:"from,omitempty"`
	Till   time.Time `json:"till,omitempty"`
	Limit  int       `json:"limit,omitempty"`
	Offset int       `json:"offset,omitempty"`
} // @name AuditPeriodDTO
