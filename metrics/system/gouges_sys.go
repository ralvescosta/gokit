package system

import (
	"go.opentelemetry.io/otel/metric"
)

func NewSysGauge(meter metric.Meter) (BasicGauges, error) {
	// ggThreads, err := meter.AsyncInt64().Gauge("go_threads", instrument.WithDescription("Number of OS threads created."))
	// if err != nil {
	// 	return nil, err
	// }

	// ggCgo, err := meter.AsyncInt64().Gauge("go_cgo", instrument.WithDescription("umber of CGO."))
	// if err != nil {
	// 	return nil, err
	// }

	// ggGRoutines, err := meter.AsyncInt64().Gauge("go_goroutines", instrument.WithDescription("Number of goroutines."))
	// if err != nil {
	// 	return nil, err
	// }

	// return &sysGauges{
	// 	ggThreads, ggCgo, ggGRoutines,
	// }, nil
	return nil, nil
}

func (s *sysGauges) Collect(meter metric.Meter) {
	// meter.RegisterCallback([]instrument.Asynchronous{}, func(ctx context.Context) {
	// 	s.ggThreads.Observe(ctx, int64(runtime.NumCPU()))
	// 	s.ggCgo.Observe(ctx, runtime.NumCgoCall())
	// 	s.ggGRoutines.Observe(ctx, int64(runtime.NumGoroutine()))
	// })
}
