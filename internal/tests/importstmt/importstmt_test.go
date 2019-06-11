package importstmt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestImportStmt(t *testing.T) {
	i, err := Flow()
	assert.NoError(t, err)
	assert.Equal(t, 123, i)
}
