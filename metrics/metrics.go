package metrics

import (
	"github.com/ralvescosta/gokit/env"
	"github.com/ralvescosta/gokit/logging"
)

type (
	Headers            map[string]string
	MetricExporterType int32
	OTLPCompression    string

	basicMetricBuilder struct {
		logger  logging.Logger
		cfg     *env.Configs
		appName string
	}
)
