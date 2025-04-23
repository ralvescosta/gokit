// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package system

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

// NewSysGauge creates a new system metrics collector that monitors
// OS threads, CGO calls, and active goroutines.
//
// Parameters:
//   - meter: The OpenTelemetry meter used to create gauge instruments.
//
// Returns:
//   - A BasicGauges implementation for system metrics collection.
//   - An error if any gauge creation fails.
func NewSysGauge(meter metric.Meter) (BasicGauges, error) {
	ggThreads, err := meter.Int64ObservableGauge("go_threads", metric.WithDescription("Number of OS threads created."))
	if err != nil {
		return nil, err
	}

	ggCgo, err := meter.Int64ObservableGauge("go_cgo", metric.WithDescription("umber of CGO."))
	if err != nil {
		return nil, err
	}

	ggGRoutines, err := meter.Int64ObservableGauge("go_goroutines", metric.WithDescription("Number of goroutines."))
	if err != nil {
		return nil, err
	}

	return &sysGauges{
		ggThreads, ggCgo, ggGRoutines,
	}, nil
}

// Collect registers callbacks for system metrics collection.
// It reads statistics from the Go runtime about CPU cores, CGO calls,
// and goroutines and reports them through the observable gauges.
//
// Parameters:
//   - meter: The OpenTelemetry meter used to register callbacks.
func (s *sysGauges) Collect(meter metric.Meter) {
	cb := func(_ context.Context, observer metric.Observer) error {
		observer.ObserveInt64(s.ggThreads, int64(runtime.NumCPU()))
		observer.ObserveInt64(s.ggCgo, int64(runtime.NumCgoCall()))
		observer.ObserveInt64(s.ggGRoutines, int64(runtime.NumGoroutine()))
		return nil
	}

	_, _ = meter.RegisterCallback(cb)
}
