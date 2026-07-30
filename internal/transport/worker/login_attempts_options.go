package worker

import (
	"time"

	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
)

const (
	defaultBatchSize        int           = 50
	defaultBatchReadTimeout time.Duration = 2 * time.Second
)

type LoginAttemptsOption func(*LoginAttemptsOptions)

type LoginAttemptsOptions struct {
	DispatcherOpts   *worker.BaseSchedulerDispatcherConfig
	ReceiveOpts      *amqp.ReceiveOptions
	BatchSize        int
	BatchReadTimeout time.Duration
	logger           logger.Logger
}

func NewLoginAttemptsOptions() *LoginAttemptsOptions {
	return &LoginAttemptsOptions{
		BatchSize:        defaultBatchSize,
		BatchReadTimeout: defaultBatchReadTimeout,
	}
}

func (lao LoginAttemptsOptions) Validate() error {
	if utils.IsNil(lao.logger) {
		return errs.NewCommonError("logger not provided (nil)", nil)
	}
	if utils.IsNil(lao.DispatcherOpts) {
		return errs.NewCommonError("dispatcher options not provided (nil)", nil)
	}
	if !(lao.BatchSize > 0) {
		return errs.NewCommonError("batch size must be greater than 0", nil)
	}
	if !(lao.BatchReadTimeout > 0) {
		return errs.NewCommonError("batch read timeout must be greater than 0", nil)
	}

	return nil
}

func WithLAODispatcherOpts(dispatcherOpts *worker.BaseSchedulerDispatcherConfig) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.DispatcherOpts = dispatcherOpts
	}
}

func WithLAOReceiveOpts(receiveOpts *amqp.ReceiveOptions) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.ReceiveOpts = receiveOpts
	}
}

func WithLAOBatchSize(batchSize int) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.BatchSize = batchSize
	}
}

func WithLAOBatchReadTimeout(batchReadTimeout time.Duration) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.BatchReadTimeout = batchReadTimeout
	}
}

func WithLAOLogger(logger logger.Logger) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.logger = logger
	}
}
