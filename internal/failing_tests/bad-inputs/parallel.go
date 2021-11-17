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

// ParallelSliceBadContextPosition has a context.Context outside of the first
// argument of the cff.Slice execution function.
func ParallelSliceBadContextPosition() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ int, c context.Context, _ string) error {
				return nil
			},
			[]string{"some", "thing"},
		),
	)
}

// ParallelSliceFirstArgNotIndex is a cff.Slice where the first non-context
// argument is not an index.
func ParallelSliceFirstArgNotIndex() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ context.Context, _ string, _ int) {},
			[]string{"some", "thing"},
		),
	)
}

// ParallelSliceMissingIndex is a cff.Slice execution function that is missing
// the argument for slice index.
func ParallelSliceMissingIndex() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ context.Context, _ string) error {
				return nil
			},
			[]string{"some", "thing"},
		),
	)
}

// ParallelSliceMissingValue is a cff.Slice execution function that is missing
// the argument for slice value.
func ParallelSliceMissingValue() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ context.Context, _ int) {},
			[]string{"some", "thing"},
		),
	)
}

// ParallelSliceWithMap is a cff.Slice that is given a map to iterate over.
func ParallelSliceWithMap() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ int, _ string) error {
				return nil
			},
			map[string]struct{}{},
		),
	)
}

// ParallelSliceFuncTooManyArgs is a cff.Slice function that has too many arguments.
func ParallelSliceFuncTooManyArgs() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ int, _ string, _ bool) {},
			[]string{"some", "thing"},
		),
	)
}

// ParallelSliceElemTypeMismatch is a cff.Slice function whose value argument type
// does not match the type of the slice elements provided to cff.Slice.
func ParallelSliceElemTypeMismatch() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ int, _ string) error {
				return nil
			},
			[]bool{true, false},
		),
	)
}

// ParallelSliceNonErrorReturn is a cff.Slice function that has a non-error
// return.
func ParallelSliceNonErrorReturn() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ int, _ string) int {
				return 1
			},
			[]string{"some", "thing"},
		),
	)
}

// ParallelSliceTooManyReturn is a cff.Slice function that has too many
// return arguments.
func ParallelSliceTooManyReturn() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ int, _ string) (int, error) {
				return 1, nil
			},
			[]string{"some", "thing"},
		),
	)
}

// ParallelSliceNonLastError is a cff.Slice function that has an error as a
// outside of the last return value.
func ParallelSliceNonLastError() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ int, _ string) (error, int) {
				return nil, 1
			},
			[]string{"some", "thing"},
		),
	)
}

func chanSend(s string, c chan<- string) {
	c <- s
}
