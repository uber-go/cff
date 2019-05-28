package nestedparent_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	nested_parent_gen "go.uber.org/cff/internal/tests/nested_parent_gen"
)

func TestNestedFlow(t *testing.T) {
	s, err := nested_parent_gen.Parent(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "1", s)
}
