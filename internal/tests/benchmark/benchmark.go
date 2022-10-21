//go:build cff
// +build cff

package benchmark

import (
	"context"
	"sync"

	"go.uber.org/cff"
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

// Baseline is a flow that has two concurrent tasks that do almost nothing, that is designed to try to measure
// the overhead incurred by cff.Flow
func Baseline() float64 {
	var res float64
	cff.Flow(
		context.Background(),
		cff.Results(&res),
		cff.Task(a),
		cff.Task(b),
		cff.Task(c),
	)
	return res
}

// BaselineNative is a re-implementation of the Simple flow that makes the most optimal use of Go synchronization primitives
// while still running the two tasks in parallel. It should serve as a baseline as comparison to the Simple function.
func BaselineNative() float64 {
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
