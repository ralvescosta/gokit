package metrics

import (
	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
)

type (
	Headers            map[string]string
	MetricExporterType int32
	OTLPCompression    string

	basicMetricsAttr struct {
		logger  logging.Logger
		cfg     *configs.Configs
		appName string
	}
)
