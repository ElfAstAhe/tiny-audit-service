package dto

import (
	"time"
)

// AuthAuditDTO является транспортным объектом (Data Transfer Object) доменного уровня,
// инкапсулирующим исчерпывающие данные о событиях контура аутентификации и авторизации.
//
// Структура полностью адаптирована для высоконагруженных сценариев сериализации, поддерживает
// сквозную трассировку и оптимизирована под стандарты Go 1.26 для предотвращения утечек дефолтных значений времени.
type AuthAuditDTO struct {
	Source       string    `json:"source,omitempty"`        // Имя микросервиса-источника, сгенерировавшего лог авторизации
	EventDate    time.Time `json:"event_date,omitzero"`     // Точная временная метка наступления события (полностью исключается из JSON, если пустая)
	Event        string    `json:"event,omitempty"`         // Тип совершенного действия (login/logoff/refresh/revoke из enums.go)
	Status       string    `json:"status,omitempty"`        // Итоговый статус выполнения прикладной операции (success/fail)
	RequestID    string    `json:"request_id,omitempty"`    // Уникальный идентификатор входящего прикладного HTTP/gRPC запроса
	TraceID      string    `json:"trace_id,omitempty"`      // Идентификатор сквозной распределенной трассировки OpenTelemetry
	Username     string    `json:"username,omitempty"`      // Имя или идентификатор пользователя, инициировавшего действие
	AccessToken  string    `json:"access_token,omitempty"`  // Системный токен доступа (или его хеш-слепок) для паспортизации сессии
	RefreshToken string    `json:"refresh_token,omitempty"` // Токен обновления сессии для долгоживущих соединений
}
