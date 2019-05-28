package fallbackwith_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	fallbackwith_gen "go.uber.org/cff/internal/tests/fallbackwith_gen"
)

func TestSerialRecovery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := fallbackwith_gen.Serial(nil, "bar")
		assert.NoError(t, err)
		assert.Equal(t, "foo", got)
	})

	t.Run("failure", func(t *testing.T) {
		s, err := fallbackwith_gen.Serial(errors.New("great sadness"), "bar")
		assert.NoError(t, err)
		assert.Equal(t, "bar", s)
	})
}
