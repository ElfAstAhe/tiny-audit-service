package usecase

import (
	"context"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

// DataAuditUseCase определяет контракт бизнес-логики (Use Case)
// для регистрации и обработки событий аудита изменений данных.
type DataAuditUseCase interface {
	Audit(ctx context.Context, data *domain.DataAudit) error
}

// DataAuditInteractor реализует интерфейс DataAuditUseCase, обеспечивая
// атомарное сохранение записей аудита данных с использованием паттерна Unit of Work.
type DataAuditInteractor struct {
	uw       libdom.UnitOfWork
	dataRepo domain.DataAuditRepository
}

var _ DataAuditUseCase = (*DataAuditInteractor)(nil)

// NewDataAudit создает и инициализирует новый экземпляр интерактора DataAuditInteractor.
func NewDataAudit(
	uw libdom.UnitOfWork,
	dataRepo domain.DataAuditRepository,
) *DataAuditInteractor {
	return &DataAuditInteractor{
		uw:       uw,
		dataRepo: dataRepo,
	}
}

// Audit выполняет валидацию входящего лога аудита и атомарно сохраняет его в репозиторий
// в рамках изолированного контекста транзакции Unit of Work.
func (dai *DataAuditInteractor) Audit(ctx context.Context, data *domain.DataAudit) error {
	if err := dai.validate(data); err != nil {
		return errs.NewBllValidateError("DataAuditInteractor.Audit", "validate income failed", err)
	}
	err := dai.uw.Execute(ctx, func(txCtx context.Context) error {
		_, err := dai.dataRepo.Create(txCtx, data)

		return err
	})
	if err != nil {
		return errs.NewBllValidateError("DataAuditInteractor.Audit", "add data audit failed", err)
	}

	return nil
}

func (dai *DataAuditInteractor) validate(data *domain.DataAudit) error {
	if utils.IsNil(data) {
		return errs.NewInvalidArgumentError("data", "data is nil")
	}

	return nil
}
