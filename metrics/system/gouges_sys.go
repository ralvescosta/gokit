package system

import (
	"context"
	"runtime"

	"go.opentelemetry.io/otel/metric"
)

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

func (s *sysGauges) Collect(meter metric.Meter) {
	cb := func(_ context.Context, observer metric.Observer) error {
		observer.ObserveInt64(s.ggThreads, int64(runtime.NumCPU()))
		observer.ObserveInt64(s.ggCgo, int64(runtime.NumCgoCall()))
		observer.ObserveInt64(s.ggGRoutines, int64(runtime.NumGoroutine()))
		return nil
	}

	_, _ = meter.RegisterCallback(cb)
}
