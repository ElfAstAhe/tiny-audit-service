package facade

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/mapper"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

type DataAudit interface {
	Audit(ctx context.Context, data *dto.DataAuditDTO) error
	ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.DataAuditDTO, error)
	ListByInstance(ctx context.Context, auditInstance *dto.AuditInstanceDTO) ([]*dto.DataAuditDTO, error)
}

type DataAuditImpl struct {
	dataAuditUC          usecase.DataAuditUseCase
	dataListByPeriodUC   usecase.DataListByPeriodUseCase
	dataListByInstanceUC usecase.DataListByInstanceUseCase
}

var _ DataAudit = (*DataAuditImpl)(nil)

func NewDataAudit(
	dataAuditUC usecase.DataAuditUseCase,
	dataListByPeriod usecase.DataListByPeriodUseCase,
	dataListByInstance usecase.DataListByInstanceUseCase,
) *DataAuditImpl {
	return &DataAuditImpl{
		dataAuditUC:          dataAuditUC,
		dataListByPeriodUC:   dataListByPeriod,
		dataListByInstanceUC: dataListByInstance,
	}
}

func (daf *DataAuditImpl) Audit(ctx context.Context, data *dto.DataAuditDTO) error {
	// validate
	if utils.IsNil(data) {
		return errs.NewInvalidArgumentError("data", "data is nil")
	}

	// logic
	err := daf.dataAuditUC.Audit(ctx, mapper.MapDataAuditDTOToModel(data))
	if err != nil {
		return errs.NewBllError("DataAuditFacadeImpl.Audit", "write audit data", err)
	}

	return nil
}

func (daf *DataAuditImpl) ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.DataAuditDTO, error) {
	// validate
	// pass to bll

	// logic
	res, err := daf.dataListByPeriodUC.List(ctx, auditPeriod.From, auditPeriod.Till, auditPeriod.Limit, auditPeriod.Offset)
	if err != nil {
		return nil, err
	}

	return mapper.MapDataAuditModelsToDTOs(res), nil
}

func (daf *DataAuditImpl) ListByInstance(ctx context.Context, auditInstance *dto.AuditInstanceDTO) ([]*dto.DataAuditDTO, error) {
	// validate
	// pass to bll

	// logic
	res, err := daf.dataListByInstanceUC.List(ctx, auditInstance.TypeName, auditInstance.InstanceID, auditInstance.Limit, auditInstance.Offset)
	if err != nil {
		return nil, err
	}

	return mapper.MapDataAuditModelsToDTOs(res), nil
}
