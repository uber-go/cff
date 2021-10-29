package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// MissingCFFLoggerAndMetrics is a flow that wants instrumentation but doesn't provide
// a cff.Metrics nor cff.Logger.
func MissingCFFLoggerAndMetrics() {
	cff.Flow(context.Background(),
		cff.Task(
			func() error {
				return nil
			},
			cff.Invoke(true),
			cff.Instrument("foo"),
		),
	)
}
