package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadMQTTConfigs() (*configs.MQTTConfigs, error) {
	mqttConfigs := &configs.MQTTConfigs{}

	mqttConfigs.Protocol = os.Getenv(keys.MQTT_PROTOCOL_ENV_KEY)
	mqttConfigs.Host = os.Getenv(keys.MQTT_HOST_ENV_KEY)
	port, err := strconv.Atoi(os.Getenv(keys.MQTT_PORT_ENV_KEY))
	if err != nil {
		return nil, err
	}

	mqttConfigs.Port = port
	mqttConfigs.User = os.Getenv(keys.MQTT_USER_ENV_KEY)
	mqttConfigs.Password = os.Getenv(keys.MQTT_PASSWORD_ENV_KEY)

	return nil, nil
}
