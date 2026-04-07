package tiny_audit_service

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"

	libdb "github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/pressly/goose/v3"
)

const (
	sqlAlterTableAuditAddTraceID string = `
alter table if exists %table_name add column if not exists trace_id varchar(100) null
`
	sqlAlterTableAuditDropTraceID string = `
alter table if exists %table_name drop column if exists trace_id
`
)

func up0003(ctx context.Context, db *sql.DB) error {
	if err := alterTablesAuditAddTraceID(ctx, db, "auth_audit"); err != nil {
		return err
	}
	if err := alterTablesAuditAddTraceID(ctx, db, "data_audit"); err != nil {
		return err
	}

	return nil
}

func down0003(ctx context.Context, db *sql.DB) error {
	if err := alterTablesAuditDropTraceID(ctx, db, "auth_audit"); err != nil {
		return err
	}
	if err := alterTablesAuditDropTraceID(ctx, db, "data_audit"); err != nil {
		return err
	}

	return nil
}

func alterTablesAuditAddTraceID(ctx context.Context, db *sql.DB, tableNane string) error {
	var errs []error
	for i := 1; i <= 5; i++ {
		tabName := libdb.BuildRepeatedObjectName(tableNane, strconv.Itoa(i))
		errs = append(errs, alterTableAuditAddTraceID(ctx, db, tabName))
	}

	return errors.Join(errs...)
}

func alterTableAuditAddTraceID(ctx context.Context, db *sql.DB, tableName string) error {
	sqlAlter := strings.Replace(sqlAlterTableAuditAddTraceID, `%table_name`, tableName, 1)
	_, err := db.ExecContext(ctx, sqlAlter)

	return err
}

func alterTablesAuditDropTraceID(ctx context.Context, db *sql.DB, tableName string) error {
	var errs []error
	for i := 1; i <= 5; i++ {
		tabName := libdb.BuildRepeatedObjectName(tableName, strconv.Itoa(i))
		errs = append(errs, alterTableAuditDropTraceID(ctx, db, tabName))
	}

	return errors.Join(errs...)
}

func alterTableAuditDropTraceID(ctx context.Context, db *sql.DB, tableName string) error {
	sqlAlter := strings.Replace(sqlAlterTableAuditDropTraceID, `%table_name`, tableName, 1)
	_, err := db.ExecContext(ctx, sqlAlter)

	return err
}

func init() {
	goose.AddMigrationNoTxContext(up0003, down0003)
}
