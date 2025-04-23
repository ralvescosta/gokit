// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

// ReadMQTTConfigs retrieves MQTT connection configuration from environment variables.
// Loads protocol, host, port, credentials, and other settings needed to connect to an
// MQTT broker. Uses default values for optional parameters when not specified.
func ReadMQTTConfigs() (*configs.MQTTConfigs, error) {
	mqttConfigs := &configs.MQTTConfigs{}

	// Get MQTT protocol (mqtt, mqtts)
	mqttConfigs.Protocol = os.Getenv(keys.MQTTProtocolEnvKey)
	// Get MQTT broker host
	mqttConfigs.Host = os.Getenv(keys.MQTTHostEnvKey)

	// Get MQTT port with default fallback to 1883
	portEnv := os.Getenv(keys.MQTTPortEnvKey)
	if portEnv == "" {
		portEnv = "1883"
	}

	// Parse port string to integer
	port, err := strconv.Atoi(portEnv)
	if err != nil {
		return nil, err
	}

	mqttConfigs.Port = port
	// Get MQTT authentication credentials (optional)
	mqttConfigs.User = os.Getenv(keys.MQTTUserEnvKey)
	mqttConfigs.Password = os.Getenv(keys.MQTTPasswordEnvKey)

	return mqttConfigs, nil
}
