package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadMQTTConfigs() (*configs.MQTTConfigs, error) {
	mqttConfigs := &configs.MQTTConfigs{}

	mqttConfigs.Protocol = os.Getenv(keys.MQTTProtocolEnvKey)
	mqttConfigs.Host = os.Getenv(keys.MQTTHostEnvKey)

	portEnv := os.Getenv(keys.MQTTPortEnvKey)
	if portEnv == "" {
		portEnv = "1883"
	}

	port, err := strconv.Atoi(portEnv)
	if err != nil {
		return nil, err
	}

	mqttConfigs.Port = port
	mqttConfigs.User = os.Getenv(keys.MQTTUserEnvKey)
	mqttConfigs.Password = os.Getenv(keys.MQTTPasswordEnvKey)

	return mqttConfigs, nil
}
