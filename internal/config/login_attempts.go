package config

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
)

type LoginAttemptsConfig struct {
	ReceiverConf       *config.AMQPReceiverConfig `mapstructure:"receiver_conf" json:"receiver_conf" yaml:"receiver_conf"`
	StartInterval      time.Duration              `mapstructure:"start_interval" json:"start_interval,omitempty" yaml:"start_interval,omitempty"`
	ScheduleInterval   time.Duration              `mapstructure:"schedule_interval" json:"schedule_interval,omitempty" yaml:"schedule_interval,omitempty"`
	WorkerCount        int                        `mapstructure:"worker_count" json:"worker_count,omitempty" yaml:"worker_count,omitempty"`
	DataCapacity       int                        `mapstructure:"data_capacity" json:"data_capacity,omitempty" yaml:"data_capacity,omitempty"`
	CompleteProcessing bool                       `mapstructure:"complete_processing" json:"complete_processing,omitempty" yaml:"complete_processing,omitempty"`
	ShutdownTimeout    time.Duration              `mapstructure:"shutdown_timeout" json:"shutdown_timeout,omitempty" yaml:"shutdown_timeout,omitempty"`
	BatchSize          int                        `mapstructure:"batch_size" json:"batch_size,omitempty" yaml:"batch_size"`
	BatchReadTimeout   time.Duration              `mapstructure:"batch_read_timeout" json:"batch_read_timeout,omitempty" yaml:"batch_read_timeout,omitempty"`
	AcknowledgeTimeout time.Duration              `mapstructure:"ack_timeout" json:"ack_timeout,omitempty" yaml:"ack_timeout,omitempty"`
}

func NewLoginAttemptsConfig(
	receiverConf *config.AMQPReceiverConfig,
) *LoginAttemptsConfig {
	return &LoginAttemptsConfig{
		ReceiverConf: receiverConf,
	}
}

func NewDefaultLoginAttemptsConfig() *LoginAttemptsConfig {
	return NewLoginAttemptsConfig(
		config.NewDefaultAMQPReceiverConfig(),
	)
}

func (lac *LoginAttemptsConfig) Validate() error {
	if utils.IsNil(lac.ReceiverConf) {
		return errs.NewConfigValidateError("login_attempts", "ReceiverConf", "absent", nil)
	}
	if err := lac.ReceiverConf.Validate(); err != nil {
		return errs.NewConfigValidateError("login_attempts", "ReceiverConf", "validate failed", err)
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

	return nil
}
