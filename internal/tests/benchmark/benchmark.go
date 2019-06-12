// +build cff

package benchmark

import (
	"context"
	"sync"

	"go.uber.org/cff"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

func a() int64 {
	return 1
}

func b() int {
	return 1
}

func c(i int64, j int) float64 {
	return float64(i + int64(j))
}

// Simple is a flow that has two concurrent tasks that do almost nothing, that is designed to try to measure
// the overhead incurred by cff.Flow
func Simple() float64 {
	var res float64
	cff.Flow(
		context.Background(),
		cff.Results(&res),
		cff.Tasks(
			a,
			b,
			c,
		),
	)
	return res
}

// SimpleNative is a re-implementation of the Simple flow that makes the most optimal use of Go synchronization primitives
// while still running the two tasks in parallel. It should serve as a baseline as comparison to the Simple function.
func SimpleNative() float64 {
	var aReturn int64
	var bReturn int
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		aReturn = a()
		wg.Done()
	}()
	go func() {
		bReturn = b()
		wg.Done()
	}()
	wg.Wait()
	return c(aReturn, bReturn)
}

// Metrics is the same flow as Simple but with instrumentation added.
func Metrics(logger *zap.Logger, scope tally.Scope) float64 {
	var res float64
	cff.Flow(
		context.Background(),
		cff.InstrumentFlow("Metrics"),
		cff.Metrics(scope),
		cff.Logger(logger),
		cff.Results(&res),
		cff.Task(
			a,
			cff.Instrument("a"),
		),
		cff.Task(
			b,
			cff.Instrument("b"),
		),
		cff.Task(
			c,
			cff.Instrument("c"),
		),
	)
	return res
}
