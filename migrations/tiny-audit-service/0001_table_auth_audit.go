package tiny_audit_service

import (
	"context"
	"database/sql"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/pressly/goose/v3"
)

func up0001(ctx context.Context, db *sql.DB) error {
	if err := createTableAuthAudit(ctx, db); err != nil {
		return err
	}

	return createIndexListByUsername(ctx, db)
}

func down0001(ctx context.Context, db *sql.DB) error {
	if err := dropIndexListByUsername(ctx, db); err != nil {
		return err
	}

	return dropTableAuthAudit(ctx, db)
}

func createTableAuthAudit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateTableAuthAudit)
	if err != nil {
		return errs.NewDBMigrationError("create table auth_audit", err)
	}

	return nil
}

func createIndexListByUsername(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlCreateIndexListByUsername)
	if err != nil {
		return errs.NewDBMigrationError("create index idx_auth_audit_username", err)
	}

	return nil
}

func dropTableAuthAudit(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropTableAuthAudit)
	if err != nil {
		return errs.NewDBMigrationError("drop table auth_audit", err)
	}

	return nil
}

func dropIndexListByUsername(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, sqlDropIndexListByUsername)
	if err != nil {
		return errs.NewDBMigrationError("drop index idx_auth_audit_username", err)
	}

	return nil
}

func init() {
	goose.AddMigrationNoTxContext(up0001, down0001)
}
