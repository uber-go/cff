// +build cff

package cycles

import (
	"go.uber.org/cff"
	"context"
)

// EasyCycle is a flow that has a simple cycle.
func EasyCycle() {
	var out string
	cff.Flow(
		context.Background(),
		cff.Results(&out),
		cff.Tasks(
			func(string) int64 {
				return int64(0)
			},
			func(int64) string {
				return ""
			},
		))
}
