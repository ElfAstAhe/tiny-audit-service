package domain

type AuditField struct {
	Value       string
	Description string
}

func NewAuditField(value string, description string) *AuditField {
	return &AuditField{
		Value:       value,
		Description: description,
	}
}

type Auditable interface {
	GetInternalTypeName() string
	GetTypeName() string
	GetTypeDescription() string
	GetInstanceID() string
	GetInstanceName() string

	HashCode() uint32
	// ToAuditMap build map, key - field name, value - field value + field description (look for AuditField structure)
	ToAuditMap() map[string]*AuditField
}
