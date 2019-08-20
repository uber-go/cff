// +build cff

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// UnusedInputs is a flow that has an input that is unused.
func UnusedInputs() {
	var s string
	cff.Flow(context.Background(),
		cff.Params(s),
	)
}
