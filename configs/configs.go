package configs

type Configs struct {
	Custom map[string]string

	AppConfigs      *AppConfigs
	HTTPConfigs     *HTTPConfigs
	MetricsConfigs  *MetricsConfigs
	TracingConfigs  *TracingConfigs
	SQLConfigs      *SQLConfigs
	IdentityConfigs *IdentityConfigs
	Auth0Configs    *Auth0Configs
	MQTTConfigs     *MQTTConfigs
	RabbitMQConfigs *RabbitMQConfigs
	KafkaConfigs    *KafkaConfigs
	AWSConfigs      *AWSConfigs
	DynamoDBConfigs *DynamoDBConfigs
}
