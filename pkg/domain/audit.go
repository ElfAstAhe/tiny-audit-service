package domain

type AuditField struct {
	Description string
	Value       string
}

func NewAuditField(description string, value string) *AuditField {
	return &AuditField{
		Description: description,
		Value:       value,
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
