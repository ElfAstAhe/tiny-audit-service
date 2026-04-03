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
	apprepo "github.com/ElfAstAhe/tiny-audit-service/internal/repository"
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

func (da *DataAuditPgRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.DataAudit, error) {
	if err := da.validateListByPeriod(from, till, limit, offset); err != nil {
		return nil, err
	}

	res, err := da.GetHelper().List(ctx, apprepo.SourceLabelListByPeriod, sqlDataAuditListByPeriod,
		from,
		till,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (da *DataAuditPgRepository) validateListByPeriod(from, till time.Time, limit, offset int) error {
	if from.IsZero() {
		return errs.NewInvalidArgumentError("from", "field is required")
	}
	if till.IsZero() {
		return errs.NewInvalidArgumentError("till", "field is required")
	}
	if !(limit > 0) {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if !(offset >= 0) {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}

func (da *DataAuditPgRepository) ListByInstance(ctx context.Context, typeName string, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	if err := da.validateListByInstance(typeName, instanceID, limit, offset); err != nil {
		return nil, err
	}

	res, err := da.GetHelper().List(ctx, apprepo.SourceLabelListByInstance, sqlDataAuditListByInstance,
		typeName,
		instanceID,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (da *DataAuditPgRepository) validateListByInstance(typeName, instanceID string, limit, offset int) error {
	if typeName == "" {
		return errs.NewInvalidArgumentError("typeName", "typeName is required")
	}
	if instanceID == "" {
		return errs.NewInvalidArgumentError("instanceID", "instanceID is required")
	}
	if !(limit > 0) {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if !(offset >= 0) {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}

func (da *DataAuditPgRepository) entityScanner(scanner repository.Scannable, sourceLabel string, entity *domain.DataAudit, params ...any) error {
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

func (da *DataAuditPgRepository) afterListYield(entity *domain.DataAudit, params ...any) (*domain.DataAudit, bool, error) {
	return entity, true, nil
}

func (da *DataAuditPgRepository) validateCreate(entity *domain.DataAudit, params ...any) error {
	if entity == nil {
		return errs.NewInvalidArgumentError("entity", "cannot be nil")
	}

	return entity.ValidateCreate()
}

func (da *DataAuditPgRepository) beforeCreate(entity *domain.DataAudit, params ...any) error {
	if err := entity.BeforeCreate(); err != nil {
		return errs.NewDalError("DataAuditPgRepository.beforeCreate", "before create entity", err)
	}

	return nil
}

func (da *DataAuditPgRepository) creator(ctx context.Context, querier db.Querier, entity *domain.DataAudit, params ...any) (*sql.Row, error) {
	var valuesRaw []byte
	var err error
	if len(entity.Values) > 0 {
		valuesRaw, err = json.Marshal(entity.Values)
		if err != nil {
			return nil, errs.NewDalError("DataAuditPgRepository.creator", "marshal values json", err)
		}
	}

	return querier.QueryRowContext(ctx, da.GetQueryBuilders().GetCreate()(),
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

func (da *DataAuditPgRepository) validateChange(entity *domain.DataAudit, params ...any) error {
	if entity == nil {
		return errs.NewInvalidArgumentError("entity", "cannot be nil")
	}

	return entity.ValidateChange()
}

func (da *DataAuditPgRepository) beforeChange(entity *domain.DataAudit, params ...any) error {
	if err := entity.BeforeChange(); err != nil {
		return errs.NewDalError("DataAuditPgRepository.beforeChange", "before change entity", err)
	}

	return nil
}

func (da *DataAuditPgRepository) changer(ctx context.Context, querier db.Querier, entity *domain.DataAudit, params ...any) (*sql.Row, error) {
	var valuesRaw []byte
	var err error
	if len(entity.Values) > 0 {
		valuesRaw, err = json.Marshal(entity.Values)
		if err != nil {
			return nil, errs.NewDalError("DataAuditPgRepository.changer", "marshal values json", err)
		}
	}

	return querier.QueryRowContext(ctx, da.GetQueryBuilders().GetChange()(),
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
