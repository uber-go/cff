package benchmark

import (
	"testing"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// BenchmarkSimple is a flow that has two concurrent tasks that do almost nothing, that is designed to try to measure
// the overhead incurred by cff.Flow
func BenchmarkSimple(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Simple()
	}
}

// BenchmarkSimpleNative is a re-implementation of the Simple flow that makes the most optimal use of Go synchronization primitives
// while still running the two tasks in parallel. It should serve as a baseline as comparison to the Simple function.
func BenchmarkSimpleNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SimpleNative()
	}
}

type metricsTestFn func(*zap.Logger, tally.Scope) float64

// BenchmarkMetrics is the same flow as Simple but with instrumentation added.
func BenchmarkMetrics(b *testing.B) {
	logger := zap.NewNop()
	scope := tally.NoopScope

	metricsCases := []struct {
		name string
		fn   metricsTestFn
	}{
		{
			"Metrics", Metrics,
		},
		{
			"Metrics100", Metrics100,
		},
		{
			"Metrics500", Metrics500,
		},
		{
			"Metrics1000", Metrics1000,
		},
	}

	for _, metricsCase := range metricsCases {
		metricsCaseClosure := metricsCase
		b.Run(metricsCase.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				metricsCaseClosure.fn(logger, scope)
			}
		})
	}
}
