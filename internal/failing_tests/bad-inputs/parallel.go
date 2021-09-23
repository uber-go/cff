// +build cff

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// ParallelInvalidParamsType is a Parallel with an invalid task parameters type.
func ParallelInvalidParamsType() {
	cff.Parallel(
		context.Background(),
		cff.Tasks(
			func(s string) bool {
				return s == "goal"
			},
		),
	)
}

// ParallelInvalidParamsMultiple is a Parallel with more than one task
// parameters.
func ParallelInvalidParamsMultiple() {
	cff.Parallel(
		context.Background(),
		cff.Tasks(
			func(_ context.Context, _ context.Context) bool {
				return s == "goal"
			},
		),
	)
}

// ParallelInvalidReturnType is a Parallel with a non-error task return value.
func ParallelInvalidReturnType() {
	cff.Parallel(
		context.Background(),
		cff.Tasks(
			func(_ context.Context) bool {
				return true
			},
		),
	)
}

// ParallelInvalidReturnMultiple is a Parallel with more than one return value.
func ParallelInvalidReturnMultiple() {
	cff.Parallel(
		context.Background(),
		cff.Tasks(
			func(_ context.Context) (error, error) {
				return true
			},
		),
	)
}

// ParallelInvalidFuncVar is a Parallel with an invalid function
// variable.
func ParallelInvalidFuncVar() {
	cff.Parallel(
		context.Background(),
		cff.Tasks(
			chanSend,
		),
	)
}

func chanSend(s string, c chan<- string) {
	c <- s
}
