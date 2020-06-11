package noresults

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwallowSuccess(t *testing.T) {
	h := &H{}
	err := h.Swallow(context.TODO(), "food")
	assert.NoError(t, err)
}

func TestSwallowError(t *testing.T) {
	h := &H{}
	err := h.Swallow(context.TODO(), "tide pods")
	assert.EqualError(t, err, "can not Swallow")
}

func TestTripleSwallow(t *testing.T) {
	h := &H{}
	err := h.TripleSwallow(context.TODO(), "tide pods")
	assert.NoError(t, err)
}

func TestUnusedInputInvoke(t *testing.T) {
	assert.NoError(t, UnusedInputInvoke())
}
