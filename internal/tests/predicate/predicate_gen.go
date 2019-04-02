// +build !cff
// @generated

package predicate

import (
	"context"
	"fmt"
	"sync"
)

// Simple runs the provided function in a task flow if the provided boolean
// is true.
func Simple(f func(), pred bool) error {
	var s string
	return func(ctx context.Context) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)

		var v1 string

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

			if func() bool { return pred }() {

				v1 = func() string {
					f()
					return "foo"
				}()

			}

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v1
		)

		*(&s) = v1

		return err
	}(context.Background())
}

// SimpleWithContextTask is a task flow which checks that context can be passed into Task w/out
// errors.
func SimpleWithContextTask() error {
	var s string
	return func(ctx context.Context, v2 int64) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)

		var v1 string

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

			if func(int64) bool {
				return false
			}(v2) {

				v1 = func(ctx context.Context) string {
					return "foo"
				}(ctx)

			}

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v1
		)

		*(&s) = v1

		return err
	}(context.Background(), int64(2))
}

// SimpleWithContextPredicate is a task flow which checks that context can be passed into
// Predicate but not Task.
func SimpleWithContextPredicate() error {
	var s string
	return func(ctx context.Context, v2 int64) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)

		var v1 string

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

			if func(context.Context, int64) bool {
				return false
			}(ctx, v2) {

				v1 = func() string {
					return "foo"
				}()

			}

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v1
		)

		*(&s) = v1

		return err
	}(context.Background(), int64(2))
}

// SimpleWithContextTaskAndPredicate is a task flow which checks that context can be passed into
// Predicate and Task.
func SimpleWithContextTaskAndPredicate() error {
	var s string
	return func(ctx context.Context, v2 int64) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)

		var v1 string

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

			if func(context.Context, int64) bool {
				return false
			}(ctx, v2) {

				v1 = func(ctx context.Context) string {
					return "foo"
				}(ctx)

			}

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v1
		)

		*(&s) = v1

		return err
	}(context.Background(), int64(2))
}

// ExtraDependencies is a task flow where the predicate has more dependencies
// than the task.
func ExtraDependencies() error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}

	var out t3
	return func(ctx context.Context, v3 int) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(2)

		var v4 t1

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

			v4 = func(int) t1 { return t1{} }(v3)

		}()

		var v5 t2

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

			v5 = func() t2 { return t2{} }()

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v4
			_ = &v5
		)

		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(2)

		var v1 string

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

			v1 = func(int) string { return "foo" }(v3)

		}()

		var v6 t3

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

			if func(int, t1) bool {
				return true
			}(v3, v4) {

				v6 = func(t2) t3 { return t3{} }(v5)

			}

		}()

		wg1.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v1
			_ = &v6
		)

		*(&out) = v6

		return err
	}(context.Background(), int(42))
}
