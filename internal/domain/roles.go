package domain

// Константы, определяющие прикладные роли пользователей (RBAC) в системе аудита.
const (
	// RoleWriter предоставляет права на запись и регистрацию новых логов аудита.
	RoleWriter string = "audit-writer"

	// RoleReader предоставляет права на чтение, поиск и просмотр архивных записей аудита.
	RoleReader string = "audit-reader"

	// RoleAdmin предоставляет полные административные права на управление системой аудита.
	RoleAdmin string = "audit-admin"
)
