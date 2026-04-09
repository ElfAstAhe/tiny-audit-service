package tiny_audit_service

const (
	sqlCreateTableAuthAudit = `
create table if not exists auth_audit (
    id varchar(50) not null,
    source varchar(100) null,
    event_date timestamptz null,
    event varchar(50) null,
    status varchar(50) null,
    request_id varchar(100) null,
    trace_id varchar(100) null,
    username varchar(100) null,
    access_token varchar(4000) null,
    refresh_token varchar(4000) null,
    created_at timestamptz null default now(),
    updated_at timestamptz null default now(),
    constraint auth_audit_pk primary key (id)
)
`
	sqlDropTableAuthAudit = `
drop table if exists auth_audit cascade
`
	sqlCreateIndexListByUsername = `
create index if not exists idx_auth_audit_username on auth_audit(username asc, event_date desc)
`
	sqlDropIndexListByUsername = `
drop index if exists idx_auth_audit_username
`
)
