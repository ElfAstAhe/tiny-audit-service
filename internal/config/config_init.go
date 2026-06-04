package config

import (
	"errors"
	"fmt"
	"os"

	conf "github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
}

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
	}

	// Парсинг
	err = res.Parse(os.Args[1:])
	if err != nil {
		return nil, err
	}

	return res, err
}

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
	)
	if err != nil {
		return errs.NewConfigError("bind flags with keys", err)
	}

	return nil
}
