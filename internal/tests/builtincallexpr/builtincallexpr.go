//go:build cff
// +build cff

package builtincallexpr

import (
	"context"
	"fmt"
	"io"
	"strconv"

	"go.uber.org/cff"
)

// Flow is some code that has a call expression on a type in the core language (error)
func Flow(s string, buf io.Writer) {
	var i int
	err := cff.Flow(context.Background(),
		cff.Params(s),
		cff.Results(&i),
		cff.Task(func(s string) (int, error) { return strconv.Atoi(s) }))
	if err != nil {
		buf.Write([]byte(fmt.Sprintf("error: %s", err.Error())))
	}
}
