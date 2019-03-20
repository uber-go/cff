// +build !cff

package basic

import (
	"bytes"
	"context"
	"fmt"
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
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v2 int64
		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v2 = func(i int) int64 {
				return int64(i)
			}(v1)

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v2
		)

		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(2)
		var v3 *foo
		var err1 error
		go func() {
			defer wg1.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once1.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v3, err1 = func(i int) (*foo, error) {
				return &foo{i}, nil
			}(v1)
			if err1 != nil {

				once1.Do(func() {
					err = err1
				})
			}

		}()
		var v4 *bar
		var err2 error
		go func() {
			defer wg1.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once1.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v4, err2 = func(i int64) (*bar, error) {
				return &bar{i}, nil
			}(v2)
			if err2 != nil {

				once1.Do(func() {
					err = err2
				})
			}

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
		var (
			wg2   sync.WaitGroup
			once2 sync.Once
		)

		wg2.Add(1)
		var v5 string
		var err3 error
		go func() {
			defer wg2.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once2.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v5, err3 = func(*foo, *bar) (string, error) {
				return "hello world", nil
			}(v3, v4)
			if err3 != nil {

				once2.Do(func() {
					err = err3
				})
			}

		}()

		wg2.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once2
			_ = &v5
		)

		*(&message) = v5

		return err
	}(context.Background(), 1)
	return message, err
}

// SimpleFlowNested has a cff.Task task within cff.Tasks.
func SimpleFlowNested() (string, error) {
	var message string
	err := func(ctx context.Context, v1 int) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v2 int64
		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v2 = func() int64 {
				return int64(1)
			}()

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v2
		)

		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(1)
		var v5 string
		go func() {
			defer wg1.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once1.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v5 = func(int64, int) string {
				return "foo"
			}(v2, v1)

		}()

		wg1.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v5
		)

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
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v6 *bytes.Buffer
		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v6 = func() *bytes.Buffer {
				return bytes.NewBufferString("hello world")
			}()

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v6
		)

		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(1)
		var v7 io.Reader
		go func() {
			defer wg1.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once1.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v7 = func(b *bytes.Buffer) io.Reader { return b }(v6)

		}()

		wg1.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v7
		)

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
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v8 t1
		var err8 error
		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v8, err8 = func() (t1, error) {
				return t1{}, f1()
			}()
			if err8 != nil {

				once0.Do(func() {
					err = err8
				})
			}

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v8
		)

		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(1)
		var v9 t2
		var err9 error
		go func() {
			defer wg1.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once1.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v9, err9 = func(t1) (t2, error) {
				return t2{}, f2()
			}(v8)
			if err9 != nil {

				once1.Do(func() {
					err = err9
				})
			}

		}()

		wg1.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v9
		)

		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg2   sync.WaitGroup
			once2 sync.Once
		)

		wg2.Add(1)
		var v10 t3
		go func() {
			defer wg2.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once2.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v10 = func(t2) t3 {
				return t3{}
			}(v9)

		}()

		wg2.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once2
			_ = &v10
		)

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
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var (
			v12 t2
			v13 t3
		)
		go func() {
			defer wg0.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v12, v13 = func(t1) (t2, t3) {
				return t2{}, t3{}
			}(v11)

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v12
			_ = &v13
		)

		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(1)
		var v14 t4
		go func() {
			defer wg1.Done()

			defer func() {
				recovered := recover()
				if recovered != nil {
					once1.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v14 = func(t2, t3) t4 {
				return t4{}
			}(v12, v13)

		}()

		wg1.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v14
		)

		*(&out) = v14

		return err
	}(context.Background(), t1{})
}
