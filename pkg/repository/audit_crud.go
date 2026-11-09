package repository

import (
	"context"
	"fmt"
	"time"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	libextdom "github.com/ElfAstAhe/go-service-template/pkg/ext/domain"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	pkgdom "github.com/ElfAstAhe/tiny-audit-service/pkg/domain"
	pkgerrs "github.com/ElfAstAhe/tiny-audit-service/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
)

// AuditableMapperFunc определяет сигнатуру функции маппинга сущности репозитория в интерфейс аудита pkgdom.Auditable.
type AuditableMapperFunc[E libdom.Entity[ID], ID comparable] func(entity E) pkgdom.Auditable

// OutboxBuilderFunc определяет сигнатуру функции конструирования сущности outbox на основе объекта dto.DataAuditDTO.
type OutboxBuilderFunc[E libextdom.Outbox[ID], ID comparable] func(context.Context, *dto.DataAuditDTO) (E, error)

// AuditCRUD представляет собой универсальный (generic) декоратор CRUD-репозитория, перехватывающий
// операции изменения данных для генерации и сохранения записей аудита через паттерн Transactional Outbox.
type AuditCRUD[E pkgdom.AuditableEntity[ID], ID comparable] struct {
	next         libdom.CRUDRepository[E, ID]
	source       string
	mapAuditable AuditableMapperFunc[E, ID]
	buildOutbox  OutboxBuilderFunc[libextdom.Outbox[ID], ID]
	outboxRepo   libextdom.OutboxRepository[libextdom.Outbox[ID], ID]
	nilInstance  E
}

// NewAuditCRUD создает и инициализирует новый экземпляр декоратора AuditCRUD.
func NewAuditCRUD[E pkgdom.AuditableEntity[ID], ID comparable](
	next libdom.CRUDRepository[E, ID],
	source string,
	mapAuditable AuditableMapperFunc[E, ID],
	buildOutbox OutboxBuilderFunc[libextdom.Outbox[ID], ID],
	outboxRepo libextdom.OutboxRepository[libextdom.Outbox[ID], ID],
) *AuditCRUD[E, ID] {
	return &AuditCRUD[E, ID]{
		next:         next,
		source:       source,
		mapAuditable: mapAuditable,
		buildOutbox:  buildOutbox,
		outboxRepo:   outboxRepo,
	}
}

// Find выполняет сквозной поиск сущности по её уникальному идентификатору через вложенный репозиторий.
func (acr *AuditCRUD[E, ID]) Find(ctx context.Context, id ID) (E, error) {
	return acr.next.Find(ctx, id)
}

// List возвращает список сущностей с учетом ограничений лимита и смещения через вложенный репозиторий.
func (acr *AuditCRUD[E, ID]) List(ctx context.Context, limit, offset int) ([]E, error) {
	return acr.next.List(ctx, limit, offset)
}

// Create создает новую сущность, формирует запись аудит-лога о создании и регистрирует её в таблице outbox.
func (acr *AuditCRUD[E, ID]) Create(ctx context.Context, entity E) (E, error) {
	// action
	res, err := acr.next.Create(ctx, entity)
	if err != nil {
		return res, err
	}

	// transform
	auditEntity := acr.mapAuditable(res)

	// builder, data
	builder := acr.builder(ctx, auditEntity, err).
		WithEvent(dto.DataEventCreate).
		WithValues(utils.BuildSingleDataAuditValues(auditEntity, false))

	// audit
	err = acr.audit(ctx, builder)
	if err != nil {
		return acr.nilInstance, pkgerrs.NewAuditError("AuditCRUD[E, ID].Create", fmt.Sprintf("audit create entity with id [%v] failed", res.GetID()), err)
	}

	return res, err
}

// Change фиксирует состояние сущности до модификации, выполняет обновление, фиксирует измененное состояние
// и регистрирует запись аудит-лога со срезом изменений (diff) в таблице outbox.
func (acr *AuditCRUD[E, ID]) Change(ctx context.Context, entity E) (E, error) {
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
	auditBeforeEntity := acr.mapAuditable(beforeEntity)
	auditAfterEntity := acr.mapAuditable(res)

	// builder
	builder := acr.builder(ctx, auditAfterEntity, err).WithEvent(dto.DataEventChange)

	// data
	// case 1 - all data in our hands
	builder.WithValues(utils.BuildBothDataAuditValues(auditBeforeEntity, auditAfterEntity))
	// audit
	err = acr.audit(ctx, builder)
	if err != nil {
		return acr.nilInstance, pkgerrs.NewAuditError("AuditCRUD[E, ID].Change", fmt.Sprintf("audit change entity with id [%v] failed", res.GetID()), err)
	}

	return res, nil
}

// Delete перехватывает фазу удаления сущности, фиксирует ее финальное состояние до очистки записей БД и отправляет лог.
func (acr *AuditCRUD[E, ID]) Delete(ctx context.Context, id ID) error {
	// get before entity
	beforeEntity, beforeErr := acr.next.Find(ctx, id)
	if beforeErr != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].Delete", fmt.Sprintf("unable to find entity with id [%v]", id), beforeErr)
	}
	// action
	err := acr.next.Delete(ctx, id)
	// auditable
	auditEntity := acr.mapAuditable(beforeEntity)
	// builder
	builder := acr.builder(ctx, auditEntity, err).WithEvent(dto.DataEventRemove)
	// data
	builder.WithValues(utils.BuildSingleDataAuditValues(auditEntity, true))
	// audit
	err = acr.audit(ctx, builder)
	if err != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].Delete", fmt.Sprintf("audit delete entity with id [%v] failed", id), beforeErr)
	}

	return err
}

func (acr *AuditCRUD[E, ID]) builder(ctx context.Context, auditEntity pkgdom.Auditable, err error) *utils.DataAuditBuilder {
	return utils.NewDataAuditBuilder().
		// common
		WithSource(acr.source).
		WithEventDate(time.Now()).
		WithStatus(utils.GetAuditStatus(err)).
		// entity info
		WithInternalTypeName(auditEntity.GetInternalTypeName()).
		WithTypeName(auditEntity.GetTypeName()).
		WithTypeDescription(auditEntity.GetTypeDescription()).
		WithInstanceID(auditEntity.GetInstanceID()).
		WithInstanceName(auditEntity.GetInstanceName()).
		// request
		WithUsername(utils.UsernameFromContext(ctx)).
		WithTraceID(utils.TraceIDFromContext(ctx)).
		WithRequestID(utils.RequestIDFromContext(ctx))
}

func (acr *AuditCRUD[E, ID]) audit(ctx context.Context, builder *utils.DataAuditBuilder) error {
	auditDTO := builder.Build()

	entity, err := acr.buildOutbox(ctx, auditDTO)
	if err != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].audit", "build outbox entity failed", err)
	}

	err = acr.outboxRepo.Save(ctx, entity)
	if err != nil {
		return pkgerrs.NewAuditError("AuditCRUD[E, ID].audit", "save audit outbox failed", err)
	}

	return nil
}

// GetSource возвращает строковый идентификатор источника событий (компонента), настроенного для данного репозитория.
func (acr *AuditCRUD[E, ID]) GetSource() string {
	return acr.source
}

// Unwrap возвращает исходный декорируемый экземпляр CRUD-репозитория, находящийся под слоем аудита.
func (acr *AuditCRUD[E, ID]) Unwrap() libdom.CRUDRepository[E, ID] {
	return acr.next
}

// GetOutboxRepository возвращает используемый репозиторий таблицы outbox, в который записываются логи телеметрии изменений.
func (acr *AuditCRUD[E, ID]) GetOutboxRepository() libextdom.OutboxRepository[libextdom.Outbox[ID], ID] {
	return acr.outboxRepo
}

// GetNilInstance возвращает пустое (дефолтное) значение сущности E для безопасного возврата при возникновении критических ошибок.
func (acr *AuditCRUD[E, ID]) GetNilInstance() E {
	return acr.nilInstance
}
