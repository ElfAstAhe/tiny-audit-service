package domain

import (
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type DataAuditValue struct {
	Name        string
	Description string
	Before      string
	After       string
}

func (dav *DataAuditValue) ValidateCreate() error {
	if dav.Name == "" {
		return errs.NewBllValidateError("DataAuditValue.ValidateCreate", "Name must not be empty", nil)
	}

	return nil
}

func (dav *DataAuditValue) ValidateChange() error {
	if dav.Name == "" {
		return errs.NewBllValidateError("DataAuditValue.ValidateChange", "Name must not be empty", nil)
	}

	return nil
}
