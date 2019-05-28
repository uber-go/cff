package nestedchild_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	nested_child_gen "go.uber.org/cff/internal/tests/nested_child_gen"
)

func TestNestedFlow(t *testing.T) {
	s, err := nested_child_gen.Itoa(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "1", s)
}
