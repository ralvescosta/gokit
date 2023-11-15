package keys

const (
	GO_ENV_KEY                 = "GO_ENV"
	LOG_LEVEL_ENV_KEY          = "LOG_LEVEL"
	LOG_PATH_ENV_KEY           = "LOG_PATH"
	APP_NAME_ENV_KEY           = "APP_NAME"
	USE_SECRET_MANAGER_ENV_KEY = "USE_SECRET_MANAGER"
	SECRET_KEY_ENV_KEY         = "SECRET_KEY"

	SQL_DB_HOST_ENV_KEY            = "SQL_DB_HOST"
	SQL_DB_PORT_ENV_KEY            = "SQL_DB_PORT"
	SQL_DB_USER_ENV_KEY            = "SQL_DB_USER"
	SQL_DB_PASSWORD_ENV_KEY        = "SQL_DB_PASSWORD"
	SQL_DB_NAME_ENV_KEY            = "SQL_DB_NAME"
	SQL_DB_SECONDS_TO_PING_ENV_KEY = "SQL_DB_SECONDS_TO_PING"

	RABBIT_HOST_ENV_KEY     = "RABBIT_HOST"
	RABBIT_PORT_ENV_KEY     = "RABBIT_PORT"
	RABBIT_USER_ENV_KEY     = "RABBIT_USER"
	RABBIT_PASSWORD_ENV_KEY = "RABBIT_PASSWORD"
	RABBIT_VHOST_ENV_KEY    = "RABBIT_VHOST"

	MQTT_PROTOCOL_ENV_KEY = "MQTT_PROTOCOL"
	MQTT_HOST_ENV_KEY     = "MQTT_HOST"
	MQTT_PORT_ENV_KEY     = "MQTT_PORT"
	MQTT_USER_ENV_KEY     = "MQTT_USER"
	MQTT_PASSWORD_ENV_KEY = "MQTT_PASSWORD"

	KAFKA_HOST_ENV_KEY     = "KAFKA_HOST"
	KAFKA_PORT_ENV_KEY     = "KAFKA_PORT"
	KAFKA_USER_ENV_KEY     = "KAFKA_USER"
	KAFKA_PASSWORD_ENV_KEY = "KAFKA_PASSWORD"

	DEFAULT_APP_NAME = "app"
	DEFAULT_LOG_PATH = "/logs/"

	TRACING_ENABLED_ENV_KEY       = "TRACING_ENABLED"
	METRICS_ENABLED_ENV_KEY       = "METRICS_ENABLE"
	TRACING_OTLP_ENDPOINT_ENV_KEY = "TRACING_OTLP_ENDPOINT"
	TRACING_OTLP_API_KEY_ENV_KEY  = "TRACING_OTLP_API_KEY"

	METRICS_OTLP_ENDPOINT_ENV_KEY = "METRICS_OTLP_ENDPOINT"
	METRICS_OTLP_API_KEY_ENV_KEY  = "METRICS_OTLP_API_KEY"
	METRICS_KIND_ENV_KEY          = "METRICS_KIND"

	JAEGER_SERVICE_NAME_KEY       = "JAEGER_SERVICE_NAME"
	JAEGER_AGENT_HOST_KEY         = "JAEGER_AGENT_HOST"
	JAEGER_SAMPLER_TYPE_KEY       = "JAEGER_SAMPLER_TYPE"
	JAEGER_SAMPLER_PARAM_KEY      = "JAEGER_SAMPLER_PARAM"
	JAEGER_REPORTER_LOG_SPANS_KEY = "JAEGER_REPORTER_LOG_SPANS"
	JAEGER_RPC_METRICS_KEY        = "JAEGER_RPC_METRICS"

	HTTP_PORT_ENV_KEY             = "HTTP_PORT"
	HTTP_HOST_ENV_KEY             = "HTTP_HOST"
	HTTP_ENABLE_PROFILING_ENV_KEY = "HTTP_ENABLE_PROFILING"

	IDENTITY_DOMAIN_ENV_KEY    = "IDENTITY_DOMAIN"
	IDENTITY_AUDIENCE_ENV_KEY  = "IDENTITY_AUDIENCE"
	IDENTITY_ISSUER_ENV_KEY    = "IDENTITY_ISSUER"
	IDENTITY_SIGNATURE_ENV_KEY = "IDENTITY_SIGNATURE"

	AUTH0_CLIENT_ID_ENV_KEY                = "AUTH0_CLIENT_ID"
	AUTH0_CLIENT_SECRET_ENV_KEY            = "AUTH0_CLIENT_SECRET"
	AUTH0_GRANT_TYPE_ENV_KEY               = "AUTH0_GRANT_TYPE"
	AUTH0_MILLISECONDS_BETWEEN_JWK_ENV_KEY = "AUTH0_MILLISECONDS_BETWEEN_JWK"
)
