// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// ExpectsFunctionCallExpression is a flow that doesn't provide a function to cff.Task
func ExpectsFunctionCallExpression() {
	cff.Flow(context.Background(), cff.Task(
		true,
	))
}
