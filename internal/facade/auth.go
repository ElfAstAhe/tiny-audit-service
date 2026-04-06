package facade

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/auth"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	domerrs "github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/mapper"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

type AuthAuditFacade interface {
	Audit(ctx context.Context, data *dto.AuthAuditDTO) error
	ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.AuthAuditDTO, error)
	ListByUsername(ctx context.Context, auditUser *dto.AuditUserDTO) ([]*dto.AuthAuditDTO, error)
}

type AuthAuditFacadeImpl struct {
	authHelper           auth.Helper
	authAuditUC          usecase.AuthAuditUseCase
	authListByPeriodUC   usecase.AuthListByPeriodUseCase
	authListByUsernameUC usecase.AuthListByUsernameUseCase
}

var _ AuthAuditFacade = (*AuthAuditFacadeImpl)(nil)

func NewAuthAuditFacade(
	authHelper auth.Helper,
	authAuditUC usecase.AuthAuditUseCase,
	authListByPeriodUC usecase.AuthListByPeriodUseCase,
	authListByUsernameUC usecase.AuthListByUsernameUseCase,
) *AuthAuditFacadeImpl {
	return &AuthAuditFacadeImpl{
		authHelper:           authHelper,
		authAuditUC:          authAuditUC,
		authListByPeriodUC:   authListByPeriodUC,
		authListByUsernameUC: authListByUsernameUC,
	}
}

func (aaf *AuthAuditFacadeImpl) Audit(ctx context.Context, data *dto.AuthAuditDTO) error {
	// subject
	subj, err := aaf.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return domerrs.NewBllForbiddenError("AuthAuditFacadeImpl.Audit", "retrieve subject", err)
	}
	// rbac
	if !(subj.HasRole(domain.RoleWriter) || subj.HasRole(domain.RoleAdmin)) {
		return domerrs.NewBllForbiddenError("AuthAuditFacadeImpl.Audit", "subject is not audit-writer", nil)
	}
	// validate
	if data == nil {
		return errs.NewInvalidArgumentError("data", "data is nil")
	}

	// logic
	err = aaf.authAuditUC.Audit(ctx, mapper.MapAuthAuditDTOToModel(data))
	if err != nil {
		return domerrs.NewBllError("AuthAuditFacadeImpl.Audit", "write audit data", err)
	}

	return nil
}

func (aaf *AuthAuditFacadeImpl) ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.AuthAuditDTO, error) {
	// subject
	subj, err := aaf.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, domerrs.NewBllForbiddenError("AuthAuditFacadeImpl.ListByPeriod", "retrieve subject", err)
	}
	// rbac
	if !(subj.HasRole(domain.RoleReader) || subj.HasRole(domain.RoleAdmin)) {
		return nil, domerrs.NewBllForbiddenError("AuthAuditFacadeImpl.ListByPeriod", "subject is not audit-reader", nil)
	}
	// validate
	// pass to bll

	// logic
	res, err := aaf.authListByPeriodUC.List(ctx, auditPeriod.From, auditPeriod.Till, auditPeriod.Limit, auditPeriod.Offset)
	if err != nil {
		return nil, err
	}

	return mapper.MapAuthAuditModelsToDTOs(res), nil
}

func (aaf *AuthAuditFacadeImpl) ListByUsername(ctx context.Context, auditUser *dto.AuditUserDTO) ([]*dto.AuthAuditDTO, error) {
	// subject
	subj, err := aaf.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, domerrs.NewBllForbiddenError("AuthAuditFacadeImpl.ListByUsername", "retrieve subject", err)
	}
	// rbac
	if !(subj.HasRole(domain.RoleReader) || subj.HasRole(domain.RoleAdmin)) {
		return nil, domerrs.NewBllForbiddenError("AuthAuditFacadeImpl.ListByUsername", "subject is not audit-reader", nil)
	}
	// validate
	// pass to bll

	// logic
	res, err := aaf.authListByUsernameUC.List(ctx, auditUser.Username, auditUser.Limit, auditUser.Offset)
	if err != nil {
		return nil, err
	}

	return mapper.MapAuthAuditModelsToDTOs(res), nil
}
