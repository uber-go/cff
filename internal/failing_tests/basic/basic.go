// +build cff

package basic

import (
	"go.uber.org/cff"
	"context"
)

// TypeArgumentNotSpecified is a flow that is missing a type to a function's arguments
func TypeArgumentNotSpecified() (string, error) {
	var message string
	err := cff.Flow(context.Background(),
		cff.Results(&message),
		cff.Tasks(
			func(float64) string {
				return ""
			},
		),
	)
	return message, err
}
