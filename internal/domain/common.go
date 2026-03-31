package domain

import (
	"fmt"
	"tiny-audit-service/internal/domain/errs"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/google/uuid"
)

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
	newID, err := uuid.NewRandom()
	if err != nil {
		return errs.NewBllError("defaultBeforeCreate", "generate new id", err)
	}

	entity.SetID(newID.String())

	return nil
}
