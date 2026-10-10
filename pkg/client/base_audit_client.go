package client

import (
	"context"
	"fmt"
	"sync/atomic"

	"github.com/ElfAstAhe/go-service-template/pkg/container"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
)

const (
	baseAuditClientNameTemplate string = "base-audit-client-%s"
)

type BaseAuditClient[D any] struct {
	name          string
	pool          worker.Pool[D]
	opts          *BaseAuditClientOptions[D]
	totalLost     *atomic.Int32
	tokenProvider auth.TokenProvider
	auditAction   AuditAction[D]
	log           logger.Logger
}

var _ AuditClient[*dto.AuthAuditDTO] = (*BaseAuditClient[*dto.AuthAuditDTO])(nil)
var _ AuditClient[*dto.DataAuditDTO] = (*BaseAuditClient[*dto.DataAuditDTO])(nil)
var _ container.Runner = (*BaseAuditClient[*dto.DataAuditDTO])(nil)

func NewBaseAuditClient[D any](options ...BaseAuditClientOption[D]) (*BaseAuditClient[D], error) {
	opts := NewBaseAuditClientOptions[D]()
	for _, option := range options {
		option(opts)
	}
	if err := opts.Validate(); err != nil {
		return nil, errs.NewCommonError("base audit client options validation failed", err)
	}
	// instance
	res := &BaseAuditClient[D]{
		name:          fmt.Sprintf(baseAuditClientNameTemplate, opts.Pool.Name),
		opts:          opts,
		totalLost:     new(atomic.Int32),
		tokenProvider: opts.TokenProvider,
		auditAction:   opts.AuditAction,
		log:           opts.Pool.Logger.GetLogger(fmt.Sprintf(baseAuditClientNameTemplate, opts.Pool.Name)),
	}
	// pool
	workerPool, err := worker.NewBasePool[D](
		worker.WithPoolName[D](opts.Pool.Name),
		worker.WithPoolWorkerCount[D](opts.Pool.WorkerCount),
		worker.WithPoolDataCapacity[D](opts.Pool.DataCapacity),
		worker.WithPoolCompleteProcess[D](opts.Pool.CompleteProcess),
		worker.WithPoolStopTimeout[D](opts.Pool.StopTimeout),
		worker.WithPoolJobHandler[D](res.jobHandler),
		worker.WithPoolLogger[D](opts.Pool.Logger),
	)
	if err != nil {
		return nil, errs.NewCommonError("failed to create worker pool", err)
	}
	// setup
	res.pool = workerPool

	return res, nil
}

func (bac *BaseAuditClient[D]) Start(ctx context.Context) error {
	return bac.pool.Start(ctx)
}

func (bac *BaseAuditClient[D]) Stop(ctx context.Context) error {
	return bac.pool.Stop(ctx)
}

func (bac *BaseAuditClient[D]) Audit(data D) error {
	bac.GetLogger().Debugf("Audit start")
	defer bac.GetLogger().Debugf("Audit end")

	res := bac.pool.TryPush(data)
	if !res {
		bac.totalLost.Add(1)
		return errs.NewCommonError("failed to push auth audit data", nil)
	}

	return nil
}

func (bac *BaseAuditClient[D]) TotalLost() int32 {
	return bac.totalLost.Load()
}

func (bac *BaseAuditClient[D]) jobHandler(ctx context.Context, workerIndex int, data D) error {
	bac.GetLogger().Debugf("jobHandler start")
	defer bac.GetLogger().Debugf("jobHandler end")

	// token acquire
	token, err := bac.tokenProvider.GetAccessToken()
	if err != nil {
		bac.totalLost.Add(1)

		return errs.NewCommonError("acquire jwt token failed", err)
	}
	// request
	err = bac.auditAction(ctx, workerIndex, data, token)
	if err != nil {
		// try push for retry audit
		if pushSuccess := bac.pool.TryPush(data); !pushSuccess {
			bac.GetLogger().Warn("failed to enqueue audit data for retry; data has been lost")
			bac.GetLogger().Warnf("lost audit data [%v]", data)
			bac.totalLost.Add(1)
		}

		return errs.NewCommonError("audit failed", err)
	}

	return nil
}

func (bac *BaseAuditClient[D]) GetOpts() *BaseAuditClientOptions[D] {
	return bac.opts
}

func (bac *BaseAuditClient[D]) GetLogger() logger.Logger {
	return bac.log
}

func (bac *BaseAuditClient[D]) IncLostCounter() {
	bac.totalLost.Add(1)
}

func (bac *BaseAuditClient[D]) GetName() string {
	return bac.pool.GetName()
}

func (bac *BaseAuditClient[D]) IsRunning() bool {
	return bac.pool.IsRunning()
}
