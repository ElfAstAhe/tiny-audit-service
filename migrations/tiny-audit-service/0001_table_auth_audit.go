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
	sqlTemplateCreateTableAuthAudit = `
create table if not exists %table_name (
    id varchar(50) not null,
    source varchar(100) null,
    event_date timestamptz null,
    event varchar(50) null,
    status varchar(50) null,
    request_id varchar(100) null,
    username varchar(100) null,
    access_token varchar(4000) null,
    refresh_token varchar(4000) null,
    created_at timestamptz null default now(),
    updated_at timestamptz null default now(),
    constraint %constraint_pk primary key (id)
)
`
	sqlTemplateDropTableAuthAudit = `
drop table if exists %table_name cascade
`
	sqlTemplateCreateIndexListByUsername = `
create index if not exists %idx_auth_audit_username on %table_name(username asc, created_at desc)
`
	sqlTemplateDropIndexListByUsername = `
drop index if exists %idx_auth_audit_username
`
)

func up0001(ctx context.Context, db *sql.DB) error {
	if err := createTablesAuthAudit(ctx, db); err != nil {
		return err
	}

	return createIndexesListByUsername(ctx, db)
}

func down0001(ctx context.Context, db *sql.DB) error {
	if err := dropIndexesListByUsername(ctx, db); err != nil {
		return err
	}

	return dropTablesAuthAudit(ctx, db)
}

func createTablesAuthAudit(ctx context.Context, db *sql.DB) error {
	var errs []error
	for index := 1; index < 6; index++ {
		errs = append(errs, createTableAuthAudit(ctx, db, index))
	}

	return errors.Join(errs...)
}

func createTableAuthAudit(ctx context.Context, db *sql.DB, index int) error {
	tableName := libdb.BuildRepeatedObjectName("auth_audit", strconv.Itoa(index))
	sqlCreate := strings.Replace(
		sqlTemplateCreateTableAuthAudit,
		"%table_name",
		tableName,
		1,
	)
	sqlCreate = strings.Replace(
		sqlCreate,
		"%constraint_pk",
		libdb.BuildPKConstraintName(tableName),
		1,
	)
	_, err := db.ExecContext(
		ctx,
		sqlCreate,
	)

	return err
}

func createIndexesListByUsername(ctx context.Context, db *sql.DB) error {
	var errs []error
	for index := 1; index < 6; index++ {
		errs = append(errs, createIndexListByUsername(ctx, db, index))
	}

	return errors.Join(errs...)
}

func createIndexListByUsername(ctx context.Context, db *sql.DB, index int) error {
	sqlCreate := strings.Replace(
		sqlTemplateCreateIndexListByUsername,
		"%idx_auth_audit_username",
		libdb.BuildRepeatedObjectName("idx_auth_audit_username", strconv.Itoa(index)),
		1,
	)
	sqlCreate = strings.Replace(
		sqlCreate,
		"%table_name",
		libdb.BuildRepeatedObjectName("auth_audit", strconv.Itoa(index)),
		1,
	)
	_, err := db.ExecContext(ctx, sqlCreate)
	if err != nil {
		return err
	}

	return nil
}

func dropTablesAuthAudit(ctx context.Context, db *sql.DB) error {
	var errs []error
	for index := 1; index < 6; index++ {
		errs = append(errs, dropTableAuthAudit(ctx, db, index))
	}

	return errors.Join(errs...)
}

func dropTableAuthAudit(ctx context.Context, db *sql.DB, index int) error {
	_, err := db.ExecContext(
		ctx,
		strings.Replace(
			sqlTemplateDropTableAuthAudit,
			"%table_name",
			libdb.BuildRepeatedObjectName("auth_audit", strconv.Itoa(index)),
			1,
		),
	)

	return err
}

func dropIndexesListByUsername(ctx context.Context, db *sql.DB) error {
	var errs []error
	for index := 1; index < 6; index++ {
		errs = append(errs, dropIndexListByUsername(ctx, db, index))
	}

	return errors.Join(errs...)
}

func dropIndexListByUsername(ctx context.Context, db *sql.DB, index int) error {
	_, err := db.ExecContext(
		ctx,
		strings.Replace(
			sqlTemplateDropIndexListByUsername,
			"%idx_auth_audit_username",
			libdb.BuildRepeatedObjectName("idx_auth_audit_username", strconv.Itoa(index)),
			1,
		),
	)

	return err
}

func init() {
	goose.AddMigrationNoTxContext(up0001, down0001)
}
