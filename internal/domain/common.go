package domain

import (
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// commonAudit определяет внутренний контракт для извлечения сквозных метаданных логов аудита.
type commonAudit interface {
	GetSource() string
	GetEventDate() time.Time
	GetEvent() string
	GetStatus() string
	GetRequestID() string
	GetTraceID() string
	GetUsername() string
}

// Константы, определяющие прикладные статусы выполнения операций аудита.
const (
	AuditStatusSuccess string = "success"
	AuditStatusFail    string = "fail"
)

// Константы, определяющие типы регистрируемых событий аутентификации.
const (
	AuthEventLogin   string = "login"
	AuthEventLogoff  string = "logoff"
	AuthEventRefresh string = "refresh"
	AuthEventRevoke  string = "revoke"
)

// Константы, определяющие типы регистрируемых событий изменения состояния данных.
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

// validateAuditStatus проверяет строковое представление статуса аудита на вхождение в белый список допустимых значений.
func validateAuditStatus(status string) error {
	if _, ok := auditStatuses[status]; ok {
		return nil
	}

	return errs.NewBllValidateError("validateAuditStatus", fmt.Sprintf("unknown audit status [%s]", status), nil)
}

// validateAuthEvent проверяет строковое представление события аутентификации на вхождение в белый список допустимых значений.
func validateAuthEvent(event string) error {
	if _, ok := authEvents[event]; ok {
		return nil
	}

	return errs.NewBllValidateError("validateAuthEvent", fmt.Sprintf("unknown auth event '%s'", event), nil)
}

// validateDataEvent проверяет строковое представление события изменений данных на вхождение в белый список допустимых значений.
func validateDataEvent(event string) error {
	if _, ok := dataEvents[event]; ok {
		return nil
	}

	return errs.NewBllValidateError("validateDataEvent", fmt.Sprintf("unknown data event '%s'", event), nil)
}

// validateCommon выполняет комплексную валидацию общих обязательных полей для любого типа лога аудита.
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
