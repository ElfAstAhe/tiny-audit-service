package usecase

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type TailListUseCase[ID comparable] interface {
	GetTail(ctx context.Context, tail time.Time) ([]ID, error)
}

type TailListInteractor[ID comparable] struct {
	tailRepo domain.TailRepository[ID]
}

func NewTailListUseCase[ID comparable](tailRepo domain.TailRepository[ID]) *TailListInteractor[ID] {
	return &TailListInteractor[ID]{
		tailRepo: tailRepo,
	}
}

func (tl *TailListInteractor[ID]) GetTail(ctx context.Context, tail time.Time) ([]ID, error) {
	if tail.IsZero() {
		return nil, errs.NewInvalidArgumentError("tail", "must not be zero")
	}

	return tl.tailRepo.GetTail(ctx, tail)
}
