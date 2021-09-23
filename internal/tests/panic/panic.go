package panic

import (
	"context"

	"go.uber.org/zap"

	"go.uber.org/cff"
	"github.com/uber-go/tally"
)

// Panicker is exported to be used by tests.
type Panicker struct {
	Scope  tally.Scope
	Logger *zap.Logger
}

// FlowPanicsParallel runs tasks in parallel.
func (p *Panicker) FlowPanicsParallel() error {
	var b bool

	err := cff.Flow(
		context.Background(),
		cff.WithEmitter(cff.TallyEmitter(p.Scope)),
		cff.WithEmitter(cff.LogEmitter(p.Logger)),
		cff.InstrumentFlow("PanicParallel"),
		cff.Results(&b),
		cff.Task(
			func() string {
				panic("panic")
			},
			cff.Instrument("T1"),
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
		cff.WithEmitter(cff.TallyEmitter(p.Scope)),
		cff.WithEmitter(cff.LogEmitter(p.Logger)),
		cff.InstrumentFlow("FlowPanicsSerial"),
		cff.Task(
			func() string {
				panic("panic")
			},
			cff.Instrument("T2"),
		),
	)

	return err
}
