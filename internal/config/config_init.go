package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/ElfAstAhe/go-service-template/pkg/config"
	"github.com/ElfAstAhe/go-service-template/pkg/errs"
	"github.com/ElfAstAhe/tiny-audit-service/internal/usecase"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func applyDefaults(v *viper.Viper) {
	// App
	v.SetDefault(keyAppEnv, defaultAppEnv)
	v.SetDefault(keyAppMaxListLimit, defaultMaxListLimit)
	v.SetDefault(keyAppTokenIssuer, defaultTokenIssuer)

	// Auth
	v.SetDefault(config.KeyAuthJWTSigningMethod, config.DefaultAuthSigningMethod)
	v.SetDefault(config.KeyAuthAccessTokenTTL, config.DefaultAuthAccessTokenTTL)
	v.SetDefault(config.KeyAuthRefreshTokenTTL, config.DefaultAuthRefreshTokenTTL)

	// HTTP
	v.SetDefault(config.KeyHTTPAddress, config.DefaultHTTPAddress)
	v.SetDefault(config.KeyHTTPReadTimeout, config.DefaultHTTPReadTimeout)
	v.SetDefault(config.KeyHTTPWriteTimeout, config.DefaultHTTPWriteTimeout)
	v.SetDefault(config.KeyHTTPIdleTimeout, config.DefaultHTTPIdleTimeout)
	v.SetDefault(config.KeyHTTPShutdownTimeout, config.DefaultHTTPShutdownTimeout)
	v.SetDefault(config.KeyHTTPSecure, config.DefaultHTTPSecure)
	v.SetDefault(config.KeyHTTPMaxRequestBodySize, config.DefaultHTTPMaxRequestBodySize)

	// gRPC
	v.SetDefault(config.KeyGRPCAddress, config.DefaultGRPCAddress)
	v.SetDefault(config.KeyGRPCMaxConnIdle, config.DefaultGRPCMaxConnIdle)
	v.SetDefault(config.KeyGRPCMaxConnAge, config.DefaultGRPCMaxConnAge)
	v.SetDefault(config.KeyGRPCMaxConnAgeGrace, config.DefaultGRPCMaxConnAgeGrace)
	v.SetDefault(config.KeyGRPCTimeout, config.DefaultGRPCTimeout)
	v.SetDefault(config.KeyGRPCKeepAliveTime, config.DefaultGRPCKeepAliveTime)
	v.SetDefault(config.KeyGRPCKeepAliveTimeout, config.DefaultGRPCKeepAliveTimeout)
	v.SetDefault(config.KeyGRPCShutdownTimeout, config.DefaultGRPCShutdownTimeout)

	// DB
	v.SetDefault(config.KeyDBDriver, config.DefaultDBDriver)
	v.SetDefault(config.KeyDBDSN, config.DefaultDBDSN)
	v.SetDefault(config.KeyDBMaxOpenConns, config.DefaultDBMaxOpenConns)
	v.SetDefault(config.KeyDBMaxIdleConns, config.DefaultDBMaxIdleConns)
	v.SetDefault(config.KeyDBConnMaxIdleLifetime, config.DefaultDBConnMaxIdleLifetime)
	v.SetDefault(config.KeyDBConnTimeout, config.DefaultDBConnTimeout)

	// Log
	v.SetDefault(config.KeyLogLevel, config.DefaultLogLevel)
	v.SetDefault(config.KeyLogFormat, config.DefaultLogFormat)

	// Telemetry
	v.SetDefault(config.KeyTelemetryEnabled, config.DefaultTelemetryEnabled)
	v.SetDefault(config.KeyTelemetryExporterEndpoint, config.DefaultTelemetryExporterEndpoint)
	v.SetDefault(config.KeyTelemetrySampleRate, config.DefaultTelemetrySampleRate)
	v.SetDefault(config.KeyTelemetryTimeout, config.DefaultTelemetryTimeout)
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
		res.String(FlagAppEnv, string(defaultAppEnv), "application environment")
		res.Int(FlagAppMaxListLimit, usecase.DefaultMaxLimit, "max list limit")
		res.String(FlagAppCipherKey, "", "cipher key")
		res.String(FlagAppCipherKey, "", "cipher key")
		res.StringSlice(FlagAppAcceptTokenIssuers, []string{}, `accept token issuers like: "issuer1 issuer2 issuer3"`)

		// Auth
		res.String(FlagAuthJWTSecret, "", "JWT secret")
		res.String(FlagAuthJWTSigningMethod, config.DefaultAuthSigningMethod, "JWT signing method")
		res.Duration(FlagAuthAccessTokenTTL, config.DefaultAuthAccessTokenTTL, "JWT access token TTL")
		res.Duration(FlagAuthRefreshTokenTTL, config.DefaultAuthRefreshTokenTTL, "JWT refresh token TTL")
		res.String(FlagAuthRSAPrivateKeyPath, "", "RSA private key path")
		res.String(FlagAuthMasterPasswordSalt, "", "master password salt")

		// HTTP
		res.String(FlagHTTPAddress, config.DefaultHTTPAddress, "http address")
		res.Duration(FlagHTTPReadTimeout, config.DefaultHTTPReadTimeout, "http read timeout")
		res.Duration(FlagHTTPWriteTimeout, config.DefaultHTTPWriteTimeout, "http write timeout")
		res.Duration(FlagHTTPIdleTimeout, config.DefaultHTTPIdleTimeout, "http idle timeout")
		res.Duration(FlagHTTPShutdownTimeout, config.DefaultHTTPShutdownTimeout, "http shutdown timeout")
		res.String(FlagHTTPPrivateKeyPath, "", "http private key path")
		res.String(FlagHTTPCertificatePath, "", "http certificate path")
		res.Bool(FlagHTTPSecure, config.DefaultHTTPSecure, "http secure mode")
		res.Int(FlagHTTPMaxRequestBodySize, config.DefaultHTTPMaxRequestBodySize, "http max request body size")

		// gRPC
		res.String(FlagGRPCAddress, config.DefaultGRPCAddress, "gRPC address")
		res.Duration(FlagGRPCMaxConnIdle, config.DefaultGRPCMaxConnIdle, "gRPC max connection idle timeout")
		res.Duration(FlagGRPCMaxConnAge, config.DefaultGRPCMaxConnAge, "gRPC max connection age timeout")
		res.Duration(FlagGRPCMaxConnAgeGrace, config.DefaultGRPCMaxConnAgeGrace, "gRPC max connection age grace timeout")
		res.Duration(FlagGRPCTimeout, config.DefaultGRPCTimeout, "gRPC timeout")
		res.Duration(FlagGRPCKeepAliveTime, config.DefaultGRPCKeepAliveTime, "gRPC keep alive timeout")
		res.Duration(FlagGRPCKeepAliveTimeout, config.DefaultGRPCKeepAliveTimeout, "gRPC keep alive timeout")
		res.Duration(FlagGRPCShutdownTimeout, config.DefaultGRPCShutdownTimeout, "gRPC shutdown timeout")

		// DB
		res.String(FlagDBDSN, config.DefaultDBDSN, "database dsn")
		res.String(FlagDBDriver, config.DefaultDBDriver, "database driver name/alias")
		res.Int(FlagDBMaxOpenConns, config.DefaultDBMaxOpenConns, "db max open connections")
		res.Int(FlagDBMaxIdleConns, config.DefaultDBMaxIdleConns, "db max idle connections")
		res.Duration(FlagDBMaxIdleLifetime, config.DefaultDBConnMaxIdleLifetime, "db max idle connection lifetime")
		res.Duration(FlagDBConnTimeout, config.DefaultDBConnTimeout, "db connection timeout)")

		// Log
		res.String(FlagLogLevel, config.DefaultLogLevel, "log level")
		res.String(FlagLogFormat, config.DefaultLogFormat, "log format")

		// Telemetry
		res.Bool(FlagTelemetryEnabled, config.DefaultTelemetryEnabled, "telemetry enabled")
		res.String(FlagTelemetryServiceName, "", "telemetry service name")
		res.String(FlagTelemetryExporterEndpoint, config.DefaultTelemetryExporterEndpoint, "telemetry exporter endpoint")
		res.Float64(FlagTelemetrySampleRate, config.DefaultTelemetrySampleRate, "telemetry sample rate")
		res.Duration(FlagTelemetryTimeout, config.DefaultTelemetryTimeout, "telemetry timeout")
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
		v.BindPFlag(keyAppEnv, flags.Lookup(FlagAppEnv)),
		v.BindPFlag(keyAppMaxListLimit, flags.Lookup(FlagAppMaxListLimit)),
		v.BindPFlag(keyAppTokenIssuer, flags.Lookup(FlagAppTokenIssuer)),
		v.BindPFlag(keyAppCipherKey, flags.Lookup(FlagAppCipherKey)),
		v.BindPFlag(keyAppAcceptTokenIssuers, flags.Lookup(FlagAppAcceptTokenIssuers)),
		// Auth
		v.BindPFlag(config.KeyAuthJWTSecret, flags.Lookup(FlagAuthJWTSecret)),
		v.BindPFlag(config.KeyAuthJWTSigningMethod, flags.Lookup(FlagAuthJWTSigningMethod)),
		v.BindPFlag(config.KeyAuthAccessTokenTTL, flags.Lookup(FlagAuthAccessTokenTTL)),
		v.BindPFlag(config.KeyAuthRefreshTokenTTL, flags.Lookup(FlagAuthRefreshTokenTTL)),
		v.BindPFlag(config.KeyAuthRSAPrivateKeyPath, flags.Lookup(FlagAuthRSAPrivateKeyPath)),
		v.BindPFlag(config.KeyAuthMasterPasswordSalt, flags.Lookup(FlagAuthMasterPasswordSalt)),
		// HTTP
		v.BindPFlag(config.KeyHTTPAddress, flags.Lookup(FlagHTTPAddress)),
		v.BindPFlag(config.KeyHTTPReadTimeout, flags.Lookup(FlagHTTPReadTimeout)),
		v.BindPFlag(config.KeyHTTPWriteTimeout, flags.Lookup(FlagHTTPWriteTimeout)),
		v.BindPFlag(config.KeyHTTPIdleTimeout, flags.Lookup(FlagHTTPIdleTimeout)),
		v.BindPFlag(config.KeyHTTPShutdownTimeout, flags.Lookup(FlagHTTPShutdownTimeout)),
		v.BindPFlag(config.KeyHTTPPrivateKeyPath, flags.Lookup(FlagHTTPPrivateKeyPath)),
		v.BindPFlag(config.KeyHTTPCertificatePath, flags.Lookup(FlagHTTPCertificatePath)),
		v.BindPFlag(config.KeyHTTPSecure, flags.Lookup(FlagHTTPSecure)),
		v.BindPFlag(config.KeyHTTPMaxRequestBodySize, flags.Lookup(FlagHTTPMaxRequestBodySize)),
		// gRPC
		v.BindPFlag(config.KeyGRPCAddress, flags.Lookup(FlagGRPCAddress)),
		v.BindPFlag(config.KeyGRPCMaxConnIdle, flags.Lookup(FlagGRPCMaxConnIdle)),
		v.BindPFlag(config.KeyGRPCMaxConnAge, flags.Lookup(FlagGRPCMaxConnAge)),
		v.BindPFlag(config.KeyGRPCMaxConnAgeGrace, flags.Lookup(FlagGRPCMaxConnAgeGrace)),
		v.BindPFlag(config.KeyGRPCTimeout, flags.Lookup(FlagGRPCTimeout)),
		v.BindPFlag(config.KeyGRPCKeepAliveTime, flags.Lookup(FlagGRPCKeepAliveTime)),
		v.BindPFlag(config.KeyGRPCKeepAliveTimeout, flags.Lookup(FlagGRPCKeepAliveTimeout)),
		v.BindPFlag(config.KeyGRPCShutdownTimeout, flags.Lookup(FlagGRPCShutdownTimeout)),
		// Log
		v.BindPFlag(config.KeyLogLevel, flags.Lookup(FlagLogLevel)),
		v.BindPFlag(config.KeyLogFormat, flags.Lookup(FlagLogFormat)),
		// DB
		v.BindPFlag(config.KeyDBDriver, flags.Lookup(FlagDBDriver)),
		v.BindPFlag(config.KeyDBDSN, flags.Lookup(FlagDBDSN)),
		v.BindPFlag(config.KeyDBMaxOpenConns, flags.Lookup(FlagDBMaxOpenConns)),
		v.BindPFlag(config.KeyDBMaxIdleConns, flags.Lookup(FlagDBMaxIdleConns)),
		v.BindPFlag(config.KeyDBConnMaxIdleLifetime, flags.Lookup(FlagDBMaxIdleLifetime)),
		v.BindPFlag(config.KeyDBConnTimeout, flags.Lookup(FlagDBConnTimeout)),
		// Telemetry
		v.BindPFlag(config.KeyTelemetryEnabled, flags.Lookup(FlagTelemetryEnabled)),
		v.BindPFlag(config.KeyTelemetryExporterEndpoint, flags.Lookup(FlagTelemetryExporterEndpoint)),
		v.BindPFlag(config.KeyTelemetryServiceName, flags.Lookup(FlagTelemetryServiceName)),
		v.BindPFlag(config.KeyTelemetrySampleRate, flags.Lookup(FlagTelemetrySampleRate)),
		v.BindPFlag(config.KeyTelemetryTimeout, flags.Lookup(FlagTelemetryTimeout)),
	)
	if err != nil {
		return errs.NewConfigError("bind flags with keys", err)
	}

	return nil
}
