package config

import (
	"fmt"
	"os"

	conf "github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"github.com/spf13/pflag"
)

func initFLags() (res *pflag.FlagSet, err error) {
	defer func() {
		if r := recover(); r != nil {
			// Проверяем, является ли r ошибкой
			recoveryErr, ok := r.(error)
			if !ok {
				// Если это строка или что-то другое, приводим к виду error вручную
				recoveryErr = errs.NewConfigError(fmt.Sprintf("panic [%v] recovery", r), nil)
			}
			res = nil
			err = errs.NewConfigError("parse cli flags panic", recoveryErr)
		}
	}()

	res = pflag.NewFlagSet("cmd flags", pflag.PanicOnError)

	// Используем константы Flag...
	{
		// app
		res.String(FlagConfig, "config/config.yaml", "path to config file")
		res.String(conf.FlagAppEnv, string(defaultAppEnv), "application environment")
		res.Duration(conf.FlagAppInitTimeout, conf.DefaultAppInitTimeout, "application init timeout")
		res.Duration(conf.FlagAppStopTimeout, conf.DefaultAppStopTimeout, "application stop timeout")
		res.Duration(conf.FlagAppCloseTimeout, conf.DefaultAppCloseTimeout, "application close timeout")
		res.String(FlagAppNodeName, defaultAppNodeName, "application node name")
		res.Int(FlagAppMaxListLimit, usecase.DefaultMaxLimit, "max list limit")
		res.String(FlagAppCipherKey, "", "cipher key")
		res.String(FlagAppTokenIssuer, "", "token issuer")
		res.StringSlice(FlagAppAcceptTokenIssuers, []string{}, `accept token issuers separated by comma, like: "issuer1,issuer2,issuer3"`)
		// auth tc
		res.Duration(FlagAuthTCStartInterval, defaultAuthTCStartInterval, "auth tail cutter worker start interval")
		res.Duration(FlagAuthTCScheduleInterval, defaultAuthTCScheduleInterval, "auth tail cutter worker schedule interval")
		res.Int(FlagAuthTCWorkerCount, defaultAuthTCWorkerCount, "auth tail cutter workers count")
		res.Int(FlagAuthTCDataCapacity, defaultAuthTCDataCapacity, "auth tail cutter data capacity")
		res.Bool(FlagAuthTCCompleteProcessing, defaultAuthTCCompleteProcessing, "auth tail cutter try complete data queue on shutdown")
		res.Duration(FlagAuthTCShutdownTimeout, defaultAuthTCShutdownTimeout, "auth tail cutter shutdown timeout")
		res.Duration(FlagAuthTCTailInterval, defaultAuthTCTailInterval, "auth tail cutter tail interval")
		res.Bool(FlagAuthTCTailCut, defaultAuthTCTailCut, "auth tail cutter enabler")
		// data tc
		res.Duration(FlagDataTCStartInterval, defaultDataTCStartInterval, "data tail cutter worker start interval")
		res.Duration(FlagDataTCScheduleInterval, defaultDataTCScheduleInterval, "data tail cutter worker schedule interval")
		res.Int(FlagDataTCWorkerCount, defaultDataTCWorkerCount, "data tail cutter workers count")
		res.Int(FlagDataTCDataCapacity, defaultDataTCDataCapacity, "data tail cutter data capacity")
		res.Bool(FlagDataTCCompleteProcessing, defaultDataTCCompleteProcessing, "data tail cutter try complete data queue on shutdown")
		res.Duration(FlagDataTCShutdownTimeout, defaultDataTCShutdownTimeout, "data tail cutter shutdown timeout")
		res.Duration(FlagDataTCTailInterval, defaultDataTCTailInterval, "data tail cutter tail interval")
		res.Bool(FlagDataTCTailCut, defaultDataTCTailCut, "data tail cutter enabler")
		// Auth
		res.String(FlagAuthJWTSecret, "", "JWT secret")
		res.String(FlagAuthJWTSigningMethod, conf.DefaultAuthSigningMethod, "JWT signing method")
		res.Duration(FlagAuthAccessTokenTTL, conf.DefaultAuthAccessTokenTTL, "JWT access token TTL")
		res.Duration(FlagAuthRefreshTokenTTL, conf.DefaultAuthRefreshTokenTTL, "JWT refresh token TTL")
		res.String(FlagAuthRSAPrivateKeyPath, "", "RSA private key path")
		res.String(FlagAuthMasterPasswordSalt, "", "master password salt")
		// HTTP
		res.String(FlagHTTPAddress, conf.DefaultHTTPAddress, "http address")
		res.Duration(FlagHTTPReadTimeout, conf.DefaultHTTPReadTimeout, "http read timeout")
		res.Duration(FlagHTTPWriteTimeout, conf.DefaultHTTPWriteTimeout, "http write timeout")
		res.Duration(FlagHTTPIdleTimeout, conf.DefaultHTTPIdleTimeout, "http idle timeout")
		res.Duration(FlagHTTPShutdownTimeout, conf.DefaultHTTPShutdownTimeout, "http shutdown timeout")
		res.String(FlagHTTPPrivateKeyPath, "", "http private key path")
		res.String(FlagHTTPCertificatePath, "", "http certificate path")
		res.Bool(FlagHTTPSecure, conf.DefaultHTTPSecure, "http secure mode")
		res.Int(FlagHTTPMaxRequestBodySize, conf.DefaultHTTPMaxRequestBodySize, "http max request body size")
		// gRPC
		res.String(FlagGRPCAddress, conf.DefaultGRPCAddress, "gRPC address")
		res.Duration(FlagGRPCMaxConnIdle, conf.DefaultGRPCMaxConnIdle, "gRPC max connection idle timeout")
		res.Duration(FlagGRPCMaxConnAge, conf.DefaultGRPCMaxConnAge, "gRPC max connection age timeout")
		res.Duration(FlagGRPCMaxConnAgeGrace, conf.DefaultGRPCMaxConnAgeGrace, "gRPC max connection age grace timeout")
		res.Duration(FlagGRPCTimeout, conf.DefaultGRPCTimeout, "gRPC timeout")
		res.Duration(FlagGRPCKeepAliveTime, conf.DefaultGRPCKeepAliveTime, "gRPC keep alive timeout")
		res.Duration(FlagGRPCKeepAliveTimeout, conf.DefaultGRPCKeepAliveTimeout, "gRPC keep alive timeout")
		res.Duration(FlagGRPCShutdownTimeout, conf.DefaultGRPCShutdownTimeout, "gRPC shutdown timeout")
		// DB
		res.String(FlagDBDSN, conf.DefaultDBDSN, "database dsn")
		res.String(FlagDBDriver, conf.DefaultDBDriver, "database driver name/alias")
		res.Int(FlagDBMaxOpenConns, conf.DefaultDBMaxOpenConns, "db max open connections")
		res.Int(FlagDBMaxIdleConns, conf.DefaultDBMaxIdleConns, "db max idle connections")
		res.Duration(FlagDBMaxIdleLifetime, conf.DefaultDBConnMaxIdleLifetime, "db max idle connection lifetime")
		res.Duration(FlagDBConnTimeout, conf.DefaultDBConnTimeout, "db connection timeout)")
		// Log
		res.String(FlagLogLevel, conf.DefaultLogLevel, "log level")
		res.String(FlagLogFormat, conf.DefaultLogFormat, "log format")
		// Telemetry
		res.Bool(FlagTelemetryEnabled, conf.DefaultTelemetryEnabled, "telemetry enabled")
		res.String(FlagTelemetryServiceName, "", "telemetry service name")
		res.String(FlagTelemetryExporterEndpoint, conf.DefaultTelemetryExporterEndpoint, "telemetry exporter endpoint")
		res.Float64(FlagTelemetrySampleRate, conf.DefaultTelemetrySampleRate, "telemetry sample rate")
		res.Duration(FlagTelemetryTimeout, conf.DefaultTelemetryTimeout, "telemetry timeout")
		// amqp connector
		res.String(FlagAMQPConnectorURL, conf.DefaultAMQPConnectorURL, "amqp connector url")
		res.String(FlagAMQPConnectorUsername, defaultAMQPConnectorUsername, "amqp connector username")
		res.String(FlagAMQPConnectorPassword, defaultAMQPConnectorPassword, "amqp connector password")
		res.Duration(FlagAMQPConnectorConnectTimeout, conf.DefaultAMQPConnectorConnectTimeout, "amqp connector connect timeout")
		res.Duration(FlagAMQPConnectorIdleTimeout, conf.DefaultAMQPConnectorIdleTimeout, "amqp connector idle timeout")
		res.Duration(FlagAMQPConnectorWriteTimeout, conf.DefaultAMQPConnectorWriteTimeout, "amqp connector write timeout")
		res.Duration(FlagAMQPConnectorShutdownTimeout, conf.DefaultAMQPConnectorShutdownTimeout, "amqp connector shutdown timeout")
		// login attempts receiver
		res.String(FlagLoginAttemptsReceiverTargetName, defaultLoginAttemptsReceiverTargetName, "login attempts receiver queue/topic name")
		res.Duration(FlagLoginAttemptsReceiverConnectTimeout, defaultLoginAttemptsReceiverConnectTimeout, "login attempts receiver connect timeout")
		res.Duration(FlagLoginAttemptsReceiverShutdownTimeout, defaultLoginAttemptsReceiverShutdownTimeout, "login attempts receiver shutdown timeout")
		res.Int(FLagLoginAttemptsReceiverPrefetchCredit, conf.DefaultAMQPReceiverPrefetchCredit, "login attempts receiver prefetch credit")
		res.Duration(FlagLoginAttemptsStartInterval, defaultLoginAttemptsStartInterval, "login attempts scheduler start interval")
		res.Duration(FlagLoginAttemptsScheduleInterval, defaultLoginAttemptsScheduleInterval, "login attempts scheduler schedule interval")
		res.Duration(FlagLoginAttemptsShutdownTimeout, defaultLoginAttemptsShutdownTimeout, "login attempts scheduler shutdown timeout")
		res.Int(FlagLoginAttemptsWorkerCount, defaultLoginAttemptsWorkerCount, "login attempts pool worker count")
		res.Int(FlagLoginAttemptsDataCapacity, defaultLoginAttemptsDataCapacity, "login attempts pool data capacity")
		res.Bool(FlagLoginAttemptsCompleteProcessing, defaultLoginAttemptsCompleteProcessing, "login attempts pool force complete processing")
		res.Int(FlagLoginAttemptsBatchSize, defaultLoginAttemptsBatchSize, "login attempts batch size")
		res.Duration(FlagLoginAttemptsBatchReadTimeout, defaultLoginAttemptsBatchReadTimeout, "login attempts batch read timeout")
		res.Duration(FlagLoginAttemptsAcknowledgeTimeout, defaultLoginAttemptsAcknowledgeTimeout, "login attempts receiver acknowledge timeout")
	}

	// Парсинг
	err = res.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	return res, err
}
