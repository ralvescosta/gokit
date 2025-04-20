// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package configs

type TracingConfigs struct {
	Enabled bool

	OtlpEndpoint string
	OtlpAPIKey   string
}
