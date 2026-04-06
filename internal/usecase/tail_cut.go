package usecase

import (
	"context"

	usecase "github.com/ElfAstAhe/go-service-template/pkg/db"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain/errs"
)

type TailCutUseCase[ID comparable] interface {
	Cut(ctx context.Context, id ID) error
}

type TailCutInteractor[ID comparable] struct {
	tm       usecase.TransactionManager
	tailRepo domain.TailRepository[ID]
}

var _ TailCutUseCase[string] = (*TailCutInteractor[string])(nil)

func NewTailCutUseCase[ID comparable](tailRepo domain.TailRepository[ID]) *TailCutInteractor[ID] {
	return &TailCutInteractor[ID]{
		tailRepo: tailRepo,
	}
}

func (tc *TailCutInteractor[ID]) Cut(ctx context.Context, id ID) error {
	err := tc.tm.WithinTransaction(ctx, nil, func(ctx context.Context) error {
		txErr := tc.tailRepo.Delete(ctx, id)
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
