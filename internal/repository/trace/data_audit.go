package trace

import (
	"context"
	"fmt"
	"time"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/repository/trace"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

// DataAuditTraceRepository представляет собой декоратор репозитория, предназначенный
// для автоматического сбора распределенной телеметрии и трассировки операций аудита данных.
type DataAuditTraceRepository struct {
	*trace.CRUD[*domain.DataAudit, string]
	next domain.DataAuditRepository
}

var _ libdom.CRUDRepository[*domain.DataAudit, string] = (*DataAuditTraceRepository)(nil)
var _ domain.DataAuditRepository = (*DataAuditTraceRepository)(nil)

// NewDataAuditRepository создает новый экземпляр трассирующего декоратора DataAuditTraceRepository.
func NewDataAuditRepository(next domain.DataAuditRepository) *DataAuditTraceRepository {
	return &DataAuditTraceRepository{
		next: next,
		CRUD: trace.NewCRUD[*domain.DataAudit, string]("DataAuditRepository", next),
	}
}

// ListByPeriod возвращает список записей аудита данных за указанный временной интервал,
// фиксируя в атрибутах спана временные границы, параметры лимитов, события ошибок и статус операции.
//
//goland:noinspection DuplicatedCode
func (dat *DataAuditTraceRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.DataAudit, error) {
	repoCtx, span := dat.StartSpan(ctx, fmt.Sprintf("%s.ListByPeriod", dat.GetRepositoryName()))
	defer span.End()

	span.SetAttributes(
		attribute.String("param.from", from.Format(time.DateTime)),
		attribute.String("param.till", till.Format(time.DateTime)),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := dat.next.ListByPeriod(repoCtx, from, till, limit, offset)
	if err != nil {
		span.AddEvent("ListByPeriod_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}

// ListByInstance возвращает список записей аудита для конкретного экземпляра сущности по её типу и ID,
// фиксируя параметры типа, инстанса, пагинации, события ошибок и статус операции в текущем спане.
func (dat *DataAuditTraceRepository) ListByInstance(ctx context.Context, typeName, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	repoCtx, span := dat.StartSpan(ctx, fmt.Sprintf("%s.ListByInstance", dat.GetRepositoryName()))
	defer span.End()

	span.SetAttributes(
		attribute.String("param.typeName", typeName),
		attribute.String("param.instanceID", instanceID),
		attribute.Int("param.limit", limit),
		attribute.Int("param.offset", offset),
	)

	res, err := dat.next.ListByInstance(repoCtx, typeName, instanceID, limit, offset)
	if err != nil {
		span.AddEvent("ListByInstance_failed")
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	return res, err
}
