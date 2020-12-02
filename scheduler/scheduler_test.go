package scheduler

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
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
		sched := cfg.Begin()

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
	sched := cfg.Begin()

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
	sched := cfg.Begin()

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
	sched := cfg.Begin()

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
