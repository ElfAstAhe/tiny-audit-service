package tiny_audit_service

const (
	sqlCreateTableDataAudit = `
create table if not exists data_audit (
    id varchar(50) not null,
    source varchar(100) null,
    event_date timestamptz null,
    event varchar(50) null,
    status varchar(50) null,
    request_id varchar(100) null,
    trace_id varchar(100) null,
    username varchar(100) null,
    internal_type_name varchar(1024) null,
    type_name varchar(256) null,
    type_description varchar(512) null,
    instance_id varchar(50) null,
    instance_name varchar(100) null,
    values text null,
    created_at timestamptz null default now(),
    updated_at timestamptz null default now(),
    constraint data_audit_pk primary key (id)
)
`
	sqlDropTableDataAudit = `
drop table if exists data_audit cascade
`
	sqlCreateIndexListByInstance = `
create index if not exists idx_data_audit_instance on data_audit (type_name asc, instance_id asc, created_at desc)
`
	sqlDropIndexListByInstance = `
drop index if exists idx_data_audit_instance
`
)
