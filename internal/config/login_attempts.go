package config

import (
	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/go-service-template/pkg/utils"
)

type LoginAttemptsConfig struct {
	ReceiverConf *config.AMQPReceiverConfig `mapstructure:"receiver_conf" json:"receiver_conf" yaml:"receiver_conf"`
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
