package repository

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/db"
	libpg "github.com/ElfAstAhe/go-service-template/pkg/db/postgres"
	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	libextdom "github.com/ElfAstAhe/go-service-template/pkg/ext/domain"
	librepo "github.com/ElfAstAhe/go-service-template/pkg/repository/v2"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/repository/postgres"
)

const (
	statementAuthListByPeriod   string = "stetement-auth-list-by-period"
	statementAuthListByUsername string = "stetement-auth-list-by-username"
	statementAuthGetTail        string = "stetement-auth-get-tail"
)

type AuthAudit struct {
	*librepo.CRUD[*domain.AuthAudit, string]
}

var _ libdom.CRUDRepository[*domain.AuthAudit, string] = (*AuthAudit)(nil)
var _ domain.AuthAuditRepository = (*AuthAudit)(nil)
var _ libextdom.TailRepository[string] = (*AuthAudit)(nil)

func NewAuthAudit(
	driverName string,
	executor db.Executor,
	errDecipher db.ErrorDecipher,
) (*AuthAudit, error) {
	// instance
	res := &AuthAudit{}
	// embed
	embed, err := librepo.NewCRUD[*domain.AuthAudit, string](
		librepo.WithExecutor[*domain.AuthAudit, string](executor),
		librepo.WithErrorDecipher[*domain.AuthAudit, string](errDecipher),
		librepo.WithHooks[*domain.AuthAudit, string](
			librepo.WithHookNewEntityFactory[*domain.AuthAudit, string](domain.NewEmptyAuthAudit),
			librepo.WithHookCreator[*domain.AuthAudit, string](res.creator),
			librepo.WithHookChanger[*domain.AuthAudit, string](res.changer),
		),
		librepo.WithEntityInfo[*domain.AuthAudit, string]("auth_audit", "AuthAudit"),
		librepo.WithDriverName[*domain.AuthAudit, string](driverName),
		librepo.WithStatementsBuilder[*domain.AuthAudit, string](libpg.DriverName, res.pgStatementsBuilder),
		librepo.WithStatementsBuilder[*domain.AuthAudit, string]("pgx", res.pgStatementsBuilder),
		librepo.WithStatementsBuilder[*domain.AuthAudit, string]("postgres", res.pgStatementsBuilder),
	)
	if err != nil {
		return nil, errs.NewDalError("NewAuthAudit", "create embed CRUD failed", err)
	}
	// setup
	res.CRUD = embed

	return res, nil
}

func (aa *AuthAudit) GetTail(ctx context.Context, tailDate time.Time) ([]string, error) {
	// fail fast
	if tailDate.IsZero() {
		return []string{}, errs.NewInvalidArgumentError("TailDate", "tail date is required")
	}
	// act
	res, err := aa.GetHelper().List(ctx, statementAuthGetTail, tailDate)
	if err != nil {
		return []string{}, errs.NewDalError("AuthAudit.GetTail", "List failed", err)
	}

	return libdom.MapEntitiesToIDs(res), nil
}

func (aa *AuthAudit) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.AuthAudit, error) {
	// fail fast
	if err := aa.validateListByPeriod(from, till, limit, offset); err != nil {
		return nil, errs.NewDalError("AuthAudit.ListByPeriod", "validate failed", err)
	}
	// act
	res, err := aa.GetHelper().List(ctx, statementAuthListByPeriod, from, till, limit, offset)
	if err != nil {
		return nil, errs.NewDalError("AuthAudit.ListByPeriod", "list failed", err)
	}

	return res, nil
}

func (aa *AuthAudit) validateListByPeriod(from, till time.Time, limit, offset int) error {
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

func (aa *AuthAudit) ListByUsername(ctx context.Context, username string, offset, limit int) ([]*domain.AuthAudit, error) {
	// fail fast
	if err := aa.validateListByUsername(username, offset, limit); err != nil {
		return nil, errs.NewDalError("AuthAudit.ListByUsername", "validate failed", err)
	}
	// act
	res, err := aa.GetHelper().List(ctx, statementAuthListByUsername, username, offset, limit)
	if err != nil {
		return nil, errs.NewDalError("AuthAudit.ListByUsername", "list failed", err)
	}

	return res, nil
}

func (aa *AuthAudit) validateListByUsername(username string, offset, limit int) error {
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

func (aa *AuthAudit) creator(ctx context.Context, querier db.Querier, entity *domain.AuthAudit, params ...any) (librepo.Scannable, error) {
	query, err := aa.GetHelper().GetStatements().GetSQL(librepo.StatementCreate)
	if err != nil {
		return nil, errs.NewDalError("AuthAudit.creator", "create statement not found", err)
	}

	return querier.QueryRowContext(ctx, query,
		entity.ID,
		entity.Source, entity.EventDate, entity.Event, entity.Status,
		entity.RequestID, entity.TraceID,
		entity.Username, entity.AccessToken, entity.RefreshToken,
		entity.CreatedAt,
	), nil
}

func (aa *AuthAudit) changer(ctx context.Context, querier db.Querier, entity *domain.AuthAudit, params ...any) (librepo.Scannable, error) {
	query, err := aa.GetHelper().GetStatements().GetSQL(librepo.StatementChange)
	if err != nil {
		return nil, errs.NewDalError("AuthAudit.changer", "change statement not found", err)
	}

	return querier.QueryRowContext(ctx, query,
		entity.ID,
		entity.Source, entity.EventDate, entity.Event, entity.Status,
		entity.RequestID, entity.TraceID,
		entity.Username, entity.AccessToken, entity.RefreshToken,
		entity.UpdatedAt,
	), nil
}

func (aa *AuthAudit) mapDefault(ctx context.Context, scannable librepo.Scannable, params ...any) (*domain.AuthAudit, error) {
	res := aa.GetHelper().GetHooks().NewEntityFactory()
	err := scannable.Scan(
		&res.ID,
		&res.Source, &res.EventDate, &res.Event, &res.Status,
		&res.RequestID, &res.TraceID,
		&res.Username, &res.AccessToken, &res.RefreshToken,
		&res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return nil, errs.NewDalError("AuthAudit.mapDefault", "default mapper failed", err)
	}

	return res, nil
}

func (aa *AuthAudit) mapTail(ctx context.Context, scannable librepo.Scannable, params ...any) (*domain.AuthAudit, error) {
	res := aa.GetHelper().GetHooks().NewEntityFactory()
	err := scannable.Scan(
		&res.ID,
	)
	if err != nil {
		return nil, errs.NewDalError("AuthAudit.mapTail", "tail mapper failed", err)
	}

	return res, nil
}

func (aa *AuthAudit) pgStatementsBuilder() *librepo.Statements[*domain.AuthAudit, string] {
	return librepo.NewStatements[*domain.AuthAudit, string]().
		WithStatement(librepo.StatementFind, librepo.NewStatement(postgres.SqlAuthAuditFind, aa.mapDefault)).
		WithStatement(librepo.StatementList, librepo.NewStatement(postgres.SqlAuthAuditList, aa.mapDefault)).
		WithStatement(librepo.StatementCreate, librepo.NewStatement(postgres.SqlAuthAuditCreate, aa.mapDefault)).
		WithStatement(librepo.StatementChange, librepo.NewStatement(postgres.SqlAuthAuditChange, aa.mapDefault)).
		WithStatement(librepo.StatementDelete, librepo.NewStatement[*domain.AuthAudit, string](postgres.SqlAuthAuditDelete, nil)).
		WithStatement(statementAuthGetTail, librepo.NewStatement(postgres.SqlAuthAuditTailList, aa.mapTail)).
		WithStatement(statementAuthListByPeriod, librepo.NewStatement(postgres.SqlAuthAuditListByPeriod, aa.mapDefault)).
		WithStatement(statementAuthListByUsername, librepo.NewStatement(postgres.SqlAuthAuditListByUsername, aa.mapDefault))
}
