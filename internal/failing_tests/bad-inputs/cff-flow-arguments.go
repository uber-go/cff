// +build cff

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// ExpectsAtLeastOneArgument is a function that doesn't have enough arguments to cff.Flow.
func ExpectsAtLeastOneArgument() {
	cff.Flow(context.Background())
}

// FlowArgumentCallExpression2 is a function that has the wrong arguments to cff.Flow.
func FlowArgumentCallExpression2() {
	cff.Flow(context.Background(),
		(cff.FlowOption)(nil),
	)
}

// FlowArgumentCallExpression is a function that has the wrong arguments to cff.Flow.
func FlowArgumentCallExpression() {
	bad := (cff.FlowOption)(nil)
	cff.Flow(context.Background(),
		bad,
	)
}

// FlowArgumentNonCFF is a function that has the wrong arguments to cff.Flow.
func FlowArgumentNonCFF() {
	badProvider := struct{ ProvidesBad func() cff.FlowOption }{ProvidesBad: func() cff.FlowOption { return cff.Params() }}
	cff.Flow(context.Background(),
		badProvider.ProvidesBad(),
	)
}
