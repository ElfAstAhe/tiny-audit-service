package rest

import (
	"net/url"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
)

type AuditClientConfig struct {
	Host        string
	Scheme      string
	BasePath    string
	ReadTimeout time.Duration
	poolConf    *worker.BasePoolConfig
}

func NewAuditClientConfig(
	baseURL string,
	readTimeout time.Duration,
	poolConf *worker.BasePoolConfig,
) (*AuditClientConfig, error) {
	res := &AuditClientConfig{
		ReadTimeout: readTimeout,
		poolConf:    poolConf,
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, errs.NewCommonError("baseURL parse failed", err)
	}
	res.Host = u.Host
	res.Scheme = u.Scheme
	res.BasePath = u.Path

	return res, nil
}
