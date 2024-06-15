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

	cfgs.Host = os.Getenv(keys.KafkaHostEnvKey)
	if cfgs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.KafkaHostEnvKey)
	}

	port := os.Getenv(keys.KafkaPortEnvKey)
	if port == "" {
		return nil, errors.NewErrRequiredConfig(keys.KafkaPortEnvKey)
	}

	cfgs.Port, _ = strconv.Atoi(port)

	cfgs.SecurityProtocol = os.Getenv(keys.KafkaSecurityProtocolEnvKey)
	cfgs.SASLMechanisms = os.Getenv(keys.KafkaSASLMechanismsEnvKey)
	cfgs.User = os.Getenv(keys.KafkaUserEnvKey)
	cfgs.Password = os.Getenv(keys.KafkaPasswordEnvKey)

	return cfgs, nil
}
