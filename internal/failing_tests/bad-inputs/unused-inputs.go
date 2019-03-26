// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// UnusedInputs is a flow that has an input that is unused.
func UnusedInputs() {
	var s string
	cff.Flow(context.Background(),
		cff.Params(s),
	)
}
