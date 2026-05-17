package benchmark

import (
	"testing"
	"time"
)

// BenchmarkPredicateFalseWithSlowDep measures consumer-dispatch latency
// when a predicate returns false and a sibling input task sleeps.
// Per-iteration wall time is bounded by Wait() draining the slow
// producer; the interesting metric is consumer_ns/op (the time from
// flow start to consumer firing), reported as a custom metric.
func BenchmarkPredicateFalseWithSlowDep(b *testing.B) {
	slowDelay := 5 * time.Millisecond
	var totalConsumer time.Duration
	for i := 0; i < b.N; i++ {
		totalConsumer += PredicateFalseWithSlowDep(slowDelay)
	}
	b.ReportMetric(float64(totalConsumer.Nanoseconds())/float64(b.N), "consumer_ns/op")
}
