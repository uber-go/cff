// +build cff

package basic

import (
	"bytes"
	"context"
	"io"

	"go.uber.org/cff"
)

type foo struct{ i int }
type bar struct{ i int64 }

// SimpleFlow is a very simple flow with some inputs and outputs.
func SimpleFlow() (string, error) {
	var message string
	err := cff.Flow(context.Background(),
		cff.Provide(1),
		cff.Result(&message),
		cff.Tasks(
			func(i int) int64 {
				return int64(i)
			},
			func(i int) (*foo, error) {
				return &foo{i}, nil
			},
			func(i int64) *bar {
				return &bar{i}
			},
			func(*foo, *bar) (string, error) {
				return "hello world", nil
			},
		),
	)
	return message, err
}

// NoParamsFlow is a flow that does not accept any parameters.
func NoParamsFlow(ctx context.Context) (io.Reader, error) {
	var r io.Reader
	err := cff.Flow(ctx,
		cff.Result(&r),
		cff.Tasks(
			func() *bytes.Buffer {
				return bytes.NewBufferString("hello world")
			},
			func(b *bytes.Buffer) io.Reader { return b },
		),
	)
	return r, err
}

// SerialFailableFlow runs the provided function in-order using a flow.
func SerialFailableFlow(f1, f2 func() error) error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}

	// We use the types to make f2 depend on f1.
	var out t3
	return cff.Flow(
		context.Background(),
		cff.Result(&out),
		cff.Tasks(
			func() (t1, error) {
				return t1{}, f1()
			},
			func(t1) (t2, error) {
				return t2{}, f2()
			},
			func(t2) t3 {
				return t3{}
			},
		),
	)
}

// ProduceMultiple has a task which produces multiple values.
func ProduceMultiple() error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}
	type t4 struct{}

	//   t1
	//   /\
	//  v   v
	// t2   t3
	//   \ /
	//    v
	//   t4

	var out t4
	return cff.Flow(
		context.Background(),
		cff.Provide(t1{}),
		cff.Result(&out),
		cff.Tasks(
			func(t1) (t2, t3) {
				return t2{}, t3{}
			},
			func(t2, t3) t4 {
				return t4{}
			},
		),
	)
}
