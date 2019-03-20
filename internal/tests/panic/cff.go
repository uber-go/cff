// +build cff

package panic

import (
	"context"
	"go.uber.org/zap"

	"go.uber.org/cff"
	"github.com/uber-go/tally"
)

type panicker struct {
	scope  tally.Scope
	logger *zap.Logger
}

func (p *panicker) FlowPanicsParallel() error {
	var b bool

	err := cff.Flow(
		context.Background(),
		cff.Metrics(p.scope),
		cff.Logger(p.logger),
		cff.InstrumentFlow("PanicParallel"),
		cff.Results(&b),
		cff.Task(
			func() string {
				panic("panic")
				return ""
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

func (p *panicker) FlowPanicsSerial() error {
	var r string

	err := cff.Flow(
		context.Background(),
		cff.Results(&r),
		cff.Metrics(p.scope),
		cff.Logger(p.logger),
		cff.InstrumentFlow("FlowPanicsSerial"),
		cff.Task(
			func() string {
				panic("panic")
				return ""
			},
			cff.Instrument("T1"),
		),
	)

	return err
}
