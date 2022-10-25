//go:build cff
// +build cff

package cffintest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff"
)

func TestIsOdd(t *testing.T) {
	err := cff.Parallel(
		context.Background(),
		cff.Concurrency(4),
		cff.Tasks(
			func() error {
				assert.Equal(t, true, isOdd(1))
				return nil
			},
			func() error {
				assert.Equal(t, false, isOdd(2))
				return nil
			},
			func() error {
				assert.Equal(t, true, isOdd(3))
				return nil
			},
			func() error {
				assert.Equal(t, false, isOdd(4))
				return nil
			},
		),
	)
	require.NoError(t, err)
}
