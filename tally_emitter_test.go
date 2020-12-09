package cff

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/uber-go/tally"
)

// TestTallyEmitter_CacheFlow verifies that we get back an initial cached object.
func TestTallyEmitter_CacheFlow(t *testing.T) {
	t.Run("same object", func(t *testing.T) {
		scope := tally.NewTestScope("", nil)
		e := TallyEmitter(scope)
		fe := e.FlowInit(&FlowInfo{
			Name:   "flow",
			File:   "location/flow.go",
			Line:   42,
			Column: 84,
		})
		fe2 := e.FlowInit(&FlowInfo{
			Name:   "flow",
			File:   "location/flow.go",
			Line:   42,
			Column: 84,
		})
		assert.True(t, fe == fe2)
	})
}

// TestTallyEmitter_CacheTask verifies that we get back an initial cached object.
func TestTallyEmitter_CacheTask(t *testing.T) {
	t.Run("same object", func(t *testing.T) {
		scope := tally.NewTestScope("", nil)
		e := TallyEmitter(scope)
		te := e.TaskInit(
			&TaskInfo{
				Name:   "task",
				File:   "location/task.go",
				Line:   42,
				Column: 84,
			},
			&FlowInfo{
				Name:   "flow",
				File:   "location/flow.go",
				Line:   42,
				Column: 84,
			})
		te2 := e.TaskInit(
			&TaskInfo{
				Name:   "another task",
				File:   "location/task.go",
				Line:   42,
				Column: 84,
			},
			&FlowInfo{
				Name:   "flow",
				File:   "location/flow.go",
				Line:   42,
				Column: 84,
			})
		te2same := e.TaskInit(
			&TaskInfo{
				Name:   "task",
				File:   "location/task.go",
				Line:   42,
				Column: 84,
			},
			&FlowInfo{
				Name:   "flow",
				File:   "location/flow.go",
				Line:   42,
				Column: 84,
			})

		assert.True(t, te == te2same)
		assert.False(t, te2 == te)
		assert.False(t, te2 == te2same)
	})
}

// TestTallyEmitter_CacheReadRace tests contention around initial caching of TaskEmitter obj.
func TestTallyEmitter_CacheTaskReadRace(t *testing.T) {
	const N = 1000
	scope := tally.NewTestScope("", nil)
	e := TallyEmitter(scope)
	start := make(chan struct{})
	results := make([]TaskEmitter, N)
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			<-start
			results[i] = e.TaskInit(
				&TaskInfo{
					Name:   "task",
					File:   "location/task.go",
					Line:   42,
					Column: 84,
				},
				&FlowInfo{
					Name:   "flow",
					File:   "location/flow.go",
					Line:   42,
					Column: 84,
				})
			// Don't need a lock because we're filling in independent indexes of
			// a pre-allocated slice.
		}(i)
	}
	close(start)
	wg.Wait()
	want := results[0]
	for _, got := range results[1:] {
		require.True(t, got == want)
		// require because we don't want to log a 1000 messages
	}
}

// TestTallyEmitter_CacheReadRace tests contention around initial caching of FlowEmitter obj.
func TestTallyEmitter_CacheFlowReadRace(t *testing.T) {
	const N = 1000
	scope := tally.NewTestScope("", nil)
	e := TallyEmitter(scope)
	start := make(chan struct{})
	results := make([]FlowEmitter, N)
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			<-start
			results[i] = e.FlowInit(
				&FlowInfo{
					Name:   "flow",
					File:   "location/flow.go",
					Line:   42,
					Column: 84,
				})
			// Don't need a lock because we're filling in independent indexes of
			// a pre-allocated slice.
		}(i)
	}
	close(start)
	wg.Wait()
	want := results[0]
	for _, got := range results[1:] {
		require.True(t, got == want)
		// require because we don't want to log a 1000 messages
	}
}
