package externalpackage

import (
	"context"
	"testing"

	"github.com/gofrs/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSimpleFlow(t *testing.T) {
	id, _ := uuid.NewV1()
	err := NestedType(context.Background(), id)
	assert.NoError(t, err)
}
