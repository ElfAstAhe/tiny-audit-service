package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/db"
	libdomain "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type AuthAuditPgRepository struct {
	*repository.BaseCRUDRepository[*domain.AuthAudit, string]
}

var _ libdomain.CRUDRepository[*domain.AuthAudit, string] = (*AuthAuditPgRepository)(nil)
var _ domain.AuthAuditRepository = (*AuthAuditPgRepository)(nil)

func NewAuthAuditPgRepository(executor db.Executor, errDecipher db.ErrorDecipher) (*AuthAuditPgRepository, error) {
	res := &AuthAuditPgRepository{}

	// sql builders
	queryBuilders := repository.NewBaseCRUDQueryBuildersBuilder().NewInstance().
		WithFind(func() string {
			return sqlAuthAuditFind
		}).
		WithList(func() string {
			return sqlAuthAuditList
		}).
		WithCreate(func() string {
			return sqlAuthAuditCreate
		}).
		WithChange(func() string {
			return sqlAuthAuditChange
		}).
		WithDelete(func() string {
			return sqlAuthAuditDelete
		}).
		Build()
	// callbacks
	callbacks, err := repository.NewBaseRepositoryCallbacksBuilder[*domain.AuthAudit, string]().NewInstance().
		WithEntityScanner(res.entityScanner).
		WithNewEntityFactory(domain.NewEmptyAuthAudit).
		WithAfterListYield(res.afterListYield).
		WithValidateCreate(res.validateCreate).
		WithBeforeCreate(res.beforeCreate).
		WithCreator(res.creator).
		WithValidateChange(res.validateChange).
		WithBeforeChange(res.beforeChange).
		WithChanger(res.changer).
		Build()
	// base CRUD
	base, err := repository.NewBaseCRUDRepository[*domain.AuthAudit, string](
		executor,
		errDecipher,
		repository.NewEntityInfo("auth_audit_[index]", "AuditAudit"),
		queryBuilders,
		callbacks,
	)
	if err != nil {
		return nil, errs.NewCommonError("error create AuthAuditPgRepository", err)
	}

	res.BaseCRUDRepository = base

	return res, nil
}

func (aa *AuthAuditPgRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.AuthAudit, error) {
	// ToDo: implement

	return nil, nil
}

func (aa *AuthAuditPgRepository) ListByUsername(ctx context.Context, username string, offset, limit int) ([]*domain.AuthAudit, error) {
	// ToDo: implement

	return nil, nil
}

func (aa *AuthAuditPgRepository) entityScanner(scanner repository.Scannable, sourceLabel string, entity *domain.AuthAudit, params ...any) error {
	return scanner.Scan(
		&entity.ID,
		&entity.Source,
		&entity.EventDate,
		&entity.Event,
		&entity.Status,
		&entity.RequestID,
		&entity.Username,
		&entity.AccessToken,
		&entity.RefreshToken,
		&entity.CreatedAt,
		&entity.UpdatedAt,
	)
}

func (aa *AuthAuditPgRepository) afterListYield(entity *domain.AuthAudit, params ...any) (*domain.AuthAudit, bool, error) {
	return entity, true, nil
}

func (aa *AuthAuditPgRepository) validateCreate(entity *domain.AuthAudit, params ...any) error {
	if entity == nil {
		return errs.NewInvalidArgumentError("entity", "auth audit entity is nil")
	}

	return entity.ValidateCreate()
}

func (aa *AuthAuditPgRepository) beforeCreate(entity *domain.AuthAudit, params ...any) error {
	if err := entity.BeforeCreate(); err != nil {
		return errs.NewDalError("AuthAuditPgRepository.beforeCreate", "before create entity", err)
	}

	return nil
}

func (aa *AuthAuditPgRepository) creator(ctx context.Context, querier db.Querier, entity *domain.AuthAudit, params ...any) (*sql.Row, error) {
	return querier.QueryRowContext(ctx, aa.GetQueryBuilders().GetCreate()(),
		entity.ID,
		entity.Source,
		entity.EventDate,
		entity.Event,
		entity.Status,
		entity.RequestID,
		entity.Username,
		entity.AccessToken,
		entity.RefreshToken,
		entity.CreatedAt,
	), nil
}

func (aa *AuthAuditPgRepository) validateChange(entity *domain.AuthAudit, params ...any) error {
	if entity == nil {
		return errs.NewInvalidArgumentError("entity", "auth audit entity is nil")
	}

	return entity.ValidateChange()
}

func (aa *AuthAuditPgRepository) beforeChange(entity *domain.AuthAudit, params ...any) error {
	if err := entity.BeforeChange(); err != nil {
		return errs.NewDalError("AuthAuditPgRepository.beforeChange", "before change entity", err)
	}

	return nil
}

func (aa *AuthAuditPgRepository) changer(ctx context.Context, querier db.Querier, entity *domain.AuthAudit, params ...any) (*sql.Row, error) {
	return querier.QueryRowContext(ctx, aa.GetQueryBuilders().GetChange()(),
		entity.ID,
		entity.Source,
		entity.EventDate,
		entity.Event,
		entity.Status,
		entity.RequestID,
		entity.Username,
		entity.AccessToken,
		entity.RefreshToken,
		entity.UpdatedAt,
	), nil
}
