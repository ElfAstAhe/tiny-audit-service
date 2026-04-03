package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/db"
	libdomain "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type DataAuditPgRepository struct {
	*repository.BaseCRUDRepository[*domain.DataAudit, string]
}

var _ libdomain.CRUDRepository[*domain.DataAudit, string] = (*DataAuditPgRepository)(nil)
var _ domain.DataAuditRepository = (*DataAuditPgRepository)(nil)

func NewDataAuditPgRepository(executor db.Executor, errDecipher db.ErrorDecipher) (*DataAuditPgRepository, error) {
	res := &DataAuditPgRepository{}

	// sql builders
	queryBuilders := repository.NewBaseCRUDQueryBuildersBuilder().NewInstance().
		WithFind(func() string {
			return sqlDataAuditFind
		}).
		WithList(func() string {
			return sqlDataAuditList
		}).
		WithCreate(func() string {
			return sqlDataAuditCreate
		}).
		WithChange(func() string {
			return sqlDataAuditChange
		}).
		WithDelete(func() string {
			return sqlDataAuditDelete
		}).
		Build()
	// callbacks
	callbacks, err := repository.NewBaseRepositoryCallbacksBuilder[*domain.DataAudit, string]().NewInstance().
		WithEntityScanner(res.entityScanner).
		WithNewEntityFactory(domain.NewEmptyDataAudit).
		WithAfterListYield(res.afterListYield).
		WithValidateCreate(res.validateCreate).
		WithBeforeCreate(res.beforeCreate).
		WithCreator(res.creator).
		WithValidateChange(res.validateChange).
		WithBeforeChange(res.beforeChange).
		WithChanger(res.changer).
		Build()
	// base CRUD
	base, err := repository.NewBaseCRUDRepository[*domain.DataAudit, string](
		executor,
		errDecipher,
		repository.NewEntityInfo("data_audit_[index]", "DataAudit"),
		queryBuilders,
		callbacks,
	)
	if err != nil {
		return nil, errs.NewCommonError("error create DataAuditPgRepository", err)
	}

	res.BaseCRUDRepository = base

	return res, nil
}

func (dr *DataAuditPgRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.DataAudit, error) {
	// ToDo: implement

	return nil, nil
}

func (dr *DataAuditPgRepository) ListByInstance(ctx context.Context, typeName string, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	// ToDo: implement

	return nil, nil
}

func (dr *DataAuditPgRepository) entityScanner(scanner repository.Scannable, sourceLabel string, entity *domain.DataAudit, params ...any) error {
	var valuesRaw []byte
	if err := scanner.Scan(
		&entity.ID,
		&entity.Source,
		&entity.EventDate,
		&entity.Event,
		&entity.Status,
		&entity.RequestID,
		&entity.Username,
		&entity.TypeName,
		&entity.TypeDescription,
		&entity.InstanceID,
		&entity.InstanceName,
		&valuesRaw,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	); err != nil {
		return errs.NewDalError("DataAuditPgRepository.entityScanner", "scan row", err)
	}

	if len(valuesRaw) > 0 {
		if err := json.Unmarshal(valuesRaw, &entity.Values); err != nil {
			return errs.NewDalError("DataAuditPgRepository.entityScanner", "unmarshal values json", err)
		}
	}

	return nil
}

func (dr *DataAuditPgRepository) afterListYield(entity *domain.DataAudit, params ...any) (*domain.DataAudit, bool, error) {
	return entity, true, nil
}

func (dr *DataAuditPgRepository) validateCreate(entity *domain.DataAudit, params ...any) error {
	if entity == nil {
		return errs.NewInvalidArgumentError("entity", "cannot be nil")
	}

	return entity.ValidateCreate()
}

func (dr *DataAuditPgRepository) beforeCreate(entity *domain.DataAudit, params ...any) error {
	if err := entity.BeforeCreate(); err != nil {
		return errs.NewDalError("DataAuditPgRepository.beforeCreate", "before create entity", err)
	}

	return nil
}

func (dr *DataAuditPgRepository) creator(ctx context.Context, querier db.Querier, entity *domain.DataAudit, params ...any) (*sql.Row, error) {
	var valuesRaw []byte
	var err error
	if len(entity.Values) > 0 {
		valuesRaw, err = json.Marshal(entity.Values)
		if err != nil {
			return nil, errs.NewDalError("DataAuditPgRepository.creator", "marshal values json", err)
		}
	}

	return querier.QueryRowContext(ctx, dr.GetQueryBuilders().GetCreate()(),
		&entity.ID,
		&entity.Source,
		&entity.EventDate,
		&entity.Event,
		&entity.Status,
		&entity.RequestID,
		&entity.Username,
		&entity.TypeName,
		&entity.TypeDescription,
		&entity.InstanceID,
		&entity.InstanceName,
		&valuesRaw,
		&entity.CreatedAt,
	), nil
}

func (dr *DataAuditPgRepository) validateChange(entity *domain.DataAudit, params ...any) error {
	if entity == nil {
		return errs.NewInvalidArgumentError("entity", "cannot be nil")
	}

	return entity.ValidateChange()
}

func (dr *DataAuditPgRepository) beforeChange(entity *domain.DataAudit, params ...any) error {
	if err := entity.BeforeChange(); err != nil {
		return errs.NewDalError("DataAuditPgRepository.beforeChange", "before change entity", err)
	}

	return nil
}

func (dr *DataAuditPgRepository) changer(ctx context.Context, querier db.Querier, entity *domain.DataAudit, params ...any) (*sql.Row, error) {
	var valuesRaw []byte
	var err error
	if len(entity.Values) > 0 {
		valuesRaw, err = json.Marshal(entity.Values)
		if err != nil {
			return nil, errs.NewDalError("DataAuditPgRepository.changer", "marshal values json", err)
		}
	}

	return querier.QueryRowContext(ctx, dr.GetQueryBuilders().GetChange()(),
		&entity.ID,
		&entity.Source,
		&entity.EventDate,
		&entity.Event,
		&entity.Status,
		&entity.RequestID,
		&entity.Username,
		&entity.TypeName,
		&entity.TypeDescription,
		&entity.InstanceID,
		&entity.InstanceName,
		&valuesRaw,
		&entity.UpdatedAt,
	), nil
}
