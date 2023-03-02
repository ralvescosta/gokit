package system

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
)

type (
	BasicGauges interface {
		Collect(meter metric.Meter)
	}

	memGauges struct {
		ggSysBytes          instrument.Int64ObservableGauge
		ggAllocBytesTotal   instrument.Int64ObservableGauge
		ggHeapAllocBytes    instrument.Int64ObservableGauge
		ggFreesTotal        instrument.Int64ObservableGauge
		ggGcSysBytes        instrument.Int64ObservableGauge
		ggHeapIdleBytes     instrument.Int64ObservableGauge
		ggInuseBytes        instrument.Int64ObservableGauge
		ggHeapObjects       instrument.Int64ObservableGauge
		ggHeapReleasedBytes instrument.Int64ObservableGauge
		ggHeapSysBytes      instrument.Int64ObservableGauge
		ggLastGcTimeSeconds instrument.Int64ObservableGauge
		ggLookupsTotal      instrument.Int64ObservableGauge
		ggMallocsTotal      instrument.Int64ObservableGauge
		ggMCacheInuseBytes  instrument.Int64ObservableGauge
		ggMCacheSysBytes    instrument.Int64ObservableGauge
		ggMspanInuseBytes   instrument.Int64ObservableGauge
		ggMspanSysBytes     instrument.Int64ObservableGauge
		ggNextGcBytes       instrument.Int64ObservableGauge
		ggOtherSysBytes     instrument.Int64ObservableGauge
		ggStackInuseBytes   instrument.Int64ObservableGauge
		ggGcCompletedCycle  instrument.Int64ObservableGauge
		ggGcPauseTotal      instrument.Int64ObservableGauge
	}

	sysGauges struct {
		ggThreads   instrument.Int64ObservableGauge
		ggCgo       instrument.Int64ObservableGauge
		ggGRoutines instrument.Int64ObservableGauge
	}
)
