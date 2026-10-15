package usecase

import (
	"context"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type DataAuditUseCase interface {
	Audit(ctx context.Context, data *domain.DataAudit) error
}

type DataAuditInteractor struct {
	uw       libdom.UnitOfWork
	dataRepo domain.DataAuditRepository
}

var _ DataAuditUseCase = (*DataAuditInteractor)(nil)

func NewDataAuditUseCase(
	uw libdom.UnitOfWork,
	dataRepo domain.DataAuditRepository,
) *DataAuditInteractor {
	return &DataAuditInteractor{
		uw:       uw,
		dataRepo: dataRepo,
	}
}

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
