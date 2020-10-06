package setconcurrency

import (
	"fmt"
	"runtime"
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

func TestGoMaxProcs(t *testing.T) {
	// Testing that if GOMAXPROCS is less than 4, we set 4 at minimum.
	tests := []struct {
		give int
		want int
	}{
		{
			give: 2,
			want: 4,
		},
		{
			give: 6,
			want: 6,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("GOMAXPROCS=%v", tt.give), func(t *testing.T) {
			old := runtime.GOMAXPROCS(0)
			runtime.GOMAXPROCS(tt.give)
			t.Cleanup(func() { runtime.GOMAXPROCS(old) })
			got, err := NumWorkersNoArg()
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
