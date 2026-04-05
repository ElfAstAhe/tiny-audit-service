package domain

import (
	"context"
	"time"
)

type TailRepository[ID comparable] interface {
	GetTail(ctx context.Context, tailDate time.Time) ([]ID, error)
	Delete(ctx context.Context, id ID) error
}
