// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

import "go.uber.org/zap"

type Configs struct {
	Logger *zap.Logger

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
