package domain

import (
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type DataAudit struct {
	ID               string
	Source           string
	EventDate        time.Time
	Event            string
	Status           string
	RequestID        string
	TraceID          string
	Username         string
	InternalTypeName string
	TypeName         string
	TypeDescription  string
	InstanceID       string
	InstanceName     string
	Values           []*DataAuditValue
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

var _ domain.Entity[string] = (*DataAudit)(nil)
var _ commonAudit = (*DataAudit)(nil)

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
	if err := validateCommon(da); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "common audit validate", err)
	}
	if err := validateDataEvent(da.Event); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateCreate", "event validate", err)
	}
	if len(da.Values) > 0 {
		for index, value := range da.Values {
			if err := value.ValidateCreate(); err != nil {
				return errs.NewBllValidateError("DataAudit.ValidateCreate", fmt.Sprintf("values[%d] validate failed", index), err)
			}
		}
	}

	return nil
}

func (da *DataAudit) ValidateChange() error {
	if da.ID == "" {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "id must not be empty", nil)
	}
	if err := validateCommon(da); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "common audit validate", err)
	}
	if err := validateDataEvent(da.Event); err != nil {
		return errs.NewBllValidateError("DataAudit.ValidateChange", "event validate", err)
	}
	if len(da.Values) > 0 {
		for index, value := range da.Values {
			if err := value.ValidateChange(); err != nil {
				return errs.NewBllValidateError("DataAudit.ValidateChange", fmt.Sprintf("values[%d] validate failed", index), err)
			}
		}
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

func (da *DataAudit) GetSource() string {
	return da.Source
}

func (da *DataAudit) GetEventDate() time.Time {
	return da.EventDate
}

func (da *DataAudit) GetEvent() string {
	return da.Event
}

func (da *DataAudit) GetStatus() string {
	return da.Status
}

func (da *DataAudit) GetRequestID() string {
	return da.RequestID
}

func (da *DataAudit) GetTraceID() string {
	return da.TraceID
}

func (da *DataAudit) GetUsername() string {
	return da.Username
}
