package rbac

import (
	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

// IsSubjectAdmin проверяет, обладает ли переданный субъект авторизации (пользователь)
// глобальными административными правами доступа RoleAdmin.
func IsSubjectAdmin(subject *auth.Subject) bool {
	if subject == nil {
		return false
	}

	return subject.HasRole(domain.RoleAdmin)
}
