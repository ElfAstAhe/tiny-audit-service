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
	sqlTemplateCreateTableDataAudit = `
create table if not exists %table_name (
    id varchar(50) not null,
    source varchar(100) null,
    event_date timestamptz null,
    event varchar(50) null,
    status varchar(50) null,
    request_id varchar(100) null,
    username varchar(100) null,
    type_name varchar(256) null,
    type_description varchar(512) null,
    instance_id varchar(50) null,
    instance_name varchar(100) null,
    values text null,
    created_at timestamptz null default now(),
    updated_at timestamptz null default now(),
    constraint %constraint_pk primary key (id)
)
`
	sqlTemplateDropTableDataAudit = `
drop table if exists %table_name cascade
`
	sqlTemplateCreateIndexListByInstance = `
create index if not exists %idx_data_audit_instance on %table_name(type_name asc, instance_id asc, created_at desc)
`
	sqlTemplateDropIndexListByInstance = `
drop index if exists %idx_data_audit_instance
`
)

func up0002(ctx context.Context, db *sql.DB) error {
	if err := createTablesDataAudit(ctx, db); err != nil {
		return err
	}

	return createIndexesListByInstance(ctx, db)
}

func down0002(ctx context.Context, db *sql.DB) error {
	if err := dropIndexesListByInstance(ctx, db); err != nil {
		return err
	}

	return dropTablesDataAudit(ctx, db)
}

func createTablesDataAudit(ctx context.Context, db *sql.DB) error {
	var errs []error
	for index := 1; index < 6; index++ {
		errs = append(errs, createTableDataAudit(ctx, db, index))
	}

	return errors.Join(errs...)
}

func createTableDataAudit(ctx context.Context, db *sql.DB, index int) error {
	tableName := libdb.BuildRepeatedObjectName("data_audit", strconv.Itoa(index))
	sqlCreate := strings.Replace(
		sqlTemplateCreateTableDataAudit,
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

func createIndexesListByInstance(ctx context.Context, db *sql.DB) error {
	var errs []error
	for index := 1; index < 6; index++ {
		errs = append(errs, createIndexListByInstance(ctx, db, index))
	}

	return errors.Join(errs...)
}

func createIndexListByInstance(ctx context.Context, db *sql.DB, index int) error {
	sqlCreate := strings.Replace(
		sqlTemplateCreateIndexListByInstance,
		"%idx_data_audit_instance",
		libdb.BuildRepeatedObjectName("idx_data_audit_instance", strconv.Itoa(index)),
		1,
	)
	sqlCreate = strings.Replace(
		sqlCreate,
		"%table_name",
		libdb.BuildRepeatedObjectName("data_audit", strconv.Itoa(index)),
		1,
	)
	_, err := db.ExecContext(ctx, sqlCreate)
	if err != nil {
		return err
	}

	return nil
}

func dropTablesDataAudit(ctx context.Context, db *sql.DB) error {
	var errs []error
	for index := 1; index < 6; index++ {
		errs = append(errs, dropTableDataAudit(ctx, db, index))
	}

	return errors.Join(errs...)
}

func dropTableDataAudit(ctx context.Context, db *sql.DB, index int) error {
	_, err := db.ExecContext(
		ctx,
		strings.Replace(
			sqlTemplateDropTableDataAudit,
			"%table_name",
			libdb.BuildRepeatedObjectName("data_audit", strconv.Itoa(index)),
			1,
		),
	)

	return err
}

func dropIndexesListByInstance(ctx context.Context, db *sql.DB) error {
	var errs []error
	for index := 1; index < 6; index++ {
		errs = append(errs, dropIndexListByInstance(ctx, db, index))
	}

	return errors.Join(errs...)
}

func dropIndexListByInstance(ctx context.Context, db *sql.DB, index int) error {
	_, err := db.ExecContext(
		ctx,
		strings.Replace(
			sqlTemplateDropIndexListByInstance,
			"%idx_data_audit_instance",
			libdb.BuildRepeatedObjectName("idx_data_audit_instance", strconv.Itoa(index)),
			1,
		),
	)

	return err
}

func init() {
	goose.AddMigrationNoTxContext(up0002, down0002)
}
