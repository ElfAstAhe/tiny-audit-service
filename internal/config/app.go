package config

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/errs"
)

// AppConfig — метаданные сервиса
type AppConfig struct {
	Env                       AppEnv        `mapstructure:"env" json:"env,omitempty" yaml:"env,omitempty"` // dev, prod, test
	MaxListLimit              int           `mapstructure:"max_list_limit" json:"max_list_limit,omitempty" yaml:"max_list_limit,omitempty"`
	TokenIssuer               string        `mapstructure:"token_issuer" json:"token_issuer,omitempty" yaml:"token_issuer,omitempty"`
	AcceptTokenIssuers        []string      `mapstructure:"accept_token_issuers" json:"accept_token_issuers,omitempty" yaml:"accept_token_issuers,omitempty"`
	CipherKey                 string        `mapstructure:"cipher_key" json:"cipher_key,omitempty" yaml:"cipher_key,omitempty"`
	AuthTailJobRepeatDuration time.Duration `mapstructure:"auth-tail-job-repeat-duration" json:"auth-tail-job-repeat-duration,omitempty" yaml:"auth-tail-job-repeat-duration,omitempty"`
	AuthTailDuration          time.Duration `mapstructure:"auth_tail_duration" json:"auth_tail_duration,omitempty" yaml:"auth_tail_duration,omitempty"`
	AuthTailCut               bool          `mapstructure:"auth_tail_cut" json:"auth_tail_cut,omitempty" yaml:"auth_tail_cut,omitempty"`
	DataTailJobRepeatDuration time.Duration `mapstructure:"data-tail-job-repeat-duration" json:"data-tail-job-repeat-duration,omitempty" yaml:"data-tail-job-repeat-duration,omitempty"`
	DataTailDuration          time.Duration `mapstructure:"data_tail_duration" json:"data_tail_duration,omitempty" yaml:"data_tail_duration,omitempty"`
	DataTailCut               bool          `mapstructure:"data_tail_cut" json:"data_tail_cut,omitempty" yaml:"data_tail_cut,omitempty"`
}

func NewAppConfig(
	env AppEnv,
	maxListLimit int,
	acceptTokenIssuers []string,
	cipherKey string,
	authTailCut bool,
	authTailDuration time.Duration,
	dataTailCut bool,
	dataTailDuration time.Duration,
) *AppConfig {
	return &AppConfig{
		Env:                env,
		MaxListLimit:       maxListLimit,
		AcceptTokenIssuers: acceptTokenIssuers,
		CipherKey:          cipherKey,
		AuthTailCut:        authTailCut,
		AuthTailDuration:   authTailDuration,
		DataTailCut:        dataTailCut,
		DataTailDuration:   dataTailDuration,
	}
}

func NewDefaultAppConfig() *AppConfig {
	return NewAppConfig(
		defaultAppEnv,
		defaultMaxListLimit,
		[]string{},
		"",
		defaultAuthTailCut,
		defaultAuthTailDuration,
		defaultDataTailCut,
		defaultDataTailDuration,
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
