// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

// FlowWithoutScope is a function that uses Instrument, but has no scope.
func FlowWithoutScope() {
	logger := zap.NewNop()

	cff.Flow(context.Background(),
		cff.Logger(logger),
		cff.InstrumentFlow("flow"),
		cff.Task(
			func() bool {
				return false
			},
			cff.Instrument("task"),
		),
	)
}

// FlowWithoutLogger is a function that uses Instrument, but has no logger.
func FlowWithoutLogger() {
	scope := tally.NewTestScope("", nil)

	cff.Flow(context.Background(),
		cff.Metrics(scope),
		cff.InstrumentFlow("flow"),
		cff.Task(
			func() bool {
				return false
			},
			cff.Instrument("task"),
		),
	)
}
