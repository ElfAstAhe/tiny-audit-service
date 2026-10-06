package rbac

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
)

type DataAudit struct {
	authHelper auth.Helper
	next       facade.DataAudit
}

var _ facade.DataAudit = (*DataAudit)(nil)

func NewDataAudit(helper auth.Helper, next facade.DataAudit) *DataAudit {
	return &DataAudit{
		authHelper: helper,
		next:       next,
	}
}

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

func (da *DataAudit) ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.DataAuditDTO, error) {
	// subject
	subj, err := da.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, errs.NewBllForbiddenError("DataAuditFacadeImpl.ListByPeriod", "retrieve subject", err)
	}
	// rbac
	if !subj.HasRole(domain.RoleReader) && subj.HasRole(domain.RoleAdmin) {
		return nil, errs.NewBllForbiddenError("DataAuditFacadeImpl.ListByPeriod", "subject is not audit-reader", nil)
	}

	return da.next.ListByPeriod(ctx, auditPeriod)
}

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
