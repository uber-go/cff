package noresults

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSwallowSuccess(t *testing.T) {
	h := &h{}
	err := h.swallow(context.TODO(), "food")
	assert.NoError(t, err)
}

func TestSwallowError(t *testing.T) {
	h := &h{}
	err := h.swallow(context.TODO(), "tide pods")
	assert.EqualError(t, err, "can not swallow")
}

func TestTripleSwallow(t *testing.T) {
	h := &h{}
	err := h.tripleSwallow(context.TODO(), "tide pods")
	assert.NoError(t, err)
}
