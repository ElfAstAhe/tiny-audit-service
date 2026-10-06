package rbac

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

type AuthAudit struct {
	authHelper auth.Helper
	next       facade.AuthAudit
}

var _ facade.AuthAudit = (*AuthAudit)(nil)

func NewAuthAudit(
	helper auth.Helper,
	next facade.AuthAudit,
) *AuthAudit {
	return &AuthAudit{
		authHelper: helper,
		next:       next,
	}
}

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
