package config

import (
	"time"
)

const (
	defaultAppEnv                    AppEnv        = AppEnvDevelopment
	defaultMaxListLimit              int           = 100
	defaultTokenIssuer               string        = "tiny-auth-service"
	defaultAuthTailJobRepeatDuration time.Duration = 1 * time.Minute
	defaultAuthTailCut               bool          = false
	defaultAuthTailDuration          time.Duration = 182 * 24 * time.Hour // 182 days
	defaultDataTailJobRepeatDuration time.Duration = 1 * time.Minute
	defaultDataTailCut               bool          = false
	defaultDataTailDuration          time.Duration = 365 * 24 * time.Hour // 1 year
)

const (
	keyAppEnv                       string = "app.env"
	keyAppMaxListLimit              string = "app.max_list_limit"
	keyAppTokenIssuer               string = "app.token_issuer"
	keyAppCipherKey                 string = "app.cipher_key"
	keyAppAcceptTokenIssuers        string = "app.accept_token_issuers"
	keyAppAuthTailJobRepeatDuration string = "app.auth_tail_job_repeat_duration"
	keyAppAuthTailCut               string = "app.auth_tail_cut"
	keyAppAuthTailDuration          string = "app.auth_tail_duration"
	keyAppDataTailJobRepeatDuration string = "app.data_tail_job_repeat_duration"
	keyAppDataTailCut               string = "app.data_tail_cut"
	keyAppDataTailDuration          string = "app.data_tail_duration"
)
