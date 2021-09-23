package parallel

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimpleParallel(t *testing.T) {
	var m sync.Map
	require.NoError(t, Simple(&m))
	var count int
	m.Range(func(_, v_ interface{}) bool {
		count++
		return true
	})
	assert.Equal(t, 2, count)
}

func TestSimpleParallelWithError(t *testing.T) {
	ch := make(chan<- string, 1)
	require.Error(t, SimpleWithError(ch))
	assert.LessOrEqual(t, len(ch), 1)
}

func TestSimpleParallelWithPanic(t *testing.T) {
	err := SimpleWithPanic()
	require.Error(t, err)
	assert.Equal(t, "parallel function panic: sad times", err.Error())
}

func TestMultipleTasks(t *testing.T) {
	ch := make(chan<- string, 2)
	require.NoError(t, MultipleTasks(ch))
	assert.Len(t, ch, 2)
}
