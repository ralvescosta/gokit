// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package internal

import (
	"os"

	"github.com/ralvescosta/gokit/configs"
	keys "github.com/ralvescosta/gokit/configs_builder/keys"
)

func ReadTracingConfigs() (*configs.TracingConfigs, error) {
	enabled := os.Getenv(keys.TracingEnabledEnvKey)

	configs := configs.TracingConfigs{}

	if enabled != "" || enabled == "true" {
		configs.Enabled = true
	}

	configs.OtlpEndpoint = os.Getenv(keys.TracingOtlpEndpointEnvKey)
	configs.OtlpAPIKey = os.Getenv(keys.TracingOtlpAPIKeyEnvKey)

	return &configs, nil
}
