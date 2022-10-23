//go:build cff && failing
// +build cff,failing

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
// does not comply with cff's Tasks function validation.
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

// DisallowInstrumentFlow is a Parallel that provides a cff.InstrumentFlow.
func DisallowInstrumentFlow() {
	cff.Parallel(
		context.Background(),
		cff.InstrumentFlow("sad"),
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

// ParallelSliceEndWithContinueOnError is a cff.Slice that uses cff.ContinueOnError
// and cff.SliceEnd.
func ParallelSliceEndWithContinueOnError() {
	cff.Parallel(
		context.Background(),
		cff.ContinueOnError(true),
		cff.Slice(
			func(int, string) error {
				return nil
			},
			[]string{"some", "thing"},
			cff.SliceEnd(func(context.Context) {}),
		),
	)
}

// ParallelSliceEndTooManyArguments is a cff.Slice function that has too many arguments.
func ParallelSliceEndTooManyArguments() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(_ int, _ string) error {
				return nil
			},
			[]string{"some", "thing"},
			cff.SliceEnd(func(context.Context, int) {}),
		),
	)
}

// ParallelSliceEndWithInvalidArgument is a cff.Slice function that has invalid
// argument in cff.SliceEnd.
func ParallelSliceEndWithInvalidArgument() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(int, string) error {
				return nil
			},
			[]string{"some", "thing"},
			cff.SliceEnd(func(int) {}),
		),
	)
}

// ParallelSliceEndWithMultipleReturns is a cff.Slice function that has multiple
// returns in cff.SliceEnd.
func ParallelSliceEndWithMultipleReturns() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(int, string) error {
				return nil
			},
			[]string{"some", "thing"},
			cff.SliceEnd(func() (error, error) {
				return nil, nil
			}),
		),
	)
}

// ParallelSliceEndWithInvalidReturn is a cff.Slice function that has an invalid
// return in cff.SliceEnd.
func ParallelSliceEndWithInvalidReturn() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(int, string) error {
				return nil
			},
			[]string{"some", "thing"},
			cff.SliceEnd(func() int {
				return 0
			}),
		),
	)
}

// ParallelSliceEndWithContinueOnErrorAndInstrument is a cff.Slice function
// that uses SliceEnd, ContinueOnError, and Instrument.
//
// This tests for a regression where SliceEnd-ContinueOnError incompatibility
// was not verified if Instrument was used.
func ParallelSliceEndWithContinueOnErrorAndInstrument() {
	cff.Parallel(context.Background(),
		cff.WithEmitter(cff.NopEmitter()),
		cff.InstrumentParallel("myparallel"),
		cff.ContinueOnError(true),
		cff.Slice(
			func(int, string) {
				// stuff
			},
			[]string{"foo", "bar"},
			cff.SliceEnd(func() error {
				return nil
			}),
		),
	)
}

// ParallelSliceEndWithContinueOnErrorAndInstrument is a cff.Slice function
// that uses SliceEnd, ContinueOnError, and Instrument.
//
// This tests for a regression where SliceEnd-ContinueOnError incompatibility
// was not verified if Instrument was used.
func ParallelMapEndWithContinueOnErrorAndInstrument() {
	cff.Parallel(context.Background(),
		cff.WithEmitter(cff.NopEmitter()),
		cff.InstrumentParallel("myparallel"),
		cff.ContinueOnError(true),
		cff.Map(
			func(string, string) {
				// stuff
			},
			map[string]string{"foo": "bar"},
			cff.MapEnd(func() error {
				return nil
			}),
		),
	)
}

// ParallelSliceWithTwoSliceEnds is a cff.Slice function that has more than
// one cff.SliceEnd.
func ParallelSliceWithTwoSliceEnds() {
	cff.Parallel(
		context.Background(),
		cff.Slice(
			func(int, string) error {
				return nil
			},
			[]string{"some", "thing"},
			cff.SliceEnd(func() error {
				return nil
			}),
			cff.SliceEnd(func() error {
				return nil
			}),
		),
	)
}

// ParallelMapNilFunction is a cff.Map with a nil value func.
func ParallelMapNilFunction() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			nil,
			map[string]string{"one": "one"},
		),
	)
}

// ParallelMapInvalidReturnType is a cff.Map with an invalid return type.
func ParallelMapInvalidReturnType() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func(k, v string) (s string) { return },
			map[string]string{"one": "one"},
		),
	)
}

// ParallelMapTooManyReturns is a cff.Map with too many return values.
func ParallelMapTooManyReturns() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func(k, v string) (s string, e error) { return },
			map[string]string{"one": "one"},
		),
	)
}

// ParallelMapNoArguments is a cff.Map with no arguments.
func ParallelMapNoArguments() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func() (e error) { return },
			map[string]string{"one": "one"},
		),
	)
}

// ParallelMapWithSlice is a cff.Map with a slice.
func ParallelMapWithSlice() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func(k, v string) (e error) { return },
			[]string{"test"},
		),
	)
}

// ParallelMapWithDifferentKeyType is a cff.Map with a different key type.
func ParallelMapWithDifferentKeyType() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func(k, v string) (e error) { return },
			map[bool]string{true: "true"},
		),
	)
}

// ParallelMapWithDifferentValueType is a cff.Map with a different value type.
func ParallelMapWithDifferentValueType() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func(k, v string) (e error) { return },
			map[string]bool{"true": true},
		),
	)
}

// ParallelMapWithMultipleMapEnds has a cff.Map call with multiple cff.MapEnd
// options.
func ParallelMapWithMultipleMapEnds() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func(k string, v bool) (e error) { return },
			map[string]bool{"true": true},
			cff.MapEnd(func() {}),
			cff.MapEnd(func(context.Context) {}),
		),
	)
}

// ParallelMapEndWithNonContextArgument has a cff.MapEnd call
// with a function that accepts a non-context argument.
func ParallelMapEndWithNonContextArgument() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func(k string, v bool) (e error) { return },
			map[string]bool{"true": true},
			cff.MapEnd(func(string) {}),
		),
	)
}

// ParallelMapEndWithNonErrorResult has a cff.MapEnd call
// with a function that returns a non-error result.
func ParallelMapEndWithNonErrorResult() {
	cff.Parallel(
		context.Background(),
		cff.Map(
			func(k string, v bool) (e error) { return },
			map[string]bool{"true": true},
			cff.MapEnd(func() int { return 0 }),
		),
	)
}

// ParallelMapEndWithContinueOnError has a cff.MapEnd call
// with a ContinueOnError.
func ParallelMapEndWithContinueOnError() {
	cff.Parallel(
		context.Background(),
		cff.ContinueOnError(true),
		cff.Map(
			func(k string, v bool) (e error) { return },
			map[string]bool{"true": true},
			cff.MapEnd(func() {}),
		),
	)
}

func chanSend(s string, c chan<- string) {
	c <- s
}
