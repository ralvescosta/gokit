package basic

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

func NewMemGauges(meter metric.Meter) (BasicGauges, error) {
	ggSysBytes, err := meter.AsyncFloat64().Gauge("go_memstats_sys_bytes", instrument.WithDescription("Number of bytes obtained from system."))
	if err != nil {
		return nil, err
	}

	ggAllocBytesTotal, err := meter.AsyncFloat64().Gauge("go_memstats_alloc_bytes_total", instrument.WithDescription("Total number of bytes allocated, even if freed."))
	if err != nil {
		return nil, err
	}

	ggHeapAllocBytes, err := meter.AsyncFloat64().Gauge("go_memstats_heap_alloc_bytes", instrument.WithDescription("Number of heap bytes allocated and still in use."))
	if err != nil {
		return nil, err
	}

	ggFreesTotal, err := meter.AsyncFloat64().Gauge("go_memstats_frees_total", instrument.WithDescription("Total number of frees."))
	if err != nil {
		return nil, err
	}

	ggGcSysBytes, err := meter.AsyncFloat64().Gauge("go_memstats_gc_sys_bytes", instrument.WithDescription("Number of bytes used for garbage collection system metadata."))
	if err != nil {
		return nil, err
	}

	ggHeapIdleBytes, err := meter.AsyncFloat64().Gauge("go_memstats_heap_idle_bytes", instrument.WithDescription("Number of heap bytes waiting to be used."))
	if err != nil {
		return nil, err
	}

	ggInuseBytes, err := meter.AsyncFloat64().Gauge("go_memstats_heap_inuse_bytes", instrument.WithDescription("Number of heap bytes that are in use."))
	if err != nil {
		return nil, err
	}

	ggHeapObjects, err := meter.AsyncFloat64().Gauge("go_memstats_heap_objects", instrument.WithDescription("Number of allocated objects."))
	if err != nil {
		return nil, err
	}

	ggHeapReleasedBytes, err := meter.AsyncFloat64().Gauge("go_memstats_heap_released_bytes", instrument.WithDescription("Number of heap bytes released to OS."))
	if err != nil {
		return nil, err
	}

	ggHeapSysBytes, err := meter.AsyncFloat64().Gauge("go_memstats_heap_sys_bytes", instrument.WithDescription("Number of heap bytes obtained from system."))
	if err != nil {
		return nil, err
	}

	ggLastGcTimeSeconds, err := meter.AsyncFloat64().Gauge("go_memstats_last_gc_time_seconds", instrument.WithDescription("Number of seconds since 1970 of last garbage collection."))
	if err != nil {
		return nil, err
	}

	ggLookupsTotal, err := meter.AsyncFloat64().Gauge("go_memstats_lookups_total", instrument.WithDescription("Total number of pointer lookups."))
	if err != nil {
		return nil, err
	}

	ggMallocsTotal, err := meter.AsyncFloat64().Gauge("go_memstats_mallocs_total", instrument.WithDescription("Total number of mallocs."))
	if err != nil {
		return nil, err
	}

	ggMCacheInuseBytes, err := meter.AsyncFloat64().Gauge("go_memstats_mcache_inuse_bytes", instrument.WithDescription("Number of bytes in use by mcache structures."))
	if err != nil {
		return nil, err
	}

	ggMCacheSysBytes, err := meter.AsyncFloat64().Gauge("go_memstats_mcache_sys_bytes", instrument.WithDescription("Number of bytes used for mcache structures obtained from system."))
	if err != nil {
		return nil, err
	}

	ggMspanInuseBytes, err := meter.AsyncFloat64().Gauge("go_memstats_mspan_inuse_bytes", instrument.WithDescription("Number of bytes in use by mspan structures."))
	if err != nil {
		return nil, err
	}

	ggMspanSysBytes, err := meter.AsyncFloat64().Gauge("go_memstats_mspan_sys_bytes", instrument.WithDescription("Number of bytes used for mspan structures obtained from system."))
	if err != nil {
		return nil, err
	}

	ggNextGcBytes, err := meter.AsyncFloat64().Gauge("go_memstats_next_gc_bytes", instrument.WithDescription("Number of heap bytes when next garbage collection will take place."))
	if err != nil {
		return nil, err
	}

	ggOtherSysBytes, err := meter.AsyncFloat64().Gauge("go_memstats_other_sys_bytes", instrument.WithDescription("Number of bytes used for other system allocations."))
	if err != nil {
		return nil, err
	}

	ggStackInuseBytes, err := meter.AsyncFloat64().Gauge("go_memstats_stack_inuse_bytes", instrument.WithDescription("Number of bytes in use by the stack allocator."))
	if err != nil {
		return nil, err
	}

	ggGcCompletedCycle, err := meter.AsyncFloat64().Gauge("go_memstats_gc_completed_cycle", instrument.WithDescription("Number of GC cycle completed."))
	if err != nil {
		return nil, err
	}

	ggGcPauseTotal, err := meter.AsyncFloat64().Gauge("go_memstats_gc_pause_total", instrument.WithDescription("Number of GC-stop-the-world caused in Nanosecond."))
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

func (m *memGauges) Collect(ctx context.Context) {
	var stats runtime.MemStats

	m.ggSysBytes.Observe(ctx, float64(stats.Sys))
	m.ggAllocBytesTotal.Observe(ctx, float64(stats.TotalAlloc))
	m.ggHeapAllocBytes.Observe(ctx, float64(stats.HeapAlloc))
	m.ggFreesTotal.Observe(ctx, float64(stats.Frees))
	m.ggGcSysBytes.Observe(ctx, float64(stats.GCSys))
	m.ggHeapIdleBytes.Observe(ctx, float64(stats.HeapIdle))
	m.ggInuseBytes.Observe(ctx, float64(stats.HeapInuse))
	m.ggHeapObjects.Observe(ctx, float64(stats.HeapObjects))
	m.ggHeapReleasedBytes.Observe(ctx, float64(stats.HeapReleased))
	m.ggHeapSysBytes.Observe(ctx, float64(stats.HeapSys))
	m.ggLastGcTimeSeconds.Observe(ctx, float64(stats.LastGC))
	m.ggLookupsTotal.Observe(ctx, float64(stats.Lookups))
	m.ggMallocsTotal.Observe(ctx, float64(stats.Mallocs))
	m.ggMCacheInuseBytes.Observe(ctx, float64(stats.MCacheInuse))
	m.ggMCacheSysBytes.Observe(ctx, float64(stats.MCacheSys))
	m.ggMspanInuseBytes.Observe(ctx, float64(stats.MSpanInuse))
	m.ggMspanSysBytes.Observe(ctx, float64(stats.MSpanSys))
	m.ggNextGcBytes.Observe(ctx, float64(stats.NextGC))
	m.ggOtherSysBytes.Observe(ctx, float64(stats.OtherSys))
	m.ggStackInuseBytes.Observe(ctx, float64(stats.StackSys))
	m.ggGcCompletedCycle.Observe(ctx, float64(stats.NumGC))
	m.ggGcPauseTotal.Observe(ctx, float64(stats.PauseTotalNs))
}
