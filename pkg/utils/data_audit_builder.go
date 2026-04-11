package utils

import (
	"time"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

type DataAuditBuilder struct {
	instance *dto.DataAuditDTO
}

func NewDataAuditBuilder() *DataAuditBuilder {
	return &DataAuditBuilder{
		instance: new(dto.DataAuditDTO),
	}
}

func (dab *DataAuditBuilder) WithSource(source string) *DataAuditBuilder {
	dab.instance.Source = source
	return dab
}

func (dab *DataAuditBuilder) WithEventDate(eventDate time.Time) *DataAuditBuilder {
	dab.instance.EventDate = eventDate
	return dab
}

func (dab *DataAuditBuilder) WithEvent(event string) *DataAuditBuilder {
	dab.instance.Event = event
	return dab
}

func (dab *DataAuditBuilder) WithStatus(status string) *DataAuditBuilder {
	dab.instance.Status = status
	return dab
}

func (dab *DataAuditBuilder) WithRequestID(requestID string) *DataAuditBuilder {
	dab.instance.RequestID = requestID
	return dab
}

func (dab *DataAuditBuilder) WithTraceID(traceID string) *DataAuditBuilder {
	dab.instance.TraceID = traceID
	return dab
}

func (dab *DataAuditBuilder) WithUsername(username string) *DataAuditBuilder {
	dab.instance.Username = username
	return dab
}

func (dab *DataAuditBuilder) WithInternalTypeName(internalTypeName string) *DataAuditBuilder {
	dab.instance.InternalTypeName = internalTypeName
	return dab
}

func (dab *DataAuditBuilder) WithTypeName(typeName string) *DataAuditBuilder {
	dab.instance.TypeName = typeName
	return dab
}

func (dab *DataAuditBuilder) WithTypeDescription(typeDescription string) *DataAuditBuilder {
	dab.instance.TypeDescription = typeDescription
	return dab
}

func (dab *DataAuditBuilder) WithInstanceID(instanceID string) *DataAuditBuilder {
	dab.instance.InstanceID = instanceID
	return dab
}

func (dab *DataAuditBuilder) WithInstanceName(instanceName string) *DataAuditBuilder {
	dab.instance.InstanceName = instanceName
	return dab
}

func (dab *DataAuditBuilder) WithValues(values []*dto.DataAuditValueDTO) *DataAuditBuilder {
	dab.instance.Values = values
	return dab
}

func (dab *DataAuditBuilder) Build() *dto.DataAuditDTO {
	return dab.instance
}
