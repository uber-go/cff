// +build !cff

package predicate

import (
	"context"
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

// ExtraDependencies is a task flow where the predicate has more dependencies
// than the task.
func ExtraDependencies() error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}

	var out t3
	return func(ctx context.Context, v2 int) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(2)
		var v3 t1
		go func() {
			defer wg0.Done()

			v3 = func(int) t1 { return t1{} }(v2)

		}()
		var v4 t2
		go func() {
			defer wg0.Done()

			v4 = func() t2 { return t2{} }()

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v3
			_ = &v4
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

			v1 = func(int) string { return "foo" }(v2)

		}()
		var v5 t3
		go func() {
			defer wg1.Done()

			if func(int, t1) bool {
				return true
			}(v2, v3) {

				v5 = func(t2) t3 { return t3{} }(v4)

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
			_ = &v5
		)

		*(&out) = v5

		return err
	}(context.Background(), int(42))
}
