// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package system

import (
	"github.com/ralvescosta/gokit/logging"
	"go.opentelemetry.io/otel"
)

// BasicMetricsCollector initializes and configures basic system metrics collection.
// It sets up memory and system gauges and starts collecting metrics.
//
// Parameters:
//   - logger: A logger instance for logging metrics-related messages.
//
// Returns:
//   - An error if metrics collection could not be initialized.
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
