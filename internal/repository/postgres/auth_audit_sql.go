package postgres

const (
	sqlAuthAuditFind string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
from
    auth_audit_1
where
    id = $1
`
	sqlAuthAuditList string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
from
    auth_audit_1
order by
    event_date desc,
    username asc
offset $2
limit $1
`
	sqlAuthAuditListByPeriod string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
from
    auth_audit_1
where
    event_date >= $1
and event_date < $2
order by
    event_date desc,
    username asc
offset $4
limit $3
`
	sqlAuthAuditListByUsername string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
from
    auth_audit_1
where
    username = $1
order by
    event_date desc
offset $3
limit $2
`
	sqlAuthAuditCreate string = `
insert into auth_audit_1 (
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    access_token,
    refresh_token,
    created_at
)
values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10
)
returning
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
`
	sqlAuthAuditChange string = `
update
    auth_audit_1
set
    source = $2,
    event_date = $3,
    event = $4,
    status = $5,
    request_id = $6,
    username = $7,
    access_token = $8,
    refresh_token = $9,
    updated_at = $10
where
    id = $1
returning
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
`
	sqlAuthAuditDelete string = `
delete
from
    auth_audit_1
where
    id = $1
`
)
