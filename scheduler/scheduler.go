// Package scheduler implements a runtime scheduler for CFF with support for
// interdependent jobs.
//
// To use the scheduler, build one with New, providing the desired maximum
// number of goroutines. This defaults to the number of CPUs available.
//
//	sched := scheduler.New(n)
//
// With a scheduler available, enqueue jobs into it with the Enqueue method.
//
//	j1 := sched.Enqueue(ctx, Job{..})
//
// The scheduler will begin running this job as soon as a worker is available.
//
// Enqueue returns a reference to the scheduled job. Use this reference in
// other Enqueue calls to specify dependencies for jobs.
//
//	j3 := sched.Enqueue(ctx, Job{
//	  ...,
//	  Dependencies: []*scheduler.ScheduledJob{j1, j2},
//	})
//
// j3 will not be run until j1 and j2 have finished successfully.
//
// Dependencies must be enqueued before jobs that depend on them. This adds
// the burden of dependency order resolution on the caller.
//
// After enqueuing all jobs, await completion with sched.Wait. This is
// comparable to WaitGroup.Wait().
//
//	err := sched.Wait(ctx)
//
// If any of the enqueued jobs failed, the remaining jobs will be aborted and
// sched.Wait will return the error.
package scheduler

import (
	"container/list"
	"context"
	"errors"
	"runtime"
	"time"

	"go.uber.org/multierr"
)

const (
	// _minDefaultWorkers sets minimum number of workers we'll spawn by default if
	// not explicitly specified by the user.
	_minDefaultWorkers = 4

	// defaultStateFlushFrequency sets the frequency at which the scheduler emits
	// metrics on its state.
	_defaultStateFlushFrequency = 100 * time.Millisecond
)

// --------------------
// IMPLEMENTATION NOTES
// --------------------
//
// There are three kinds of goroutines at play here.
//
// Caller
//   This is the goroutine that calls scheduler.New(n), Scheduler.Enqueue,
//   and Scheduler.Wait. This is the point in the code where fan-out begins
//   (Scheduler.Enqueue) and ends (Scheduler.Wait).
//
// Workers
//   One or more worker goroutines run the scheduled jobs. These are the
//   simplest component: they pull jobs off a channel, run them, and post
//   results to a different channel.
//
// Scheduler Loop
//   The Scheduler Loop runs in the background, manages internal state, and
//   orchestrates the workers. The Scheduler Loop decides which jobs are ready
//   to be executed, posts them to workers, and processes results coming from
//   these workers.
//
// We can keep the core scheduler logic lockless because all state management
// is deferred to the same goroutine: the Scheduler Loop. DO NOT read or write
// internal state outside that goroutine.
//
//  .-----.
//  |     |----------.
//  |     |          |
//  |     |          v
//  |  C  |     +--------------------------------------------+
//  |  A  |     | enqueuec := make(chan *ScheduledJob, 1)    |----------.
//  |  L  |     +--------------------------------------------+          |
//  |  L  |                                                             v
//  |  E  |                                                          .-----.
//  |  R  |     +--------------------------------------------+       |     |
//  |     |     | readyc := make(chan *ScheduledJob)         |<------|  S  |
//  |     |     +--------------------------------------------+       |  C  |
//  |     |          |              |                  |             |  H  |
//  '-----'          |              |                  |             |  E  |
//                   v              v                  v             |  D  |
//              .----------.   .----------.       .----------.       |  U  |
//              | Worker 1 |   | Worker 2 |  ...  | Worker N |       |  L  |
//              '----------'   '----------'       '----------'       |  E  |
//                   |              |                  |             |  R  |
//                   |              |                  |             |     |
//                   v              v                  v             |  L  |
//              +--------------------------------------------+       |  O  |
//              | donec := make(chan jobResult, N)           |------>|  O  |
//              +--------------------------------------------+       |  P  |
//                                                                   |     |
//                                                                   '-----'

type jobResult struct {
	Job *ScheduledJob // job that was executed
	Err error         // failure, if any
}

// errJobInvalid is a sentinel error to capture jobs that are invalidated
// because their dependencies have errored or are marked invalid.
var errJobInvalid = errors.New("job invalid")

// worker implements the logic for a worker goroutine. Workers may read from
// the ScheduledJob but they MUST NOT modify it. All output from workers should
// be sent up to the scheduler via jobResult. The scheduler loop is the only
// goroutine permitted to modify ScheduledJob.
//
// NOTE: If you rename this function, update _workerFunction in
// internal/tests/setconcurrency/setconcurrency.go.
func worker(readyc <-chan *ScheduledJob, donec chan<- jobResult) {
	var (
		currentJob  *ScheduledJob
		exitCleanly bool
	)
	defer func() {
		if exitCleanly {
			return
		}
		donec <- jobResult{Job: currentJob, Err: errors.New("job exited unexpectedly")}
		go worker(readyc, donec)
	}()

	for j := range readyc {
		res := jobResult{Job: j}
		currentJob = j

		if err := j.ctx.Err(); err != nil {
			// Don't run if context already cancelled.
			res.Err = err
		} else if j.invalid {
			// Don't run if marked as invalid.
			res.Err = errJobInvalid
		} else {
			res.Err = j.run(j.ctx)
		}
		currentJob = nil
		donec <- res
	}
	exitCleanly = true
}

// Scheduler schedules jobs for a CFF flow based on their dependencies.
type Scheduler struct {
	// Closed when the Scheduler Loop exits.
	finishedc chan struct{}

	// Error encountered while running the jobs, if any.
	err error

	// Enqueue posts partially initialized ScheduledJobs to this channel
	// without inspecting any internal state.
	//
	// The Scheduler Loop initializes the object fully, with access to
	// internal state.
	enqueuec chan *ScheduledJob

	// The Scheduler Loop posts jobs that are ready to be executed by
	// workers to this channel.
	readyc chan<- *ScheduledJob

	// Workers post results of executed jobs to this channel.
	donec <-chan jobResult

	// Concurrency is the number of workers the scheduler can process tasks
	// with.
	concurrency int

	// If true when a job fails, directs the scheduler to record its failure,
	// invalidate all jobs that depend on the failed job, and keep running.
	continueOnError bool
}

// Config stores parameters the scheduler should run with and is the
// entry point for running the scheduler.
type Config struct {
	// Concurrency is the number of concurrent workers to schedule tasks to.
	Concurrency int
	// Emitter provides a hook into the state of the scheduler.
	Emitter Emitter
	// StateFlushFrequency is how often the scheduler will emit metrics with the
	// emitter. This defaults to 100 milliseconds.
	StateFlushFrequency time.Duration
	// ContinueOnError, if true when a job fails, directs the scheduler to
	// record its failure, invalidate all jobs that depend on the failed job,
	// and keep running.
	ContinueOnError bool
}

// New begins execution of a flow with the provided number of
// goroutines. Concurrency defaults to max(GOMAXPROCS, 4) if zero.
//
// Enqueue jobs into the returned scheduler using the Enqueue method, and wait
// for the result with Wait.
func (c Config) New() *Scheduler {
	if c.Concurrency == 0 {
		c.Concurrency = runtime.GOMAXPROCS(0)
		if c.Concurrency < _minDefaultWorkers {
			c.Concurrency = _minDefaultWorkers
		}
	}

	if c.StateFlushFrequency == 0 {
		c.StateFlushFrequency = _defaultStateFlushFrequency
	}

	// Channel size 1: Support enqueuing one additional job when the
	// scheduler is busy.
	enqueuec := make(chan *ScheduledJob, 1)

	// Unbuffered channel: If all workers are busy, don't schedule more
	// work. This ensures that if we quit early, workers don't do
	// additional work that will be thrown away. For example, given jobs
	// [A, B, C], if A failed, we shouldn't post B and C to a buffered
	// channel because if we do that, the workers will not exit until
	// after they've run B and C, the results for which will be discarded
	// anyway because A failed.
	readyc := make(chan *ScheduledJob)

	// Channel size should match concurrency: Workers should always be
	// able to post their results, even if the Scheduler Loop is busy.
	donec := make(chan jobResult, c.Concurrency)

	// Start the workers.
	go func() {
		// TODO(abg): Maybe we should spawn workers on demand as
		// needed with a maximum of N workers instead of spawning them
		// in advance.
		for i := 0; i < c.Concurrency; i++ {
			go worker(readyc, donec)
		}
	}()

	sched := &Scheduler{
		enqueuec:        enqueuec,
		readyc:          readyc,
		donec:           donec,
		finishedc:       make(chan struct{}),
		concurrency:     c.Concurrency,
		continueOnError: c.ContinueOnError,
	}

	// We lie to the caller about the number of goroutines. Spawn one
	// extra goroutine for the Scheduler Loop.
	go sched.run(c.Emitter, c.StateFlushFrequency)

	return sched
}

// Job is an independent executable unit meant to be executed by the
// scheduler.
type Job struct {
	// Run executes the job and returns the error it encountered, if any.
	Run func(context.Context) error

	// Dependencies are previously enqueued jobs that must run before this
	// job.
	Dependencies []*ScheduledJob
}

// ScheduledJob is a job that has been scheduled for execution by the
// scheduler.
type ScheduledJob struct {
	// The following fields are initialized in Scheduler.Enqueue. These
	// are read-only. They MUST NOT be changed once initialized.

	ctx  context.Context
	run  func(context.Context) error
	deps []*ScheduledJob

	// The following fields track the internal state of the job. These are
	// read-write, but only within Scheduler.run. DO NOT read or write
	// them outside scheduler.run, as that will introduce a data race.

	remaining int             // number of jobs we're waiting for
	consumers []*ScheduledJob // jobs waiting for this job
	done      bool            // whether this was run, regardless of success or failure
	err       error           // the job error, if encountered when the job ran
	invalid   bool            // whether the job is marked invalid and should not run

	// NOTE: DO NOT add methods to ScheduledJob. There's danger of using
	// methods that read or write internal state outside the Scheduler.run
	// function which, as discussed above, introduces a data race.
}

// Enqueue queues up a job for execution with the scheduler.
// The returned object may be used as a dependency for other jobs.
//
// Enqueue will panic if called after calling Wait.
func (s *Scheduler) Enqueue(ctx context.Context, j Job) *ScheduledJob {
	// Enqueue is invoked from the Caller goroutine, which is running at
	// the same time as the Scheduler Loop. To avoid data races here,
	// Enqueue MUST NOT access any internal state. To that end, Enqueue
	// places a partially initialized object into the enqueuec channel,
	// and the Scheduler Loop initializes the rest of it.
	pj := &ScheduledJob{
		ctx:  ctx,
		run:  j.Run,
		deps: j.Dependencies,
	}
	s.enqueuec <- pj // panics if closed
	return pj
}

// run implements the Scheduler Loop. The Scheduler Loop works by maintaining
// the ready list, which contains jobs ready to be run, with no outstanding dependencies.
//
// Each tick of the loop runs one of the following branches:
//
//   - Attempt to schedule a job if `ready` is non-empty and a worker is
//     available.
//   - Process a newly Enqueued job, placing it in `ready` if it's ready to be executed.
//   - If a job finished running, signal jobs that were awaiting
//     its completion. Those that have no more dependencies outstanding are
//     moved to the `ready` list.
func (s *Scheduler) run(emitter Emitter, freq time.Duration) {
	defer close(s.finishedc) // unblock Wait()
	defer close(s.readyc)    // kill workers

	// Upon exit, drain enqueuec. This is necessary because the caller
	// goroutine will roughly take the following form, where tasks begin
	// executing as soon as possible.
	//
	//  sched.Enqueue(ctx, j1)
	//  sched.Enqueue(ctx, j2)
	//  // ...
	//  sched.Enqueue(ctx, jN)
	//  err := sched.Wait(ctx)
	//
	// If j1 fails and causes the Scheduler Loop to exit early, we still
	// need to process the remaining Enqueue invocations so that we get
	// to sched.Wait.
	defer func() {
		for range s.enqueuec {
		}
	}()

	var tickerC <-chan time.Time
	if emitter != nil {
		// Note: Phab marks this block as untested, but we believe this is
		// tested (GM-876).
		ticker := time.NewTicker(freq)
		defer ticker.Stop()

		tickerC = ticker.C
	}

	// Jobs ready to be thrown into the ready channel.
	ready := list.New() // []*ScheduledJob
	// TODO(abg): Use a maxheap here based on the number of consumers.
	// That way, we'll run jobs that unblock the most consumers first.

	// Number of jobs that are executing.
	ongoing := 0

	// Total number of jobs in flight. This includes jobs that are
	// executing or waiting to be executed.
	pending := 0

	// Number of jobs waiting for other jobs to finish.
	waiting := 0

	// Tracks whether we're still expecting new Enqueue calls. After this
	// is set to nil, we don't expect new Enqueue requests.
	enqueuec := s.enqueuec

	for {
		// If there's at least one job ready to be executed, grab it.
		// If no jobs are ready, this leaves `readyc` as nil. Trying
		// to insert into a nil channel never resolves so the select
		// will never pick that path.
		readyc := s.readyc
		var (
			nextEl *list.Element
			next   *ScheduledJob
		)
		if ready.Len() > 0 {
			nextEl = ready.Front()
			next = nextEl.Value.(*ScheduledJob)
		} else {
			readyc = nil
		}

		select {
		case readyc <- next:
			// Remove from the ready queue only if we scheduled in
			// this iteration.
			ready.Remove(nextEl)

			ongoing++

		case job, ok := <-enqueuec:
			// Wait was called and the enqueue channel was closed.
			// Make sure we never hit this branch of the select
			// again. (A nil channel never resolves.)
			if !ok {
				enqueuec = nil
				break
			}

			// Ask to be notified when dependencies are run --
			// unless they've already been run.
			//
			// If a dependency has already run and errored, mark the job as
			// invalid.
			for _, dep := range job.deps {
				if dep.done {
					if dep.err != nil {
						job.invalid = true
					}
					continue
				}
				dep.consumers = append(dep.consumers, job)
				job.remaining++
			}

			pending++

			// No outstanding dependencies. Ready to run.
			if job.remaining == 0 {
				ready.PushBack(job)
			} else {
				waiting++
			}

		case res := <-s.donec:
			job := res.Job
			job.done = true

			pending--
			ongoing--

			if err := res.Err; err != nil {
				job.err = err

				// Record the failure and return early if the job
				// failed.
				if !s.continueOnError {
					s.err = err
					return
				}
				// With continueOnError, mark invalid directly dependent jobs,
				// append non-sentinel errors, and continue the scheduler loop.
				if !errors.Is(err, errJobInvalid) {
					s.err = multierr.Append(s.err, err)
				}
				for _, consumer := range job.consumers {
					consumer.invalid = true
				}
			}

			// Notify jobs waiting on this job, adding them to
			// ready if this was their last outstanding
			// dependency.
			for _, consumer := range job.consumers {
				consumer.remaining--
				if consumer.remaining == 0 {
					waiting--
					ready.PushBack(consumer)
				}
			}

		case <-tickerC:
			// If emitter is nil, tickerC will be a nil channel that
			// never resolves.
			// Note: Phab marks this line as untested, but we believe this is
			// tested (GM-876).
			emitter.Emit(
				State{
					Pending:     pending,
					Ready:       ready.Len(),
					Waiting:     waiting,
					IdleWorkers: idleWorkers(s.concurrency, ongoing),
					Concurrency: s.concurrency,
				},
			)
		}

		// If all enqueued jobs have been finished and no new enqueues
		// are allowed, we can exit.
		if pending == 0 && enqueuec == nil {
			return
		}
	}
}

// Wait waits for all scheduled jobs to finish by default it returns the first
// error encountered, if any.
//
// If ContinueOnError was specified, the returned error will combine all
// errors from jobs failures.
//
// No new jobs may be enqueued once Wait is called.
func (s *Scheduler) Wait(ctx context.Context) error {
	close(s.enqueuec) // disallow new Enqueues
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-s.finishedc: // wait for Scheduler Loop to exit
		err := s.err
		// If both channels are ready to read from, select will pick
		// one randomly. In that case, if there was no job failure,
		// pick the context-level failure in case the context timed
		// out at the same time the job finished.
		if err == nil {
			err = ctx.Err()
		}
		return err
	}
}

// idleWorkers tracks the difference of scheduler workers and ongoing jobs.
// Fewer jobs than workers mean there are idle workers.
func idleWorkers(concurrency, ongoing int) int {
	idle := concurrency - ongoing
	if idle < 0 {
		// It's impossible to have more ongoing jobs than available workers,
		// but we should guard against it.
		idle = 0
	}
	return idle
}
