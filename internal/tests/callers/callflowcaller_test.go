package callers_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.uber.org/cff/internal/tests/callers"
)

func TestSandwichFlowCall(t *testing.T) {
	aFlowOut := "nested doll"
	bFlowOut := "nested doll"
	actualAFlow, actualBFlow := callers.PackageCall()
	assert.Equal(t, aFlowOut, actualAFlow)
	assert.Equal(t, bFlowOut, actualBFlow)
}
