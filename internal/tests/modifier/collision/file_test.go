package collision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFlows tests cases where we have cff modifiers whose
// positions in their respective files collide with each other.
// This test compiles cff.Flows from two different files that
// will collide in line/column pair, to make sure that the
// modifiers generate unique names for colliding primitives.
func TestFlows(t *testing.T) {
	i, err := Flow1()
	assert.NoError(t, err)
	j, err := Flow2()
	assert.NoError(t, err)
	assert.Equal(t, 1, i)
	assert.Equal(t, 2, j)
}
