package repository

import (
	"strings"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	libextdom "github.com/ElfAstAhe/go-service-template/pkg/ext/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	pkgdom "github.com/ElfAstAhe/tiny-audit-service/pkg/domain"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/errs"
)

// AuditCRUDOption определяет сигнатуру функциональной опции для конфигурирования структуры AuditCRUDOptions.
type AuditCRUDOption[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable] func(*AuditCRUDOptions[E, O, ID])

// AuditCRUDOptions объединяет зависимости и параметры конфигурации, необходимые для работы
// декоратора CRUD-репозитория с поддержкой аудита изменений и паттерна Transactional Outbox.
type AuditCRUDOptions[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable] struct {
	Source          string
	Next            libdom.CRUDRepository[E, ID]
	AuditableMapper AuditableMapperFunc[E, ID]
	PayloadBuilder  OutboxPayloadBuilderFunc
	OutboxBuilder   OutboxBuilderFunc[O, ID]
	OutboxRepo      libextdom.OutboxRepository[O, ID]
}

// NewAuditCRUDOptions создает и инициализирует новый пустой экземпляр настроек AuditCRUDOptions.
func NewAuditCRUDOptions[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable]() *AuditCRUDOptions[E, O, ID] {
	return &AuditCRUDOptions[E, O, ID]{}
}

// Validate проверяет корректность и наличие всех обязательных зависимостей в AuditCRUDOptions.
func (aco *AuditCRUDOptions[E, O, ID]) Validate() error {
	if strings.TrimSpace(aco.Source) == "" {
		return errs.NewAuditError("AuditCRUDOptions[E, O, ID].Validate", "source required", nil)
	}
	if utils.IsNil(aco.Next) {
		return errs.NewAuditError("AuditCRUDOptions[E, O, ID].Validate", "next repository required", nil)
	}
	if utils.IsNil(aco.PayloadBuilder) {
		return errs.NewAuditError("AuditCRUDOptions[E, O, ID].Validate", "payload builder required", nil)
	}
	if utils.IsNil(aco.OutboxBuilder) {
		return errs.NewAuditError("AuditCRUDOptions[E, O, ID].Validate", "outbox builder required", nil)
	}

	return nil
}

// ================================================================
// Fluent API
// ================================================================

// WithAuditSource возвращает опцию для установки имени источника событий аудита.
func WithAuditSource[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable](source string) AuditCRUDOption[E, O, ID] {
	return func(options *AuditCRUDOptions[E, O, ID]) {
		options.Source = source
	}
}

// WithAuditNext возвращает опцию для передачи декорируемого базового CRUD-репозитория.
func WithAuditNext[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable](next libdom.CRUDRepository[E, ID]) AuditCRUDOption[E, O, ID] {
	return func(options *AuditCRUDOptions[E, O, ID]) {
		options.Next = next
	}
}

// WithAuditAuditableMapper возвращает опцию для настройки функции маппинга сущности в модель аудита.
func WithAuditAuditableMapper[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable](mapper AuditableMapperFunc[E, ID]) AuditCRUDOption[E, O, ID] {
	return func(options *AuditCRUDOptions[E, O, ID]) {
		options.AuditableMapper = mapper
	}
}

// WithAuditPayloadBuilder возвращает опцию для передачи конструктора полезной нагрузки (payload) сообщения Outbox.
func WithAuditPayloadBuilder[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable](builder OutboxPayloadBuilderFunc) AuditCRUDOption[E, O, ID] {
	return func(options *AuditCRUDOptions[E, O, ID]) {
		options.PayloadBuilder = builder
	}
}

// WithAuditOutboxBuilder возвращает опцию для настройки фабрики создания сущности Outbox.
func WithAuditOutboxBuilder[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable](builder OutboxBuilderFunc[O, ID]) AuditCRUDOption[E, O, ID] {
	return func(options *AuditCRUDOptions[E, O, ID]) {
		options.OutboxBuilder = builder
	}
}

// WithAuditOutboxRepository возвращает опцию для привязки репозитория сохранения записей Transactional Outbox.
func WithAuditOutboxRepository[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable](repo libextdom.OutboxRepository[O, ID]) AuditCRUDOption[E, O, ID] {
	return func(options *AuditCRUDOptions[E, O, ID]) {
		options.OutboxRepo = repo
	}
}
