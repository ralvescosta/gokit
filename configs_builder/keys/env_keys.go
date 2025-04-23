// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package keys provides constants for environment variable keys used
// throughout the configuration system. These keys are used to retrieve
// configuration values from environment variables and .env files.
package keys

// Environment variable keys for various configuration components
const (
	// Core application settings
	GoEnvKey               = "GO_ENV"             // Application environment (development, production, etc.)
	LogLevelEnvKey         = "LOG_LEVEL"          // Logging level (debug, info, warn, error, etc.)
	LogPathEnvKey          = "LOG_PATH"           // Path for log files
	AppNameEnvKey          = "APP_NAME"           // Application name
	UseSecretManagerEnvKey = "USE_SECRET_MANAGER" // Whether to use a secret manager for sensitive configs
	SecretKeyEnvKey        = "SECRET_KEY"         // Secret key for encryption/security

	// SQL database configuration
	SQLDbHostEnvKey          = "SQL_DB_HOST"            // SQL database host
	SQLDbPortEnvKey          = "SQL_DB_PORT"            // SQL database port
	SQLDbUserEnvKey          = "SQL_DB_USER"            // SQL database username
	SQLDbPasswordEnvKey      = "SQL_DB_PASSWORD"        // SQL database password
	SQLDbNameEnvKey          = "SQL_DB_NAME"            // SQL database name
	SQLDbSecondsToPingEnvKey = "SQL_DB_SECONDS_TO_PING" // Interval for database health check pings

	// RabbitMQ configuration
	RabbitSchemaEnvKey   = "RABBIT_SCHEMA"   // RabbitMQ connection schema (amqp, amqps)
	RabbitHostEnvKey     = "RABBIT_HOST"     // RabbitMQ host
	RabbitPortEnvKey     = "RABBIT_PORT"     // RabbitMQ port
	RabbitUserEnvKey     = "RABBIT_USER"     // RabbitMQ username
	RabbitPasswordEnvKey = "RABBIT_PASSWORD" // RabbitMQ password
	RabbitVHostEnvKey    = "RABBIT_VHOST"    // RabbitMQ virtual host

	// MQTT configuration
	MQTTProtocolEnvKey = "MQTT_PROTOCOL" // MQTT protocol (mqtt, mqtts)
	MQTTHostEnvKey     = "MQTT_HOST"     // MQTT broker host
	MQTTPortEnvKey     = "MQTT_PORT"     // MQTT broker port
	MQTTUserEnvKey     = "MQTT_USER"     // MQTT username
	MQTTPasswordEnvKey = "MQTT_PASSWORD" // MQTT password

	// Kafka configuration
	KafkaHostEnvKey             = "KAFKA_HOST"              // Kafka broker host
	KafkaPortEnvKey             = "KAFKA_PORT"              // Kafka broker port
	KafkaSecurityProtocolEnvKey = "KAFKA_SECURITY_PROTOCOL" // Kafka security protocol (plaintext, ssl, etc.)
	KafkaSASLMechanismsEnvKey   = "KAFKA_SASL_MECHANISMS"   // Kafka SASL mechanism (PLAIN, SCRAM, etc.)
	KafkaUserEnvKey             = "KAFKA_USER"              // Kafka username
	KafkaPasswordEnvKey         = "KAFKA_PASSWORD"          // Kafka password

	// Default values
	DefaultAppName = "app"    // Default application name if not specified
	DefaultLogPath = "/logs/" // Default log file path if not specified

	// Tracing and metrics configuration
	TracingEnabledEnvKey      = "TRACING_ENABLED"       // Whether tracing is enabled
	MetricsEnabledEnvKey      = "METRICS_ENABLE"        // Whether metrics collection is enabled
	TracingOtlpEndpointEnvKey = "TRACING_OTLP_ENDPOINT" // OpenTelemetry endpoint for tracing
	TracingOtlpAPIKeyEnvKey   = "TRACING_OTLP_API_KEY"  // API key for OpenTelemetry tracing

	MetricsOtlpEndpointEnvKey = "METRICS_OTLP_ENDPOINT" // OpenTelemetry endpoint for metrics
	MetricsOtlpAPIKeyEnvKey   = "METRICS_OTLP_API_KEY"  // API key for OpenTelemetry metrics
	MetricsKindEnvKey         = "METRICS_KIND"          // Kind of metrics system (prometheus, otlp)

	// Jaeger tracing configuration
	JaegerServiceNameKey      = "JAEGER_SERVICE_NAME"       // Service name for Jaeger tracing
	JaegerAgentHostKey        = "JAEGER_AGENT_HOST"         // Jaeger agent host
	JaegerSamplerTypeKey      = "JAEGER_SAMPLER_TYPE"       // Jaeger sampler type (const, probabilistic, etc.)
	JaegerSamplerParamKey     = "JAEGER_SAMPLER_PARAM"      // Jaeger sampler parameter
	JaegerReporterLogSpansKey = "JAEGER_REPORTER_LOG_SPANS" // Whether to log spans in Jaeger
	JaegerRPCMetricsKey       = "JAEGER_RPC_METRICS"        // Whether to collect RPC metrics

	// HTTP server configuration
	HTTPPortEnvKey            = "HTTP_PORT"             // HTTP server port
	HTTPHostEnvKey            = "HTTP_HOST"             // HTTP server host
	HTTPEnableProfilingEnvKey = "HTTP_ENABLE_PROFILING" // Whether to enable pprof profiling endpoints

	// Identity/Auth configuration
	IdentityClientIDEnvKey               = "IDENTITY_CLIENT_ID"                // OAuth client ID
	IdentityClientSecretEnvKey           = "IDENTITY_CLIENT_SECRET"            // OAuth client secret
	IdentityGrantTypeEnvKey              = "IDENTITY_GRANT_TYPE"               // OAuth grant type
	IdentityMillisecondsBetweenJwkEnvKey = "IDENTITY_MILLISECONDS_BETWEEN_JWK" // Interval for refreshing JWK
	IdentityDomainEnvKey                 = "IDENTITY_DOMAIN"                   // Identity provider domain
	IdentityAudienceEnvKey               = "IDENTITY_AUDIENCE"                 // JWT audience
	IdentityIssuerEnvKey                 = "IDENTITY_ISSUER"                   // JWT issuer
	IdentitySignatureEnvKey              = "IDENTITY_SIGNATURE"                // JWT signature algorithm
)
