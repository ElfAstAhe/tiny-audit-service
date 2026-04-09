package domain

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type AuthAudit struct {
	ID           string
	Source       string
	EventDate    time.Time
	Event        string
	Status       string
	RequestID    string
	TraceID      string
	Username     string
	AccessToken  string
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var _ domain.Entity[string] = (*AuthAudit)(nil)
var _ commonAudit = (*AuthAudit)(nil)

func NewEmptyAuthAudit() *AuthAudit {
	return &AuthAudit{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func (aa *AuthAudit) GetID() string {
	return aa.ID
}

func (aa *AuthAudit) SetID(id string) {
	aa.ID = id
}

func (aa *AuthAudit) IsExists() bool {
	return aa.ID != ""
}

func (aa *AuthAudit) BeforeCreate() error {
	if err := defaultBeforeCreate(aa); err != nil {
		return errs.NewBllError("AuthAudit.BeforeCreate", "default before create failed", err)
	}

	if aa.CreatedAt.IsZero() {
		aa.CreatedAt = time.Now()
	}
	aa.UpdatedAt = time.Now()

	return nil
}

func (aa *AuthAudit) BeforeChange() error {
	aa.UpdatedAt = time.Now()

	return nil
}

func (aa *AuthAudit) ValidateCreate() error {
	if aa.ID != "" {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "id must be empty", nil)
	}
	if err := validateAuthEvent(aa.Event); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "event validate", err)
	}
	if err := validateCommon(aa); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "common.go audit validate", err)
	}

	return nil
}

func (aa *AuthAudit) ValidateChange() error {
	if aa.ID == "" {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "id cannot be empty", nil)
	}
	if err := validateAuthEvent(aa.Event); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "event validate", err)
	}
	if err := validateCommon(aa); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "common.go audit validate", err)
	}

	return nil
}

func (aa *AuthAudit) GetSource() string {
	return aa.Source
}

func (aa *AuthAudit) GetEventDate() time.Time {
	return aa.EventDate
}

func (aa *AuthAudit) GetEvent() string {
	return aa.Event
}

func (aa *AuthAudit) GetStatus() string {
	return aa.Status
}

func (aa *AuthAudit) GetRequestID() string {
	return aa.RequestID
}

func (aa *AuthAudit) GetTraceID() string {
	return aa.TraceID
}

func (aa *AuthAudit) GetUsername() string {
	return aa.Username
}
