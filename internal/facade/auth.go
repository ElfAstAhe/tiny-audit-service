package facade

import (
	"context"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/dto"
	"github.com/ElfAstAhe/tiny-audit-service/internal/facade/mapper"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

// AuthAudit определяет интерфейс фасада для управления и выборки
// логов аудита аутентификации на уровне адаптеров приложения.
type AuthAudit interface {
	Audit(ctx context.Context, data *dto.AuthAuditDTO) error
	ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.AuthAuditDTO, error)
	ListByUsername(ctx context.Context, auditUser *dto.AuditUserDTO) ([]*dto.AuthAuditDTO, error)
}

// AuthAuditFacade реализует интерфейс AuthAudit, агрегируя различные UseCase сценарии
// и обеспечивая маппинг между структурами DTO слоя представления и доменными моделями.
type AuthAuditFacade struct {
	authAuditUC          usecase.AuthAuditUseCase
	authListByPeriodUC   usecase.AuthListByPeriodUseCase
	authListByUsernameUC usecase.AuthListByUsernameUseCase
}

var _ AuthAudit = (*AuthAuditFacade)(nil)

// NewAuthAudit создает и инициализирует новый экземпляр фасада AuthAuditFacade.
func NewAuthAudit(
	authAuditUC usecase.AuthAuditUseCase,
	authListByPeriodUC usecase.AuthListByPeriodUseCase,
	authListByUsernameUC usecase.AuthListByUsernameUseCase,
) *AuthAuditFacade {
	return &AuthAuditFacade{
		authAuditUC:          authAuditUC,
		authListByPeriodUC:   authListByPeriodUC,
		authListByUsernameUC: authListByUsernameUC,
	}
}

// Audit выполняет первичную валидацию DTO, преобразует его в доменную модель и передает в UseCase для регистрации записи аудита.
func (aaf *AuthAuditFacade) Audit(ctx context.Context, data *dto.AuthAuditDTO) error {
	// validate
	if utils.IsNil(data) {
		return errs.NewInvalidArgumentError("data", "data is nil")
	}

	// logic
	err := aaf.authAuditUC.Audit(ctx, mapper.MapAuthAuditDTOToModel(data))
	if err != nil {
		return errs.NewBllError("AuthAuditFacadeImpl.Audit", "write audit data", err)
	}

	return nil
}

// ListByPeriod запрашивает архивные логи аудита аутентификации за указанный период времени и возвращает их в виде DTO-слайса.
func (aaf *AuthAuditFacade) ListByPeriod(ctx context.Context, auditPeriod *dto.AuditPeriodDTO) ([]*dto.AuthAuditDTO, error) {
	// validate
	// pass to bll

	// logic
	res, err := aaf.authListByPeriodUC.List(ctx, auditPeriod.From, auditPeriod.Till, auditPeriod.Limit, auditPeriod.Offset)
	if err != nil {
		return nil, err
	}

	return mapper.MapAuthAuditModelsToDTOs(res), nil
}

// ListByUsername запрашивает логи аудита аутентификации для конкретного пользователя и возвращает их в виде DTO-слайса.
func (aaf *AuthAuditFacade) ListByUsername(ctx context.Context, auditUser *dto.AuditUserDTO) ([]*dto.AuthAuditDTO, error) {
	// validate
	// pass to bll

	// logic
	res, err := aaf.authListByUsernameUC.List(ctx, auditUser.Username, auditUser.Limit, auditUser.Offset)
	if err != nil {
		return nil, err
	}

	return mapper.MapAuthAuditModelsToDTOs(res), nil
}
