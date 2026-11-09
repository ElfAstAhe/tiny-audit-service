package repository

import (
	"context"
	"encoding/json"
	"strings"
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
	statementDataListByPeriod   string = "statement-data-list-by-period"
	statementDataListByInstance string = "statement-data-list-by-instance"
	statementDataGetTail        string = "statement-data-get-tail"
)

// DataAudit представляет собой базовый репозиторий для выполнения CRUD-операций
// и специализированных выборок сущности domain.DataAudit в базе данных.
type DataAudit struct {
	*librepo.CRUD[*domain.DataAudit, string]
}

var _ libdom.CRUDRepository[*domain.DataAudit, string] = (*DataAudit)(nil)
var _ domain.DataAuditRepository = (*DataAudit)(nil)
var _ libextdom.TailRepository[string] = (*DataAudit)(nil)

// NewDataAudit создает и инициализирует новый экземпляр репозитория DataAudit,
// регистрируя необходимые мапперы, хуки и SQL-инструкции для драйверов базы данных.
func NewDataAudit(
	driverName string,
	executor db.Executor,
	errDecipher db.ErrorDecipher,
) (*DataAudit, error) {
	// instance
	res := &DataAudit{}
	// embed
	embed, err := librepo.NewCRUD(
		librepo.WithExecutor[*domain.DataAudit, string](executor),
		librepo.WithErrorDecipher[*domain.DataAudit, string](errDecipher),
		librepo.WithHooks[*domain.DataAudit, string](
			librepo.WithHookNewEntityFactory[*domain.DataAudit, string](domain.NewEmptyDataAudit),
			librepo.WithHookCreator[*domain.DataAudit, string](res.creator),
			librepo.WithHookChanger[*domain.DataAudit, string](res.changer),
		),
		librepo.WithEntityInfo[*domain.DataAudit, string]("data_audit", "DataAudit"),
		librepo.WithDriverName[*domain.DataAudit, string](driverName),
		librepo.WithStatementsBuilder[*domain.DataAudit, string](libpg.DriverName, res.pgStatementsBuilder),
		librepo.WithStatementsBuilder[*domain.DataAudit, string]("pgx", res.pgStatementsBuilder),
		librepo.WithStatementsBuilder[*domain.DataAudit, string]("postgres", res.pgStatementsBuilder),
	)
	if err != nil {
		return nil, errs.NewDalError("NewDataAudit", "create embed CRUD failed", err)
	}
	// setup
	res.CRUD = embed

	return res, nil
}

// GetTail возвращает список идентификаторов старых записей аудита данных, созданных до указанной даты.
func (da *DataAudit) GetTail(ctx context.Context, tailDate time.Time) ([]string, error) {
	// fail fast
	if tailDate.IsZero() {
		return []string{}, errs.NewInvalidArgumentError("TailDate", "tail date is required")
	}
	// act
	res, err := da.GetHelper().List(ctx, statementDataGetTail, tailDate)
	if err != nil {
		return []string{}, errs.NewDalError("DataAudit.GetTail", "list failed", err)
	}

	return libdom.MapEntitiesToIDs(res), nil
}

// CutTail выполняет физическое удаление записи аудита данных из "хвоста" таблицы по её идентификатору.
func (da *DataAudit) CutTail(ctx context.Context, id string) error {
	// fail fast
	if strings.TrimSpace(id) == "" {
		return errs.NewInvalidArgumentError("ID", "id can not be empty")
	}

	return da.Delete(ctx, id)
}

// ListByPeriod возвращает срез записей аудита данных за указанный временной период с фильтрацией по лимиту и смещению.
func (da *DataAudit) ListByPeriod(ctx context.Context, from, till time.Time, limit, offset int) ([]*domain.DataAudit, error) {
	// fail fast
	if err := da.validateListByPeriod(from, till, limit, offset); err != nil {
		return nil, errs.NewDalError("DataAudit.ListByPeriod", "validate failed", err)
	}
	// act
	res, err := da.GetHelper().List(ctx, statementDataListByPeriod, from, till, limit, offset)
	if err != nil {
		return nil, errs.NewDalError("DataAudit.ListByPeriod", "list failed", err)
	}

	return res, nil
}

func (da *DataAudit) validateListByPeriod(from, till time.Time, limit, offset int) error {
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

// ListByInstance возвращает срез записей аудита для конкретного экземпляра сущности по её типу и уникальному ID.
func (da *DataAudit) ListByInstance(ctx context.Context, typeName string, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	// fail fast
	if err := da.validateListByInstance(typeName, instanceID, offset, limit); err != nil {
		return nil, errs.NewDalError("AuthAudit.ListByUsername", "validate failed", err)
	}
	// act
	res, err := da.GetHelper().List(ctx, statementDataListByInstance, typeName, instanceID, limit, offset)
	if err != nil {
		return nil, errs.NewDalError("AuthAudit.ListByUsername", "list failed", err)
	}

	return res, nil
}

//goland:noinspection DuplicatedCode
func (da *DataAudit) validateListByInstance(typeName, instanceID string, limit, offset int) error {
	if strings.TrimSpace(typeName) == "" {
		return errs.NewInvalidArgumentError("typeName", "typeName is required")
	}
	if strings.TrimSpace(instanceID) == "" {
		return errs.NewInvalidArgumentError("instanceID", "instanceID is required")
	}
	if limit <= 0 {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if offset < 0 {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}

func (da *DataAudit) creator(ctx context.Context, querier db.Querier, entity *domain.DataAudit, params ...any) (librepo.Scannable, error) {
	query, err := da.GetHelper().GetStatements().GetSQL(librepo.StatementCreate)
	if err != nil {
		return nil, errs.NewDalError("DataAudit.creator", "create statement not found", err)
	}
	var valuesRaw []byte
	if len(entity.Values) > 0 {
		valuesRaw, err = json.Marshal(entity.Values)
		if err != nil {
			return nil, errs.NewDalError("DataAudit.creator", "marshal values json", err)
		}
	}

	return querier.QueryRowContext(ctx, query,
		entity.ID,
		entity.Source, entity.EventDate, entity.Event, entity.Status,
		entity.RequestID, entity.TraceID, entity.Username,
		entity.InternalTypeName, entity.TypeName, entity.TypeDescription, entity.InstanceID, entity.InstanceName,
		valuesRaw,
		entity.CreatedAt,
	), nil
}

func (da *DataAudit) changer(ctx context.Context, querier db.Querier, entity *domain.DataAudit, params ...any) (librepo.Scannable, error) {
	query, err := da.GetHelper().GetStatements().GetSQL(librepo.StatementChange)
	if err != nil {
		return nil, errs.NewDalError("DataAudit.changer", "change statement not found", err)
	}
	var valuesRaw []byte
	if len(entity.Values) > 0 {
		valuesRaw, err = json.Marshal(entity.Values)
		if err != nil {
			return nil, errs.NewDalError("DataAudit.changer", "marshal values json", err)
		}
	}

	return querier.QueryRowContext(ctx, query,
		entity.ID,
		entity.Source, entity.EventDate, entity.Event, entity.Status,
		entity.RequestID, entity.TraceID, entity.Username,
		entity.InternalTypeName, entity.TypeName, entity.TypeDescription, entity.InstanceID, entity.InstanceName,
		valuesRaw,
		entity.UpdatedAt,
	), nil
}

//goland:noinspection DuplicatedCode
func (da *DataAudit) mapDefault(ctx context.Context, scannable librepo.Scannable, params ...any) (*domain.DataAudit, error) {
	res := da.GetHelper().GetHooks().NewEntityFactory()
	var valuesRaw []byte
	err := scannable.Scan(
		&res.ID,
		&res.Source, &res.EventDate, &res.Event, &res.Status,
		&res.RequestID, &res.TraceID, &res.Username,
		&res.InternalTypeName, &res.TypeName, &res.TypeDescription, &res.InstanceID, &res.InstanceName,
		&valuesRaw,
		&res.CreatedAt, &res.UpdatedAt,
	)
	if err != nil {
		return nil, errs.NewDalError("DataAudit.mapDefault", "default mapper failed", err)
	}
	if len(valuesRaw) > 0 {
		if err := json.Unmarshal(valuesRaw, &res.Values); err != nil {
			return nil, errs.NewDalError("DataAudit.mapDefault", "unmarshal values json", err)
		}
	}

	return res, nil
}

func (da *DataAudit) mapTail(ctx context.Context, scannable librepo.Scannable, params ...any) (*domain.DataAudit, error) {
	res := da.GetHelper().GetHooks().NewEntityFactory()
	err := scannable.Scan(
		&res.ID,
	)
	if err != nil {
		return nil, errs.NewDalError("DataAudit.mapTail", "tail mapper failed", err)
	}

	return res, nil
}

//goland:noinspection DuplicatedCode
func (da *DataAudit) pgStatementsBuilder() *librepo.Statements[*domain.DataAudit, string] {
	return librepo.NewStatements[*domain.DataAudit, string]().
		WithStatement(librepo.StatementFind, librepo.NewStatement(postgres.SqlDataAuditFind, da.mapDefault)).
		WithStatement(librepo.StatementList, librepo.NewStatement(postgres.SqlDataAuditList, da.mapDefault)).
		WithStatement(librepo.StatementCreate, librepo.NewStatement(postgres.SqlDataAuditCreate, da.mapDefault)).
		WithStatement(librepo.StatementChange, librepo.NewStatement(postgres.SqlDataAuditChange, da.mapDefault)).
		WithStatement(librepo.StatementDelete, librepo.NewStatement[*domain.DataAudit, string](postgres.SqlDataAuditDelete, nil)).
		WithStatement(statementDataGetTail, librepo.NewStatement(postgres.SqlDataAuditTailList, da.mapTail)).
		WithStatement(statementDataListByPeriod, librepo.NewStatement(postgres.SqlDataAuditListByPeriod, da.mapDefault)).
		WithStatement(statementDataListByInstance, librepo.NewStatement(postgres.SqlDataAuditListByInstance, da.mapDefault))
}
