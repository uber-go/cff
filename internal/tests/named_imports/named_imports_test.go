package namedimports

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNestedFlow(t *testing.T) {
	assert.Error(t, run(context.Background()))
}
