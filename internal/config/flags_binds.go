package config

import (
	"errors"

	conf "github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func bindFlags(flags *pflag.FlagSet, v *viper.Viper) error {
	err := errors.Join(
		// App
		v.BindPFlag(conf.KeyAppEnv, flags.Lookup(conf.FlagAppEnv)),
		v.BindPFlag(conf.KeyAppInitTimeout, flags.Lookup(conf.FlagAppInitTimeout)),
		v.BindPFlag(conf.KeyAppStopTimeout, flags.Lookup(conf.FlagAppStopTimeout)),
		v.BindPFlag(conf.KeyAppCloseTimeout, flags.Lookup(conf.FlagAppCloseTimeout)),
		v.BindPFlag(keyAppNodeName, flags.Lookup(FlagAppNodeName)),
		v.BindPFlag(keyAppMaxListLimit, flags.Lookup(FlagAppMaxListLimit)),
		v.BindPFlag(keyAppTokenIssuer, flags.Lookup(FlagAppTokenIssuer)),
		v.BindPFlag(keyAppCipherKey, flags.Lookup(FlagAppCipherKey)),
		v.BindPFlag(keyAppAcceptTokenIssuers, flags.Lookup(FlagAppAcceptTokenIssuers)),
		// auth tc
		v.BindPFlag(keyAuthTCStartInterval, flags.Lookup(FlagAuthTCStartInterval)),
		v.BindPFlag(keyAuthTCScheduleInterval, flags.Lookup(FlagAuthTCScheduleInterval)),
		v.BindPFlag(keyAuthTCWorkerCount, flags.Lookup(FlagAuthTCWorkerCount)),
		v.BindPFlag(keyAuthTCDataCapacity, flags.Lookup(FlagAuthTCDataCapacity)),
		v.BindPFlag(keyAuthTCCompleteProcessing, flags.Lookup(FlagAuthTCCompleteProcessing)),
		v.BindPFlag(keyAuthTCShutdownTimeout, flags.Lookup(FlagAuthTCShutdownTimeout)),
		v.BindPFlag(keyAuthTCTailInterval, flags.Lookup(FlagAuthTCTailInterval)),
		v.BindPFlag(keyAuthTCTailCut, flags.Lookup(FlagAuthTCTailCut)),
		// data tc
		v.BindPFlag(keyDataTCStartInterval, flags.Lookup(FlagDataTCStartInterval)),
		v.BindPFlag(keyDataTCScheduleInterval, flags.Lookup(FlagDataTCScheduleInterval)),
		v.BindPFlag(keyDataTCWorkerCount, flags.Lookup(FlagDataTCWorkerCount)),
		v.BindPFlag(keyDataTCDataCapacity, flags.Lookup(FlagDataTCDataCapacity)),
		v.BindPFlag(keyDataTCCompleteProcessing, flags.Lookup(FlagDataTCCompleteProcessing)),
		v.BindPFlag(keyDataTCShutdownTimeout, flags.Lookup(FlagDataTCShutdownTimeout)),
		v.BindPFlag(keyDataTCTailInterval, flags.Lookup(FlagDataTCTailInterval)),
		v.BindPFlag(keyDataTCTailCut, flags.Lookup(FlagDataTCTailCut)),
		// Auth
		v.BindPFlag(conf.KeyAuthJWTSecret, flags.Lookup(FlagAuthJWTSecret)),
		v.BindPFlag(conf.KeyAuthJWTSigningMethod, flags.Lookup(FlagAuthJWTSigningMethod)),
		v.BindPFlag(conf.KeyAuthAccessTokenTTL, flags.Lookup(FlagAuthAccessTokenTTL)),
		v.BindPFlag(conf.KeyAuthRefreshTokenTTL, flags.Lookup(FlagAuthRefreshTokenTTL)),
		v.BindPFlag(conf.KeyAuthRSAPrivateKeyPath, flags.Lookup(FlagAuthRSAPrivateKeyPath)),
		v.BindPFlag(conf.KeyAuthMasterPasswordSalt, flags.Lookup(FlagAuthMasterPasswordSalt)),
		// HTTP
		v.BindPFlag(conf.KeyHTTPAddress, flags.Lookup(FlagHTTPAddress)),
		v.BindPFlag(conf.KeyHTTPReadTimeout, flags.Lookup(FlagHTTPReadTimeout)),
		v.BindPFlag(conf.KeyHTTPWriteTimeout, flags.Lookup(FlagHTTPWriteTimeout)),
		v.BindPFlag(conf.KeyHTTPIdleTimeout, flags.Lookup(FlagHTTPIdleTimeout)),
		v.BindPFlag(conf.KeyHTTPShutdownTimeout, flags.Lookup(FlagHTTPShutdownTimeout)),
		v.BindPFlag(conf.KeyHTTPPrivateKeyPath, flags.Lookup(FlagHTTPPrivateKeyPath)),
		v.BindPFlag(conf.KeyHTTPCertificatePath, flags.Lookup(FlagHTTPCertificatePath)),
		v.BindPFlag(conf.KeyHTTPSecure, flags.Lookup(FlagHTTPSecure)),
		v.BindPFlag(conf.KeyHTTPMaxRequestBodySize, flags.Lookup(FlagHTTPMaxRequestBodySize)),
		// gRPC
		v.BindPFlag(conf.KeyGRPCAddress, flags.Lookup(FlagGRPCAddress)),
		v.BindPFlag(conf.KeyGRPCMaxConnIdle, flags.Lookup(FlagGRPCMaxConnIdle)),
		v.BindPFlag(conf.KeyGRPCMaxConnAge, flags.Lookup(FlagGRPCMaxConnAge)),
		v.BindPFlag(conf.KeyGRPCMaxConnAgeGrace, flags.Lookup(FlagGRPCMaxConnAgeGrace)),
		v.BindPFlag(conf.KeyGRPCTimeout, flags.Lookup(FlagGRPCTimeout)),
		v.BindPFlag(conf.KeyGRPCKeepAliveTime, flags.Lookup(FlagGRPCKeepAliveTime)),
		v.BindPFlag(conf.KeyGRPCKeepAliveTimeout, flags.Lookup(FlagGRPCKeepAliveTimeout)),
		v.BindPFlag(conf.KeyGRPCShutdownTimeout, flags.Lookup(FlagGRPCShutdownTimeout)),
		// Log
		v.BindPFlag(conf.KeyLogLevel, flags.Lookup(FlagLogLevel)),
		v.BindPFlag(conf.KeyLogFormat, flags.Lookup(FlagLogFormat)),
		// DB
		v.BindPFlag(conf.KeyDBDriver, flags.Lookup(FlagDBDriver)),
		v.BindPFlag(conf.KeyDBDSN, flags.Lookup(FlagDBDSN)),
		v.BindPFlag(conf.KeyDBMaxOpenConns, flags.Lookup(FlagDBMaxOpenConns)),
		v.BindPFlag(conf.KeyDBMaxIdleConns, flags.Lookup(FlagDBMaxIdleConns)),
		v.BindPFlag(conf.KeyDBConnMaxIdleLifetime, flags.Lookup(FlagDBMaxIdleLifetime)),
		v.BindPFlag(conf.KeyDBConnTimeout, flags.Lookup(FlagDBConnTimeout)),
		// Telemetry
		v.BindPFlag(conf.KeyTelemetryEnabled, flags.Lookup(FlagTelemetryEnabled)),
		v.BindPFlag(conf.KeyTelemetryExporterEndpoint, flags.Lookup(FlagTelemetryExporterEndpoint)),
		v.BindPFlag(conf.KeyTelemetryServiceName, flags.Lookup(FlagTelemetryServiceName)),
		v.BindPFlag(conf.KeyTelemetrySampleRate, flags.Lookup(FlagTelemetrySampleRate)),
		v.BindPFlag(conf.KeyTelemetryTimeout, flags.Lookup(FlagTelemetryTimeout)),
		// amqp connector
		v.BindPFlag(keyAMQPConnectorURL, flags.Lookup(FlagAMQPConnectorURL)),
		v.BindPFlag(keyAMQPConnectorUsername, flags.Lookup(FlagAMQPConnectorUsername)),
		v.BindPFlag(keyAMQPConnectorPassword, flags.Lookup(FlagAMQPConnectorPassword)),
		v.BindPFlag(keyAMQPConnectorConnectTimeout, flags.Lookup(FlagAMQPConnectorConnectTimeout)),
		v.BindPFlag(keyAMQPConnectorWriteTimeout, flags.Lookup(FlagAMQPConnectorWriteTimeout)),
		v.BindPFlag(keyAMQPConnectorIdleTimeout, flags.Lookup(FlagAMQPConnectorIdleTimeout)),
		v.BindPFlag(keyAMQPConnectorShutdownTimeout, flags.Lookup(FlagAMQPConnectorShutdownTimeout)),
		// amqp login attempts receiver
		v.BindPFlag(keyLoginAttemptsReceiverTargetName, flags.Lookup(FlagLoginAttemptsReceiverTargetName)),
		v.BindPFlag(keyLoginAttemptsReceiverConnectTimeout, flags.Lookup(FlagLoginAttemptsReceiverConnectTimeout)),
		v.BindPFlag(keyLoginAttemptsReceiverShutdownTimeout, flags.Lookup(FlagLoginAttemptsReceiverShutdownTimeout)),
		v.BindPFlag(keyLoginAttemptsReceiverPrefetchCredit, flags.Lookup(FLagLoginAttemptsReceiverPrefetchCredit)),
		v.BindPFlag(keyLoginAttemptsStartInterval, flags.Lookup(FlagLoginAttemptsStartInterval)),
		v.BindPFlag(keyLoginAttemptsScheduleInterval, flags.Lookup(FlagLoginAttemptsScheduleInterval)),
		v.BindPFlag(keyLoginAttemptsWorkerCount, flags.Lookup(FlagLoginAttemptsWorkerCount)),
		v.BindPFlag(keyLoginAttemptsDataCapacity, flags.Lookup(FlagLoginAttemptsDataCapacity)),
		v.BindPFlag(keyLoginAttemptsCompleteProcessing, flags.Lookup(FlagLoginAttemptsCompleteProcessing)),
		v.BindPFlag(keyLoginAttemptsShutdownTimeout, flags.Lookup(FlagLoginAttemptsShutdownTimeout)),
		v.BindPFlag(keyLoginAttemptsBatchSize, flags.Lookup(FlagLoginAttemptsBatchSize)),
		v.BindPFlag(keyLoginAttemptsBatchReadTimeout, flags.Lookup(FlagLoginAttemptsBatchReadTimeout)),
		v.BindPFlag(keyLoginAttemptsAcknowledgeTimeout, flags.Lookup(FlagLoginAttemptsAcknowledgeTimeout)),
	)
	if err != nil {
		return errs.NewConfigError("bind flags with keys", err)
	}

	return nil
}
