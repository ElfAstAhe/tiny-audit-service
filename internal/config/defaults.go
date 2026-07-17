package config

import (
	"time"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	conf "github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/spf13/viper"
)

// app
const (
	defaultAppEnv       config.AppEnv = config.AppEnvDevelopment
	defaultAppNodeName  string        = ApplicationName
	defaultMaxListLimit int           = 100
	defaultTokenIssuer  string        = "tiny-auth-service"
)

// auth tail cutter
const (
	defaultAuthTCStartInterval      time.Duration = 5 * time.Second
	defaultAuthTCScheduleInterval   time.Duration = 1 * time.Minute
	defaultAuthTCWorkerCount        int           = 2
	defaultAuthTCDataCapacity       int           = 128
	defaultAuthTCCompleteProcessing bool          = false
	defaultAuthTCShutdownTimeout    time.Duration = 15 * time.Second
	defaultAuthTCTailInterval       time.Duration = 182 * 24 * time.Hour // 182 days
	defaultAuthTCTailCut            bool          = true
)

// data tail cutter
const (
	defaultDataTCStartInterval      time.Duration = 5 * time.Second
	defaultDataTCScheduleInterval   time.Duration = 1 * time.Minute
	defaultDataTCWorkerCount        int           = 2
	defaultDataTCDataCapacity       int           = 128
	defaultDataTCCompleteProcessing bool          = false
	defaultDataTCShutdownTimeout    time.Duration = 15 * time.Second
	defaultDataTCTailInterval       time.Duration = 365 * 24 * time.Hour // 1 year
	defaultDataTCTailCut            bool          = true
)

// amqp connector
const (
	defaultAMQPConnectorUsername string = "svc-audit"
	defaultAMQPConnectorPassword string = "test"
)

// amqp login attempts receiver (FQQN artemis style)
const (
	defaultLoginAttemptsReceiverTargetName string = "tiny.auth::login.attempts"
)

// app
const (
	keyAppNodeName           string = "app.node_name"
	keyAppMaxListLimit       string = "app.max_list_limit"
	keyAppTokenIssuer        string = "app.token_issuer"
	keyAppCipherKey          string = "app.cipher_key"
	keyAppAcceptTokenIssuers string = "app.accept_token_issuers"
)

// auth tail cutter
const (
	keyAuthTCStartInterval      string = "auth_tc.start_interval"
	keyAuthTCScheduleInterval   string = "auth_tc.schedule_interval"
	keyAuthTCWorkerCount        string = "auth_tc.worker_count"
	keyAuthTCDataCapacity       string = "auth_tc.data_capacity"
	keyAuthTCCompleteProcessing string = "auth_tc.complete_processing"
	keyAuthTCShutdownTimeout    string = "auth_tc.shutdown_timeout"
	keyAuthTCTailInterval       string = "auth_tc.tail_interval"
	keyAuthTCTailCut            string = "auth_tc.tail_cut"
)

// data tail cutter
const (
	keyDataTCStartInterval      string = "data_tc.start_interval"
	keyDataTCScheduleInterval   string = "data_tc.schedule_interval"
	keyDataTCWorkerCount        string = "data_tc.worker_count"
	keyDataTCDataCapacity       string = "data_tc.data_capacity"
	keyDataTCCompleteProcessing string = "data_tc.complete_processing"
	keyDataTCShutdownTimeout    string = "data_tc.shutdown_timeout"
	keyDataTCTailInterval       string = "data_tc.tail_interval"
	keyDataTCTailCut            string = "data_tc.tail_cut"
)

// amqp connector
const (
	keyAMQPConnectorURL             string = "amqp_connector.url"
	keyAMQPConnectorUsername        string = "amqp_connector.username"
	keyAMQPConnectorPassword        string = "amqp_connector.password"
	keyAMQPConnectorConnectTimeout  string = "amqp_connector.connect_timeout"
	keyAMQPConnectorWriteTimeout    string = "amqp_connector.write_timeout"
	keyAMQPConnectorIdleTimeout     string = "amqp_connector.idle_timeout"
	keyAMQPConnectorShutdownTimeout string = "amqp_connector.shutdown_timeout"
)

// amqp login attempts sender
const (
	keyLoginAttemptsReceiverTargetName      string = "login_attempts_receiver.target_name"
	keyLoginAttemptsReceiverConnectTimeout  string = "login_attempts_receiver.connect_timeout"
	keyLoginAttemptsReceiverShutdownTimeout string = "login_attempts_receiver.shutdown_timeout"
)

//goland:noinspection DuplicatedCode
func applyDefaults(v *viper.Viper) {
	// App
	v.SetDefault(conf.KeyAppEnv, defaultAppEnv)
	v.SetDefault(conf.KeyAppInitTimeout, conf.DefaultAppInitTimeout)
	v.SetDefault(conf.KeyAppStopTimeout, conf.DefaultAppStopTimeout)
	v.SetDefault(conf.KeyAppCloseTimeout, conf.DefaultAppCloseTimeout)
	v.SetDefault(keyAppNodeName, defaultAppNodeName)
	v.SetDefault(keyAppMaxListLimit, defaultMaxListLimit)
	v.SetDefault(keyAppTokenIssuer, defaultTokenIssuer)
	// auth tc
	v.SetDefault(keyAuthTCStartInterval, defaultAuthTCStartInterval)
	v.SetDefault(keyAuthTCScheduleInterval, defaultAuthTCScheduleInterval)
	v.SetDefault(keyAuthTCWorkerCount, defaultAuthTCWorkerCount)
	v.SetDefault(keyAuthTCDataCapacity, defaultAuthTCDataCapacity)
	v.SetDefault(keyAuthTCCompleteProcessing, defaultAuthTCCompleteProcessing)
	v.SetDefault(keyAuthTCShutdownTimeout, defaultAuthTCShutdownTimeout)
	v.SetDefault(keyAuthTCTailInterval, defaultAuthTCTailInterval)
	v.SetDefault(keyAuthTCTailCut, defaultAuthTCTailCut)
	// data tc
	v.SetDefault(keyDataTCStartInterval, defaultDataTCStartInterval)
	v.SetDefault(keyDataTCScheduleInterval, defaultDataTCScheduleInterval)
	v.SetDefault(keyDataTCWorkerCount, defaultDataTCWorkerCount)
	v.SetDefault(keyDataTCDataCapacity, defaultDataTCDataCapacity)
	v.SetDefault(keyDataTCCompleteProcessing, defaultDataTCCompleteProcessing)
	v.SetDefault(keyDataTCShutdownTimeout, defaultDataTCShutdownTimeout)
	v.SetDefault(keyDataTCTailInterval, defaultDataTCTailInterval)
	v.SetDefault(keyDataTCTailCut, defaultDataTCTailCut)
	// Auth
	v.SetDefault(conf.KeyAuthJWTSigningMethod, conf.DefaultAuthSigningMethod)
	v.SetDefault(conf.KeyAuthAccessTokenTTL, conf.DefaultAuthAccessTokenTTL)
	v.SetDefault(conf.KeyAuthRefreshTokenTTL, conf.DefaultAuthRefreshTokenTTL)
	// HTTP
	v.SetDefault(conf.KeyHTTPAddress, conf.DefaultHTTPAddress)
	v.SetDefault(conf.KeyHTTPReadTimeout, conf.DefaultHTTPReadTimeout)
	v.SetDefault(conf.KeyHTTPWriteTimeout, conf.DefaultHTTPWriteTimeout)
	v.SetDefault(conf.KeyHTTPIdleTimeout, conf.DefaultHTTPIdleTimeout)
	v.SetDefault(conf.KeyHTTPShutdownTimeout, conf.DefaultHTTPShutdownTimeout)
	v.SetDefault(conf.KeyHTTPSecure, conf.DefaultHTTPSecure)
	v.SetDefault(conf.KeyHTTPMaxRequestBodySize, conf.DefaultHTTPMaxRequestBodySize)
	// gRPC
	v.SetDefault(conf.KeyGRPCAddress, conf.DefaultGRPCAddress)
	v.SetDefault(conf.KeyGRPCMaxConnIdle, conf.DefaultGRPCMaxConnIdle)
	v.SetDefault(conf.KeyGRPCMaxConnAge, conf.DefaultGRPCMaxConnAge)
	v.SetDefault(conf.KeyGRPCMaxConnAgeGrace, conf.DefaultGRPCMaxConnAgeGrace)
	v.SetDefault(conf.KeyGRPCTimeout, conf.DefaultGRPCTimeout)
	v.SetDefault(conf.KeyGRPCKeepAliveTime, conf.DefaultGRPCKeepAliveTime)
	v.SetDefault(conf.KeyGRPCKeepAliveTimeout, conf.DefaultGRPCKeepAliveTimeout)
	v.SetDefault(conf.KeyGRPCShutdownTimeout, conf.DefaultGRPCShutdownTimeout)
	// DB
	v.SetDefault(conf.KeyDBDriver, conf.DefaultDBDriver)
	v.SetDefault(conf.KeyDBDSN, conf.DefaultDBDSN)
	v.SetDefault(conf.KeyDBMaxOpenConns, conf.DefaultDBMaxOpenConns)
	v.SetDefault(conf.KeyDBMaxIdleConns, conf.DefaultDBMaxIdleConns)
	v.SetDefault(conf.KeyDBConnMaxIdleLifetime, conf.DefaultDBConnMaxIdleLifetime)
	v.SetDefault(conf.KeyDBConnTimeout, conf.DefaultDBConnTimeout)
	// Log
	v.SetDefault(conf.KeyLogLevel, conf.DefaultLogLevel)
	v.SetDefault(conf.KeyLogFormat, conf.DefaultLogFormat)
	// Telemetry
	v.SetDefault(conf.KeyTelemetryEnabled, conf.DefaultTelemetryEnabled)
	v.SetDefault(conf.KeyTelemetryExporterEndpoint, conf.DefaultTelemetryExporterEndpoint)
	v.SetDefault(conf.KeyTelemetrySampleRate, conf.DefaultTelemetrySampleRate)
	v.SetDefault(conf.KeyTelemetryTimeout, conf.DefaultTelemetryTimeout)
	// amqp connector
	v.SetDefault(keyAMQPConnectorURL, conf.DefaultAMQPConnectorURL)
	v.SetDefault(keyAMQPConnectorUsername, defaultAMQPConnectorUsername)
	v.SetDefault(keyAMQPConnectorPassword, defaultAMQPConnectorPassword)
	v.SetDefault(keyAMQPConnectorConnectTimeout, conf.DefaultAMQPSenderConnectTimeout)
	v.SetDefault(keyAMQPConnectorWriteTimeout, conf.DefaultAMQPConnectorWriteTimeout)
	v.SetDefault(keyAMQPConnectorIdleTimeout, conf.DefaultAMQPConnectorIdleTimeout)
	v.SetDefault(keyAMQPConnectorShutdownTimeout, conf.DefaultAMQPConnectorShutdownTimeout)
	// amqp login attempts receiver
	v.SetDefault(keyLoginAttemptsReceiverTargetName, defaultLoginAttemptsReceiverTargetName)
}
