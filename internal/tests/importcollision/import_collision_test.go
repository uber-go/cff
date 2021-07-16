package importcollision

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNoImportCollision(t *testing.T) {
	_, err := Flow()
	assert.NoError(t, err)
}
