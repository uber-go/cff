package setconcurrency

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConcurrencyDirective(t *testing.T) {
	for _, tt := range []int{1, 2, 4, 8, 16} {
		t.Run(strconv.Itoa(tt), func(t *testing.T) {
			got, err := NumWorkers(tt)
			require.NoError(t, err)

			assert.Equal(t, tt, got, "number of workers does not match")
		})
	}
}
