package simple

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/cff/internal/tests/modifier/external"
)

func TestFlow(t *testing.T) {
	iRes, sRes, err := Flow()
	assert.NoError(t, err)
	assert.Equal(t, iRes, 1)
	assert.Equal(t, sRes, "non-zero")
}

func TestModifyVarInScope(t *testing.T) {
	res, side, err := ModifyVarInScope()
	assert.NoError(t, err)
	assert.Equal(t, res, true)
	assert.Equal(t, side, []int{1, 2, 3})
}

func TestExternal(t *testing.T) {
	res, err := External()
	assert.NoError(t, err)
	assert.Equal(t, res, true)
}

func TestParams(t *testing.T) {
	sRes, eRes, err := Params()
	assert.NoError(t, err)
	assert.Equal(t, sRes, "true")
	assert.Equal(t, eRes, external.A(1))
}
