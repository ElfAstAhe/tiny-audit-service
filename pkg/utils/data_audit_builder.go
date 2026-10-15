package utils

import (
	"time"

	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
)

// DataAuditBuilder provides a fluent builder API to incrementally assemble and instantiate
// decoupled, strongly-typed dto.DataAuditDTO payload structures.
type DataAuditBuilder struct {
	instance *dto.DataAuditDTO // Underlying core data transfer object instance state payload
}

// NewDataAuditBuilder acts as a factory constructor allocating a fresh builder state holding an empty data target.
func NewDataAuditBuilder() *DataAuditBuilder {
	return &DataAuditBuilder{
		instance: new(dto.DataAuditDTO),
	}
}

// WithSource maps the system node or originating framework context tag descriptor string onto the object.
func (dab *DataAuditBuilder) WithSource(source string) *DataAuditBuilder {
	dab.instance.Source = source
	return dab
}

// WithEventDate records the precise temporal execution clock boundaries configuration parameter payload.
func (dab *DataAuditBuilder) WithEventDate(eventDate time.Time) *DataAuditBuilder {
	dab.instance.EventDate = eventDate
	return dab
}

// WithEvent sets the specific domain operation or audit transaction hook event type category string.
func (dab *DataAuditBuilder) WithEvent(event string) *DataAuditBuilder {
	dab.instance.Event = event
	return dab
}

// WithStatus maps the target final processing completion code or transport state outcome flag.
func (dab *DataAuditBuilder) WithStatus(status string) *DataAuditBuilder {
	dab.instance.Status = status
	return dab
}

// WithRequestID binds the correlation networking unique tracking sequence identifier token string.
func (dab *DataAuditBuilder) WithRequestID(requestID string) *DataAuditBuilder {
	dab.instance.RequestID = requestID
	return dab
}

// WithTraceID injects active distributed observability trace indicators parameters values.
func (dab *DataAuditBuilder) WithTraceID(traceID string) *DataAuditBuilder {
	dab.instance.TraceID = traceID
	return dab
}

// WithUsername sets caller identifier properties managing non-anonymous audit security contexts.
func (dab *DataAuditBuilder) WithUsername(username string) *DataAuditBuilder {
	dab.instance.Username = username
	return dab
}

// WithInternalTypeName maps metadata structural reflection names detailing underlying domains taxonomy models.
func (dab *DataAuditBuilder) WithInternalTypeName(internalTypeName string) *DataAuditBuilder {
	dab.instance.InternalTypeName = internalTypeName
	return dab
}

// WithTypeName populates flat canonical string alias schemas defining targeted application configurations nodes.
func (dab *DataAuditBuilder) WithTypeName(typeName string) *DataAuditBuilder {
	dab.instance.TypeName = typeName
	return dab
}

// WithTypeDescription logs short descriptive contextual human-readable explanation metrics labels.
func (dab *DataAuditBuilder) WithTypeDescription(typeDescription string) *DataAuditBuilder {
	dab.instance.TypeDescription = typeDescription
	return dab
}

// WithInstanceID pins the unique record index entity state handle allocated inside downstream СУБД systems.
func (dab *DataAuditBuilder) WithInstanceID(instanceID string) *DataAuditBuilder {
	dab.instance.InstanceID = instanceID
	return dab
}

// WithInstanceName attaches target literal labels representation for centralized operational reporting metrics loops.
func (dab *DataAuditBuilder) WithInstanceName(instanceName string) *DataAuditBuilder {
	dab.instance.InstanceName = instanceName
	return dab
}

// WithValues injects a structured array collection listing detailed properties variations DTO schemas.
func (dab *DataAuditBuilder) WithValues(values []*dto.DataAuditValueDTO) *DataAuditBuilder {
	dab.instance.Values = values
	return dab
}

// Build finalizes current setup parameters configurations steps to yield the fully prepared data payload wrapper structure pointer.
func (dab *DataAuditBuilder) Build() *dto.DataAuditDTO {
	return dab.instance
}
