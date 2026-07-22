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

	return nil
}
