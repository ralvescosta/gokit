package keys

const (
	GoEnvKey               = "GO_ENV"
	LogLevelEnvKey         = "LOG_LEVEL"
	LogPathEnvKey          = "LOG_PATH"
	AppNameEnvKey          = "APP_NAME"
	UseSecretManagerEnvKey = "USE_SECRET_MANAGER"
	SecretKeyEnvKey        = "SECRET_KEY"

	SQLDbHostEnvKey          = "SQL_DB_HOST"
	SQLDbPortEnvKey          = "SQL_DB_PORT"
	SQLDbUserEnvKey          = "SQL_DB_USER"
	SQLDbPasswordEnvKey      = "SQL_DB_PASSWORD"
	SQLDbNameEnvKey          = "SQL_DB_NAME"
	SQLDbSecondsToPingEnvKey = "SQL_DB_SECONDS_TO_PING"

	RabbitHostEnvKey     = "RABBIT_HOST"
	RabbitPortEnvKey     = "RABBIT_PORT"
	RabbitUserEnvKey     = "RABBIT_USER"
	RabbitPasswordEnvKey = "RABBIT_PASSWORD"
	RabbitVHostEnvKey    = "RABBIT_VHOST"

	MQTTProtocolEnvKey = "MQTT_PROTOCOL"
	MQTTHostEnvKey     = "MQTT_HOST"
	MQTTPortEnvKey     = "MQTT_PORT"
	MQTTUserEnvKey     = "MQTT_USER"
	MQTTPasswordEnvKey = "MQTT_PASSWORD"

	KafkaHostEnvKey             = "KAFKA_HOST"
	KafkaPortEnvKey             = "KAFKA_PORT"
	KafkaSecurityProtocolEnvKey = "KAFKA_SECURITY_PROTOCOL"
	KafkaSASLMechanismsEnvKey   = "KAFKA_SASL_MECHANISMS"
	KafkaUserEnvKey             = "KAFKA_USER"
	KafkaPasswordEnvKey         = "KAFKA_PASSWORD"

	DefaultAppName = "app"
	DefaultLogPath = "/logs/"

	TracingEnabledEnvKey      = "TRACING_ENABLED"
	MetricsEnabledEnvKey      = "METRICS_ENABLE"
	TracingOtlpEndpointEnvKey = "TRACING_OTLP_ENDPOINT"
	TracingOtlpAPIKeyEnvKey   = "TRACING_OTLP_API_KEY"

	MetricsOtlpEndpointEnvKey = "METRICS_OTLP_ENDPOINT"
	MetricsOtlpAPIKeyEnvKey   = "METRICS_OTLP_API_KEY"
	MetricsKindEnvKey         = "METRICS_KIND"

	JaegerServiceNameKey      = "JAEGER_SERVICE_NAME"
	JaegerAgentHostKey        = "JAEGER_AGENT_HOST"
	JaegerSamplerTypeKey      = "JAEGER_SAMPLER_TYPE"
	JaegerSamplerParamKey     = "JAEGER_SAMPLER_PARAM"
	JaegerReporterLogSpansKey = "JAEGER_REPORTER_LOG_SPANS"
	JaegerRPCMetricsKey       = "JAEGER_RPC_METRICS"

	HTTPPortEnvKey            = "HTTP_PORT"
	HTTPHostEnvKey            = "HTTP_HOST"
	HTTPEnableProfilingEnvKey = "HTTP_ENABLE_PROFILING"

	IdentityClientIDEnvKey               = "IDENTITY_CLIENT_ID"
	IdentityClientSecretEnvKey           = "IDENTITY_CLIENT_SECRET"
	IdentityGrantTypeEnvKey              = "IDENTITY_GRANT_TYPE"
	IdentityMillisecondsBetweenJwkEnvKey = "IDENTITY_MILLISECONDS_BETWEEN_JWK"
	IdentityDomainEnvKey                 = "IDENTITY_DOMAIN"
	IdentityAudienceEnvKey               = "IDENTITY_AUDIENCE"
	IdentityIssuerEnvKey                 = "IDENTITY_ISSUER"
	IdentitySignatureEnvKey              = "IDENTITY_SIGNATURE"
)
