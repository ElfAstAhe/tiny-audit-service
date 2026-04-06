package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	domerrs "github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type AuthListByPeriodUseCase interface {
	List(ctx context.Context, from, till time.Time, limit int, offset int) ([]*domain.AuthAudit, error)
}

type AuthListByPeriodInteractor struct {
	authAuditRepo domain.AuthAuditRepository
}

var _ AuthListByPeriodUseCase = (*AuthListByPeriodInteractor)(nil)

func NewAuthListByPeriodUseCase(authAuditRepo domain.AuthAuditRepository) *AuthListByPeriodInteractor {
	return &AuthListByPeriodInteractor{
		authAuditRepo: authAuditRepo,
	}
}

func (alp *AuthListByPeriodInteractor) List(ctx context.Context, from, till time.Time, limit int, offset int) ([]*domain.AuthAudit, error) {
	if err := alp.validate(from, till, limit, offset); err != nil {
		return nil, domerrs.NewBllValidateError("AuthListByPeriodInteractor.List", "validate income failed", err)
	}

	res, err := alp.authAuditRepo.ListByPeriod(ctx, from, till, limit, offset)
	if err != nil {
		return nil, domerrs.NewBllError("AuthListByPeriodInteractor.List", fmt.Sprintf("list auth audit data by period from [%s] till [%s] with limit [%v] and offset [%v] failed", from.Format(time.DateTime), till.Format(time.DateTime), limit, offset), err)
	}

	return res, nil
}

func (alp *AuthListByPeriodInteractor) validate(from, till time.Time, limit int, offset int) error {
	if from.IsZero() {
		return errs.NewInvalidArgumentError("from", "field is required")
	}
	if till.IsZero() {
		return errs.NewInvalidArgumentError("till", "field is required")
	}
	if !(limit > 0) {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if !(offset >= 0) {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}
