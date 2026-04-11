package client

import (
	"context"
	"time"
)

type AuditAction[D any] func(ctx context.Context, workerIndex int, data D, token string) error

type AuditClient[A any] interface {
	Start(ctx context.Context) error
	Stop(stopTimeout time.Duration) error
	Audit(data A) error
	TotalLost() int32
}
