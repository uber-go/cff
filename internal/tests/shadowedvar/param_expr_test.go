package shadowedvar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlowParamOrder(t *testing.T) {
	order := map[string]int{
		"ctx":    0,
		"param1": 1,
		"param2": 2,
	}
	check := orderCheck{
		t:     t,
		order: order,
	}

	assert.NoError(t, ParamOrder(&check))
	assert.Equal(
		t,
		len(order),
		check.counter,
		"%d param expression evaluations expected; got %d",
		len(order),
		check.counter,
	)
}
