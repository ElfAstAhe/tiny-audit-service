package domain

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// AuthAuditOutbox представляет собой доменную модель и системную сущность транзакционного журнала Outbox для событий авторизации.
//
// Структура полностью сопоставляется со схемой таблицы СУБД, инкапсулирует сериализованные байты авторизационных логов
// и реализует полный жизненный цикл хуков валидации и генерации первичных ключей UUIDv7 фреймворка.
type AuthAuditOutbox struct {
	ID          string
	Username    string
	ContentType string
	Payload     []byte
	Status      string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Контролируем реализацию базового контракта Entity фреймворка на этапе компиляции
var _ domain.Entity[string] = (*AuthAuditOutbox)(nil)

// NewEmptyAuthAuditOutbox инициализирует пустой экземпляр сущности со стартовым статусом PENDING.
func NewEmptyAuthAuditOutbox() *AuthAuditOutbox {
	return &AuthAuditOutbox{
		Status: OutboxStatusPending,
	}
}

// NewAuthAuditOutbox — параметризованный конструктор для явного формирования записи Outbox-журнала авторизации.
func NewAuthAuditOutbox(ID string, Username string, ContentType string, Payload []byte, status string) *AuthAuditOutbox {
	return &AuthAuditOutbox{
		ID:          ID,
		Username:    Username,
		ContentType: ContentType,
		Payload:     Payload,
		Status:      status,
	}
}

// GetID возвращает уникальный первичный ключ текущей записи Outbox-журнала.
func (a *AuthAuditOutbox) GetID() string {
	return a.ID
}

// SetID выполняет мутацию и принудительную инициализацию первичного ключа записи.
func (a *AuthAuditOutbox) SetID(id string) {
	a.ID = id
}

// IsExists проверяет персистентный статус сущности (задан ли первичный ключ записи).
func (a *AuthAuditOutbox) IsExists() bool {
	return a.ID != ""
}

// BeforeCreate — инфраструктурный хук (интерцептор), автоматически запускаемый ORM перед фазой вставки (INSERT).
// Выполняет генерацию и присвоение хронологически сортируемого идентификатора UUIDv7 и выравнивает временные метки.
func (a *AuthAuditOutbox) BeforeCreate() error {
	if err := domain.AssignUUIDv7(a); err != nil {
		return errs.NewBllError("AuthAuditOutbox.BeforeCreate", "default before create failed", err)
	}

	if a.CreatedAt.IsZero() {
		a.CreatedAt = time.Now()
	}
	a.UpdatedAt = time.Now()

	return nil
}

// BeforeChange — infrastructure-хук (интерцептор), автоматически запускаемый ORM перед фазой обновления статусов (UPDATE).
func (a *AuthAuditOutbox) BeforeChange() error {
	a.UpdatedAt = time.Now()

	return nil
}

// ValidateCreate выполняет доменную валидацию инвариантов и ограничений модели на этапе создания записи.
func (a *AuthAuditOutbox) ValidateCreate() error {
	if a.ID != "" {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateCreate", "id must be empty", nil)
	}
	if a.Username == "" {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateCreate", "username required", nil)
	}
	if a.ContentType == "" {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateCreate", "content_type required", nil)
	}
	if len(a.Payload) == 0 {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateCreate", "payload required", nil)
	}
	if err := ValidateOutboxStatus(a.Status); err != nil {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateCreate", "status validate failed", err)
	}

	return nil
}

// ValidateChange выполняет доменную валидацию инвариантов и ограничений модели на этапе модификации полей и перевода статусов доставки.
func (a *AuthAuditOutbox) ValidateChange() error {
	if a.ID == "" {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateChange", "id cannot be empty", nil)
	}
	if a.Username == "" {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateChange", "username required", nil)
	}
	if a.ContentType == "" {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateChange", "content_type required", nil)
	}
	if len(a.Payload) == 0 {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateChange", "payload required", nil)
	}
	if err := ValidateOutboxStatus(a.Status); err != nil {
		return errs.NewBllValidateError("AuthAuditOutbox.ValidateChange", "status validate failed", err)
	}

	return nil
}
