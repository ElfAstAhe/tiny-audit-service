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

// login attempts common
const (
	// FlagLoginAttemptsReceiverKind - receiver kind, accepted values: amqp, kafka
	FlagLoginAttemptsReceiverKind string = "login-attempts-receiver-kind"
)

// login attempts worker
const (
	FlagLoginAttemptsWorkerStartInterval      string = "login-attempts-worker-start-interval"
	FlagLoginAttemptsWorkerScheduleInterval   string = "login-attempts-worker-schedule-interval"
	FlagLoginAttemptsWorkerWorkerCount        string = "login-attempts-worker-worker-count"
	FlagLoginAttemptsWorkerDataCapacity       string = "login-attempts-worker-data-capacity"
	FlagLoginAttemptsWorkerCompleteProcessing string = "login-attempts-worker-complete-processing"
	FlagLoginAttemptsWorkerShutdownTimeout    string = "login-attempts-worker-shutdown-timeout"
	FlagLoginAttemptsWorkerBatchSize          string = "login-attempts-worker-batch-size"
	FlagLoginAttemptsWorkerBatchReadTimeout   string = "login-attempts-worker-batch-read-timeout"
	FlagLoginAttemptsWorkerAcknowledgeTimeout string = "login-attempts-worker-ack-timeout"
)

// amqp login attempts receiver
const (
	FlagLoginAttemptsReceiverAMQPConfigTargetName      string = "login-attempts-receiver-amqp-target-name"
	FlagLoginAttemptsReceiverAMQPConfigConnectTimeout  string = "login-attempts-receiver-amqp-connect-timeout"
	FlagLoginAttemptsReceiverAMQPConfigShutdownTimeout string = "login-attempts-receiver-amqp-shutdown-timeout"
	FLagLoginAttemptsReceiverAMQPConfigPrefetchCredit  string = "login-attempts-receiver-amqp-prefetch-credit"
)

// kafka login attempts receiver
const (
	FlagLoginAttemptsReceiverKafkaConfigBrokers           string = "login-attempts-receiver-kafka-brokers"
	FlagLoginAttemptsReceiverKafkaConfigTargetName        string = "login-attempts-receiver-kafka-target-name"
	FlagLoginAttemptsReceiverKafkaConfigGroupID           string = "login-attempts-receiver-kafka-group-id"
	FlagLoginAttemptsReceiverKafkaConfigPartition         string = "login-attempts-receiver-kafka-partition"
	FlagLoginAttemptsReceiverKafkaConfigConnectTimeout    string = "login-attempts-receiver-kafka-connect-timeout"
	FlagLoginAttemptsReceiverKafkaConfigShutdownTimeout   string = "login-attempts-receiver-kafka-shutdown-timeout"
	FlagLoginAttemptsReceiverKafkaConfigMinBytes          string = "login-attempts-receiver-kafka-min-bytes"
	FlagLoginAttemptsReceiverKafkaConfigMaxBytes          string = "login-attempts-receiver-kafka-max-bytes"
	FlagLoginAttemptsReceiverKafkaConfigMaxWait           string = "login-attempts-receiver-kafka-max-wait"
	FlagLoginAttemptsReceiverKafkaConfigUsername          string = "login-attempts-receiver-kafka-username"
	FlagLoginAttemptsReceiverKafkaConfigPassword          string = "login-attempts-receiver-kafka-password"
	FlagLoginAttemptsReceiverKafkaConfigHeartbeatInterval string = "login-attempts-receiver-kafka-heartbeat-interval"
	FlagLoginAttemptsReceiverKafkaConfigSessionTimeout    string = "login-attempts-receiver-kafka-session-timeout"
	FlagLoginAttemptsReceiverKafkaConfigRebalanceTimeout  string = "login-attempts-receiver-kafka-rebalance-timeout"
	FlagLoginAttemptsReceiverKafkaConfigReadTimeout       string = "login-attempts-receiver-kafka-read-timeout"
	FlagLoginAttemptsReceiverKafkaConfigMaxAttempts       string = "login-attempts-receiver-kafka-max-attempts"
	FlagLoginAttemptsReceiverKafkaConfigQueueCapacity     string = "login-attempts-receiver-kafka-queue-capacity"
	FlagLoginAttemptsReceiverKafkaConfigStartOffset       string = "login-attempts-receiver-kafka-start-offset"
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

// login attempts worker
const (
	keyLoginAttemptsWorkerStartInterval      string = "login_attempts_receiver.worker_start_interval"
	keyLoginAttemptsWorkerScheduleInterval   string = "login_attempts_receiver.worker_schedule_interval"
	keyLoginAttemptsWorkerWorkerCount        string = "login_attempts_receiver.worker_worker_count"
	keyLoginAttemptsWorkerDataCapacity       string = "login_attempts_receiver.worker_data_capacity"
	keyLoginAttemptsWorkerCompleteProcessing string = "login_attempts_receiver.worker_complete_processing"
	keyLoginAttemptsWorkerShutdownTimeout    string = "login_attempts_receiver.worker_shutdown_timeout"
	keyLoginAttemptsWorkerBatchSize          string = "login_attempts_receiver.worker_batch_size"
	keyLoginAttemptsWorkerBatchReadTimeout   string = "login_attempts_receiver.worker_batch_read_timeout"
	keyLoginAttemptsWorkerAcknowledgeTimeout string = "login_attempts_receiver.worker_ack_timeout"
)

// login attempts common
const (
	keyLoginAttemptsReceiverKind string = "login_attempts_receiver.receiver_kind"
)

// amqp login attempts receiver
const (
	keyLoginAttemptsReceiverAMQPConfigTargetName      string = "login_attempts_receiver.amqp_config.target_name"
	keyLoginAttemptsReceiverAMQPConfigConnectTimeout  string = "login_attempts_receiver.amqp_config.connect_timeout"
	keyLoginAttemptsReceiverAMQPConfigShutdownTimeout string = "login_attempts_receiver.amqp_config.shutdown_timeout"
	keyLoginAttemptsReceiverAMQPConfigPrefetchCredit  string = "login_attempts_receiver.amqp_config.prefetch_credit"
)

const (
	keyLoginAttemptsReceiverKafkaConfigBrokers           string = "login_attempts_receiver.kafka_config.brokers"
	keyLoginAttemptsReceiverKafkaConfigTargetName        string = "login_attempts_receiver.kafka_config.target_name"
	keyLoginAttemptsReceiverKafkaConfigGroupID           string = "login_attempts_receiver.kafka_config.group_id"
	keyLoginAttemptsReceiverKafkaConfigPartition         string = "login_attempts_receiver.kafka_config.partition"
	keyLoginAttemptsReceiverKafkaConfigConnectTimeout    string = "login_attempts_receiver.kafka_config.connect_timeout"
	keyLoginAttemptsReceiverKafkaConfigShutdownTimeout   string = "login_attempts_receiver.kafka_config.shutdown_timeout"
	keyLoginAttemptsReceiverKafkaConfigMinBytes          string = "login_attempts_receiver.kafka_config.min_bytes"
	keyLoginAttemptsReceiverKafkaConfigMaxBytes          string = "login_attempts_receiver.kafka_config.max_bytes"
	keyLoginAttemptsReceiverKafkaConfigMaxWait           string = "login_attempts_receiver.kafka_config.max_wait"
	keyLoginAttemptsReceiverKafkaConfigUsername          string = "login_attempts_receiver.kafka_config.username"
	keyLoginAttemptsReceiverKafkaConfigPassword          string = "login_attempts_receiver.kafka_config.password"
	keyLoginAttemptsReceiverKafkaConfigHeartbeatInterval string = "login_attempts_receiver.kafka_config.heartbeat_interval"
	keyLoginAttemptsReceiverKafkaConfigSessionTimeout    string = "login_attempts_receiver.kafka_config.session_timeout"
	keyLoginAttemptsReceiverKafkaConfigRebalanceTimeout  string = "login_attempts_receiver.kafka_config.rebalance_timeout"
	keyLoginAttemptsReceiverKafkaConfigReadTimeout       string = "login_attempts_receiver.kafka_config.read_timeout"
	keyLoginAttemptsReceiverKafkaConfigMaxAttempts       string = "login_attempts_receiver.kafka_config.max_attempts"
	keyLoginAttemptsReceiverKafkaConfigQueueCapacity     string = "login_attempts_receiver.kafka_config.queue_capacity"
	keyLoginAttemptsReceiverKafkaConfigStartOffset       string = "login_attempts_receiver.kafka_config.start_offset"
)
