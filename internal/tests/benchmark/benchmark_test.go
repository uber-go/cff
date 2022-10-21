package benchmark

import (
	"testing"
)

// BenchmarkBaseline is a flow that has two concurrent tasks that do almost nothing, that is designed to try to measure
// the overhead incurred by cff.Flow
func BenchmarkBaseline(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Baseline()
	}
}

// BenchmarkBaselineNative is a re-implementation of the Simple flow that makes the most optimal use of Go synchronization primitives
// while still running the two tasks in parallel. It should serve as a baseline as comparison to the Simple function.
func BenchmarkBaselineNative(b *testing.B) {
	for i := 0; i < b.N; i++ {
		BaselineNative()
	}
}

// BenchmarkPredicateCombined emulates the state of cff.Flow before
// https://code.uberinternal.com/D5495165. This serves as the baseline
// implementation in comparison to BenchmarkPredicateSplit.
func BenchmarkPredicateCombined(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PredicateCombined()
	}
}

// BenchmarkPredicateSplit is a flow that deploys a cff.Task with a
// cff.Predicate option.
func BenchmarkPredicateSplit(b *testing.B) {
	for i := 0; i < b.N; i++ {
		PredicateSplit()
	}
}
