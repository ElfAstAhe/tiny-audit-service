package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

// AuthListByUsernameUseCase определяет контракт бизнес-логики (Use Case)
// для получения записей аудита событий аутентификации конкретного пользователя.
type AuthListByUsernameUseCase interface {
	List(ctx context.Context, username string, limit, offset int) ([]*domain.AuthAudit, error)
}

// AuthListByUsernameInteractor реализует интерфейс AuthListByUsernameUseCase, обеспечивая
// чтение записей аудита событий аутентификации для выбранного пользователя из репозитория.
type AuthListByUsernameInteractor struct {
	authAuditRepo domain.AuthAuditRepository
}

var _ AuthListByUsernameUseCase = (*AuthListByUsernameInteractor)(nil)

// NewAuthListByUsername создает и инициализирует новый экземпляр интерактора AuthListByUsernameInteractor.
func NewAuthListByUsername(authAuditRepo domain.AuthAuditRepository) *AuthListByUsernameInteractor {
	return &AuthListByUsernameInteractor{
		authAuditRepo: authAuditRepo,
	}
}

// List выполняет валидацию входных параметров имени пользователя и пагинации,
// после чего запрашивает список записей аудита аутентификации из нижележащего репозитория.
func (alu *AuthListByUsernameInteractor) List(ctx context.Context, username string, limit, offset int) ([]*domain.AuthAudit, error) {
	if err := alu.validate(username, limit, offset); err != nil {
		return nil, errs.NewBllValidateError("AuthListByUsernameInteractor.List", "validate income failed", err)
	}

	res, err := alu.authAuditRepo.ListByUsername(ctx, username, offset, limit)
	if err != nil {
		return nil, errs.NewBllError("AuthListByUsernameInteractor.List", fmt.Sprintf("list auth audit data for user [%s] with limit [%v] and offset [%v] failed", username, limit, offset), err)
	}

	return res, nil
}

func (alu *AuthListByUsernameInteractor) validate(username string, limit, offset int) error {
	if strings.TrimSpace(username) == "" {
		return errs.NewInvalidArgumentError("username", "field is required")
	}
	if limit <= 0 {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if offset < 0 {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}
