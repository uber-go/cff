//go:build cff && failing
// +build cff,failing

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// MissingCffLoggerAndMetrics is a flow that wants instrumentation but doesn't provide
// a cff.Metrics nor cff.Logger.
func MissingCffLoggerAndMetrics() {
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
