package usecase

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	domerrs "github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type AuthListByUsernameUseCase interface {
	List(ctx context.Context, username string, limit, offset int) ([]*domain.AuthAudit, error)
}

type AuthListByUsernameInteractor struct {
	authAuditRepo domain.AuthAuditRepository
}

var _ AuthListByUsernameUseCase = (*AuthListByUsernameInteractor)(nil)

func NewAuthListByUsernameUseCase(authAuditRepo domain.AuthAuditRepository) *AuthListByUsernameInteractor {
	return &AuthListByUsernameInteractor{
		authAuditRepo: authAuditRepo,
	}
}

func (alu *AuthListByUsernameInteractor) List(ctx context.Context, username string, limit, offset int) ([]*domain.AuthAudit, error) {
	if err := alu.valudate(username, limit, offset); err != nil {
		return nil, domerrs.NewBllValidateError("AuthListByUsernameInteractor.List", "validate income failed", err)
	}

	res, err := alu.authAuditRepo.ListByUsername(ctx, username, limit, offset)
	if err != nil {
		return nil, domerrs.NewBllError("AuthListByUsernameInteractor.List", fmt.Sprintf("list auth audit data for user [%s] with limit [%v] and offset [%v] failed", username, limit, offset), err)
	}

	return res, nil
}

func (alu *AuthListByUsernameInteractor) valudate(username string, limit, offset int) error {
	if username == "" {
		return errs.NewInvalidArgumentError("username", "field is required")
	}
	if !(limit > 0) {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if !(offset >= 0) {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}
