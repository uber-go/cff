// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// ExpectsAtLeastOneArgument is a function that doesn't have enough arguments to cff.Flow
func ExpectsAtLeastOneArgument() {
	cff.Flow(context.Background())
}
