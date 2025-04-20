// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package system

import (
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel"
)

func BasicMetricsCollector(logger logging.Logger) error {
	logger.Debug("configuring basic metrics...")

	meter := otel.Meter("github.com/ralvescosta/gokit/metric/basic")

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

	mem.Collect(meter)
	sys.Collect(meter)

	return nil
}
