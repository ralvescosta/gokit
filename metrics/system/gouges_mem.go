// Copyright (c) 2023, The GoKit Authors
// MIT License
// All rights reserved.

package system

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

// NewMemGauges creates a new memory metrics collector that monitors various aspects
// of the Go runtime memory usage and garbage collection.
//
// Parameters:
//   - meter: The OpenTelemetry meter used to create gauge instruments.
//
// Returns:
//   - A BasicGauges implementation for memory metrics collection.
//   - An error if any gauge creation fails.
func NewMemGauges(meter metric.Meter) (BasicGauges, error) {
	ggSysBytes, err := meter.Int64ObservableGauge("go_memstats_sys_bytes", metric.WithDescription("Number of bytes obtained from system."))
	if err != nil {
		return nil, err
	}

	ggAllocBytesTotal, err := meter.Int64ObservableGauge("go_memstats_alloc_bytes_total", metric.WithDescription("Total number of bytes allocated, even if freed."))
	if err != nil {
		return nil, err
	}

	ggHeapAllocBytes, err := meter.Int64ObservableGauge("go_memstats_heap_alloc_bytes", metric.WithDescription("Number of heap bytes allocated and still in use."))
	if err != nil {
		return nil, err
	}

	ggFreesTotal, err := meter.Int64ObservableGauge("go_memstats_frees_total", metric.WithDescription("Total number of frees."))
	if err != nil {
		return nil, err
	}

	ggGcSysBytes, err := meter.Int64ObservableGauge("go_memstats_gc_sys_bytes", metric.WithDescription("Number of bytes used for garbage collection system metadata."))
	if err != nil {
		return nil, err
	}

	ggHeapIdleBytes, err := meter.Int64ObservableGauge("go_memstats_heap_idle_bytes", metric.WithDescription("Number of heap bytes waiting to be used."))
	if err != nil {
		return nil, err
	}

	ggInuseBytes, err := meter.Int64ObservableGauge("go_memstats_heap_inuse_bytes", metric.WithDescription("Number of heap bytes that are in use."))
	if err != nil {
		return nil, err
	}

	ggHeapObjects, err := meter.Int64ObservableGauge("go_memstats_heap_objects", metric.WithDescription("Number of allocated objects."))
	if err != nil {
		return nil, err
	}

	ggHeapReleasedBytes, err := meter.Int64ObservableGauge("go_memstats_heap_released_bytes", metric.WithDescription("Number of heap bytes released to OS."))
	if err != nil {
		return nil, err
	}

	ggHeapSysBytes, err := meter.Int64ObservableGauge("go_memstats_heap_sys_bytes", metric.WithDescription("Number of heap bytes obtained from system."))
	if err != nil {
		return nil, err
	}

	ggLastGcTimeSeconds, err := meter.Int64ObservableGauge("go_memstats_last_gc_time_seconds", metric.WithDescription("Number of seconds since 1970 of last garbage collection."))
	if err != nil {
		return nil, err
	}

	ggLookupsTotal, err := meter.Int64ObservableGauge("go_memstats_lookups_total", metric.WithDescription("Total number of pointer lookups."))
	if err != nil {
		return nil, err
	}

	ggMallocsTotal, err := meter.Int64ObservableGauge("go_memstats_mallocs_total", metric.WithDescription("Total number of mallocs."))
	if err != nil {
		return nil, err
	}

	ggMCacheInuseBytes, err := meter.Int64ObservableGauge("go_memstats_mcache_inuse_bytes", metric.WithDescription("Number of bytes in use by mcache structures."))
	if err != nil {
		return nil, err
	}

	ggMCacheSysBytes, err := meter.Int64ObservableGauge("go_memstats_mcache_sys_bytes", metric.WithDescription("Number of bytes used for mcache structures obtained from system."))
	if err != nil {
		return nil, err
	}

	ggMspanInuseBytes, err := meter.Int64ObservableGauge("go_memstats_mspan_inuse_bytes", metric.WithDescription("Number of bytes in use by mspan structures."))
	if err != nil {
		return nil, err
	}

	ggMspanSysBytes, err := meter.Int64ObservableGauge("go_memstats_mspan_sys_bytes", metric.WithDescription("Number of bytes used for mspan structures obtained from system."))
	if err != nil {
		return nil, err
	}

	ggNextGcBytes, err := meter.Int64ObservableGauge("go_memstats_next_gc_bytes", metric.WithDescription("Number of heap bytes when next garbage collection will take place."))
	if err != nil {
		return nil, err
	}

	ggOtherSysBytes, err := meter.Int64ObservableGauge("go_memstats_other_sys_bytes", metric.WithDescription("Number of bytes used for other system allocations."))
	if err != nil {
		return nil, err
	}

	ggStackInuseBytes, err := meter.Int64ObservableGauge("go_memstats_stack_inuse_bytes", metric.WithDescription("Number of bytes in use by the stack allocator."))
	if err != nil {
		return nil, err
	}

	ggGcCompletedCycle, err := meter.Int64ObservableGauge("go_memstats_gc_completed_cycle", metric.WithDescription("Number of GC cycle completed."))
	if err != nil {
		return nil, err
	}

	ggGcPauseTotal, err := meter.Int64ObservableGauge("go_memstats_gc_pause_total", metric.WithDescription("Number of GC-stop-the-world caused in Nanosecond."))
	if err != nil {
		return nil, err
	}

	return &memGauges{
		ggSysBytes,
		ggAllocBytesTotal,
		ggHeapAllocBytes,
		ggFreesTotal,
		ggGcSysBytes,
		ggHeapIdleBytes,
		ggInuseBytes,
		ggHeapObjects,
		ggHeapReleasedBytes,
		ggHeapSysBytes,
		ggLastGcTimeSeconds,
		ggLookupsTotal,
		ggMallocsTotal,
		ggMCacheInuseBytes,
		ggMCacheSysBytes,
		ggMspanInuseBytes,
		ggMspanSysBytes,
		ggNextGcBytes,
		ggOtherSysBytes,
		ggStackInuseBytes,
		ggGcCompletedCycle,
		ggGcPauseTotal,
	}, nil
}

// Collect registers callbacks for memory metrics collection.
// It reads memory statistics from the Go runtime and reports them through the
// observable gauges.
//
// Parameters:
//   - meter: The OpenTelemetry meter used to register callbacks.
func (m *memGauges) Collect(meter metric.Meter) {

	cb := func(_ context.Context, observer metric.Observer) error {
		var stats runtime.MemStats
		runtime.ReadMemStats(&stats)

		observer.ObserveInt64(m.ggSysBytes, int64(stats.Sys))
		observer.ObserveInt64(m.ggAllocBytesTotal, int64(stats.TotalAlloc))
		observer.ObserveInt64(m.ggHeapAllocBytes, int64(stats.HeapAlloc))
		observer.ObserveInt64(m.ggFreesTotal, int64(stats.Frees))
		observer.ObserveInt64(m.ggGcSysBytes, int64(stats.GCSys))
		observer.ObserveInt64(m.ggHeapIdleBytes, int64(stats.HeapIdle))
		observer.ObserveInt64(m.ggInuseBytes, int64(stats.HeapInuse))
		observer.ObserveInt64(m.ggHeapObjects, int64(stats.HeapObjects))
		observer.ObserveInt64(m.ggHeapReleasedBytes, int64(stats.HeapReleased))
		observer.ObserveInt64(m.ggHeapSysBytes, int64(stats.HeapSys))
		observer.ObserveInt64(m.ggLastGcTimeSeconds, int64(stats.LastGC))
		observer.ObserveInt64(m.ggLookupsTotal, int64(stats.Lookups))
		observer.ObserveInt64(m.ggMallocsTotal, int64(stats.Mallocs))
		observer.ObserveInt64(m.ggMCacheInuseBytes, int64(stats.MCacheInuse))
		observer.ObserveInt64(m.ggMCacheSysBytes, int64(stats.MCacheSys))
		observer.ObserveInt64(m.ggMspanInuseBytes, int64(stats.MSpanInuse))
		observer.ObserveInt64(m.ggMspanSysBytes, int64(stats.MSpanSys))
		observer.ObserveInt64(m.ggNextGcBytes, int64(stats.NextGC))
		observer.ObserveInt64(m.ggOtherSysBytes, int64(stats.OtherSys))
		observer.ObserveInt64(m.ggStackInuseBytes, int64(stats.StackSys))
		observer.ObserveInt64(m.ggGcCompletedCycle, int64(stats.NumGC))
		observer.ObserveInt64(m.ggGcPauseTotal, int64(stats.PauseTotalNs))

		return nil
	}

	_, _ = meter.RegisterCallback(cb)
}
