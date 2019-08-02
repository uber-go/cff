package earlyresult

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEarlyResult(t *testing.T) {
	err := EarlyResult(context.Background())
	assert.NoError(t, err)
}

func TestConsumesResult(t *testing.T) {
	err := ConsumesResult()
	assert.NoError(t, err)
}
