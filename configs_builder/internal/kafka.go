package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadKafkaConfigs() (*configs.KafkaConfigs, error) {
	cfgs := &configs.KafkaConfigs{}

	cfgs.Host = os.Getenv(keys.KAFKA_HOST_ENV_KEY)
	if cfgs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.KAFKA_HOST_ENV_KEY)
	}

	port := os.Getenv(keys.KAFKA_PORT_ENV_KEY)
	if port == "" {
		return nil, errors.NewErrRequiredConfig(keys.KAFKA_PORT_ENV_KEY)
	}

	cfgs.Port, _ = strconv.Atoi(port)

	cfgs.SecurityProtocol = os.Getenv(keys.KAFKA_SECURITY_PROTOCOL_ENV_KEY)
	cfgs.SASLMechanisms = os.Getenv(keys.KAFKA_SASL_MECHANISMS_ENV_KEY)
	cfgs.User = os.Getenv(keys.KAFKA_USER_ENV_KEY)
	cfgs.Password = os.Getenv(keys.KAFKA_PASSWORD_ENV_KEY)

	return cfgs, nil
}
