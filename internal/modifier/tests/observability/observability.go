package observability

import (
	"context"

	"github.com/uber-go/tally"
	"go.uber.org/cff"
	"go.uber.org/zap"
)

// InstrumentFlow is a cff.Flow with an instrumented flow.
func InstrumentFlow(scope tally.Scope, logger *zap.Logger) (int64, error) {
	var res int64
	err := cff.Flow(context.Background(),
		cff.Concurrency(2),
		cff.Params(1),
		cff.Results(&res),
		cff.WithEmitter(cff.TallyEmitter(scope)),
		cff.WithEmitter(cff.LogEmitter(logger)),
		cff.InstrumentFlow("Instrumented"),
		cff.Task(
			func(i int) int64 {
				return int64(1)
			},
		),
	)
	return res, err
}
