package config

import (
	"time"
)

// app
const (
	defaultAppEnv       AppEnv = AppEnvDevelopment
	defaultAppNodeName  string = ApplicationName
	defaultMaxListLimit int    = 100
	defaultTokenIssuer  string = "tiny-auth-service"
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

// app
const (
	keyAppEnv                string = "app.env"
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
