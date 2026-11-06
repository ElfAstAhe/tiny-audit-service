package postgres

const (
	SqlAuthAuditFind string = `
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
    auth_audit
where
    id = $1
`
	SqlAuthAuditList string = `
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
    auth_audit
order by
    event_date desc,
    username asc
offset $2
limit $1
`
	SqlAuthAuditListByPeriod string = `
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
    auth_audit
where
    event_date >= $1
and event_date < $2
order by
    event_date desc,
    username asc
offset $4
limit $3
`
	SqlAuthAuditListByUsername string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
from
    auth_audit
where
    username = $1
order by
    event_date desc
offset $3
limit $2
`
	SqlAuthAuditCreate string = `
insert into auth_audit (
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    access_token,
    refresh_token,
    created_at
)
values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
)
returning
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
`
	SqlAuthAuditChange string = `
update
    auth_audit
set
    source = $2,
    event_date = $3,
    event = $4,
    status = $5,
    request_id = $6,
    trace_id = $7,
    username = $8,
    access_token = $9,
    refresh_token = $10,
    updated_at = $11
where
    id = $1
returning
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    access_token,
    refresh_token,
    created_at,
    updated_at
`
	SqlAuthAuditDelete string = `
delete
from
    auth_audit
where
    id = $1
`
	SqlAuthAuditTailList string = `
select
    id
from
    auth_audit
where
    event_date < $1
order by
    event_date desc
limit 50
`
)
