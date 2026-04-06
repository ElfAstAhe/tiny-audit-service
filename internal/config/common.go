package config

// FlagConfig - файл конфигурации
const FlagConfig = "config-path"

// App config flags
const (
	FlagAppEnv                       string = "env"
	FlagAppMaxListLimit              string = "app-max-list-limit"
	FlagAppTokenIssuer               string = "app-token-issuer"
	FlagAppCipherKey                 string = "app-cipher-key"
	FlagAppAcceptTokenIssuers        string = "app-accept-token-issuers"
	FlagAppAuthTailJobRepeatDuration string = "app-auth-tail-job-repeat-duration"
	FlagAppAuthTailCut               string = "app-auth-tail-cut"
	FlagAppAuthTailDuration          string = "app-auth-tail-duration"
	FlagAppDataTailJobRepeatDuration string = "app-data-tail-job-repeat-duration"
	FlagAppDataTailCut               string = "app-data-tail-cut"
	FlagAppDataTailDuration          string = "app-data-tail-duration"
)

// Auth config flags
const (
	FlagAuthJWTSecret          string = "auth-jwt-secret"
	FlagAuthJWTSigningMethod   string = "auth-jwt-signing-method"
	FlagAuthAccessTokenTTL     string = "auth-access-token-ttl"
	FlagAuthRefreshTokenTTL    string = "auth-refresh-token-ttl"
	FlagAuthRSAPrivateKeyPath  string = "auth-rsa-private-key-path"
	FlagAuthMasterPasswordSalt string = "auth-master-password-salt"
)

// DB config flags
const (
	FlagDBDSN             string = "db-dsn"
	FlagDBDriver          string = "db-driver"
	FlagDBMaxOpenConns    string = "db-max-open-conns"
	FlagDBMaxIdleConns    string = "db-max-idle-conns"
	FlagDBMaxIdleLifetime string = "db-max-idle-lifetime"
	FlagDBConnTimeout     string = "db-conn-timeout"
)

// gRPC config flags
const (
	FlagGRPCAddress          string = "grpc-address"
	FlagGRPCMaxConnIdle      string = "grpc-max-conn-idle"
	FlagGRPCMaxConnAge       string = "grpc-max-conn-age"
	FlagGRPCMaxConnAgeGrace  string = "grpc-max-conn-age-grace"
	FlagGRPCTimeout          string = "grpc-timeout"
	FlagGRPCKeepAliveTime    string = "grpc-keep-alive-time"
	FlagGRPCKeepAliveTimeout string = "grpc-keep-alive-timeout"
	FlagGRPCShutdownTimeout  string = "grpc-shutdown-timeout"
)

// http config flags
const (
	FlagHTTPAddress            string = "http-address"
	FlagHTTPReadTimeout        string = "http-read-timeout"
	FlagHTTPWriteTimeout       string = "http-write-timeout"
	FlagHTTPIdleTimeout        string = "http-idle-timeout"
	FlagHTTPShutdownTimeout    string = "http-shutdown-timeout"
	FlagHTTPPrivateKeyPath     string = "http-private-key-path"
	FlagHTTPCertificatePath    string = "http-certificate-path"
	FlagHTTPSecure             string = "http-secure"
	FlagHTTPMaxRequestBodySize string = "http-max-request-body-size"
)

// log config flags
const (
	FlagLogLevel  string = "log-level"
	FlagLogFormat string = "log-format"
)

// telemetry
const (
	FlagTelemetryEnabled          string = "telemetry-enabled"
	FlagTelemetryServiceName      string = "telemetry-service-name"
	FlagTelemetryExporterEndpoint string = "telemetry-exporter-endpoint"
	FlagTelemetrySampleRate       string = "telemetry-sample-rate"
	FlagTelemetryTimeout          string = "telemetry-timeout"
)

// EnvConfig - файл конфигурации
const EnvConfig string = "CONFIG_PATH"

type AppEnv string

func (ae AppEnv) Exists() bool {
	return appEnvs.Contains(ae)
}

type appEnvList map[AppEnv]struct{}

func (ae appEnvList) Contains(env AppEnv) bool {
	_, ok := ae[env]

	return ok
}

// app env enum
const (
	AppEnvProduction  AppEnv = "prod"
	AppEnvDevelopment AppEnv = "dev"
	AppEnvTest        AppEnv = "test"
)

var appEnvs appEnvList = map[AppEnv]struct{}{
	AppEnvProduction:  {},
	AppEnvDevelopment: {},
	AppEnvTest:        {},
}
