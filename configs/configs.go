package configs

type Configs struct {
	Custom map[string]string

	AppConfigs      *AppConfigs
	HTTPConfigs     *HTTPConfigs
	OtelConfigs     *OtelConfigs
	SqlConfigs      *SqlConfigs
	Auth0Configs    *Auth0Configs
	MQTTConfigs     *MQTTConfigs
	RabbitMQConfigs *RabbitMQConfigs
	AWSConfigs      *AWSConfigs
	DynamoDBConfigs *DynamoDBConfigs
}
