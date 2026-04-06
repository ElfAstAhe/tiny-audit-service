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
	apprepo "github.com/ElfAstAhe/tiny-audit-service/internal/repository"
)

type AuthAuditPgRepository struct {
	*repository.BaseCRUDRepository[*domain.AuthAudit, string]
}

var _ libdomain.CRUDRepository[*domain.AuthAudit, string] = (*AuthAuditPgRepository)(nil)
var _ domain.AuthAuditRepository = (*AuthAuditPgRepository)(nil)
var _ domain.TailRepository[string] = (*AuthAuditPgRepository)(nil)

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

func (aa *AuthAuditPgRepository) GetTail(ctx context.Context, tailDate time.Time) ([]string, error) {
	res, err := aa.GetHelper().List(ctx, apprepo.SourceLabelTailList, sqlAuthAuditTailList, tailDate)
	if err != nil {
		return nil, err
	}

	return libdomain.EntitiesToIDList(res), err
}

func (aa *AuthAuditPgRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.AuthAudit, error) {
	if err := aa.validateListByPeriod(from, till, limit, offset); err != nil {
		return nil, err
	}

	res, err := aa.GetHelper().List(ctx, apprepo.SourceLabelListByPeriod, sqlAuthAuditListByPeriod,
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

func (aa *AuthAuditPgRepository) validateListByPeriod(from, till time.Time, limit, offset int) error {
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

func (aa *AuthAuditPgRepository) ListByUsername(ctx context.Context, username string, offset, limit int) ([]*domain.AuthAudit, error) {
	if err := aa.validateListByUsername(username, offset, limit); err != nil {
		return nil, err
	}

	res, err := aa.GetHelper().List(ctx, apprepo.SourceLabelListByUsername, sqlAuthAuditListByUsername,
		username,
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (aa *AuthAuditPgRepository) validateListByUsername(username string, offset, limit int) error {
	if username == "" {
		return errs.NewInvalidArgumentError("username", "field is required")
	}
	if !(limit > 0) {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if !(offset >= 0) {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}

func (aa *AuthAuditPgRepository) entityScanner(scanner repository.Scannable, sourceLabel string, entity *domain.AuthAudit, params ...any) error {
	switch sourceLabel {
	case apprepo.SourceLabelTailList:
		return scanner.Scan(&entity.ID)
	default:
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
