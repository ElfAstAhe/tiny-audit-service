package rbac

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

// AuthAudit представляет собой декоратор контроля доступа (RBAC) для фасада facade.AuthAudit,
// перехватывающий вызовы для валидации ролей пользователя (Subject) в контексте выполнения.
type AuthAudit struct {
	authHelper auth.Helper
	next       facade.AuthAudit
}

var _ facade.AuthAudit = (*AuthAudit)(nil)

// NewAuthAudit создает и инициализирует новый экземпляр RBAC-декоратора AuthAudit.
func NewAuthAudit(
	helper auth.Helper,
	next facade.AuthAudit,
) *AuthAudit {
	return &AuthAudit{
		authHelper: helper,
		next:       next,
	}
}

// Audit извлекает субъект из контекста и проверяет наличие ролей Writer или Admin перед передачей вызова на запись аудит-лога.
func (aa *AuthAudit) Audit(ctx context.Context, data *dto.AuthAuditDTO) error {
	// subject
	subj, err := aa.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return errs.NewBllForbiddenError("AuthAuditFacadeImpl.Audit", "retrieve subject", err)
	}
	// rbac
	if !subj.HasRole(domain.RoleWriter) && !subj.HasRole(domain.RoleAdmin) {
		return errs.NewBllForbiddenError("AuthAuditFacadeImpl.Audit", "subject is not audit-writer", nil)
	}

	return aa.next.Audit(ctx, data)
}

// ListByPeriod извлекает субъект из контекста и проверяет наличие ролей Reader или Admin перед чтением логов за период.
func (aa *AuthAudit) ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.AuthAuditDTO, error) {
	// subject
	subj, err := aa.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, errs.NewBllForbiddenError("AuthAuditFacadeImpl.ListByPeriod", "retrieve subject", err)
	}
	// rbac
	if !subj.HasRole(domain.RoleReader) && !subj.HasRole(domain.RoleAdmin) {
		return nil, errs.NewBllForbiddenError("AuthAuditFacadeImpl.ListByPeriod", "subject is not audit-reader", nil)
	}

	return aa.next.ListByPeriod(ctx, auditPeriod)
}

// ListByUsername извлекает субъект из контекста и проверяет наличие ролей Reader или Admin перед поиском логов по имени пользователя.
func (aa *AuthAudit) ListByUsername(ctx context.Context, auditUser *dto.AuditUserDTO) ([]*dto.AuthAuditDTO, error) {
	// subject
	subj, err := aa.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, errs.NewBllForbiddenError("AuthAuditFacadeImpl.ListByUsername", "retrieve subject", err)
	}
	// rbac
	if !subj.HasRole(domain.RoleReader) && !subj.HasRole(domain.RoleAdmin) {
		return nil, errs.NewBllForbiddenError("AuthAuditFacadeImpl.ListByUsername", "subject is not audit-reader", nil)
	}

	return aa.next.ListByUsername(ctx, auditUser)
}
