package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

// DataListByPeriodUseCase определяет контракт бизнес-логики (Use Case)
// для получения архивных записей аудита изменений данных за указанный период времени.
type DataListByPeriodUseCase interface {
	List(ctx context.Context, from, till time.Time, limit int, offset int) ([]*domain.DataAudit, error)
}

// DataListByPeriodInteractor реализует интерфейс DataListByPeriodUseCase, обеспечивая
// чтение записей аудита изменений данных за выбранный временной интервал из репозитория.
type DataListByPeriodInteractor struct {
	dataAuditRepo domain.DataAuditRepository
}

var _ DataListByPeriodUseCase = (*DataListByPeriodInteractor)(nil)

// NewDataListByPeriod создает и инициализирует новый экземпляр интерактора DataListByPeriodInteractor.
func NewDataListByPeriod(dataAuditRepo domain.DataAuditRepository) *DataListByPeriodInteractor {
	return &DataListByPeriodInteractor{
		dataAuditRepo: dataAuditRepo,
	}
}

// List выполняет валидацию входных параметров времени и пагинации, после чего запрашивает
// список записей аудита изменений данных из нижележащего репозитория.
func (dlp *DataListByPeriodInteractor) List(ctx context.Context, from, till time.Time, limit int, offset int) ([]*domain.DataAudit, error) {
	if err := dlp.validate(from, till, limit, offset); err != nil {
		return nil, errs.NewBllValidateError("DataListByPeriodInteractor.List", "validate income failed", err)
	}
	res, err := dlp.dataAuditRepo.ListByPeriod(ctx, from, till, limit, offset)
	if err != nil {
		return nil, errs.NewBllError("DataListByPeriodInteractor.List", fmt.Sprintf("list auth audit data by period from [%s] till [%s] with limit [%v] and offset [%v] failed", from.Format(time.DateTime), till.Format(time.DateTime), limit, offset), err)
	}

	return res, nil
}

func (dlp *DataListByPeriodInteractor) validate(from, till time.Time, limit int, offset int) error {
	if from.IsZero() {
		return errs.NewInvalidArgumentError("from", "field is required")
	}
	if till.IsZero() {
		return errs.NewInvalidArgumentError("till", "field is required")
	}
	if limit <= 0 {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if offset < 0 {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}
