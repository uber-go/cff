package insidegeneric

import (
	"context"

	"go.uber.org/cff"
)

// Producer is a function that produces a value of the given type.
type Producer[T any] func(context.Context) (T, error)

// JoinTwo combines two producers producing different values using the provided
// function.
func JoinTwo[A, B, C any](
	pa Producer[A],
	pb Producer[B],
	fn func(A, B) C,
) (C, error) {
	var c C
	err := cff.Flow(context.Background(),
		cff.Results(&c),
		cff.Task(pa),
		cff.Task(pb),
		cff.Task(fn),
	)
	return c, err
}

// JoinMany runs the given producers and returns a slice of their results
// in-order.
func JoinMany[T any](producers ...Producer[T]) ([]T, error) {
	results := make([]T, len(producers))
	err := cff.Parallel(context.Background(),
		cff.Slice(
			func(ctx context.Context, idx int, fn Producer[T]) error {
				v, err := fn(ctx)
				results[idx] = v
				return err
			},
			producers,
		),
	)
	return results, err
}
