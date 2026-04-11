package grpc

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
)

type RawGRPCClientConfig struct {
	Secure           bool
	ConnTimeout      time.Duration
	KATime           time.Duration
	KATimeOut        time.Duration
	KAPermitWOStream bool
}

func NewRawGRPCClientConfig(
	secure bool,
	connTimeout time.Duration,
	KATime time.Duration,
	KATimeOut time.Duration,
	KAPermitWOStream bool,
) *RawGRPCClientConfig {
	return &RawGRPCClientConfig{
		Secure:           secure,
		ConnTimeout:      connTimeout,
		KATime:           KATime,
		KATimeOut:        KATimeOut,
		KAPermitWOStream: KAPermitWOStream,
	}
}

type AuditClientConfig struct {
	Target   string
	GRPCConf *RawGRPCClientConfig
	poolConf *worker.BasePoolConfig
}

func NewAuditAuditClientConfig(
	grpcConf *RawGRPCClientConfig,
	poolConf *worker.BasePoolConfig,
) *AuditClientConfig {
	return &AuditClientConfig{
		GRPCConf: grpcConf,
		poolConf: poolConf,
	}
}
