package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// ParallelInvalidParamsType is a Parallel task with an invalid parameters
// type.
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

// ParallelInvalidParamsMultiple is a Parallel task with more than one
// parameter.
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

// ParallelInvalidReturnType is a Parallel task with a non-error return value.
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

// ParallelInvalidReturnMultiple is a Parallel task with more than one
// return value.
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

// ParallelInvalidFuncVar is a Parallel task with a function reference that
// does not comply with CFF's Tasks function validation.
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
