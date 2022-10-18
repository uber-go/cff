package instrumentimplied

import (
	"context"
	"strconv"

	"github.com/uber-go/tally"
	"go.uber.org/cff"
	"go.uber.org/zap"
)

// H is used by some tests
type H struct {
	Scope  tally.Scope
	Logger *zap.Logger
}

// ImpliedName is a flow with a task with no instrument call but the flow is instrumented
func (h *H) ImpliedName(ctx context.Context, req string) (res int, err error) {
	var unsigned uint

	err = cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res, &unsigned),
		cff.Results(&unsigned),
		cff.WithEmitter(cff.TallyEmitter(h.Scope)),
		cff.WithEmitter(cff.LogEmitter(h.Logger)),
		cff.InstrumentFlow("ImpliedName"),
		cff.Task(
			strconv.Atoi,
			// Instrument annotation intentionally omitted.
		),
		cff.Task(
			func(i int) (uint, error) {
				return uint(i), nil
			},
			// Instrument annotation intentionally omitted.
		),
	)
	return
}
