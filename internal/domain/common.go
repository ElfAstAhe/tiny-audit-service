package domain

import (
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
	"github.com/google/uuid"
)

type commonAudit interface {
	GetSource() string
	GetEventDate() time.Time
	GetEvent() string
	GetStatus() string
	GetRequestID() string
	GetTraceID() string
	GetUsername() string
}

const (
	AuditStatusSuccess string = "success"
	AuditStatusFail    string = "fail"
)

const (
	AuthEventLogin   string = "login"
	AuthEventLogoff  string = "logoff"
	AuthEventRefresh string = "refresh"
	AuthEventRevoke  string = "revoke"
)

const (
	DataEventCreate string = "create"
	DataEventChange string = "change"
	DataEventRemove string = "remove"
)

var (
	auditStatuses = map[string]struct{}{
		AuditStatusSuccess: {},
		AuditStatusFail:    {},
	}
	authEvents = map[string]struct{}{
		AuthEventLogin:   {},
		AuthEventLogoff:  {},
		AuthEventRefresh: {},
		AuthEventRevoke:  {},
	}
	dataEvents = map[string]struct{}{
		DataEventCreate: {},
		DataEventChange: {},
		DataEventRemove: {},
	}
)

func validateAuditStatus(status string) error {
	if _, ok := auditStatuses[status]; ok {
		return nil
	}

	return errs.NewBllValidateError("validateAuditStatus", fmt.Sprintf("unknown audit status [%s]", status), nil)
}

func validateAuthEvent(event string) error {
	if _, ok := authEvents[event]; ok {
		return nil
	}

	return errs.NewBllValidateError("validateAuthEvent", fmt.Sprintf("unknown auth event '%s'", event), nil)
}

func validateDataEvent(event string) error {
	if _, ok := dataEvents[event]; ok {
		return nil
	}

	return errs.NewBllValidateError("validateDataEvent", fmt.Sprintf("unknown data event '%s'", event), nil)
}

func defaultBeforeCreate(entity domain.Entity[string]) error {
	newID, err := uuid.NewV7()
	if err != nil {
		return errs.NewBllError("defaultBeforeCreate", "generate new id", err)
	}

	entity.SetID(newID.String())

	return nil
}

func validateCommon(commonAudit commonAudit) error {
	if commonAudit.GetSource() == "" {
		return errs.NewBllValidateError("validateCommon", "source must not be empty", nil)
	}
	if commonAudit.GetEventDate().IsZero() {
		return errs.NewBllValidateError("validateCommon", "event_date must not be empty", nil)
	}
	if err := validateAuditStatus(commonAudit.GetStatus()); err != nil {
		return errs.NewBllValidateError("validateCommon", "status validate", err)
	}
	if commonAudit.GetUsername() == "" {
		return errs.NewBllValidateError("validateCommon", "username cannot be empty", nil)
	}

	return nil
}
