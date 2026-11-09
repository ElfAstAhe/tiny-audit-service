package postgres

// Константы, содержащие SQL-запросы для взаимодействия с таблицей data_audit в СУБД PostgreSQL.
const (
	// SqlDataAuditFind содержит SQL-запрос для выборки одной записи аудита данных по её уникальному идентификатору.
	SqlDataAuditFind string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    internal_type_name,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
from
    data_audit
where
    id = $1
`
	// SqlDataAuditList содержит SQL-запрос для постраничной выборки всех записей аудита данных.
	SqlDataAuditList string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    internal_type_name,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
from
    data_audit
order by
    event_date desc,
    type_name asc,
    instance_id asc
offset $2
limit $1
`
	// SqlDataAuditListByPeriod содержит SQL-запрос для постраничной выборки записей аудита данных за определенный период времени.
	SqlDataAuditListByPeriod string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    internal_type_name,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
from
    data_audit
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
	// SqlDataAuditListByInstance содержит SQL-запрос для выборки изменений конкретного экземпляра сущности по типу и ID.
	SqlDataAuditListByInstance string = `
select
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    internal_type_name,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
from
    data_audit
where
    type_name = $1
and instance_id = $2
order by
    event_date desc
offset $4
limit $3
`
	// SqlDataAuditCreate содержит SQL-запрос для вставки новой записи в таблицу аудита данных.
	SqlDataAuditCreate string = `
insert into data_audit (
    id,
    source,
    event_date,
    event,
    status,
    request_id,
    trace_id,
    username,
    internal_type_name,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at
)
values (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
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
    internal_type_name,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
`
	// SqlDataAuditChange содержит SQL-запрос для обновления существующей записи аудита данных.
	SqlDataAuditChange string = `
update
    data_audit
set
    source = $2,
    event_date = $3,
    event = $4,
    status = $5,
    request_id = $6,
    trace_id = $7,
    username = $8,
    internal_type_name = $9,
    type_name = $10,
    type_description = $11,
    instance_id = $12,
    instance_name = $13,
    values = $14,
    updated_at = $15
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
    internal_type_name,
    type_name,
    type_description,
    instance_id,
    instance_name,
    values,
    created_at,
    updated_at
`
	// SqlDataAuditDelete содержит SQL-запрос для удаления записи аудита данных по её идентификатору.
	SqlDataAuditDelete string = `
delete
from
    data_audit
where
    id = $1
`
	// SqlDataAuditTailList содержит SQL-запрос для извлечения идентификаторов старых записей, попадающих под очистку хвостов таблиц.
	SqlDataAuditTailList string = `
select
    id
from
    data_audit
where
    event_date < $1
order by
    event_date desc
limit 50
`
)
