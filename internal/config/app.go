package config

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// AppConfig — метаданные сервиса
type AppConfig struct {
	*config.AppConfig  `mapstructure:",squash"`
	NodeName           string   `mapstructure:"node_name" json:"node_name,omitempty" yaml:"node_name,omitempty"`
	MaxListLimit       int      `mapstructure:"max_list_limit" json:"max_list_limit,omitempty" yaml:"max_list_limit,omitempty"`
	TokenIssuer        string   `mapstructure:"token_issuer" json:"token_issuer,omitempty" yaml:"token_issuer,omitempty"`
	AcceptTokenIssuers []string `mapstructure:"accept_token_issuers" json:"accept_token_issuers,omitempty" yaml:"accept_token_issuers,omitempty"`
	CipherKey          string   `mapstructure:"cipher_key" json:"cipher_key,omitempty" yaml:"cipher_key,omitempty"`
}

func NewAppConfig(
	env config.AppEnv,
	initTimeout time.Duration,
	stopTimeout time.Duration,
	closeTimeout time.Duration,
	maxListLimit int,
	acceptTokenIssuers []string,
	cipherKey string,
) *AppConfig {
	return &AppConfig{
		AppConfig:          config.NewAppConfig(env, initTimeout, stopTimeout, closeTimeout),
		MaxListLimit:       maxListLimit,
		AcceptTokenIssuers: acceptTokenIssuers,
		CipherKey:          cipherKey,
	}
}

func NewDefaultAppConfig() *AppConfig {
	return NewAppConfig(
		defaultAppEnv,
		config.DefaultAppInitTimeout,
		config.DefaultAppStopTimeout,
		config.DefaultAppCloseTimeout,
		defaultMaxListLimit,
		[]string{},
		"",
	)
}

func (ac *AppConfig) Validate() error {
	if ac.Env == "" {
		return errs.NewConfigValidateError("app", "env", "must not be empty", nil)
	}

	if !ac.Env.Exists() {
		return errs.NewConfigValidateError("app", "env", "env value not match", nil)
	}

	if ac.MaxListLimit < 0 {
		return errs.NewConfigValidateError("app", "max_list_limit", "must be positive", nil)
	}

	if ac.CipherKey == "" {
		return errs.NewConfigValidateError("app", "cipher-key", "must not be empty", nil)
	}

	// ToDo: restore in future
	//if len(ac.AcceptTokenIssuers) == 0 {
	//	return errs.NewConfigValidateError("app", "accept-token-issuers", "must not be empty", nil)
	//}

	return nil
}
