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

	return &memGauges{
		ggSysBytes, ggAllocBytesTotal,
	}, nil
}

func (m *memGauges) Collect(ctx context.Context) {
	var stats runtime.MemStats

	m.ggSysBytes.Observe(ctx, float64(stats.Sys))
	m.ggAllocBytesTotal.Observe(ctx, float64(stats.TotalAlloc))
}
