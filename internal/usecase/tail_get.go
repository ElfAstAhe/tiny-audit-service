package usecase

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

type TailGetUseCase[ID comparable] interface {
	GetTail(ctx context.Context, tail time.Time) ([]ID, error)
}

type TailGetInteractor[ID comparable] struct {
	tailRepo domain.TailRepository[ID]
}

func NewTailGetUseCase[ID comparable](tailRepo domain.TailRepository[ID]) *TailGetInteractor[ID] {
	return &TailGetInteractor[ID]{
		tailRepo: tailRepo,
	}
}

func (tl *TailGetInteractor[ID]) GetTail(ctx context.Context, tail time.Time) ([]ID, error) {
	if tail.IsZero() {
		return nil, errs.NewInvalidArgumentError("tail", "must not be zero")
	}

	return tl.tailRepo.GetTail(ctx, tail)
}
