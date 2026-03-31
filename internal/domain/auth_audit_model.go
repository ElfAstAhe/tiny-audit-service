package domain

import (
	"time"
	"tiny-audit-service/internal/domain/errs"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
)

type AuthAudit struct {
	ID           string
	Source       string
	EventDate    time.Time
	Event        string
	Status       string
	RequestID    string
	Username     string
	AccessToken  string
	RefreshToken string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

var _ domain.Entity[string] = (*AuthAudit)(nil)

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
		return errs.NewBllError("DataAudit.BeforeCreate", "default before create failed", err)
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
	if aa.Source != "" {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "source must be empty", nil)
	}
	if aa.EventDate.IsZero() {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "event_date must not be empty", nil)
	}
	if err := validateAuthEvent(aa.Event); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "event validate", err)
	}
	if err := validateAuditStatus(aa.Status); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "status validate", err)
	}
	if aa.Username == "" {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "username cannot be empty", nil)
	}

	return nil
}

func (aa *AuthAudit) ValidateChange() error {
	if aa.ID == "" {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "id cannot be empty", nil)
	}
	if aa.EventDate.IsZero() {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "event_date must not be empty", nil)
	}
	if err := validateAuthEvent(aa.Event); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "event validate", err)
	}
	if err := validateAuditStatus(aa.Status); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "status validate", err)
	}
	if aa.Username == "" {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "username cannot be empty", nil)
	}

	return nil
}
