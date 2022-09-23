package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFlow(t *testing.T) {
	iRes, sRes, err := Flow()
	assert.NoError(t, err)
	assert.Equal(t, iRes, 1)
	assert.Equal(t, sRes, "non-zero")
}
