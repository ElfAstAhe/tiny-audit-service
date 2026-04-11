package repository

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libutils "github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	auditdomain "github.com/ElfAstAhe/tiny-audit-service/pkg/domain"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/utils"
)

type AuditableEntity[ID comparable] interface {
	domain.Entity[ID]
	auditdomain.Auditable
}

type AuditableMapper[ID comparable] func(entity domain.Entity[ID]) auditdomain.Auditable

type AuditCRUDRepository[E AuditableEntity[ID], ID comparable] struct {
	next        domain.CRUDRepository[E, ID]
	source      string
	mapper      AuditableMapper[ID]
	auditClient client.DataAuditClient
	log         logger.Logger
}

func NewAuditCRUDRepository[E AuditableEntity[ID], ID comparable](
	next domain.CRUDRepository[E, ID],
	source string,
	mapper AuditableMapper[ID],
	auditClient client.DataAuditClient,
	log logger.Logger,
) *AuditCRUDRepository[E, ID] {
	return &AuditCRUDRepository[E, ID]{
		next:        next,
		source:      source,
		mapper:      mapper,
		auditClient: auditClient,
		log:         log.GetLogger(libutils.GetTypeName(next)),
	}
}

func (acr *AuditCRUDRepository[E, ID]) Find(ctx context.Context, id ID) (E, error) {
	return acr.next.Find(ctx, id)
}

func (acr *AuditCRUDRepository[E, ID]) List(ctx context.Context, limit, offset int) ([]E, error) {
	return acr.next.List(ctx, limit, offset)
}

func (acr *AuditCRUDRepository[E, ID]) Create(ctx context.Context, entity E) (E, error) {
	acr.log.Debugf("audit Create start")
	defer acr.log.Debugf("audit Create finish")

	res, err := acr.next.Create(ctx, entity)
	if err != nil {
		acr.log.Debugf("audit original Create failed [%v]", err)
	}

	auditEntity := acr.mapper(res)
	// common
	builder := utils.NewDataAuditBuilder().
		WithSource(acr.source).
		WithEventDate(time.Now()).
		WithEvent(dto.DataEventCreate).
		WithStatus(utils.GetAuditStatus(err))
	if auditEntity != nil {
		// instance
		builder.WithInternalTypeName(auditEntity.GetInternalTypeName()).
			WithTypeName(auditEntity.GetTypeName()).
			WithTypeDescription(auditEntity.GetTypeDescription()).
			WithInstanceID(auditEntity.GetInstanceID()).
			WithInstanceName(auditEntity.GetInstanceName()).
			WithValues(utils.BuildSingleDataAuditValues(auditEntity, false))
		// request
		builder.WithUsername(utils.UsernameFromContext(ctx)).
			WithTraceID(utils.TraceIDFromContext(ctx)).
			WithRequestID(utils.RequestIDFromContext(ctx))
	}
	clientDTO := builder.Build()
	acr.log.Debugf("audit Create client dto [%v]", clientDTO)

	auditErr := acr.auditClient.Audit(clientDTO)
	if auditErr != nil {
		acr.log.Errorf("audit Create audit [%v] failed [%v]", clientDTO, auditErr)
	}

	return res, err
}

func (acr *AuditCRUDRepository[E, ID]) Change(ctx context.Context, entity E) (E, error) {
	return acr.next.Change(ctx, entity)
}

func (acr *AuditCRUDRepository[E, ID]) Delete(ctx context.Context, id ID) error {
	return acr.next.Delete(ctx, id)
}
