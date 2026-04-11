package client

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/client/dto"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
)

type BaseAuditClient[D any] struct {
	pool          worker.Pool[D]
	conf          *worker.BasePoolConfig
	totalLost     *atomic.Int32
	tokenProvider auth.TokenProvider
	auditAction   AuditAction[D]
	log           logger.Logger
}

var _ AuditClient[*dto.AuthAuditDTO] = (*BaseAuditClient[*dto.AuthAuditDTO])(nil)
var _ AuditClient[*dto.DataAuditDTO] = (*BaseAuditClient[*dto.DataAuditDTO])(nil)

func NewBaseAuditClient[D any](
	name string,
	conf *worker.BasePoolConfig,
	auditAction AuditAction[D],
	tokenProvider auth.TokenProvider,
	log logger.Logger,
) *BaseAuditClient[D] {
	res := &BaseAuditClient[D]{
		conf:          conf,
		totalLost:     new(atomic.Int32),
		tokenProvider: tokenProvider,
		auditAction:   auditAction,
		log:           log.GetLogger("BaseAuditClient"),
	}
	res.totalLost.Store(0)
	res.pool = worker.NewBasePool(
		name,
		worker.NewBasePoolConfig(
			conf.WorkerCount,
			conf.DataCapacity,
			conf.CompleteProcess,
		),
		res.jobHandler,
		log,
	)

	return res
}

func (bac *BaseAuditClient[D]) Start(ctx context.Context) error {
	return bac.pool.Start(ctx)
}

func (bac *BaseAuditClient[D]) Stop(stopTimeout time.Duration) error {
	return bac.pool.Stop(stopTimeout)
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

func (bac *BaseAuditClient[D]) GetConfig() *worker.BasePoolConfig {
	return bac.conf
}

func (bac *BaseAuditClient[D]) GetLogger() logger.Logger {
	return bac.log
}

func (bac *BaseAuditClient[D]) IncLostCounter() {
	bac.totalLost.Add(1)
}
