package trace

import (
	"context"
	"fmt"
	"time"

	libdomain "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

type DataAuditTraceRepository struct {
	*repository.BaseCRUDTraceRepository[*domain.DataAudit, string]
	next domain.DataAuditRepository
}

var _ libdomain.CRUDRepository[*domain.DataAudit, string] = (*DataAuditTraceRepository)(nil)
var _ domain.DataAuditRepository = (*DataAuditTraceRepository)(nil)

func NewDataAuditTraceRepository(next domain.DataAuditRepository) *DataAuditTraceRepository {
	return &DataAuditTraceRepository{
		next:                    next,
		BaseCRUDTraceRepository: repository.NewBaseCRUDTraceRepository[*domain.DataAudit, string]("DataAuditRepository", next),
	}
}

func (dat *DataAuditTraceRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.DataAudit, error) {
	ctx, span := dat.StartSpan(ctx, fmt.Sprintf("%s.ListByPeriod", dat.BaseCRUDTraceRepository.GetRepositoryName()))
	defer span.End()

	span.SetAttributes(
		attribute.String("param.from", from.String()),
		attribute.String("param.till", till.String()),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := dat.next.ListByPeriod(ctx, from, till, limit, offset)
	if err != nil {
		span.AddEvent("ListByPeriod_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}

func (dat *DataAuditTraceRepository) ListByInstance(ctx context.Context, typeName, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	ctx, span := dat.StartSpan(ctx, fmt.Sprintf("%s.ListByInstance", dat.BaseCRUDTraceRepository.GetRepositoryName()))
	defer span.End()

	span.SetAttributes(
		attribute.String("param.typeName", typeName),
		attribute.String("param.instanceID", instanceID),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := dat.next.ListByInstance(ctx, typeName, instanceID, limit, offset)
	if err != nil {
		span.AddEvent("ListByInstance_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
