//go:build cff
// +build cff

package importstmt

import (
	"context"
	"strconv"

	"go.uber.org/cff"
)

// Flow tests a flow where there is a non-grouped import before "go.uber.org/cff"
// which delete during codegen phase.
func Flow() (int, error) {
	var s int

	err := cff.Flow(context.Background(),
		cff.Params("123"),
		cff.Results(&s),
		cff.Task(strconv.Atoi),
	)

	return s, err
}
