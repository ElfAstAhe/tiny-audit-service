package rbac

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

// DataAudit представляет собой декоратор контроля доступа (RBAC) для фасада facade.DataAudit,
// перехватывающий вызовы для валидации ролей пользователя (Subject) в контексте выполнения.
type DataAudit struct {
	authHelper auth.Helper
	next       facade.DataAudit
}

var _ facade.DataAudit = (*DataAudit)(nil)

// NewDataAudit создает и инициализирует новый экземпляр RBAC-декоратора DataAudit.
func NewDataAudit(helper auth.Helper, next facade.DataAudit) *DataAudit {
	return &DataAudit{
		authHelper: helper,
		next:       next,
	}
}

// Audit извлекает субъект из контекста и проверяет наличие ролей Writer или Admin перед передачей вызова на запись аудит-лога данных.
func (da *DataAudit) Audit(ctx context.Context, data *dto.DataAuditDTO) error {
	// subject
	subj, err := da.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return errs.NewBllForbiddenError("DataAuditFacadeImpl.Audit", "retrieve subject", err)
	}
	// rbac
	if !subj.HasRole(domain.RoleWriter) && !subj.HasRole(domain.RoleAdmin) {
		return errs.NewBllForbiddenError("DataAuditFacadeImpl.Audit", "subject is not audit-writer", nil)
	}

	return da.next.Audit(ctx, data)
}

// ListByPeriod извлекает субъект из контекста и проверяет наличие ролей Reader или Admin перед чтением логов данных за период.
func (da *DataAudit) ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.DataAuditDTO, error) {
	// subject
	subj, err := da.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, errs.NewBllForbiddenError("DataAuditFacadeImpl.ListByPeriod", "retrieve subject", err)
	}
	// rbac
	if !subj.HasRole(domain.RoleReader) && !subj.HasRole(domain.RoleAdmin) {
		return nil, errs.NewBllForbiddenError("DataAuditFacadeImpl.ListByPeriod", "subject is not audit-reader", nil)
	}

	return da.next.ListByPeriod(ctx, auditPeriod)
}

// ListByInstance извлекает субъект из контекста и проверяет наличие ролей Reader или Admin перед поиском логов конкретного экземпляра сущности.
func (da *DataAudit) ListByInstance(ctx context.Context, auditInstance *dto.AuditInstanceDTO) ([]*dto.DataAuditDTO, error) {
	// subject
	subj, err := da.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, errs.NewBllForbiddenError("DataAuditFacadeImpl.ListByInstance", "retrieve subject", err)
	}
	// rbac
	if !subj.HasRole(domain.RoleReader) && !subj.HasRole(domain.RoleAdmin) {
		return nil, errs.NewBllForbiddenError("DataAuditFacadeImpl.ListByInstance", "subject is not audit-reader", nil)
	}

	return da.next.ListByInstance(ctx, auditInstance)
}
