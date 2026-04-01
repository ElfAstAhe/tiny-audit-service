package config

const (
	defaultAppEnv       AppEnv = AppEnvDevelopment
	defaultMaxListLimit int    = 100
	defaultTokenIssuer  string = "tiny-auth-service"
)

const (
	keyAppEnv                string = "app.env"
	keyAppMaxListLimit       string = "app.max_list_limit"
	keyAppTokenIssuer        string = "app.token_issuer"
	keyAppCipherKey          string = "app.cipher_key"
	keyAppAcceptTokenIssuers string = "app.accept_token_issuers"
)
