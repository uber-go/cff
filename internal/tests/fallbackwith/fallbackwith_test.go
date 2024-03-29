package fallbackwith

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSerialRecovery(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := Serial(nil, "bar")
		assert.NoError(t, err)
		assert.Equal(t, "foo", got)
	})

	t.Run("failure", func(t *testing.T) {
		s, err := Serial(errors.New("great sadness"), "bar")
		assert.NoError(t, err)
		assert.Equal(t, "bar", s)
	})
}

func TestNoOutputRecovery(t *testing.T) {
	err := NoOutput()
	assert.NoError(t, err)
}

func TestPanic(t *testing.T) {
	s, err := Panic()
	assert.NoError(t, err)
	assert.Equal(t, "fallback", s)
}
