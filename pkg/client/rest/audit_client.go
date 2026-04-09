package rest

import (
	"context"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/tiny-audit-service/pkg/api/http/audit/v1/client/audit"
	"github.com/ElfAstAhe/tiny-auth-service/pkg/transport/auth"
)

type AuditAction[D any] func(ctx context.Context, workerIndex int, data D, token string) error

type AuditClient[A any] interface {
	Start(ctx context.Context) error
	Stop(stopTimeout time.Duration) error
	Audit(data A) error
	TotalLost() int32
}

type AuditClientConfig struct {
	host       string
	scheme     string
	basePath   string
	timeout    time.Duration
	poolConfig *worker.BasePoolConfig
}

func NewAuditClientConfig(
	baseURL string,
	timeout time.Duration,
	poolConfig *worker.BasePoolConfig,
) (*AuditClientConfig, error) {
	res := &AuditClientConfig{
		timeout:    timeout,
		poolConfig: poolConfig,
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, errs.NewCommonError("baseURL parse failed", err)
	}
	res.host = u.Host
	res.scheme = u.Scheme
	res.basePath = u.Path

	return res, nil
}

type BaseAuditClient[D any] struct {
	client        audit.ClientService
	pool          worker.Pool[D]
	conf          *AuditClientConfig
	totalLost     *atomic.Int32
	tokenProvider auth.TokenProvider
	auditAction   AuditAction[D]
}

func NewBaseAuditClient[D any](
	name string,
	conf *AuditClientConfig,
	auditAction AuditAction[D],
	tokenProvider auth.TokenProvider,
	log logger.Logger,
) *BaseAuditClient[D] {
	res := &BaseAuditClient[D]{
		conf:          conf,
		client:        audit.NewClientWithBearerToken(conf.host, conf.basePath, conf.scheme, ""),
		totalLost:     new(atomic.Int32),
		tokenProvider: tokenProvider,
		auditAction:   auditAction,
	}
	res.totalLost.Store(0)
	res.pool = worker.NewBasePool(
		name,
		conf.poolConfig,
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
	// token acquire
	token, err := bac.tokenProvider.GetAccessToken()
	if err != nil {
		bac.totalLost.Add(1)

		return errs.NewCommonError("acquire jwt token failed", err)
	}
	// request
	err = bac.auditAction(ctx, workerIndex, data, token)
	if err != nil {
		return errs.NewCommonError("audit failed", err)
	}

	return nil
}

func (bac *BaseAuditClient[D]) GetConfig() *AuditClientConfig {
	return bac.conf
}
