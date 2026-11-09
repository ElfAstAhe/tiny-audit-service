package domain

import (
	"time"

	libdomain "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// AuthAudit представляет собой доменную модель (Entity) для хранения
// и валидации записей аудита событий аутентификации пользователей.
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

var _ libdomain.Entity[string] = (*AuthAudit)(nil)
var _ commonAudit = (*AuthAudit)(nil)

// NewEmptyAuthAudit возвращает новый экземпляр AuthAudit с инициализированными метками времени создания и обновления.
func NewEmptyAuthAudit() *AuthAudit {
	return &AuthAudit{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// GetID возвращает уникальный строковый идентификатор текущей записи аудита.
func (aa *AuthAudit) GetID() string {
	return aa.ID
}

// SetID принудительно устанавливает уникальный строковый идентификатор для текущей записи аудита.
func (aa *AuthAudit) SetID(id string) {
	aa.ID = id
}

// IsExists проверяет, была ли данная сущность уже сохранена в базе данных (имеет ли она заполненный ID).
func (aa *AuthAudit) IsExists() bool {
	return aa.ID != ""
}

// BeforeCreate выполняет хук жизненного цикла перед сохранением: генерирует UUIDv7 и выставляет метки времени.
func (aa *AuthAudit) BeforeCreate() error {
	if err := libdomain.AssignUUIDv7(aa); err != nil {
		return errs.NewBllError("AuthAudit.BeforeCreate", "default before create failed", err)
	}

	if aa.CreatedAt.IsZero() {
		aa.CreatedAt = time.Now()
	}
	aa.UpdatedAt = time.Now()

	return nil
}

// BeforeChange выполняет хук жизненного цикла перед обновлением сущности, актуализируя метку времени UpdatedAt.
func (aa *AuthAudit) BeforeChange() error {
	aa.UpdatedAt = time.Now()

	return nil
}

// ValidateCreate осуществляет проверку бизнес-правил доменного слоя, необходимых для успешного создания записи аудита.
func (aa *AuthAudit) ValidateCreate() error {
	if aa.ID != "" {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "id must be empty", nil)
	}
	if err := validateAuthEvent(aa.Event); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "event validate", err)
	}
	if err := validateCommon(aa); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateCreate", "common audit validate", err)
	}

	return nil
}

// ValidateChange осуществляет проверку бизнес-правил доменного слоя, необходимых для успешного обновления записи аудита.
func (aa *AuthAudit) ValidateChange() error {
	if aa.ID == "" {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "id cannot be empty", nil)
	}
	if err := validateAuthEvent(aa.Event); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "event validate", err)
	}
	if err := validateCommon(aa); err != nil {
		return errs.NewBllValidateError("AuthAudit.ValidateChange", "common audit validate", err)
	}

	return nil
}

// GetSource возвращает наименование компонента-источника события.
func (aa *AuthAudit) GetSource() string {
	return aa.Source
}

// GetEventDate возвращает точную дату и время фиксации события.
func (aa *AuthAudit) GetEventDate() time.Time {
	return aa.EventDate
}

// GetEvent возвращает строковый маркер типа события аутентификации.
func (aa *AuthAudit) GetEvent() string {
	return aa.Event
}

// GetStatus возвращает статус завершения операции (success/fail).
func (aa *AuthAudit) GetStatus() string {
	return aa.Status
}

// GetRequestID возвращает сквозной идентификатор HTTP/gRPC запроса.
func (aa *AuthAudit) GetRequestID() string {
	return aa.RequestID
}

// GetTraceID возвращает идентификатор распределенной трассировки OpenTelemetry.
func (aa *AuthAudit) GetTraceID() string {
	return aa.TraceID
}

// GetUsername возвращает имя пользователя, инициировавшего событие.
func (aa *AuthAudit) GetUsername() string {
	return aa.Username
}
