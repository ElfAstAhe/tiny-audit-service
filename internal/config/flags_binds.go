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
	)
	if err != nil {
		return errs.NewConfigError("bind flags with keys", err)
	}

	return nil
}
