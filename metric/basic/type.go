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
		ggSysBytes        asyncfloat64.Gauge
		ggAllocBytesTotal asyncfloat64.Gauge
	}

	sysGauges struct {
		ggThreads   asyncint64.Gauge
		ggCgo       asyncint64.Gauge
		ggGRoutines asyncint64.Gauge
	}
)
