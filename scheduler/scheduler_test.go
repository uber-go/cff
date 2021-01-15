package scheduler

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestScheduler(t *testing.T) {
	t.Parallel()

	type job struct {
		deps []int // indexes of dependencies in jobs list
		err  error // error to return, if any
		run  runLevel
	}

	type jobs []job

	type testCase struct {
		desc    string
		jobs    jobs
		wantErr error // error expected in Wait, if any
	}

	errSad := errors.New("great sadness")

	// Add comments documenting the dependency graph. Use the following
	// conventions.
	//
	//  A <- B      B depends on A
	//  {A, B}      both, A and B
	//  A!          A will fail
	//  A?          A may or may not run
	tests := []testCase{
		{
			// 0 <- {1, 2} <- 3
			desc: "diamond",
			jobs: jobs{
				{},                  // 0
				{deps: []int{0}},    // 1
				{deps: []int{0}},    // 2
				{deps: []int{1, 2}}, // 3
			},
		},
		{
			// 0! <- {1, 2} <- 3
			desc: "diamond/fail initial",
			jobs: jobs{
				{err: errSad},                        // 0
				{deps: []int{0}, run: mustNotRun},    // 1
				{deps: []int{0}, run: mustNotRun},    // 2
				{deps: []int{1, 2}, run: mustNotRun}, // 3
			},
			wantErr: errSad,
		},
		{
			// 0 <- {1?, 2!} <- 3
			// If 2 gets picked up early and aborts the run, then
			// 1 will not run.
			desc: "diamond/fail middle",
			jobs: jobs{
				{},                                   // 0
				{deps: []int{0}, run: mayRun},        // 1
				{deps: []int{0}, err: errSad},        // 2
				{deps: []int{1, 2}, run: mustNotRun}, // 3
			},
			wantErr: errSad,
		},
		{
			// 0 <- 1
			// 2 <- 3
			desc: "independent graph",
			jobs: jobs{
				{},               // 0
				{deps: []int{0}}, // 1
				{},               // 2
				{deps: []int{2}}, // 3
			},
		},
		{
			// 0! <- 1
			// 2? <- 3?
			// Based on scheduler performance, 2 and 3 may or may
			// not run.
			desc: "independent graph/fail part",
			jobs: jobs{
				{err: errSad},                     // 0
				{deps: []int{0}, run: mustNotRun}, // 1
				{run: mayRun},                     // 2
				{deps: []int{2}, run: mayRun},     // 3
			},
			wantErr: errSad,
		},
		{
			desc: "independent 100/no deps",
			jobs: make(jobs, 100),
		},
		{
			// 0 <- 1 <- 2 <- ... <- 100
			desc: "chain/100",
			jobs: func() (jobs jobs) {
				jobs = append(jobs, job{})
				for i := 0; i < 99; i++ {
					jobs = append(jobs, job{deps: []int{i}})
				}
				return jobs
			}(),
		},
		{
			// 0! <- 1 <- 2 <- ... <- 100
			desc: "chain/100/fail initial",
			jobs: func() (jobs jobs) {
				jobs = append(jobs, job{err: errSad})
				for i := 0; i < 99; i++ {
					jobs = append(jobs, job{
						deps: []int{i},
						run:  mustNotRun,
					})
				}
				return jobs
			}(),
			wantErr: errSad,
		},
	}

	runTestCase := func(t *testing.T, numWorkers int, tt testCase) {
		ctrl := newFakeJobController(t)
		defer ctrl.Verify()

		cfg := Config{Concurrency: numWorkers}
		sched := cfg.New()

		ctx := context.Background()
		jobs := make([]*ScheduledJob, len(tt.jobs))
		for i, job := range tt.jobs {
			deps := make([]*ScheduledJob, 0, len(job.deps))
			for _, dep := range job.deps {
				if dep >= i {
					t.Fatalf("job %d depends on job %d > %d", i, dep, i)
				}
				deps = append(deps, jobs[dep])
			}

			job := ctrl.NewJob(&fakeJobConfig{
				Deps:     deps,
				Run:      job.run,
				FailWith: job.err,
			})
			jobs[i] = sched.Enqueue(ctx, job)
		}

		err := sched.Wait(ctx)

		if tt.wantErr != nil {
			if err == nil {
				t.Error("expected failure, got success")
			}
			return
		}

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}

	for _, numWorkers := range []int{0, 1, 2, 4, 8} {
		numWorkers := numWorkers
		t.Run(fmt.Sprintf("workers=%d", numWorkers), func(t *testing.T) {
			t.Parallel()

			for _, tt := range tests {
				tt := tt
				t.Run(tt.desc, func(t *testing.T) {
					t.Parallel()

					runTestCase(t, numWorkers, tt)
				})
			}
		})
	}
}

func TestScheduler_Wait(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	cfg := Config{Concurrency: 0}
	sched := cfg.New()

	if err := sched.Wait(ctx); err != nil {
		t.Fatalf("Wait without enqueuing anything failed: %v", err)
	}

	t.Run("Enqueue after Wait", func(t *testing.T) {
		t.Parallel()

		defer func() {
			if recover() == nil {
				t.Error("Enqueue should panic after Wait, got success instead")
			}
		}()

		ctrl := newFakeJobController(t)
		defer ctrl.Verify()

		sched.Enqueue(ctx, ctrl.NewJob(&fakeJobConfig{Run: mustNotRun}))
	})
}

func TestScheduler_WaitAfterCanceled(t *testing.T) {
	t.Parallel()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	cfg := Config{Concurrency: 0}
	sched := cfg.New()

	if err := sched.Wait(ctx); err == nil {
		t.Error("Wait with canceled context should fail")
	}
}

func TestScheduler_EnqueueManyConcurrently(t *testing.T) {
	t.Parallel()

	const N = 100

	ctrl := newFakeJobController(t)
	defer ctrl.Verify()

	jobs := make([]Job, N)
	for i := 0; i < N; i++ {
		jobs[i] = ctrl.NewJob(&fakeJobConfig{})
	}

	ctx := context.Background()
	cfg := Config{Concurrency: 0}
	sched := cfg.New()

	// Goroutines use 'ready' to wait for each other so that we have a
	// higher chance of a race. We use `done` to wait for all these
	// goroutines to be finished.
	var ready, done sync.WaitGroup
	done.Add(N)
	ready.Add(N)
	for i := 0; i < N; i++ {
		go func(i int) {
			defer done.Done()

			ready.Done() // I'm ready
			ready.Wait() // ...but is everyone else?

			sched.Enqueue(ctx, jobs[i])
		}(i)
	}

	done.Wait()

	if err := sched.Wait(ctx); err != nil {
		t.Errorf("unexpected failure from Scheduler.Wait: %v", err)
	}
}

// Test that a nil emitter does not break the scheduler.
func TestScheduler_EmitterNil(t *testing.T) {
	t.Parallel()

	sched := Config{
		Concurrency:         1,
		StateFlushFrequency: time.Millisecond,
	}.New()

	sched.Enqueue(context.Background(), Job{
		Run: func(c context.Context) error {
			return nil
		},
	})

	assert.NoError(t, sched.Wait(context.Background()))
}

// emitterFn is a convenience type to block scheduler progress until Emit
// is called.
type emitterFn func(State)

func (e emitterFn) Emit(s State) {
	e(s)
}

// Test that the scheduler correctly emits state when there is no scheduled
// activity.
func TestScheduler_SchedulerEmpty(t *testing.T) {
	t.Parallel()

	currGoMaxProcs := runtime.GOMAXPROCS(1)
	defer runtime.GOMAXPROCS(currGoMaxProcs)

	var called bool
	done := make(chan struct{})
	emitter := emitterFn(func(s State) {
		if called {
			return
		}
		called = true
		assert.Equal(t, State{
			Concurrency: _minDefaultWorkers,
			IdleWorkers: _minDefaultWorkers,
		}, s)
		close(done)
	})

	sched := Config{
		Emitter:             emitter,
		StateFlushFrequency: time.Millisecond,
	}.New()

	<-done

	if err := sched.Wait(context.Background()); err != nil {
		t.Errorf("unexpected failure from Scheduler.Wait: %v", err)
	}
}

// Test that the scheduler correctly emits state when there is a currently
// scheduled job.
func TestScheduler_EmitSingleJob(t *testing.T) {
	t.Parallel()

	emitter, statec := newChannelEmitter()

	sched := Config{
		Concurrency:         1,
		Emitter:             emitter,
		StateFlushFrequency: time.Millisecond,
	}.New()

	blocker := newBlocker()
	sched.Enqueue(context.Background(), Job{
		Run: blocker.Run,
	})

	s := awaitStableState(t, statec)
	assert.Equal(t, 1, s.Pending)

	blocker.AwaitRunning()
	assert.Equal(t, State{
		Pending:     1,
		Ready:       0,
		IdleWorkers: 0,
		Concurrency: 1,
	}, awaitStableState(t, statec))

	blocker.UnblockAndWait()

	s = awaitStableState(t, statec)
	assert.Equal(t, 0, s.Pending)
	assert.Equal(t, 1, s.IdleWorkers)

	if err := sched.Wait(context.Background()); err != nil {
		t.Errorf("unexpected failure from Scheduler.Wait: %v", err)
	}
}

// Test that the scheduler correctly emits state when a job is waiting on
// another unrelated job to finish.
func TestScheduler_EmitTwoIndependentJobs(t *testing.T) {
	t.Parallel()

	emitter, statec := newChannelEmitter()

	sched := Config{
		Concurrency:         2,
		Emitter:             emitter,
		StateFlushFrequency: time.Millisecond,
	}.New()

	blockerA := newBlocker()
	sched.Enqueue(context.Background(), Job{
		Run: blockerA.Run,
	})

	s := awaitStableState(t, statec)
	assert.Equal(t, 1, s.Pending)

	// A must be running before we schedule B to ensure that B does not
	// get scheduled first, breaking our assertions below.
	blockerA.AwaitRunning()

	assert.Equal(t, State{
		Pending:     1,
		Ready:       0,
		IdleWorkers: 1,
		Waiting:     0,
		Concurrency: 2,
	}, awaitStableState(t, statec))

	blockerB := newBlocker()
	sched.Enqueue(context.Background(), Job{
		Run: blockerB.Run,
	})

	blockerB.AwaitRunning()
	assert.Equal(t, State{
		Pending:     2,
		Ready:       0,
		Waiting:     0,
		IdleWorkers: 0,
		Concurrency: 2,
	}, awaitStableState(t, statec))

	blockerA.UnblockAndWait()
	blockerB.UnblockAndWait()

	s = awaitStableState(t, statec)
	assert.Equal(t, 0, s.Pending)
	assert.Equal(t, 2, s.IdleWorkers)

	if err := sched.Wait(context.Background()); err != nil {
		t.Errorf("unexpected failure from Scheduler.Wait: %v", err)
	}
}

// Test that the scheduler correctly emits state when there are two scheduled
// jobs with one job waiting on the other.
func TestScheduler_EmitTwoDependentJobs(t *testing.T) {
	t.Parallel()

	emitter, statec := newChannelEmitter()

	sched := Config{
		Concurrency:         2,
		Emitter:             emitter,
		StateFlushFrequency: time.Millisecond,
	}.New()

	blockerA := newBlocker()
	scheduledA := sched.Enqueue(context.Background(), Job{
		Run: blockerA.Run,
	})

	blockerB := newBlocker()
	sched.Enqueue(context.Background(), Job{
		Run: blockerB.Run,
		Dependencies: []*ScheduledJob{
			scheduledA,
		},
	})

	s := awaitStableState(t, statec)
	assert.Equal(t, 1, s.Waiting)

	// When all dependencies for a job have run, that job could be
	// ready or running.
	blockerA.UnblockAndWait()

	blockerA.AwaitRunning()
	assert.Equal(t, State{
		Pending:     1,
		Ready:       0,
		Waiting:     0,
		IdleWorkers: 1,
		Concurrency: 2,
	}, awaitStableState(t, statec))

	blockerB.UnblockAndWait()

	s = awaitStableState(t, statec)
	assert.Equal(t, 0, s.Pending)
	assert.Equal(t, 2, s.IdleWorkers)

	if err := sched.Wait(context.Background()); err != nil {
		t.Errorf("unexpected failure from Scheduler.Wait: %v", err)
	}
}

// blocker is a testing convenience object that runs to block the execution
// of jobs until Unblock is called.
type blocker struct {
	done    chan struct{}
	proceed chan struct{}
	running chan struct{}
}

func newBlocker() blocker {
	return blocker{
		done:    make(chan struct{}),
		proceed: make(chan struct{}),
		running: make(chan struct{}),
	}
}

// Run called inside a job causes execution to block until the
// corresponding UnblockAndWait is called.
func (j blocker) Run(context.Context) error {
	close(j.running)
	defer close(j.done)
	<-j.proceed
	return nil
}

func (j blocker) UnblockAndWait() {
	close(j.proceed)
	<-j.done
}

func (j blocker) AwaitRunning() {
	<-j.running
}

// channelEmitter is an Emitter that posts scheduler state to a channel.
//
// The Emitter does not block on channel writes and drops messages if the
// receiver is slow.
type channelEmitter struct {
	statec chan State
}

func newChannelEmitter() (Emitter, <-chan State) {
	state := make(chan State, 1)
	return channelEmitter{statec: state}, state
}

func (t channelEmitter) Emit(s State) {
	select {
	case t.statec <- s:
	default:
	}
}

// awaitStableState polls a State channel until a stable State is found to
// minimize the effect of a race condition while asserting the state of
// the scheduler.
// The race is that we can't coordinate the timing between when the state of
// the scheduler changes in response to tasks finishing and when the state
// assertion is evaluated. In the tests, we only unblock the task or know
// that the defer statement has fired in a job. Here, the scheduler state
// could be updated before or after the assertion of the state is evaluated.
func awaitStableState(t *testing.T, ch <-chan State) State {
	const (
		// stableN is the threshold of consistent observations considered
		// stable. This value was found by tuning dependent tests with
		// --runs_per_test=100.
		stableN int = 5
		// maxAttempts is the num of observations allowed to find a stable
		// state. This value was found by tuning dependent tests with
		// --runs_per_test=100.
		maxAttempts int = 3
	)
	prevState := <-ch
attempt:
	for i := 0; i < maxAttempts; i++ {
		for run := 1; run < stableN; run++ {
			s := <-ch
			if prevState != s {
				prevState = s
				continue attempt
			}
		}
		return prevState
	}
	t.Fatalf("failed to find stable state after %d attempts", maxAttempts)
	return State{}
}

func TestIdleWorkers(t *testing.T) {
	tests := []struct {
		desc                       string
		concurrency, ongoing, want int
	}{
		{
			desc:        "all idle",
			concurrency: 3,
			ongoing:     0,
			want:        3,
		},
		{
			desc:        "some idle",
			concurrency: 3,
			ongoing:     1,
			want:        2,
		},
		{
			desc:        "none idle",
			concurrency: 3,
			ongoing:     3,
			want:        0,
		},
		{
			desc:        "more jobs than concurrency",
			concurrency: 3,
			ongoing:     5,
			want:        0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, idleWorkers(tt.concurrency, tt.ongoing))
		})
	}
}
