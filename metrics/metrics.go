// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package metrics

import (
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
)

type (
	// Headers is a map of string key-value pairs used for metric exporter HTTP headers.
	Headers map[string]string

	// MetricExporterType defines the type of metrics exporter being used.
	MetricExporterType int32

	// OTLPCompression defines the compression type used for OpenTelemetry Protocol.
	OTLPCompression string

	// basicMetricsAttr provides common attributes needed by metrics exporters.
	basicMetricsAttr struct {
		// logger is used for logging metrics-related events and errors.
		logger logging.Logger

		// cfg contains application configuration values.
		cfg *configs.Configs

		// appName stores the application name for metrics identification.
		appName string
	}
)
