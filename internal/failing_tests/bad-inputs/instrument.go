package badinputs

import (
	"context"

	"go.uber.org/cff"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

// MissingCFFMetrics is a flow that wants instrumentation but doesn't provide
// a cff.Metrics.
func MissingCFFMetrics() {
	core, _ := observer.New(zap.DebugLevel)
	logger := zap.New(core)
	cff.Flow(context.Background(),
		cff.WithEmitter(cff.LogEmitter(logger)),
		cff.Task(
			func() error {
				return nil
			},
			cff.Invoke(true),
			cff.Instrument("foo"),
		),
	)
}

// MissingCFFLogger is a flow that wants instrumentation but doesn't provide
// a cff.Logger.
func MissingCFFLogger() {
	cff.Flow(context.Background(),
		cff.WithEmitter(cff.TallyEmitter(tally.NewTestScope("", nil))),
		cff.Task(
			func() error {
				return nil
			},
			cff.Invoke(true),
			cff.Instrument("foo"),
		),
	)
}

// MissingCFFLoggerME is a flow that wants instrumentation but doesn't provide
// a cff.Logger.
func MissingCFFLoggerME() {
	cff.Flow(context.Background(),
		cff.WithEmitter(cff.TallyEmitter(tally.NewTestScope("", nil))),
		cff.Task(
			func() error {
				return nil
			},
			cff.Invoke(true),
			cff.Instrument("foo"),
		),
	)
}

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
