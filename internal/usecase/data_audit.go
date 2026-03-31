package usecase

import (
	"context"

	usecade "github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type DataAuditUseCase interface {
	Audit(ctx context.Context, data *domain.DataAudit) error
}

type DataAuditInteractor struct {
	tm       usecade.TransactionManager
	dataRepo domain.DataAuditRepository
}

var _ DataAuditUseCase = (*DataAuditInteractor)(nil)

func NewDataAuditUseCase(
	tm usecade.TransactionManager,
	dataRepo domain.DataAuditRepository,
) *DataAuditInteractor {
	return &DataAuditInteractor{
		tm:       tm,
		dataRepo: dataRepo,
	}
}

func (dai *DataAuditInteractor) Audit(ctx context.Context, data *domain.DataAudit) error {
	if err := dai.validate(data); err != nil {
		return errs.NewBllValidateError("DataAuditInteractor.Audit", "validate failed", err)
	}
	err := dai.tm.WithinTransaction(ctx, nil, func(ctx context.Context) error {
		_, err := dai.dataRepo.Create(ctx, data)

		return err
	})
	if err != nil {
		return errs.NewBllValidateError("DataAuditInteractor.Audit", "add data audit failed", err)
	}

	return nil
}

func (dai *DataAuditInteractor) validate(data *domain.DataAudit) error {
	// ToDo: implement

	return nil
}
