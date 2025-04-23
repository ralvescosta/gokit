// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

// Package system provides system metrics collection capabilities for monitoring
// memory usage, garbage collection, threads, and goroutines.
package system

import (
	"go.opentelemetry.io/otel/metric"
)

type (
	// BasicGauges defines an interface for metrics collectors that gather
	// system-level metrics using OpenTelemetry observable gauges.
	BasicGauges interface {
		// Collect registers callbacks for the metrics with the provided meter.
		Collect(meter metric.Meter)
	}

	// memGauges implements BasicGauges to collect memory-related metrics.
	// It contains observable gauges for various memory statistics including
	// heap allocation, garbage collection, and system memory usage.
	memGauges struct {
		// System memory metrics
		ggSysBytes          metric.Int64ObservableGauge
		ggAllocBytesTotal   metric.Int64ObservableGauge
		ggHeapAllocBytes    metric.Int64ObservableGauge
		ggFreesTotal        metric.Int64ObservableGauge
		ggGcSysBytes        metric.Int64ObservableGauge
		ggHeapIdleBytes     metric.Int64ObservableGauge
		ggInuseBytes        metric.Int64ObservableGauge
		ggHeapObjects       metric.Int64ObservableGauge
		ggHeapReleasedBytes metric.Int64ObservableGauge
		ggHeapSysBytes      metric.Int64ObservableGauge
		ggLastGcTimeSeconds metric.Int64ObservableGauge
		ggLookupsTotal      metric.Int64ObservableGauge
		ggMallocsTotal      metric.Int64ObservableGauge
		ggMCacheInuseBytes  metric.Int64ObservableGauge
		ggMCacheSysBytes    metric.Int64ObservableGauge
		ggMspanInuseBytes   metric.Int64ObservableGauge
		ggMspanSysBytes     metric.Int64ObservableGauge
		ggNextGcBytes       metric.Int64ObservableGauge
		ggOtherSysBytes     metric.Int64ObservableGauge
		ggStackInuseBytes   metric.Int64ObservableGauge
		ggGcCompletedCycle  metric.Int64ObservableGauge
		ggGcPauseTotal      metric.Int64ObservableGauge
	}

	// sysGauges implements BasicGauges to collect system-level metrics.
	// It contains observable gauges for OS threads, CGo calls, and goroutines.
	sysGauges struct {
		// OS and runtime metrics
		ggThreads   metric.Int64ObservableGauge
		ggCgo       metric.Int64ObservableGauge
		ggGRoutines metric.Int64ObservableGauge
	}
)
