package system

import (
	"go.opentelemetry.io/otel/metric"
)

type (
	BasicGauges interface {
		Collect(meter metric.Meter)
	}

	memGauges struct {
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

	sysGauges struct {
		ggThreads   metric.Int64ObservableGauge
		ggCgo       metric.Int64ObservableGauge
		ggGRoutines metric.Int64ObservableGauge
	}
)
