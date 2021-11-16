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
				return "some" == "goal"
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
				return nil, nil
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

// InstrumentParallelInvalid is a Parallel that provides a cff.InstrumentParallel
// without an emitter.
func InstrumentParallelInvalid() {
	cff.Parallel(
		context.Background(),
		cff.InstrumentParallel("some instrument"),
		cff.Task(
			func() error {
				return nil
			},
		),
	)
}

// ParallelTaskInvalidParamsType is a Parallel with an invalid task parameters type.
func ParallelTaskInvalidParamsType() {
	cff.Parallel(
		context.Background(),
		cff.Task(
			func(s string) bool {
				return s == "goal"
			},
		),
	)
}

// ParallelTaskInvalidParamsMultiple is a Parallel with more than one task
// parameters.
func ParallelTaskInvalidParamsMultiple() {
	cff.Parallel(
		context.Background(),
		cff.Task(
			func(_ context.Context, _ context.Context) bool {
				return "some" == "goal"
			},
		),
	)
}

// ParallelTaskInvalidReturnType is a Parallel with a non-error task return value.
func ParallelTaskInvalidReturnType() {
	cff.Parallel(
		context.Background(),
		cff.Task(
			func(_ context.Context) bool {
				return true
			},
		),
	)
}

// ParallelTaskInvalidReturnMultiple is a Parallel with more than one return value.
func ParallelTaskInvalidReturnMultiple() {
	cff.Parallel(
		context.Background(),
		cff.Task(
			func(_ context.Context) (error, error) {
				return nil, nil
			},
		),
	)
}

// ParallelTaskInvalidFuncVar is a Parallel with an invalid function
// variable.
func ParallelTaskInvalidFuncVar() {
	cff.Parallel(
		context.Background(),
		cff.Task(
			chanSend,
		),
	)
}

// InstrumentParallelTaskInvalid is a Parallel that instruments a cff.Task
// without an emitter.
func InstrumentParallelTaskInvalid() {
	cff.Parallel(
		context.Background(),
		cff.Task(
			func() error {
				return nil
			},
			cff.Instrument("BadTask"),
		),
	)
}

func chanSend(s string, c chan<- string) {
	c <- s
}
