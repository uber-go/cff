package fallbackwith

import (
	"context"
	"errors"

	"go.uber.org/cff"
)

// Serial executes a flow that fails with the given error, if any and recovers
// with the given string.
func Serial(e error, r string) (string, error) {
	var s string
	err := cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Task(func() (string, error) {
			return "foo", e
		}, cff.FallbackWith(r)),
	)
	return s, err
}

// NoOutput is a flow that uses a no-output task (task of no return values) but uses FallbackWith.
func NoOutput() error {
	err := cff.Flow(
		context.Background(),
		cff.Task(
			func() error {
				return errors.New("always errors")
			},
			cff.FallbackWith(),
			cff.Invoke(true),
		),
	)

	return err
}

// Panic executes a flow that will panic but recovers with FallbackWith
func Panic() (string, error) {
	var rv string
	err := cff.Flow(
		context.Background(),
		cff.Results(&rv),
		cff.Task(func() (string, error) {
			panic("always panics")
		}, cff.FallbackWith("fallback")),
	)
	return rv, err
}
