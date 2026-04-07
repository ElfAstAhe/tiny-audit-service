package utils

import (
	"time"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/api/http/audit/v1/models"
)

type AuthAuditBuilder struct {
	instance *models.AuthAuditDTO
}

func NewAuthAuditBuilder() *AuthAuditBuilder {
	return &AuthAuditBuilder{
		instance: new(models.AuthAuditDTO),
	}
}

func (aab *AuthAuditBuilder) NewInstance() *AuthAuditBuilder {
	aab.instance = new(models.AuthAuditDTO)
	return aab
}

func (aab *AuthAuditBuilder) WithSource(source string) *AuthAuditBuilder {
	aab.instance.Source = source
	return aab
}

func (aab *AuthAuditBuilder) WithEventDate(eventDate time.Time) *AuthAuditBuilder {
	aab.instance.EventDate = eventDate.Format(time.RFC3339)
	return aab
}

func (aab *AuthAuditBuilder) WithEvent(event string) *AuthAuditBuilder {
	aab.instance.Event = event
	return aab
}

func (aab *AuthAuditBuilder) WithStatus(status string) *AuthAuditBuilder {
	aab.instance.Status = status
	return aab
}

func (aab *AuthAuditBuilder) WithRequestID(requestID string) *AuthAuditBuilder {
	aab.instance.RequestID = requestID
	return aab
}

func (aab *AuthAuditBuilder) WithTraceID(traceID string) *AuthAuditBuilder {
	aab.instance.TraceID = traceID
	return aab
}

func (aab *AuthAuditBuilder) WithUsername(username string) *AuthAuditBuilder {
	aab.instance.Username = username
	return aab
}

func (aab *AuthAuditBuilder) WithAccessToken(accessToken string) *AuthAuditBuilder {
	aab.instance.AccessToken = accessToken
	return aab
}

func (aab *AuthAuditBuilder) WithRefreshToken(refreshToken string) *AuthAuditBuilder {
	aab.instance.RefreshToken = refreshToken
	return aab
}

func (aab *AuthAuditBuilder) Build() *models.AuthAuditDTO {
	return aab.instance
}
