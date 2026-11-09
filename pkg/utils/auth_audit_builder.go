package utils

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

// AuthAuditBuilder представляет собой строитель (Builder) для пошагового
// и безопасного конструирования объектов dto.AuthAuditDTO.
type AuthAuditBuilder struct {
	instance *dto.AuthAuditDTO
}

var _ utils.Builder[*dto.AuthAuditDTO] = (*AuthAuditBuilder)(nil)

// NewAuthAuditBuilder создает и инициализирует новый экземпляр строителя AuthAuditBuilder.
func NewAuthAuditBuilder() *AuthAuditBuilder {
	return &AuthAuditBuilder{
		instance: new(dto.AuthAuditDTO),
	}
}

// WithSource устанавливает источник события (Source) в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithSource(source string) *AuthAuditBuilder {
	aab.instance.Source = source
	return aab
}

// WithEventDate устанавливает дату и время фиксации события (EventDate) в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithEventDate(eventDate time.Time) *AuthAuditBuilder {
	aab.instance.EventDate = eventDate
	return aab
}

// WithEvent устанавливает тип или наименование события (Event) в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithEvent(event string) *AuthAuditBuilder {
	aab.instance.Event = event
	return aab
}

// WithStatus устанавливает статус выполнения операции (Status) в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithStatus(status string) *AuthAuditBuilder {
	aab.instance.Status = status
	return aab
}

// WithRequestID устанавливает сквозной идентификатор запроса (RequestID) в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithRequestID(requestID string) *AuthAuditBuilder {
	aab.instance.RequestID = requestID
	return aab
}

// WithTraceID устанавливает идентификатор распределенной трассировки (TraceID) в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithTraceID(traceID string) *AuthAuditBuilder {
	aab.instance.TraceID = traceID
	return aab
}

// WithUsername устанавливает имя пользователя (Username), инициировавшего операцию, в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithUsername(username string) *AuthAuditBuilder {
	aab.instance.Username = username
	return aab
}

// WithAccessToken устанавливает токен доступа (AccessToken) в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithAccessToken(accessToken string) *AuthAuditBuilder {
	aab.instance.AccessToken = accessToken
	return aab
}

// WithRefreshToken устанавливает токен обновления (RefreshToken) в конструируемый объект DTO.
func (aab *AuthAuditBuilder) WithRefreshToken(refreshToken string) *AuthAuditBuilder {
	aab.instance.RefreshToken = refreshToken
	return aab
}

// Build завершает процесс конфигурации и возвращает сформированный экземпляр dto.AuthAuditDTO.
func (aab *AuthAuditBuilder) Build() *dto.AuthAuditDTO {
	return aab.instance
}
