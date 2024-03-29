//go:build cff && failing
// +build cff,failing

package cycles

import (
	"context"

	"go.uber.org/cff"
)

// EasyCycle is a flow that has a simple cycle.
func EasyCycle() {
	var out string
	cff.Flow(
		context.Background(),
		cff.Results(&out),
		cff.Task(
			func(string) int64 {
				return int64(0)
			},
		),
		cff.Task(
			func(int64) string {
				return ""
			},
		))
}
