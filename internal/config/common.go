package config

// FlagConfig - файл конфигурации
const FlagConfig = "config-path"

// App config flags
const (
	FlagAppNodeName           string = "node-name"
	FlagAppMaxListLimit       string = "app-max-list-limit"
	FlagAppTokenIssuer        string = "app-token-issuer"
	FlagAppCipherKey          string = "app-cipher-key"
	FlagAppAcceptTokenIssuers string = "app-accept-token-issuers"
)

// auth tc config flags
const (
	FlagAuthTCStartInterval      string = "auth-tc-start-interval"
	FlagAuthTCScheduleInterval   string = "auth-tc-schedule-interval"
	FlagAuthTCWorkerCount        string = "auth-tc-worker-count"
	FlagAuthTCDataCapacity       string = "auth-tc-data-capacity"
	FlagAuthTCCompleteProcessing string = "auth-tc-complete-processing"
	FlagAuthTCShutdownTimeout    string = "auth-tc-shutdown-timeout"
	FlagAuthTCTailInterval       string = "auth-tc-tail-interval"
	FlagAuthTCTailCut            string = "auth-tc-tail-cut"
)

// data tc config flags
const (
	FlagDataTCStartInterval      string = "data-tc-start-interval"
	FlagDataTCScheduleInterval   string = "data-tc-schedule-interval"
	FlagDataTCWorkerCount        string = "data-tc-worker-count"
	FlagDataTCDataCapacity       string = "data-tc-data-capacity"
	FlagDataTCCompleteProcessing string = "data-tc-complete-processing"
	FlagDataTCShutdownTimeout    string = "data-tc-shutdown-timeout"
	FlagDataTCTailInterval       string = "data-tc-tail-interval"
	FlagDataTCTailCut            string = "data-tc-tail-cut"
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

// amqp connector
const (
	FlagAMQPConnectorURL             string = "amqp-connector-url"
	FlagAMQPConnectorUsername        string = "amqp-connector-username"
	FlagAMQPConnectorPassword        string = "amqp-connector-password"
	FlagAMQPConnectorConnectTimeout  string = "amqp-connector-connect-timeout"
	FlagAMQPConnectorWriteTimeout    string = "amqp-connector-write-timeout"
	FlagAMQPConnectorIdleTimeout     string = "amqp-connector-idle-timeout"
	FlagAMQPConnectorShutdownTimeout string = "amqp-connector-shutdown-timeout"
)

// amqp login attempts receiver
const (
	FlagLoginAttemptsReceiverTargetName      string = "login-attempts-receiver-target-name"
	FlagLoginAttemptsReceiverConnectTimeout  string = "login-attempts-receiver-connect-timeout"
	FlagLoginAttemptsReceiverShutdownTimeout string = "login-attempts-receiver-shutdown-timeout"
	FLagLoginAttemptsReceiverPrefetchCredit  string = "login-attempts-receiver-prefetch-credit"
	FlagLoginAttemptsStartInterval           string = "login-attempts-start-interval"
	FlagLoginAttemptsScheduleInterval        string = "login-attempts-schedule-interval"
	FlagLoginAttemptsWorkerCount             string = "login-attempts-worker-count"
	FlagLoginAttemptsDataCapacity            string = "login-attempts-data-capacity"
	FlagLoginAttemptsCompleteProcessing      string = "login-attempts-complete-processing"
	FlagLoginAttemptsShutdownTimeout         string = "login-attempts-shutdown-timeout"
	FlagLoginAttemptsBatchSize               string = "login-attempts-batch-size"
	FlagLoginAttemptsBatchReadTimeout        string = "login-attempts-batch-read-timeout"
	FlagLoginAttemptsAcknowledgeTimeout      string = "login-attempts-ack-timeout"
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

// amqp login attempts receiver
const (
	keyLoginAttemptsReceiverTargetName      string = "login_attempts_receiver.receiver_conf.target_name"
	keyLoginAttemptsReceiverConnectTimeout  string = "login_attempts_receiver.receiver_conf.connect_timeout"
	keyLoginAttemptsReceiverShutdownTimeout string = "login_attempts_receiver.receiver_conf.shutdown_timeout"
	keyLoginAttemptsReceiverPrefetchCredit  string = "login_attempts_receiver.receiver_conf.prefetch_credit"
	keyLoginAttemptsStartInterval           string = "login_attempts_receiver.start_interval"
	keyLoginAttemptsScheduleInterval        string = "login_attempts_receiver.schedule_interval"
	keyLoginAttemptsWorkerCount             string = "login_attempts_receiver.worker_count"
	keyLoginAttemptsDataCapacity            string = "login_attempts_receiver.data_capacity"
	keyLoginAttemptsCompleteProcessing      string = "login_attempts_receiver.complete_processing"
	keyLoginAttemptsShutdownTimeout         string = "login_attempts_receiver.shutdown_timeout"
	keyLoginAttemptsBatchSize               string = "login_attempts_receiver.batch_size"
	keyLoginAttemptsBatchReadTimeout        string = "login_attempts_receiver.batch_read_timeout"
	keyLoginAttemptsAcknowledgeTimeout      string = "login_attempts_receiver.ack_timeout"
)
