package postgres

const (
	sqlDataAuditFind string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
from
    data_audit_1
where
    id = $1
`
	sqlDataAuditList string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
from
    data_audit_1
order by
    event_date desc,
    type_name asc,
    instance_id asc
offset $2
limit $1
`
	sqlDataAuditListByPeriod string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
from
    data_audit_1
where
    event_date >= $1
and event_date < $2
order by
    event_date desc,
    type_name asc,
    instance_id asc
offset $4
limit $3
`
	sqlDataAuditListByInstance string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
from
    data_audit_1
where
    type_name = $1
and instance_id = $2
order by
    event_date desc
offset $4
limit $3
`
	sqlDataAuditCreate string = `
insert into data_audit_1 (
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at
)
values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13
)
returning
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    username,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
`
	sqlDataAuditChange string = `
update
    data_audit_1
set
    source = $2,
    event_date = $3,
    event = $4,
    status = $5,
    request_id = $6,
    username = $7,
    type_name = $8,
    type_description = $9,
    instance_id = $10,
    instance_name = $11,
    values = $12,
    updated_at = $13
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
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
`
	sqlDataAuditDelete string = `
delete
from
    data_audit_1
where
    id = $1
`
	sqlDataAuditTailList string = `
select
    id
from
    data_audit_1
where
    event_date < $1
order by
    event_date desc
limit 50
`
)
