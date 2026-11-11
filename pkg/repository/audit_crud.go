package repository

import (
	"context"
	"fmt"
	"time"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	libextdom "github.com/ElfAstAhe/go-service-template/pkg/ext/domain"
	libutl "github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	pkgdom "github.com/ElfAstAhe/tiny-audit-service/pkg/domain"
	pkgerrs "github.com/ElfAstAhe/tiny-audit-service/pkg/errs"
	pkgutl "github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
)

// AuditableMapperFunc определяет сигнатуру функции маппинга сущности репозитория в интерфейс аудита pkgdom.Auditable.
type AuditableMapperFunc[E libdom.Entity[ID], ID comparable] func(entity E) (pkgdom.Auditable, error)

// OutboxBuilderFunc определяет сигнатуру функции конструирования сущности domain.Outbox.
type OutboxBuilderFunc[O libextdom.Outbox[ID], ID comparable] func(ctx context.Context, builder *pkgutl.DataAuditBuilder) (O, error)

// AuditCRUD представляет собой универсальный (generic) декоратор CRUD-репозитория, перехватывающий
// операции изменения данных для генерации и сохранения записей аудита через паттерн Transactional Outbox.
type AuditCRUD[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable] struct {
	next            libdom.CRUDRepository[E, ID]
	source          string
	auditableMapper AuditableMapperFunc[E, ID]
	outboxBuilder   OutboxBuilderFunc[O, ID]
	outboxRepo      libextdom.OutboxRepository[O, ID]
	nilInstance     E
	opts            *AuditCRUDOptions[E, O, ID]
}

// NewAuditCRUD создает и инициализирует новый экземпляр декоратора AuditCRUD.
func NewAuditCRUD[E pkgdom.AuditableEntity[ID], O libextdom.Outbox[ID], ID comparable](
	options ...AuditCRUDOption[E, O, ID],
) (*AuditCRUD[E, O, ID], error) {
	opts := NewAuditCRUDOptions[E, O, ID]()
	for _, option := range options {
		option(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, pkgerrs.NewAuditError("NewAuditCRUD", "options validation failed", err)
	}

	return &AuditCRUD[E, O, ID]{
		opts:            opts,
		source:          opts.Source,
		next:            opts.Next,
		auditableMapper: opts.AuditableMapper,
		outboxBuilder:   opts.OutboxBuilder,
		outboxRepo:      opts.OutboxRepo,
	}, nil
}

// Find выполняет сквозной поиск сущности по её уникальному идентификатору через вложенный репозиторий.
func (acr *AuditCRUD[E, O, ID]) Find(ctx context.Context, id ID) (E, error) {
	return acr.next.Find(ctx, id)
}

// List возвращает список сущностей с учетом ограничений лимита и смещения через вложенный репозиторий.
func (acr *AuditCRUD[E, O, ID]) List(ctx context.Context, limit, offset int) ([]E, error) {
	return acr.next.List(ctx, limit, offset)
}

// Create создает новую сущность, формирует запись аудит-лога о создании и регистрирует её в таблице outbox.
func (acr *AuditCRUD[E, O, ID]) Create(ctx context.Context, entity E) (E, error) {
	// action
	res, err := acr.next.Create(ctx, entity)
	if err != nil {
		return res, err
	}

	// transform
	auditEntity, err := acr.mapToAuditable(res)
	if err != nil {
		return acr.nilInstance, pkgerrs.NewAuditError("AuditCRUD[E, ID].Create", "after entity to auditable transform failed", err)
	}

	// builder, data
	builder := acr.builder(ctx, auditEntity, nil).
		WithEvent(dto.DataEventCreate).
		WithValues(pkgutl.BuildSingleDataAuditValues(auditEntity, false))

	// audit
	err = acr.audit(ctx, builder)
	if err != nil {
		return acr.nilInstance, pkgerrs.NewAuditError("AuditCRUD[E, ID].Create", fmt.Sprintf("audit create entity with id [%v] failed", res.GetID()), err)
	}

	return res, err
}

// Change фиксирует состояние сущности до модификации, выполняет обновление, фиксирует измененное состояние
// и регистрирует запись аудит-лога со срезом изменений (diff) в таблице outbox.
func (acr *AuditCRUD[E, O, ID]) Change(ctx context.Context, entity E) (E, error) {
	// get before entity
	beforeEntity, beforeErr := acr.next.Find(ctx, entity.GetID())
	if beforeErr != nil {
		return acr.nilInstance, pkgerrs.NewAuditError("AuditCRUD[E, ID].Change", fmt.Sprintf("unable to find entity with id [%v]", entity.GetID()), beforeErr)
	}

	// action
	res, err := acr.next.Change(ctx, entity)
	if err != nil {
		return res, pkgerrs.NewAuditError("AuditCRUD[E, ID].Change", fmt.Sprintf("apply change entity with id [%v] failed", entity.GetID()), err)
	}

	// transform
	auditBeforeEntity, err := acr.mapToAuditable(beforeEntity)
	if err != nil {
		return acr.nilInstance, pkgerrs.NewAuditError("AuditCRUD[E, ID].Change", "before entity to auditable transform failed", err)
	}
	auditAfterEntity, err := acr.mapToAuditable(res)
	if err != nil {
		return acr.nilInstance, pkgerrs.NewAuditError("AuditCRUD[E, ID].Change", "after entity to auditable transform failed", err)
	}

	// builder
	builder := acr.builder(ctx, auditAfterEntity, err).
		WithEvent(dto.DataEventChange)

	// data
	// case 1 - all data in our hands
	builder.WithValues(pkgutl.BuildBothDataAuditValues(auditBeforeEntity, auditAfterEntity))
	// audit
	err = acr.audit(ctx, builder)
	if err != nil {
		return acr.nilInstance, pkgerrs.NewAuditError("AuditCRUD[E, ID].Change", fmt.Sprintf("audit change entity with id [%v] failed", res.GetID()), err)
	}

	return res, nil
}

// Delete перехватывает фазу удаления сущности, фиксирует ее финальное состояние до очистки записей БД и отправляет лог.
func (acr *AuditCRUD[E, O, ID]) Delete(ctx context.Context, id ID) error {
	// get before entity
	beforeEntity, beforeErr := acr.next.Find(ctx, id)
	if beforeErr != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].Delete", fmt.Sprintf("unable to find entity with id [%v]", id), beforeErr)
	}
	// action
	err := acr.next.Delete(ctx, id)
	if err != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].Delete", fmt.Sprintf("unable to delete entity with id [%v]", id), err)
	}
	// auditable
	auditEntity, err := acr.mapToAuditable(beforeEntity)
	if err != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].Delete", "before entity to auditable transform failed", err)
	}

	// builder, data
	builder := acr.builder(ctx, auditEntity, nil).
		WithEvent(dto.DataEventRemove).
		WithValues(pkgutl.BuildSingleDataAuditValues(auditEntity, true))

	// audit
	err = acr.audit(ctx, builder)
	if err != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].Delete", fmt.Sprintf("audit delete entity with id [%v] failed", id), beforeErr)
	}

	return err
}

func (acr *AuditCRUD[E, O, ID]) mapToAuditable(entity E) (pkgdom.Auditable, error) {
	if res, ok := any(entity).(pkgdom.Auditable); ok {
		return res, nil
	}
	if libutl.IsNil(acr.auditableMapper) {
		return nil, pkgerrs.NewAuditError("AuditCRUD[E, ID].mapToAuditable", "custom auditable mapper unassigned", nil)
	}

	res, err := acr.auditableMapper(entity)
	if err != nil {
		return nil, pkgerrs.NewAuditError("AuditCRUD[E, ID].mapToAuditable", "custom auditable mapper failed", err)
	}

	return res, nil
}

func (acr *AuditCRUD[E, O, ID]) builder(ctx context.Context, auditEntity pkgdom.Auditable, err error) *pkgutl.DataAuditBuilder {
	return pkgutl.NewDataAuditBuilder().
		// common
		WithSource(acr.source).
		WithEventDate(time.Now()).
		WithStatus(pkgutl.GetAuditStatus(err)).
		// entity info
		WithInternalTypeName(auditEntity.GetInternalTypeName()).
		WithTypeName(auditEntity.GetTypeName()).
		WithTypeDescription(auditEntity.GetTypeDescription()).
		WithInstanceID(auditEntity.GetInstanceID()).
		WithInstanceName(auditEntity.GetInstanceName()).
		// request
		WithUsername(pkgutl.UsernameFromContext(ctx)).
		WithTraceID(pkgutl.TraceIDFromContext(ctx)).
		WithRequestID(pkgutl.RequestIDFromContext(ctx))
}

func (acr *AuditCRUD[E, O, ID]) audit(ctx context.Context, builder *pkgutl.DataAuditBuilder) error {
	// outbox
	entity, err := acr.outboxBuilder(ctx, builder)
	if err != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].audit", "build outbox entity failed", err)
	}

	// save
	err = acr.outboxRepo.Save(ctx, entity)
	if err != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].audit", "save audit outbox failed", err)
	}

	return nil
}

// GetSource возвращает строковый идентификатор источника событий (компонента), настроенного для данного репозитория.
func (acr *AuditCRUD[E, O, ID]) GetSource() string {
	return acr.source
}

// Unwrap возвращает исходный декорируемый экземпляр CRUD-репозитория, находящийся под слоем аудита.
func (acr *AuditCRUD[E, O, ID]) Unwrap() libdom.CRUDRepository[E, ID] {
	return acr.next
}

// GetOutboxRepository возвращает используемый репозиторий таблицы outbox, в который записываются логи телеметрии изменений.
func (acr *AuditCRUD[E, O, ID]) GetOutboxRepository() libextdom.OutboxRepository[O, ID] {
	return acr.outboxRepo
}

// GetNilInstance возвращает пустое (дефолтное) значение сущности E для безопасного возврата при возникновении критических ошибок.
func (acr *AuditCRUD[E, O, ID]) GetNilInstance() E {
	return acr.nilInstance
}

// GetOpts возвращает набор опций, примененных при инициализации декоратора AuditCRUD.
func (acr *AuditCRUD[E, O, ID]) GetOpts() *AuditCRUDOptions[E, O, ID] {
	return acr.opts
}
