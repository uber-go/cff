package badinputs

import (
	"context"
	"errors"

	"go.uber.org/cff"
)

// ExpectsAtLeastOneArgument is a function that doesn't have enough arguments to cff.Flow.
func ExpectsAtLeastOneArgument() {
	cff.Flow(context.Background())
}

// FlowArgumentCallExpression2 is a function that has the wrong arguments to cff.Flow.
func FlowArgumentCallExpression2() {
	cff.Flow(context.Background(),
		(cff.Option)(nil),
	)
}

// FlowArgumentCallExpression is a function that has the wrong arguments to cff.Flow.
func FlowArgumentCallExpression() {
	bad := (cff.Option)(nil)
	cff.Flow(context.Background(),
		bad,
	)
}

// FlowArgumentNonCFF is a function that has the wrong arguments to cff.Flow.
func FlowArgumentNonCFF() {
	badProvider := struct{ ProvidesBad func() cff.Option }{ProvidesBad: func() cff.Option { return cff.Params() }}
	cff.Flow(context.Background(),
		badProvider.ProvidesBad(),
	)
}

// DisallowContinueOnError is a function that provides cff.ContinueOnError
// to cff.Flow.
func DisallowContinueOnError() {
	cff.Flow(context.Background(),
		cff.ContinueOnError(true),
		cff.Task(
			func() error { return errors.New("sad times") },
		),
	)
}

// DisallowSlice is a function that provides cff.Slice
// to cff.Flow.
func DisallowSlice() {
	cff.Flow(context.Background(),
		cff.Slice(
			func(_ int, elem string) error { return nil },
			[]string{"sad", "times"},
		),
	)
}
