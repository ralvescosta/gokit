package basic

import (
	"context"
	"time"

	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel/metric/global"
)

func BasicMetricsCollector(logger logging.Logger) error {
	logger.Debug("configuring basic metrics...")

	meter := global.Meter("github.com/ralvescosta/gokit/metric/basic")

	//Memory stats
	mem, err := NewMemGauges(meter)
	if err != nil {
		return err
	}

	//sys
	sys, err := NewSysGauge(meter)
	if err != nil {
		return err
	}

	logger.Debug("basic metrics configured")

	for {
		time.Sleep(time.Second * 10)
		logger.Debug("collecting basic metrics...")
		ctx := context.Background()

		mem.Collect(ctx)
		sys.Collect(ctx)

		logger.Debug("basic metrics collected")
	}
}
