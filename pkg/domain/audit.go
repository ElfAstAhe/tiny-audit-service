package domain

// AuditField encapsulates a single granular snapshot of an entity field data payload,
// combining its raw serialized string value with a human-readable metadata description label.
type AuditField struct {
	Value       string
	Description string
}

// NewAuditField acts as a factory constructor allocating an immutable field snapshot tuple.
func NewAuditField(value string, description string) *AuditField {
	return &AuditField{
		Value:       value,
		Description: description,
	}
}

// Auditable defines the core behavioral contract that domain aggregate models must implement
// to support automated differential state reflection, hash integrity validation, and structured audit registry mapping.
type Auditable interface {
	// GetInternalTypeName unpacks low-level structural reflection names detailing underlying domains taxonomy models.
	GetInternalTypeName() string
	// GetTypeName outputs the flat canonical string alias representation characterizing this class configuration node.
	GetTypeName() string
	// GetTypeDescription outputs short descriptive contextual explanation metrics labels.
	GetTypeDescription() string
	// GetInstanceID extracts the unique primary identity record index allocated within downstream data layers.
	GetInstanceID() string
	// GetInstanceName maps generic literal indicators for centralized system monitoring operations.
	GetInstanceName() string

	// HashCode serializes the cumulative scalar state parameters into a stable cryptographic checksum.
	HashCode() uint32
	// ToAuditMap build map, key - field name, value - field value + field description (look for AuditField structure)
	ToAuditMap() map[string]*AuditField
}
