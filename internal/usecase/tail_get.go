package usecase

import (
	"context"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/domain"
)

// TailGetUseCase defines the generic application orchestration contract to retrieve historical tracking
// or log slice tails boundaries matching specific temporal time criteria.
type TailGetUseCase[ID comparable] interface {
	GetTail(ctx context.Context, tail time.Time) ([]ID, error)
}

// TailGetInteractor orchestrates decoupled lookup requests, managing state invariant checkpoints prior to querying underlying generic data stores.
type TailGetInteractor[ID comparable] struct {
	tailRepo domain.TailRepository[ID] // Downstream generic persistence repository layer execution implementation path
}

// Compile-time interface compliance verification for canonical string identifiers
var _ TailGetUseCase[string] = (*TailGetInteractor[string])(nil)

// NewTailGetUseCase acts as a generic factory constructor initializing data store persistence adapters dependencies.
func NewTailGetUseCase[ID comparable](tailRepo domain.TailRepository[ID]) *TailGetInteractor[ID] {
	return &TailGetInteractor[ID]{
		tailRepo: tailRepo,
	}
}

// GetTail evaluates time boundary invariants before delegating the diagnostic lookup command down the persistence stream loop.
func (tl *TailGetInteractor[ID]) GetTail(ctx context.Context, tail time.Time) ([]ID, error) {
	if tail.IsZero() {
		return nil, errs.NewInvalidArgumentError("tail", "must not be zero")
	}

	return tl.tailRepo.GetTail(ctx, tail)
}
