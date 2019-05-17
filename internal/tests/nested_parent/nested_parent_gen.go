// +build !cff
// @generated

package nestedparent

import (
	"context"
	"fmt"
	"sync"

	// When built under +cff build tag, this refers to the raw flow. After code
	// generation, this refers to the generated code.
	"go.uber.org/cff/internal/tests/nested_child"
)

// Parent is a CFF flow that uses a nested CFF flow.
func Parent(ctx context.Context, i int) (s string, err error) {
	err = func(ctx context.Context, v1 int) (err error) {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		var (
			once0 sync.Once
		)

		var v2 string
		var err0 error
		func() {

			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)

						err = recoveredErr
					})
				}
			}()

			v2, err0 = nestedchild.Itoa(ctx, v1)
			if err0 != nil {

				once0.Do(func() {
					err = err0
				})
			}

		}()

		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v2
		)

		*(&s) = v2

		return err
	}(ctx, i)

	return
}
