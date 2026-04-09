package tiny_audit_service

import (
	"context"
	"database/sql"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/pressly/goose/v3"
)

func up0002(ctx context.Context, db *sql.DB) error {
	if err := createTableDataAudit(ctx, db); err != nil {
		return err
	}

	return createIndexListByInstance(ctx, db)
}

func down0002(ctx context.Context, db *sql.DB) error {
	if err := dropIndexListByInstance(ctx, db); err != nil {
		return err
	}

	return dropTableDataAudit(ctx, db)
}

func createTableDataAudit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateTableDataAudit)
	if err != nil {
		return errs.NewDBMigrationError("create table data_audit", err)
	}

	return nil
}

func createIndexListByInstance(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateIndexListByInstance)
	if err != nil {
		return errs.NewDBMigrationError("create index idx_data_audit_instance", err)
	}

	return nil
}

func dropTableDataAudit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropTableDataAudit)
	if err != nil {
		return errs.NewDBMigrationError("drop table data_audit", err)
	}

	return nil
}

func dropIndexListByInstance(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropIndexListByInstance)
	if err != nil {
		return errs.NewDBMigrationError("drop index idx_data_audit_instance", err)
	}

	return nil
}

func init() {
	goose.AddMigrationNoTxContext(up0002, down0002)
}
