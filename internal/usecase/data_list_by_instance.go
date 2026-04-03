package usecase

import (
	"context"
	"fmt"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	domerrs "github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type DataListByInstanceUseCase interface {
	List(ctx context.Context, typeName, instanceID string, limit, offset int) ([]*domain.DataAudit, error)
}

type DataListByInstanceInteractor struct {
	dataAuditRepo domain.DataAuditRepository
}

var _ DataListByInstanceUseCase = (*DataListByInstanceInteractor)(nil)

func NewDataListByInstanceUseCase(dataAuditRepo domain.DataAuditRepository) *DataListByInstanceInteractor {
	return &DataListByInstanceInteractor{
		dataAuditRepo: dataAuditRepo,
	}
}

func (dli *DataListByInstanceInteractor) List(ctx context.Context, typeName, instanceID string, limit, offset int) ([]*domain.DataAudit, error) {
	if err := dli.validate(typeName, instanceID, limit, offset); err != nil {
		return nil, domerrs.NewBllValidateError("DataListByInstanceInteractor.List", "validate income failed", err)
	}

	res, err := dli.dataAuditRepo.ListByInstance(ctx, typeName, instanceID, limit, offset)
	if err != nil {
		return nil, domerrs.NewBllError("DataListByInstanceInteractor.List", fmt.Sprintf("data audit list for type name [%s] and instance id [%s] with limit [%v] and offset [%v] failed", typeName, instanceID, limit, offset), err)
	}

	return res, nil
}

func (dli *DataListByInstanceInteractor) validate(typeName, instanceID string, limit, offset int) error {
	if typeName == "" {
		return errs.NewInvalidArgumentError("typeName", "typeName is required")
	}
	if instanceID == "" {
		return errs.NewInvalidArgumentError("instanceID", "instanceID is required")
	}
	if !(limit > 0) {
		return errs.NewInvalidArgumentError("limit", "limit must be grater than zero")
	}
	if !(offset >= 0) {
		return errs.NewInvalidArgumentError("offset", "offset must be greater or equal than zero")
	}

	return nil
}
