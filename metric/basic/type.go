package basic

import (
	"context"

	"go.opentelemetry.io/otel/metric/instrument/asyncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/asyncint64"
)

type (
	BasicGauges interface {
		Collect(ctx context.Context)
	}

	memGauges struct {
		ggSysBytes          asyncfloat64.Gauge
		ggAllocBytesTotal   asyncfloat64.Gauge
		ggHeapAllocBytes    asyncfloat64.Gauge
		ggFreesTotal        asyncfloat64.Gauge
		ggGcSysBytes        asyncfloat64.Gauge
		ggHeapIdleBytes     asyncfloat64.Gauge
		ggInuseBytes        asyncfloat64.Gauge
		ggHeapObjects       asyncfloat64.Gauge
		ggHeapReleasedBytes asyncfloat64.Gauge
		ggHeapSysBytes      asyncfloat64.Gauge
		ggLastGcTimeSeconds asyncfloat64.Gauge
		ggLookupsTotal      asyncfloat64.Gauge
		ggMallocsTotal      asyncfloat64.Gauge
		ggMCacheInuseBytes  asyncfloat64.Gauge
		ggMCacheSysBytes    asyncfloat64.Gauge
		ggMspanInuseBytes   asyncfloat64.Gauge
		ggMspanSysBytes     asyncfloat64.Gauge
		ggNextGcBytes       asyncfloat64.Gauge
		ggOtherSysBytes     asyncfloat64.Gauge
		ggStackInuseBytes   asyncfloat64.Gauge
		ggGcCompletedCycle  asyncfloat64.Gauge
		ggGcPauseTotal      asyncfloat64.Gauge
	}

	sysGauges struct {
		ggThreads   asyncint64.Gauge
		ggCgo       asyncint64.Gauge
		ggGRoutines asyncint64.Gauge
	}
)
