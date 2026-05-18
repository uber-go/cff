//go:build cff
// +build cff

package benchmark

import (
	"context"
	"time"

	"go.uber.org/cff"
)

const (
	_workDuration = 100 * time.Millisecond
	_concurrency  = 2
)

// work is a task function that takes a pre-specifed amount of time.
func work() int {
	time.Sleep(_workDuration)
	return 0
}

// pred is a predicate function that takes a pre-specifed amount of time.
func pred() bool {
	time.Sleep(_workDuration)
	return true
}

// PredicateCombined runs a cff.Flow that exercises the function of a prior
// version of cff.Predicate that ran the predicate function within the task
// that the predicate is an option for (i.e the state of the world before
// https://code.uberinternal.com/D5495165).
func PredicateCombined() float64 {
	var res float64
	cff.Flow(
		context.Background(),
		cff.Concurrency(_concurrency),
		cff.Results(&res),
		cff.Task(
			func() int {
				return work()
			},
		),
		cff.Task(
			func(num int) (f float64) {
				if !pred() {
					return
				}
				f = float64(work())
				return
			},
		),
	)
	return res
}

// PredicateSplit runs a cff.Flow to exercise the function of the current
// predicate optimization which decouples the predicate function from the task
// it is an option for (i.e the state of the world after
// https://code.uberinternal.com/D5495165).
func PredicateSplit() float64 {
	var res float64
	cff.Flow(
		context.Background(),
		cff.Concurrency(_concurrency),
		cff.Results(&res),
		cff.Task(
			func() int {
				return work()
			},
		),
		cff.Task(
			func(num int) float64 {
				return float64(work())
			},
			cff.Predicate(
				pred,
			),
		),
	)
	return res
}

// PredicateFalseWithSlowDep runs a cff.Flow where a predicate returns false
// while a sibling input task sleeps for slowDelay. Returns the elapsed
// wall-clock time from flow start to consumer dispatch, so the benchmark
// can report consumer_ns/op as a custom metric.
//
// Before scheduler-aware predicates, the consumer waited the full
// slowDelay before its skip gate fired. After, the consumer dispatches
// immediately on predicate=false, even though Wait() still drains the
// slow producer.
func PredicateFalseWithSlowDep(slowDelay time.Duration) time.Duration {
	type slowOut struct{}
	type fastOut struct{}
	type skippedOut struct{}

	var elapsed time.Duration
	start := time.Now()
	cff.Flow(
		context.Background(),
		cff.Concurrency(_concurrency),
		cff.Results(&elapsed),
		cff.Task(func() slowOut {
			time.Sleep(slowDelay)
			return slowOut{}
		}),
		cff.Task(func() fastOut { return fastOut{} }),
		cff.Task(
			func(slowOut) (skippedOut, error) { return skippedOut{}, nil },
			cff.Predicate(func(fastOut) bool { return false }),
			cff.FallbackWith(skippedOut{}),
		),
		cff.Task(func(skippedOut) time.Duration {
			return time.Since(start)
		}),
	)
	return elapsed
}
