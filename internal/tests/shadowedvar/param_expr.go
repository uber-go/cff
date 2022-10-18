//go:build cff
// +build cff

package shadowedvar

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff"
)

// ParamOrder initializes a cff.Flow to test that the order in which user
// provided expressions are evalauted matches the order in which they were
// provided to the cff.Flow.
func ParamOrder(track *orderCheck) error {
	var res string
	return cff.Flow(
		track.ctx(),
		cff.Params(track.param1(), track.param2()),
		cff.Results(&res),
		cff.Task(func(_ int, _ bool) (string, error) {
			return "", nil
		}),
	)
}

// NilParam verifies that cff.Flow is compilable with a user provided nil.
// CFF should compile and generate this flow even if no test function
// uses it.
func NilParam() {
	var res []int
	cff.Flow(
		context.Background(),
		cff.Params(1, true),
		cff.Results(&res),
		cff.Task(
			func(_ int, _ bool) ([]int, error) {
				return nil, nil
			},
			cff.FallbackWith(nil),
		),
	)
}

// checks the call order of parameter expressions invoked by cff.Flow.
// cff.Flow assumes all non-task parameters expressions are not concurrently
// invoked.
type orderCheck struct {
	t *testing.T
	// counter tracks order in which expressions are invoked.
	counter int
	// order expectations for parameter expressions.
	order map[string]int
}

func (c *orderCheck) ctx() context.Context {
	o, ok := c.order["ctx"]
	require.True(c.t, ok)
	assert.Equal(c.t, o, c.counter)

	c.counter++
	return context.Background()
}

func (c *orderCheck) param1() int {
	o, ok := c.order["param1"]
	require.True(c.t, ok)
	assert.Equal(c.t, o, c.counter)

	c.counter++
	return 0
}

func (c *orderCheck) param2() bool {
	o, ok := c.order["param2"]
	require.True(c.t, ok)
	assert.Equal(c.t, o, c.counter)

	c.counter++
	return true
}
