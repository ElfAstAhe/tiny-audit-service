package config

import (
	"time"

	conf "github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/spf13/viper"
)

// app
const (
	defaultAppEnv       conf.AppEnv = conf.AppEnvDevelopment
	defaultAppNodeName  string      = ApplicationName
	defaultMaxListLimit int         = 100
	defaultTokenIssuer  string      = "tiny-auth-service"
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

// login attempts receiver worker
const (
	defaultLoginAttemptsWorkerStartInterval      time.Duration = 5 * time.Second
	defaultLoginAttemptsWorkerScheduleInterval   time.Duration = 1 * time.Minute
	defaultLoginAttemptsWorkerWorkerCount        int           = 2
	defaultLoginAttemptsWorkerDataCapacity       int           = 128
	defaultLoginAttemptsWorkerCompleteProcessing bool          = false
	defaultLoginAttemptsWorkerShutdownTimeout    time.Duration = 15 * time.Second
	defaultLoginAttemptsWorkerBatchSize          int           = 50
	defaultLoginAttemptsWorkerBatchReadTimeout   time.Duration = 2 * time.Second
	defaultLoginAttemptsWorkerAcknowledgeTimeout time.Duration = 1 * time.Second
)

// login attempts receiver common
const (
	defaultLoginAttemptsReceiverKind string = "amqp"
)

// amqp login attempts receiver (FQQN artemis style)
const (
	// defaultLoginAttemptsReceiverAMQPConfigTargetName - queue/topic (FQQN artemis style)
	defaultLoginAttemptsReceiverAMQPConfigTargetName      string        = "tiny.auth::login.attempts"
	defaultLoginAttemptsReceiverAMQPConfigConnectTimeout  time.Duration = 10 * time.Second
	defaultLoginAttemptsReceiverAMQPConfigShutdownTimeout time.Duration = 10 * time.Second
	defaultLoginAttemptsReceiverAMQPConfigPrefetchCredit  int           = 50
)

// kafka login attempts receiver
const (
	defaultLoginAttemptsReceiverKafkaConfigTargetName string = "tiny.auth.login.attempts"
	defaultLoginAttemptsReceiverKafkaConfigPartition  int    = -1
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
	// login attempts worker
	v.SetDefault(keyLoginAttemptsWorkerStartInterval, defaultLoginAttemptsWorkerStartInterval)
	v.SetDefault(keyLoginAttemptsWorkerScheduleInterval, defaultLoginAttemptsWorkerScheduleInterval)
	v.SetDefault(keyLoginAttemptsWorkerWorkerCount, defaultLoginAttemptsWorkerWorkerCount)
	v.SetDefault(keyLoginAttemptsWorkerDataCapacity, defaultLoginAttemptsWorkerDataCapacity)
	v.SetDefault(keyLoginAttemptsWorkerCompleteProcessing, defaultLoginAttemptsWorkerCompleteProcessing)
	v.SetDefault(keyLoginAttemptsWorkerShutdownTimeout, defaultLoginAttemptsWorkerShutdownTimeout)
	v.SetDefault(keyLoginAttemptsWorkerBatchSize, defaultLoginAttemptsWorkerBatchSize)
	v.SetDefault(keyLoginAttemptsWorkerBatchReadTimeout, defaultLoginAttemptsWorkerBatchReadTimeout)
	v.SetDefault(keyLoginAttemptsWorkerAcknowledgeTimeout, defaultLoginAttemptsWorkerAcknowledgeTimeout)
	// login attempts common
	v.SetDefault(keyLoginAttemptsReceiverKind, defaultLoginAttemptsReceiverKind)
	// amqp login attempts receiver
	v.SetDefault(keyLoginAttemptsReceiverAMQPConfigTargetName, defaultLoginAttemptsReceiverAMQPConfigTargetName)
	v.SetDefault(keyLoginAttemptsReceiverAMQPConfigPrefetchCredit, defaultLoginAttemptsReceiverAMQPConfigPrefetchCredit)
	v.SetDefault(keyLoginAttemptsReceiverAMQPConfigConnectTimeout, defaultLoginAttemptsReceiverAMQPConfigConnectTimeout)
	v.SetDefault(keyLoginAttemptsReceiverAMQPConfigShutdownTimeout, defaultLoginAttemptsReceiverAMQPConfigShutdownTimeout)
	// kafka login attempts receiver
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigBrokers, conf.DefaultKafkaBrokers)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigTargetName, defaultLoginAttemptsReceiverKafkaConfigTargetName)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigPartition, defaultLoginAttemptsReceiverKafkaConfigPartition)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigConnectTimeout, conf.DefaultKafkaReceiverConnectTimeout)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigShutdownTimeout, conf.DefaultKafkaReceiverShutdownTimeout)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigMinBytes, conf.DefaultKafkaReceiverMinBytes)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigMaxBytes, conf.DefaultKafkaReceiverMaxBytes)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigMaxWait, conf.DefaultKafkaReceiverMaxWait)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigHeartbeatInterval, conf.DefaultKafkaReceiverHeartbeatInterval)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigSessionTimeout, conf.DefaultKafkaReceiverSessionTimeout)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigRebalanceTimeout, conf.DefaultKafkaReceiverRebalanceTimeout)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigReadTimeout, conf.DefaultKafkaReceiverReadTimeout)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigMaxAttempts, conf.DefaultKafkaReceiverMaxAttempts)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigQueueCapacity, conf.DefaultKafkaReceiverQueueCapacity)
	v.SetDefault(keyLoginAttemptsReceiverKafkaConfigStartOffset, conf.DefaultKafkaReceiverStartOffset)
}
