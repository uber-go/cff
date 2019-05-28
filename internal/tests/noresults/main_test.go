package noresults_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	noresults_gen "go.uber.org/cff/internal/tests/noresults_gen"
)

func TestSwallowSuccess(t *testing.T) {
	h := &noresults_gen.H{}
	err := h.Swallow(context.TODO(), "food")
	assert.NoError(t, err)
}

func TestSwallowError(t *testing.T) {
	h := &noresults_gen.H{}
	err := h.Swallow(context.TODO(), "tide pods")
	assert.EqualError(t, err, "can not Swallow")
}

func TestTripleSwallow(t *testing.T) {
	h := &noresults_gen.H{}
	err := h.TripleSwallow(context.TODO(), "tide pods")
	assert.NoError(t, err)
}
