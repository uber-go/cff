package main

import (
	"testing"

	"go.uber.org/cff/mode"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenMode(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		m, err := genMode("source-map")
		require.NoError(t, err)
		assert.Equal(t, m, mode.SourceMap)
	})

	t.Run("error", func(t *testing.T) {
		_, err := genMode("sad")
		assert.EqualError(t, err, `"unknown" is an invalid CFF generation mode. Argument was "sad"`)
	})
}
