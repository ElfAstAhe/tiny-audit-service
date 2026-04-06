package usecase

import (
	"context"

	usecase "github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	domerrs "github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type AuthAuditUseCase interface {
	Audit(ctx context.Context, data *domain.AuthAudit) error
}

type AuthAuditInteractor struct {
	tm       usecase.TransactionManager
	authRepo domain.AuthAuditRepository
}

var _ AuthAuditUseCase = (*AuthAuditInteractor)(nil)

func NewAuthAuditUseCase(tm usecase.TransactionManager, authRepo domain.AuthAuditRepository) *AuthAuditInteractor {
	return &AuthAuditInteractor{
		tm:       tm,
		authRepo: authRepo,
	}
}

func (aai *AuthAuditInteractor) Audit(ctx context.Context, data *domain.AuthAudit) error {
	if err := aai.validate(data); err != nil {
		return domerrs.NewBllValidateError("AuthAuditInteractor.Audit", "validate failed", err)
	}

	err := aai.tm.WithinTransaction(ctx, nil, func(ctx context.Context) error {
		_, txErr := aai.authRepo.Create(ctx, data)

		return txErr
	})
	if err != nil {
		return domerrs.NewBllError("AuthAuditInteractor.Audit", "add auth audit failed", err)
	}

	return nil
}

func (aai *AuthAuditInteractor) validate(data *domain.AuthAudit) error {
	if data == nil {
		return errs.NewInvalidArgumentError("data", "data is nil")
	}

	return nil
}
