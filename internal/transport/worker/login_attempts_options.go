package worker

import (
	"strings"
	"time"

	"github.com/Azure/go-amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/logger"
	libamqp "github.com/ElfAstAhe/go-service-template/pkg/transport/amqp"
	"github.com/ElfAstAhe/go-service-template/pkg/transport/worker"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
)

const (
	defaultBatchSize        int           = 50
	defaultBatchReadTimeout time.Duration = 2 * time.Second
)

type LoginAttemptsOption func(*LoginAttemptsOptions)

type LoginAttemptsOptions struct {
	Name               string
	DispatcherOpts     *worker.BaseSchedulerDispatcherConfig
	ReceiveOpts        *amqp.ReceiveOptions
	Receiver           libamqp.Receiver[*amqp.ReceiveOptions]
	AuthAuditUC        usecase.AuthAuditUseCase
	BatchSize          int
	BatchReadTimeout   time.Duration
	AcknowledgeTimeout time.Duration
	Logger             logger.Logger
}

func NewLoginAttemptsOptions() *LoginAttemptsOptions {
	return &LoginAttemptsOptions{
		BatchSize:        defaultBatchSize,
		BatchReadTimeout: defaultBatchReadTimeout,
	}
}

func (lao LoginAttemptsOptions) Validate() error {
	if strings.TrimSpace(lao.Name) == "" {
		return errs.NewCommonError("name not applied (empty)", nil)
	}
	if utils.IsNil(lao.Logger) {
		return errs.NewCommonError("logger not provided (nil)", nil)
	}
	if utils.IsNil(lao.DispatcherOpts) {
		return errs.NewCommonError("dispatcher options not provided (nil)", nil)
	}
	if utils.IsNil(lao.Receiver) {
		return errs.NewCommonError("receiver not provided (nil)", nil)
	}
	if utils.IsNil(lao.AuthAuditUC) {
		return errs.NewCommonError("auth audit use case not provided (nil)", nil)
	}
	if !(lao.BatchSize > 0) {
		return errs.NewCommonError("batch size must be greater than 0", nil)
	}
	if !(lao.BatchReadTimeout > 0) {
		return errs.NewCommonError("batch read timeout must be greater than 0", nil)
	}
	if !(lao.AcknowledgeTimeout > 0) {
		return errs.NewCommonError("acknowledge timeout must be greater than 0", nil)
	}

	return nil
}

func WithLAOName(name string) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.Name = name
	}
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

func WithLAOReceiver(receiver libamqp.Receiver[*amqp.ReceiveOptions]) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.Receiver = receiver
	}
}

func WithLAOAuthAuditUseCase(authAuditUC usecase.AuthAuditUseCase) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.AuthAuditUC = authAuditUC
	}
}

func WithLAOBatchSize(batchSize int) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.BatchSize = batchSize
	}
}

func WithLAOBatchReadTimeout(timeout time.Duration) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.BatchReadTimeout = timeout
	}
}

func WithLAOAcknowledgeTimeout(timeout time.Duration) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.AcknowledgeTimeout = timeout
	}
}

func WithLAOLogger(logger logger.Logger) LoginAttemptsOption {
	return func(opts *LoginAttemptsOptions) {
		opts.Logger = logger
	}
}
