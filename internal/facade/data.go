package facade

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/mapper"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

// DataAudit определяет интерфейс фасада для управления и выборки
// логов аудита изменений данных на уровне адаптеров приложения.
type DataAudit interface {
	Audit(ctx context.Context, data *dto.DataAuditDTO) error
	ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.DataAuditDTO, error)
	ListByInstance(ctx context.Context, auditInstance *dto.AuditInstanceDTO) ([]*dto.DataAuditDTO, error)
}

// DataAuditImpl реализует интерфейс DataAudit, агрегируя различные UseCase сценарии
// и обеспечивая маппинг между структурами DTO слоя представления и доменными моделями аудита данных.
type DataAuditImpl struct {
	dataAuditUC          usecase.DataAuditUseCase
	dataListByPeriodUC   usecase.DataListByPeriodUseCase
	dataListByInstanceUC usecase.DataListByInstanceUseCase
}

var _ DataAudit = (*DataAuditImpl)(nil)

// NewDataAudit создает и инициализирует новый экземпляр фасада DataAuditImpl.
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

// Audit выполняет первичную валидацию DTO, преобразует его в доменную модель и передает в UseCase для регистрации записи аудита данных.
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

// ListByPeriod запрашивает архивные логи аудита изменений данных за указанный период времени и возвращает их в виде DTO-слайса.
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

// ListByInstance запрашивает логи аудита для конкретного экземпляра сущности по её типу и уникальному ID, возвращая результат в виде DTO-слайса.
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
