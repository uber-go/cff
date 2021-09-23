package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// FallbackWithNoError is a function that uses FallbackWith, but the function does not return an error.
func FallbackWithNoError() {
	cff.Flow(context.Background(),
		cff.Task(
			func() bool {
				return false
			},
			cff.FallbackWith(true),
		),
	)
}

// FallbackWithBadPositionalArguments is a function with too few positional arguments in the fallback literals.
func FallbackWithBadPositionalArguments() {
	cff.Flow(context.Background(),
		cff.Task(
			func() (bool, string, error) {
				return false, "", nil
			},
			cff.FallbackWith(true),
		),
	)
}

// FallbackWithTypeMismatch is a function whose fallback positional arguments are not assignable to the
// original function.
func FallbackWithTypeMismatch() {
	cff.Flow(context.Background(),
		cff.Task(
			func() (bool, string, error) {
				return false, "", nil
			},
			cff.FallbackWith("fb", false),
		),
	)
}
