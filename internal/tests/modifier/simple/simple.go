//go:build cff
// +build cff

package simple

import (
	"context"

	"go.uber.org/cff"
	"go.uber.org/cff/internal/tests/modifier/external"
)

type bar struct{ i int64 }

// Flow is a very simple flow with some inputs and outputs.
func Flow() (int, string, error) {
	var (
		iRes int
		sRes string
	)
	err := cff.Flow(context.Background(),
		cff.Concurrency(2),
		cff.Results(&iRes, &sRes),
		cff.Task(
			func() int64 {
				return int64(1)
			},
		),
		cff.Task(
			func(i int64) (*bar, error) {
				return &bar{i}, nil
			}),
		cff.Task(
			func(*bar) (int, error) {
				return 1, nil
			},
		),
		cff.Task(
			func(i int) (string, error) {
				if i != 0 {
					return "non-zero", nil
				}
				return "zero", nil
			},
		),
	)
	return iRes, sRes, err
}

// ModifyVarInScope is a simple flow that has a side effect of modifying a variable
// in scope.
func ModifyVarInScope() (bool, []int, error) {
	var res bool
	slc := make([]int, 3)
	err := cff.Flow(context.Background(),
		cff.Concurrency(2),
		cff.Results(&res),
		cff.Task(
			func() int64 {
				slc[0] = 1
				return int64(1)
			},
		),
		cff.Task(
			func(i int64) (*bar, error) {
				slc[1] = 2
				return &bar{i}, nil
			}),
		cff.Task(
			func(*bar) (bool, error) {
				slc[2] = 3
				return true, nil
			},
		),
	)
	return res, slc, err
}

// External is a simple flow that depends on an external package.
func External() (bool, error) {
	var res bool
	err := cff.Flow(context.Background(),
		cff.Concurrency(2),
		cff.Results(&res),
		cff.Task(
			func() external.A {
				return 1
			},
		),
		cff.Task(external.Run),
		cff.Task(
			func(b external.B) (bool, error) {
				return bool(b), nil
			},
		),
	)
	return res, err
}

// Params is a simple cff.Flow that depends on cff.Params.
func Params() (string, external.A, error) {
	var (
		res1 string
		res2 external.A
	)
	err := cff.Flow(context.Background(),
		cff.Concurrency(2),
		cff.Params(1, true),
		cff.Results(&res1, &res2),
		cff.Task(
			func(i int) int64 {
				return int64(i)
			},
		),
		cff.Task(
			func(i int64) (external.A, error) {
				return external.A(i), nil
			}),
		cff.Task(
			func(b bool) (string, error) {
				if b {
					return "true", nil
				}
				return "false", nil
			},
		),
	)
	return res1, res2, err
}
