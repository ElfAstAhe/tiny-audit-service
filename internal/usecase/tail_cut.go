package usecase

import (
	"context"

	libdom "github.com/ElfAstAhe/go-service-template/pkg/domain"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type TailCutUseCase[ID comparable] interface {
	Cut(ctx context.Context, id ID) error
}

type TailCutInteractor[ID comparable] struct {
	uw       libdom.UnitOfWork
	tailRepo domain.TailRepository[ID]
}

var _ TailCutUseCase[string] = (*TailCutInteractor[string])(nil)

func NewTailCutUseCase[ID comparable](
	uw libdom.UnitOfWork,
	tailRepo domain.TailRepository[ID],
) *TailCutInteractor[ID] {
	return &TailCutInteractor[ID]{
		uw:       uw,
		tailRepo: tailRepo,
	}
}

func (tc *TailCutInteractor[ID]) Cut(ctx context.Context, id ID) error {
	err := tc.uw.Execute(ctx, func(txCtx context.Context) error {
		txErr := tc.tailRepo.Delete(txCtx, id)
		if txErr != nil {
			return txErr
		}

		return nil
	})
	if err != nil {
		return errs.NewBllError("TailCutInteractor.Cut", "cut tail", err)
	}

	return nil
}
