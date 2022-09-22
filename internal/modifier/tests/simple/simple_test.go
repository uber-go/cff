package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	v, err := Flow()
	assert.NoError(t, err)
	assert.Equal(t, v, 1)
}
