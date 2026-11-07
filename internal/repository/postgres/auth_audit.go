package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/db"
	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	libextdom "github.com/ElfAstAhe/go-service-template/pkg/ext/domain"
	librepo "github.com/ElfAstAhe/go-service-template/pkg/repository"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

// AuthAuditPgRepository репозиторий сущности domain.AuthAudit
//
// Deprecated: AuthAuditPgRepository реализация v1 repository
type AuthAuditPgRepository struct {
	*librepo.CRUD[*domain.AuthAudit, string]
}

var _ libdom.CRUDRepository[*domain.AuthAudit, string] = (*AuthAuditPgRepository)(nil)
var _ domain.AuthAuditRepository = (*AuthAuditPgRepository)(nil)
var _ libextdom.TailRepository[string] = (*AuthAuditPgRepository)(nil)

//goland:noinspection DuplicatedCode
func NewAuthAuditPgRepository(executor db.Executor, errDecipher db.ErrorDecipher) (*AuthAuditPgRepository, error) {
	res := &AuthAuditPgRepository{}

	// sql builders
	queryBuilders := librepo.NewBaseCRUDQueryBuildersBuilder().NewInstance().
		WithFind(func() string {
			return SqlAuthAuditFind
		}).
		WithList(func() string {
			return SqlAuthAuditList
		}).
		WithCreate(func() string {
			return SqlAuthAuditCreate
		}).
		WithChange(func() string {
			return SqlAuthAuditChange
		}).
		WithDelete(func() string {
			return SqlAuthAuditDelete
		}).
		Build()
	// callbacks
	callbacks, _ := librepo.NewCallbacksBuilder[*domain.AuthAudit, string]().NewInstance().
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
	base, err := librepo.NewCRUD[*domain.AuthAudit, string](
		executor,
		errDecipher,
		librepo.NewEntityInfo("auth_audit", "AuthAudit"),
		queryBuilders,
		callbacks,
	)
	if err != nil {
		return nil, errs.NewCommonError("error create AuthAuditPgRepository", err)
	}

	res.CRUD = base

	return res, nil
}

func (aa *AuthAuditPgRepository) GetTail(ctx context.Context, tailDate time.Time) ([]string, error) {
	res, err := aa.GetHelper().List(ctx, SourceLabelTailList, SqlAuthAuditTailList, tailDate)
	if err != nil {
		return nil, err
	}

	return libdom.MapEntitiesToIDs(res), err
}

func (aa *AuthAuditPgRepository) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.AuthAudit, error) {
	if err := aa.validateListByPeriod(from, till, limit, offset); err != nil {
		return nil, err
	}

	res, err := aa.GetHelper().List(ctx, SourceLabelListByPeriod, SqlAuthAuditListByPeriod,
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
	if limit <= 0 {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if offset < 0 {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}

func (aa *AuthAuditPgRepository) ListByUsername(ctx context.Context, username string, offset, limit int) ([]*domain.AuthAudit, error) {
	if err := aa.validateListByUsername(username, offset, limit); err != nil {
		return nil, err
	}

	res, err := aa.GetHelper().List(ctx, SourceLabelListByUsername, SqlAuthAuditListByUsername,
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
	if limit <= 0 {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if offset < 0 {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}

func (aa *AuthAuditPgRepository) entityScanner(scanner librepo.Scannable, sourceLabel string, entity *domain.AuthAudit, params ...any) error {
	switch sourceLabel {
	case SourceLabelTailList:
		return scanner.Scan(&entity.ID)
	default:
		return scanner.Scan(
			&entity.ID,
			&entity.Source,
			&entity.EventDate,
			&entity.Event,
			&entity.Status,
			&entity.RequestID,
			&entity.TraceID,
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
		return errs.NewDalError("AuthAuditPgRepository.beforeCreate", "before create failed", err)
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
		entity.TraceID,
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
		return errs.NewDalError("AuthAuditPgRepository.beforeChange", "before change failed", err)
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
		entity.TraceID,
		entity.Username,
		entity.AccessToken,
		entity.RefreshToken,
		entity.UpdatedAt,
	), nil
}
