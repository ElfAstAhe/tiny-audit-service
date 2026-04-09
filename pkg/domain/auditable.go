package domain

type Auditable interface {
	GetInternalTypeName() string
	GetTypeName() string
	GetTypeDescription() string
	GetInstanceID() string
	GetInstanceName() string

	HashCode() uint32
	ToAuditMap() map[string]string
}
