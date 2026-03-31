package domain

import (
	"time"
	"tiny-audit-service/internal/domain/errs"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

type DataAudit struct {
	ID              string
	Source          string
	EventDate       time.Time
	Event           string
	Status          string
	RequestID       string
	Username        string
	TypeName        string
	TypeDescription string
	InstanceID      string
	InstanceName    string
	Values          []*DataAuditValue
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

var _ domain.Entity[string] = (*DataAudit)(nil)

func NewEmptyDataAudit() *DataAudit {
	return &DataAudit{
		Values:    make([]*DataAuditValue, 0),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (da *DataAudit) GetID() string {
	return da.ID
}

func (da *DataAudit) SetID(id string) {
	da.ID = id
}

func (da *DataAudit) IsExists() bool {
	return da.ID != ""
}

func (da *DataAudit) ValidateCreate() error {
	if da.ID != "" {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "id must be empty", nil)
	}
	if da.Source != "" {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "source must be empty", nil)
	}
	if da.EventDate.IsZero() {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "event_date must not be empty", nil)
	}
	if err := validateDataEvent(da.Event); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "event validate", err)
	}
	if err := validateAuditStatus(da.Status); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "status validate", err)
	}
	if da.Username == "" {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "username cannot be empty", nil)
	}

	return nil
}

func (da *DataAudit) ValidateChange() error {
	if da.ID == "" {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "id must be empty", nil)
	}
	if da.Source != "" {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "source must be empty", nil)
	}
	if da.EventDate.IsZero() {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "event_date must not be empty", nil)
	}
	if err := validateDataEvent(da.Event); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "event validate", err)
	}
	if err := validateAuditStatus(da.Status); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "status validate", err)
	}
	if da.Username == "" {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "username cannot be empty", nil)
	}

	return nil
}

func (da *DataAudit) BeforeCreate() error {
	if err := defaultBeforeCreate(da); err != nil {
		return errs.NewBllError("DataAudit.BeforeCreate", "default before create failed", err)
	}

	if da.CreatedAt.IsZero() {
		da.CreatedAt = time.Now()
	}
	da.UpdatedAt = time.Now()

	return nil
}

func (da *DataAudit) BeforeChange() error {
	da.UpdatedAt = time.Now()

	return nil
}
