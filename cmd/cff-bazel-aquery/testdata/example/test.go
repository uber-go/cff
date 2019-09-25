// +build cff

package example

import (
	"context"

	"go.uber.org/cff"
)

// Example ...
func Example() {
	var s string
	_ = cff.Flow(context.Background(),
		cff.Results(&s),
		cff.Task(func() string { return "" }),
	)
}
