package config

import (
	"fmt"
	"slices"
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
)

type LoginAttemptsConfig struct {
	ReceiverKind       string                      `mapstructure:"receiver_kind" json:"receiver_kind,omitempty" yaml:"receiver_kind,omitempty"`
	StartInterval      time.Duration               `mapstructure:"worker_start_interval" json:"worker_start_interval,omitempty" yaml:"worker_start_interval,omitempty"`
	ScheduleInterval   time.Duration               `mapstructure:"worker_schedule_interval" json:"worker_schedule_interval,omitempty" yaml:"worker_schedule_interval,omitempty"`
	WorkerCount        int                         `mapstructure:"worker_worker_count" json:"worker_worker_count,omitempty" yaml:"worker_worker_count,omitempty"`
	DataCapacity       int                         `mapstructure:"worker_data_capacity" json:"worker_data_capacity,omitempty" yaml:"worker_data_capacity,omitempty"`
	CompleteProcessing bool                        `mapstructure:"worker_complete_processing" json:"worker_complete_processing,omitempty" yaml:"worker_complete_processing,omitempty"`
	ShutdownTimeout    time.Duration               `mapstructure:"worker_shutdown_timeout" json:"worker_shutdown_timeout,omitempty" yaml:"worker_shutdown_timeout,omitempty"`
	BatchSize          int                         `mapstructure:"worker_batch_size" json:"worker_batch_size,omitempty" yaml:"worker_batch_size,omitempty"`
	BatchReadTimeout   time.Duration               `mapstructure:"worker_batch_read_timeout" json:"worker_batch_read_timeout,omitempty" yaml:"worker_batch_read_timeout,omitempty"`
	AcknowledgeTimeout time.Duration               `mapstructure:"worker_ack_timeout" json:"worker_ack_timeout,omitempty" yaml:"worker_ack_timeout,omitempty"`
	AMQPConfig         *config.AMQPReceiverConfig  `mapstructure:"amqp_config" json:"amqp_config,omitempty" yaml:"amqp_config,omitempty"`
	KafkaConfig        *config.KafkaReceiverConfig `mapstructure:"kafka_config" json:"kafka_config,omitempty" yaml:"kafka_config,omitempty"`
}

func NewLoginAttemptsConfig(
	receiverKind string,
	startInterval time.Duration,
	scheduleInterval time.Duration,
	workerCount int,
	dataCapacity int,
	completeProcessing bool,
	shutdownTimeout time.Duration,
	batchSize int,
	batchReadTimeout time.Duration,
	acknowledgeTimeout time.Duration,
	amqpConfig *config.AMQPReceiverConfig,
	kafkaConfig *config.KafkaReceiverConfig,
) *LoginAttemptsConfig {
	return &LoginAttemptsConfig{
		ReceiverKind:       receiverKind,
		StartInterval:      startInterval,
		ScheduleInterval:   scheduleInterval,
		WorkerCount:        workerCount,
		DataCapacity:       dataCapacity,
		CompleteProcessing: completeProcessing,
		ShutdownTimeout:    shutdownTimeout,
		BatchSize:          batchSize,
		BatchReadTimeout:   batchReadTimeout,
		AcknowledgeTimeout: acknowledgeTimeout,
		AMQPConfig:         amqpConfig,
		KafkaConfig:        kafkaConfig,
	}
}

func NewDefaultLoginAttemptsConfig() *LoginAttemptsConfig {
	return NewLoginAttemptsConfig(
		defaultLoginAttemptsReceiverKind,
		defaultLoginAttemptsWorkerStartInterval,
		defaultLoginAttemptsWorkerScheduleInterval,
		defaultLoginAttemptsWorkerWorkerCount,
		defaultLoginAttemptsWorkerDataCapacity,
		defaultLoginAttemptsWorkerCompleteProcessing,
		defaultLoginAttemptsWorkerShutdownTimeout,
		defaultLoginAttemptsWorkerBatchSize,
		defaultLoginAttemptsWorkerBatchReadTimeout,
		defaultLoginAttemptsWorkerAcknowledgeTimeout,
		config.NewDefaultAMQPReceiverConfig(),
		config.NewDefaultKafkaReceiverConfig(),
	)
}

func (lac *LoginAttemptsConfig) Validate() error {
	if !slices.Contains([]string{"amqp", "kafka"}, lac.ReceiverKind) {
		return errs.NewConfigValidateError("login_attempts", "ReceiverKind", fmt.Sprintf("unknown receiver kind: %s", lac.ReceiverKind), nil)
	}
	if !(lac.ScheduleInterval > 0) {
		return errs.NewConfigValidateError("login_attempts", "ScheduleInterval", "must be greater than zero", nil)
	}
	if !(lac.WorkerCount > 0) {
		return errs.NewConfigValidateError("login_attempts", "WorkerCount", "must be greater than zero", nil)
	}
	if !(lac.DataCapacity > 0) {
		return errs.NewConfigValidateError("login_attempts", "DataCapacity", "must be greater than zero", nil)
	}
	if !(lac.ShutdownTimeout > 0) {
		return errs.NewConfigValidateError("login_attempts", "ShutdownTimeout", "must be greater than zero", nil)
	}
	if !(lac.BatchSize > 0) {
		return errs.NewConfigValidateError("login_attempts", "BatchSize", "must be greater than zero", nil)
	}
	if !(lac.BatchReadTimeout > 0) {
		return errs.NewConfigValidateError("login_attempts", "BatchReadTimeout", "must be greater than zero", nil)
	}
	if !(lac.AcknowledgeTimeout > 0) {
		return errs.NewConfigValidateError("login_attempts", "AcknowledgeTimeout", "must be greater than zero", nil)
	}
	if lac.ReceiverKind == "amqp" && utils.IsNil(lac.AMQPConfig) {
		return errs.NewConfigValidateError("login_attempts", "AMQPConfig", "is required", nil)
	}
	if lac.ReceiverKind == "kafka" && utils.IsNil(lac.KafkaConfig) {
		return errs.NewConfigValidateError("login_attempts", "KafkaConfig", "is required", nil)
	}
	if lac.ReceiverKind == "amqp" {
		if err := lac.AMQPConfig.Validate(); err != nil {
			return errs.NewConfigValidateError("login_attempts", "AMQPConfig", "validation failed", err)
		}
	}
	if lac.ReceiverKind == "kafka" {
		if err := lac.KafkaConfig.Validate(); err != nil {
			return errs.NewConfigValidateError("login_attempts", "KafkaConfig", "validation failed", err)
		}
	}

	return nil
}
