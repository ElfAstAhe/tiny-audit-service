package repository

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libutils "github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	pkgdomain "github.com/ElfAstAhe/tiny-audit-service/pkg/domain"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
)

type AuditableEntity[ID comparable] interface {
	domain.Entity[ID]
	pkgdomain.Auditable
}

type AuditableMapper[E domain.Entity[ID], ID comparable] func(entity E) pkgdomain.Auditable

type BaseAuditCRUDRepository[E AuditableEntity[ID], ID comparable] struct {
	next        domain.CRUDRepository[E, ID]
	source      string
	mapper      AuditableMapper[E, ID]
	auditClient client.DataAuditClient
	log         logger.Logger
}

func NewBaseAuditCRUDRepository[E AuditableEntity[ID], ID comparable](
	next domain.CRUDRepository[E, ID],
	source string,
	mapper AuditableMapper[E, ID],
	auditClient client.DataAuditClient,
	log logger.Logger,
) *BaseAuditCRUDRepository[E, ID] {
	return &BaseAuditCRUDRepository[E, ID]{
		next:        next,
		source:      source,
		mapper:      mapper,
		auditClient: auditClient,
		log:         log.GetLogger(libutils.GetTypeName(next)),
	}
}

func (acr *BaseAuditCRUDRepository[E, ID]) Find(ctx context.Context, id ID) (E, error) {
	return acr.next.Find(ctx, id)
}

func (acr *BaseAuditCRUDRepository[E, ID]) List(ctx context.Context, limit, offset int) ([]E, error) {
	return acr.next.List(ctx, limit, offset)
}

func (acr *BaseAuditCRUDRepository[E, ID]) Create(ctx context.Context, entity E) (E, error) {
	acr.log.Debugf("audit Create start")
	defer acr.log.Debugf("audit Create finish")

	// action
	res, err := acr.next.Create(ctx, entity)
	if err != nil {
		acr.log.Debugf("audit original Create failed [%v]", err)
	}

	// resolve entity instance
	var auditEntity pkgdomain.Auditable
	if err == nil {
		auditEntity = acr.mapper(res)
	} else {
		auditEntity = acr.mapper(entity)
	}

	// builder, data
	builder := acr.builder(ctx, auditEntity, err).WithEvent(dto.DataEventCreate).
		WithValues(utils.BuildSingleDataAuditValues(auditEntity, false))

	// audit
	acr.audit(builder)

	return res, err
}

func (acr *BaseAuditCRUDRepository[E, ID]) Change(ctx context.Context, entity E) (E, error) {
	acr.log.Debugf("audit Change start")
	defer acr.log.Debugf("audit Change finish")

	// get before entity
	beforeEntity, beforeErr := acr.next.Find(ctx, entity.GetID())
	if beforeErr != nil {
		acr.log.Debugf("audit Change get before entity failed [%v]", beforeErr)
	}

	// action
	res, err := acr.next.Change(ctx, entity)
	if err != nil {
		acr.log.Debugf("audit original Change failed [%v]", err)
	}

	var auditBeforeEntity pkgdomain.Auditable
	if beforeErr != nil {
		auditBeforeEntity = acr.mapper(beforeEntity)
	}
	// resolve entity instance
	var auditAfterEntity pkgdomain.Auditable
	if err == nil {
		auditAfterEntity = acr.mapper(res)
	} else {
		auditAfterEntity = acr.mapper(entity)
	}
	// builder
	builder := acr.builder(ctx, auditAfterEntity, err).WithEvent(dto.DataEventChange)
	// data
	if err == nil && beforeErr == nil && utils.HasChanges(auditBeforeEntity, auditAfterEntity) {
		// case 1 - all data in our hands
		builder.WithValues(utils.BuildBothDataAuditValues(auditBeforeEntity, auditAfterEntity))
	} else if err == nil && beforeErr != nil {
		// case 2 - only after entity in our hands
		builder.WithValues(utils.BuildSingleDataAuditValues(auditAfterEntity, false))
	} else if err != nil && beforeErr == nil {
		// case 3 - all data in our hands, but got error on change
		builder.WithValues(utils.BuildBothDataAuditValues(auditBeforeEntity, auditAfterEntity))
	} else {
		// default, got all errors (get before and change)
		builder.WithValues(utils.BuildSingleDataAuditValues(auditAfterEntity, false))
	}
	// audit
	acr.audit(builder)

	return res, err
}

func (acr *BaseAuditCRUDRepository[E, ID]) Delete(ctx context.Context, id ID) error {
	acr.log.Debugf("audit Delete start")
	defer acr.log.Debugf("audit Delete finish")

	// get before entity
	beforeEntity, beforeErr := acr.next.Find(ctx, id)
	if beforeErr != nil {
		acr.log.Debugf("audit Delete get before entity failed [%v]", beforeErr)
	}
	// action
	err := acr.next.Delete(ctx, id)
	// auditable
	var auditEntity pkgdomain.Auditable = nil
	if beforeErr != nil {
		auditEntity = acr.mapper(beforeEntity)
	}
	// builder
	builder := acr.builder(ctx, auditEntity, err).WithEvent(dto.DataEventRemove)
	// data
	if beforeErr != nil {
		builder.WithValues(utils.BuildSingleDataAuditValues(auditEntity, true))
	}
	// audit
	acr.audit(builder)

	return err
}

func (acr *BaseAuditCRUDRepository[E, ID]) builder(ctx context.Context, auditEntity pkgdomain.Auditable, err error) *utils.DataAuditBuilder {
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

func (acr *BaseAuditCRUDRepository[E, ID]) audit(builder *utils.DataAuditBuilder) {
	clientDTO := builder.Build()
	acr.log.Debugf("audit Create client dto [%v]", clientDTO)

	auditErr := acr.auditClient.Audit(clientDTO)
	if auditErr != nil {
		acr.log.Errorf("audit Create audit [%v] failed [%v]", clientDTO, auditErr)
	}
}

func (acr *BaseAuditCRUDRepository[E, ID]) GetSource() string {
	return acr.source
}

// Unwrap get decorable repository
func (acr *BaseAuditCRUDRepository[E, ID]) Unwrap() domain.CRUDRepository[E, ID] {
	return acr.next
}
func (acr *BaseAuditCRUDRepository[E, ID]) GetAuditClient() client.DataAuditClient {
	return acr.auditClient
}

func (acr *BaseAuditCRUDRepository[E, ID]) GetLogger() logger.Logger {
	return acr.log
}
