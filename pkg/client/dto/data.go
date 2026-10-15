package dto

import (
	"time"
)

// DataAuditValueDTO инкапсулирует детальный паспорт изменений конкретного атрибута (поля) бизнес-сущности.
//
// Используется для фиксации срезов состояний «До» (Before) и «После» (After) в рамках аудита мутаций данных.
type DataAuditValueDTO struct {
	Name        string `json:"name"`        // Системное техническое имя изменяемого поля структуры
	Description string `json:"description"` // Понятное человекочитаемое описание атрибута для интерфейса аналитики
	Before      string `json:"before"`      // Строковое представление значения поля ДО наступления события (мутации)
	After       string `json:"after"`       // Строковое представление значения поля ПОСЛЕ наступления события (мутации)
}

// DataAuditDTO является центральным транспортным объектом (Data Transfer Object) доменного уровня,
// аккумулирующим исчерпывающую информацию о событии модификации бизнес-данных в системе.
//
// Структура полностью оптимизирована для сериализации в высоконагруженных контурах и последующего
// сохранения в аналитических хранилищах (ClickHouse / Postgres), поддерживая сквозную трассировку запросов.
type DataAuditDTO struct {
	Source           string               `json:"source,omitempty"`             // Имя микросервиса-источника, сгенерировавшего лог аудита
	EventDate        time.Time            `json:"event_date,omitzero"`          // Точная временная метка наступления события в наносекундах
	Event            string               `json:"event,omitempty"`              // Тип совершенной мутации (create/change/remove из enums.go)
	Status           string               `json:"status,omitempty"`             // Итоговый статус прикладной операции (success/fail)
	RequestID        string               `json:"request_id,omitempty"`         // Уникальный идентификатор входящего прикладного HTTP/gRPC запроса
	TraceID          string               `json:"trace_id,omitempty"`           // Идентификатор сквозной распределенной трассировки OpenTelemetry
	Username         string               `json:"username,omitempty"`           // Имя или ID пользователя, совершившего целевое действие
	InternalTypeName string               `json:"internal_type_name,omitempty"` // Системный тип данных на уровне кодовой базы Go (например, *models.User)
	TypeName         string               `json:"type_name,omitempty"`          // Бизнес-ориентированное наименование сущности (например, "Пользователь")
	TypeDescription  string               `json:"type_description,omitempty"`   // Подробное человекочитаемое описание доменной сущности
	InstanceID       string               `json:"instance_id,omitempty"`        // Уникальный первичный ключ (UUID/ID) конкретного экземпляра сущности в БД
	InstanceName     string               `json:"instance_name,omitempty"`      // Значимое текстовое имя экземпляра (например, Email или Название компании)
	Values           []*DataAuditValueDTO `json:"values,omitempty"`             // Массив измененных полей с фиксацией состояний «До/После»
}
