// +build !cff

package basic

import (
	"bytes"
	"context"
	"io"
	"sync"
)

type foo struct{ i int }
type bar struct{ i int64 }

// SimpleFlow is a very simple flow with some inputs and outputs.
func SimpleFlow() (string, error) {
	var message string
	err := func(ctx context.Context, v1 int) (err error) {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v2 int64
		v2 = func(i int) int64 {
			return int64(i)
		}(v1)

		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(2)

		var v3 *foo
		go func() {
			defer wg1.Done()

			var err1 error
			v3, err = func(i int) (*foo, error) {
				return &foo{i}, nil
			}(v1)
			if err1 != nil {

				once1.Do(func() {
					err = err1
				})
			}

		}()

		var v4 *bar
		go func() {
			defer wg1.Done()

			v4 = func(i int64) *bar {
				return &bar{i}
			}(v2)

		}()

		wg1.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v3
			_ = &v4
		)

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v5 string
		v5, err = func(*foo, *bar) (string, error) {
			return "hello world", nil
		}(v3, v4)
		if err != nil {

			return err
		}

		*(&message) = v5

		return err
	}(context.Background(), 1)
	return message, err
}

// NoParamsFlow is a flow that does not accept any parameters.
func NoParamsFlow(ctx context.Context) (io.Reader, error) {
	var r io.Reader
	err := func(ctx context.Context) (err error) {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v6 *bytes.Buffer
		v6 = func() *bytes.Buffer {
			return bytes.NewBufferString("hello world")
		}()

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v7 io.Reader
		v7 = func(b *bytes.Buffer) io.Reader { return b }(v6)

		*(&r) = v7

		return err
	}(ctx)
	return r, err
}

// SerialFailableFlow runs the provided function in-order using a flow.
func SerialFailableFlow(ctx context.Context, f1, f2 func() error) error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}

	// We use the types to make f2 depend on f1.
	var out t3
	return func(ctx context.Context) (err error) {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v8 t1
		v8, err = func() (t1, error) {
			return t1{}, f1()
		}()
		if err != nil {

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v9 t2
		v9, err = func(t1) (t2, error) {
			return t2{}, f2()
		}(v8)
		if err != nil {

			return err
		}

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v10 t3
		v10 = func(t2) t3 {
			return t3{}
		}(v9)

		*(&out) = v10

		return err
	}(ctx)
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
	return func(ctx context.Context, v11 t1) (err error) {

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var (
			v12 t2
			v13 t3
		)
		v12, v13 = func(t1) (t2, t3) {
			return t2{}, t3{}
		}(v11)

		if ctx.Err() != nil {
			return ctx.Err()
		}

		var v14 t4
		v14 = func(t2, t3) t4 {
			return t4{}
		}(v12, v13)

		*(&out) = v14

		return err
	}(context.Background(), t1{})
}
