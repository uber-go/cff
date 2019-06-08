package nestedparent

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedFlow(t *testing.T) {
	s, err := Parent(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "1", s)
}
