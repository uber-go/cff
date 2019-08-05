// build +cff

package cycles

import (
	"go.uber.org/cff"
	"context"
)

// EasyCycleNoOut is a flow that has a simple cycle with no result.
func EasyCycleNoOut() {
	cff.Flow(
		context.Background(),
		cff.Task(
			func(string) int32 {
				return int32(0)
			},
		),
		cff.Task(
			func(int32) string {
				return ""
			},
		))
}

type moo struct{}

// EasyCycleSingleNode is a flow that has a simple cycle with no result.
func EasyCycleSingleNode() {
	cff.Flow(
		context.Background(),
		cff.Task(
			func(*moo) *moo {
				return &moo{}
			},
		))
}
