// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// MissingProvider is a flow that doesn't have a provider for a type.
func MissingProvider() {
	var s string
	cff.Flow(context.Background(),
		cff.Results(&s),
		cff.Task(
			func(float64) string {
				return ""
			},
		),
	)
}
