package usecase

import (
	"context"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type AuthAuditUseCase interface {
	Audit(ctx context.Context, data *domain.AuthAudit) error
}

type AuthAuditInteractor struct {
	uw       libdom.UnitOfWork
	authRepo domain.AuthAuditRepository
}

var _ AuthAuditUseCase = (*AuthAuditInteractor)(nil)

func NewAuthAuditUseCase(
	uw libdom.UnitOfWork,
	authRepo domain.AuthAuditRepository,
) *AuthAuditInteractor {
	return &AuthAuditInteractor{
		uw:       uw,
		authRepo: authRepo,
	}
}

func (aai *AuthAuditInteractor) Audit(ctx context.Context, data *domain.AuthAudit) error {
	if err := aai.validate(data); err != nil {
		return errs.NewBllValidateError("AuthAuditInteractor.Audit", "validate failed", err)
	}

	err := aai.uw.Execute(ctx, func(txCtx context.Context) error {
		_, txErr := aai.authRepo.Create(txCtx, data)

		return txErr
	})
	if err != nil {
		return errs.NewBllError("AuthAuditInteractor.Audit", "add auth audit failed", err)
	}

	return nil
}

func (aai *AuthAuditInteractor) validate(data *domain.AuthAudit) error {
	if utils.IsNil(data) {
		return errs.NewInvalidArgumentError("data", "data is nil")
	}

	return nil
}
