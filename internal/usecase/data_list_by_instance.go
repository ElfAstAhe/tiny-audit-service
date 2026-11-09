package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

// DataListByInstanceUseCase определяет контракт бизнес-логики (Use Case)
// для получения записей аудита изменений по конкретному типу и идентификатору экземпляра сущности.
type DataListByInstanceUseCase interface {
	List(ctx context.Context, typeName, instanceID string, limit, offset int) ([]*domain.DataAudit, error)
}

// DataListByInstanceInteractor реализует интерфейс DataListByInstanceUseCase, обеспечивая
// чтение записей аудита изменений данных для выбранного инстанса из репозитория.
type DataListByInstanceInteractor struct {
	dataAuditRepo domain.DataAuditRepository
}

var _ DataListByInstanceUseCase = (*DataListByInstanceInteractor)(nil)

// NewDataListByInstance создает и инициализирует новый экземпляр интерактора DataListByInstanceInteractor.
func NewDataListByInstance(dataAuditRepo domain.DataAuditRepository) *DataListByInstanceInteractor {
	return &DataListByInstanceInteractor{
		dataAuditRepo: dataAuditRepo,
	}
}

// List выполняет валидацию входных параметров типа, идентификатора инстанса и пагинации,
// после чего запрашивает список записей аудита изменений данных из нижележащего репозитория.
func (dli *DataListByInstanceInteractor) List(ctx context.Context, typeName, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	if err := dli.validate(typeName, instanceID, limit, offset); err != nil {
		return nil, errs.NewBllValidateError("DataListByInstanceInteractor.List", "validate income failed", err)
	}

	res, err := dli.dataAuditRepo.ListByInstance(ctx, typeName, instanceID, limit, offset)
	if err != nil {
		return nil, errs.NewBllError("DataListByInstanceInteractor.List", fmt.Sprintf("data audit list for type name [%s] and instance id [%s] with limit [%v] and offset [%v] failed", typeName, instanceID, limit, offset), err)
	}

	return res, nil
}

func (dli *DataListByInstanceInteractor) validate(typeName, instanceID string, limit, offset int) error {
	if strings.TrimSpace(typeName) == "" {
		return errs.NewInvalidArgumentError("typeName", "typeName is required")
	}
	if instanceID == "" {
		return errs.NewInvalidArgumentError("instanceID", "instanceID is required")
	}
	if limit <= 0 {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if offset < 0 {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}
