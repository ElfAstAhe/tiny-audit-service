package domain

import (
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// DataAuditValue представляет собой структуру для хранения измененного значения
// конкретного поля сущности, фиксируя состояние до и после операции.
type DataAuditValue struct {
	Name        string
	Description string
	Before      string
	After       string
}

// ValidateCreate осуществляет проверку бизнес-правил доменного слоя, необходимых для создания значения аудита поля.
func (dav *DataAuditValue) ValidateCreate() error {
	if dav.Name == "" {
		return errs.NewBllValidateError("DataAuditValue.ValidateCreate", "Name must not be empty", nil)
	}

	return nil
}

// ValidateChange осуществляет проверку бизнес-правил доменного слоя, необходимых для обновления значения аудита поля.
func (dav *DataAuditValue) ValidateChange() error {
	if dav.Name == "" {
		return errs.NewBllValidateError("DataAuditValue.ValidateChange", "Name must not be empty", nil)
	}

	return nil
}
