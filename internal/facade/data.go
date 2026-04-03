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

type DataAuditFacade interface {
	Audit(ctx context.Context, data *dto.DataAuditDTO) error
	ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.DataAuditDTO, error)
	ListByInstance(ctx context.Context, auditInstance *dto.AuditInstanceDTO) ([]*dto.DataAuditDTO, error)
}

type DataAuditFacadeImpl struct {
	authHelper           auth.Helper
	dataAuditUC          usecase.DataAuditUseCase
	dataListByPeriodUC   usecase.DataListByPeriodUseCase
	dataListByInstanceUC usecase.DataListByInstanceUseCase
}

var _ DataAuditFacade = (*DataAuditFacadeImpl)(nil)

func NewDataAuditFacade(
	authHelper auth.Helper,
	dataAuditUC usecase.DataAuditUseCase,
	dataListByPeriod usecase.DataListByPeriodUseCase,
	dataListByInstance usecase.DataListByInstanceUseCase,
) *DataAuditFacadeImpl {
	return &DataAuditFacadeImpl{
		authHelper:           authHelper,
		dataAuditUC:          dataAuditUC,
		dataListByPeriodUC:   dataListByPeriod,
		dataListByInstanceUC: dataListByInstance,
	}
}

func (daf *DataAuditFacadeImpl) Audit(ctx context.Context, data *dto.DataAuditDTO) error {
	// subject
	subj, err := daf.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return domerrs.NewBllForbiddenError("DataAuditFacadeImpl.Audit", "retrieve subject", err)
	}
	// rbac
	if !(subj.HasRole(domain.RoleWriter) || subj.HasRole(domain.RoleAdmin)) {
		return domerrs.NewBllForbiddenError("DataAuditFacadeImpl.Audit", "subject is not audit-writer", nil)
	}
	// validate
	if data == nil {
		return errs.NewInvalidArgumentError("data", "data is nil")
	}

	// logic
	err = daf.dataAuditUC.Audit(ctx, mapper.MapDataAuditDTOToModel(data))
	if err != nil {
		return domerrs.NewBllError("DataAuditFacadeImpl.Audit", "write audit data", err)
	}

	return nil
}

func (daf *DataAuditFacadeImpl) ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.DataAuditDTO, error) {
	// subject
	subj, err := daf.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, domerrs.NewBllForbiddenError("DataAuditFacadeImpl.ListByPeriod", "retrieve subject", err)
	}
	// rbac
	if !(subj.HasRole(domain.RoleReader) || subj.HasRole(domain.RoleAdmin)) {
		return nil, domerrs.NewBllForbiddenError("DataAuditFacadeImpl.ListByPeriod", "subject is not audit-reader", nil)
	}
	// validate
	// pass to bll

	// logic
	res, err := daf.dataListByPeriodUC.List(ctx, auditPeriod.From, auditPeriod.Till, auditPeriod.Limit, auditPeriod.Offset)
	if err != nil {
		return nil, err
	}

	return mapper.MapDataAuditModelsToDTOs(res), nil
}

func (daf *DataAuditFacadeImpl) ListByInstance(ctx context.Context, auditInstance *dto.AuditInstanceDTO) ([]*dto.DataAuditDTO, error) {
	// subject
	subj, err := daf.authHelper.SubjectFromContext(ctx)
	if err != nil {
		return nil, domerrs.NewBllForbiddenError("DataAuditFacadeImpl.ListByInstance", "retrieve subject", err)
	}
	// rbac
	if !(subj.HasRole(domain.RoleReader) || subj.HasRole(domain.RoleAdmin)) {
		return nil, domerrs.NewBllForbiddenError("DataAuditFacadeImpl.ListByInstance", "subject is not audit-reader", nil)
	}
	// validate
	// pass to bll

	// logic
	res, err := daf.dataListByInstanceUC.List(ctx, auditInstance.TypeName, auditInstance.InstanceID, auditInstance.Limit, auditInstance.Offset)
	if err != nil {
		return nil, err
	}

	return mapper.MapDataAuditModelsToDTOs(res), nil
}
