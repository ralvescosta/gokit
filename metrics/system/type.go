package system

import (
	"go.opentelemetry.io/otel/metric"
)

type (
	BasicGauges interface {
		Collect(meter metric.Meter)
	}

	memGauges struct {
		// ggSysBytes instrument.Int64ObservableGauge
		// ggAllocBytesTotal   asyncint64.Gauge
		// ggHeapAllocBytes    asyncint64.Gauge
		// ggFreesTotal        asyncint64.Gauge
		// ggGcSysBytes        asyncint64.Gauge
		// ggHeapIdleBytes     asyncint64.Gauge
		// ggInuseBytes        asyncint64.Gauge
		// ggHeapObjects       asyncint64.Gauge
		// ggHeapReleasedBytes asyncint64.Gauge
		// ggHeapSysBytes      asyncint64.Gauge
		// ggLastGcTimeSeconds asyncint64.Gauge
		// ggLookupsTotal      asyncint64.Gauge
		// ggMallocsTotal      asyncint64.Gauge
		// ggMCacheInuseBytes  asyncint64.Gauge
		// ggMCacheSysBytes    asyncint64.Gauge
		// ggMspanInuseBytes   asyncint64.Gauge
		// ggMspanSysBytes     asyncint64.Gauge
		// ggNextGcBytes       asyncint64.Gauge
		// ggOtherSysBytes     asyncint64.Gauge
		// ggStackInuseBytes   asyncint64.Gauge
		// ggGcCompletedCycle  asyncint64.Gauge
		// ggGcPauseTotal      asyncint64.Gauge
	}

	sysGauges struct {
		// ggThreads   asyncint64.Gauge
		// ggCgo       asyncint64.Gauge
		// ggGRoutines asyncint64.Gauge
	}
)
