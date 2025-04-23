// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"
	"strconv"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/configs_builder/errors"
	"github.com/ralvescosta/gokit/configs_builder/keys"
)

// ReadKafkaConfigs retrieves Kafka connection configuration from environment variables.
// Validates required connection parameters and returns an error if any required
// configuration is missing.
func ReadKafkaConfigs() (*configs.KafkaConfigs, error) {
	cfgs := &configs.KafkaConfigs{}

	// Get and validate Kafka broker host
	cfgs.Host = os.Getenv(keys.KafkaHostEnvKey)
	if cfgs.Host == "" {
		return nil, errors.NewErrRequiredConfig(keys.KafkaHostEnvKey)
	}

	// Get and validate Kafka broker port
	port := os.Getenv(keys.KafkaPortEnvKey)
	if port == "" {
		return nil, errors.NewErrRequiredConfig(keys.KafkaPortEnvKey)
	}

	// Parse port string to integer
	cfgs.Port, _ = strconv.Atoi(port)

	// Get optional security settings
	cfgs.SecurityProtocol = os.Getenv(keys.KafkaSecurityProtocolEnvKey)
	cfgs.SASLMechanisms = os.Getenv(keys.KafkaSASLMechanismsEnvKey)
	cfgs.User = os.Getenv(keys.KafkaUserEnvKey)
	cfgs.Password = os.Getenv(keys.KafkaPasswordEnvKey)

	return cfgs, nil
}
