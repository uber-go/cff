package badinputs

import (
	"go.uber.org/cff"
)

// BadTopLevelFunction is a bad call to a cff.X top level function.
func BadTopLevelFunction() {
	// cff.Predicate cannot be at the top level.
	cff.Predicate(func() {})
}
