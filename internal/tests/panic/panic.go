//go:build cff
// +build cff

package panic

import (
	"context"

	"go.uber.org/cff"
)

// Panicker is exported to be used by tests.
type Panicker struct{}

// FlowPanicsParallel runs tasks in parallel.
func (p *Panicker) FlowPanicsParallel() error {
	var b bool

	err := cff.Flow(
		context.Background(),
		cff.Results(&b),
		cff.Task(
			func() string {
				panic("panic")
			},
		),
		// This task is necessary so that task 1 and 2 are run in parallel, which necessitates running them
		// in separate goroutines.
		cff.Task(
			func() int64 {
				return 0
			},
		),
		cff.Task(
			func(string, int64) bool {
				return true
			},
		),
	)

	return err
}

// FlowPanicsSerial runs a single flow.
func (p *Panicker) FlowPanicsSerial() error {
	var r string

	err := cff.Flow(
		context.Background(),
		cff.Results(&r),
		cff.Task(
			func() string {
				panic("panic")
			},
		),
	)

	return err
}
