package cff

import (
	"context"
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
			Flow:   "flow",
			File:   "location/flow.go",
			Line:   42,
			Column: 84,
		})
		fe2 := e.FlowInit(&FlowInfo{
			Flow:   "flow",
			File:   "location/flow.go",
			Line:   42,
			Column: 84,
		})
		assert.True(t, fe == fe2)
	})
	t.Run("mutate object and not get it back", func(t *testing.T) {
		scope := tally.NewTestScope("", nil)
		e := TallyEmitter(scope)
		fe := e.FlowInit(&FlowInfo{
			Flow:   "flow",
			File:   "location/flow.go",
			Line:   42,
			Column: 84,
		})
		feMutated := fe.FlowFailedTask(context.Background(), "foo", nil)
		fe2 := e.FlowInit(&FlowInfo{
			Flow:   "flow",
			File:   "location/flow.go",
			Line:   42,
			Column: 84,
		})
		assert.False(t, feMutated == fe)
		assert.False(t, feMutated == fe2)
		assert.True(t, fe == fe2) // Already tested in "same object".
	})

}

// TestTallyEmitter_CacheTask verifies that we get back an initial cached object.
func TestTallyEmitter_CacheTask(t *testing.T) {
	t.Run("same object", func(t *testing.T) {
		scope := tally.NewTestScope("", nil)
		e := TallyEmitter(scope)
		te := e.TaskInit(
			&TaskInfo{
				Task:   "task",
				File:   "location/task.go",
				Line:   42,
				Column: 84,
			},
			&FlowInfo{
				Flow:   "flow",
				File:   "location/flow.go",
				Line:   42,
				Column: 84,
			})
		te2 := e.TaskInit(
			&TaskInfo{
				Task:   "another task",
				File:   "location/task.go",
				Line:   42,
				Column: 84,
			},
			&FlowInfo{
				Flow:   "flow",
				File:   "location/flow.go",
				Line:   42,
				Column: 84,
			})
		te2same := e.TaskInit(
			&TaskInfo{
				Task:   "task",
				File:   "location/task.go",
				Line:   42,
				Column: 84,
			},
			&FlowInfo{
				Flow:   "flow",
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
					Task:   "task",
					File:   "location/task.go",
					Line:   42,
					Column: 84,
				},
				&FlowInfo{
					Flow:   "flow",
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
	mutatedResults := make([]FlowEmitter, N/2)
	var wg sync.WaitGroup
	wg.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer wg.Done()
			<-start
			results[i] = e.FlowInit(
				&FlowInfo{
					Flow:   "flow",
					File:   "location/flow.go",
					Line:   42,
					Column: 84,
				})
			// Verifying that original cached FlowEmitter isn't mutated and we
			// return a new object.
			if i%2 == 0 {
				results[i].FlowFailedTask(context.Background(), "foo", nil)
			}
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
	wantMutated := mutatedResults[0]
	for i, got := range mutatedResults[1:] {
		require.True(t, got == wantMutated)
		require.False(t, got == results[i])
	}
}
