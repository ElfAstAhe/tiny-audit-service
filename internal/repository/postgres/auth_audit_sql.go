package postgres

// Константы, содержащие SQL-запросы для взаимодействия с таблицей auth_audit в СУБД PostgreSQL.
const (
	// SqlAuthAuditFind содержит SQL-запрос для выборки одной записи аудита аутентификации по её уникальному идентификатору.
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
	// SqlAuthAuditList содержит SQL-запрос для постраничной выборки всех записей аудита аутентификации.
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
	// SqlAuthAuditListByPeriod содержит SQL-запрос для постраничной выборки записей аудита аутентификации за определенный период времени.
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
	// SqlAuthAuditListByUsername содержит SQL-запрос для выборки записей событий аудита для конкретного пользователя.
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
	// SqlAuthAuditCreate содержит SQL-запрос для вставки новой записи в таблицу аудита аутентификации.
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
	// SqlAuthAuditChange содержит SQL-запрос для обновления существующей записи аудита аутентификации.
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
	// SqlAuthAuditDelete содержит SQL-запрос для удаления записи аудита аутентификации по её идентификатору.
	SqlAuthAuditDelete string = `
delete
from
    auth_audit
where
    id = $1
`
	// SqlAuthAuditTailList содержит SQL-запрос для извлечения идентификаторов старых записей, попадающих под очистку хвостов таблиц auth_audit.
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
