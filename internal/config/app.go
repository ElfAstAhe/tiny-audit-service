package config

import (
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// AppConfig — метаданные сервиса
type AppConfig struct {
	Env                AppEnv   `mapstructure:"env" json:"env,omitempty" yaml:"env,omitempty"` // dev, prod, test
	MaxListLimit       int      `mapstructure:"max_list_limit" json:"max_list_limit,omitempty" yaml:"max_list_limit,omitempty"`
	TokenIssuer        string   `mapstructure:"token_issuer" json:"token_issuer,omitempty" yaml:"token_issuer,omitempty"`
	AcceptTokenIssuers []string `mapstructure:"accept_token_issuers" json:"accept_token_issuers,omitempty" yaml:"accept_token_issuers,omitempty"`
	CipherKey          string   `mapstructure:"cipher_key" json:"cipher_key,omitempty" yaml:"cipher_key,omitempty"`
}

func NewAppConfig(env AppEnv, maxListLimit int) *AppConfig {
	return &AppConfig{
		Env:                env,
		MaxListLimit:       maxListLimit,
		AcceptTokenIssuers: make([]string, 0),
	}
}

func NewDefaultAppConfig() *AppConfig {
	return NewAppConfig(defaultAppEnv, defaultMaxListLimit)
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
