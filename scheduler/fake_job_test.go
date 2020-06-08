package scheduler

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"go.uber.org/atomic"
)

// fakeJobController provides a customizable mock-like API to build jobs to
// test the scheduler.
//
//  ctrl := newFakeJobController(t)
//  defer ctrl.Verify()
//
//  job := ctrl.NewJob(...)
type fakeJobController struct {
	t    *testing.T
	jobs []*fakeJob
}

func newFakeJobController(t *testing.T) *fakeJobController {
	return &fakeJobController{t: t}
}

// Verify expectations set on the jobs.
func (c *fakeJobController) Verify() {
	for _, j := range c.jobs {
		j.verify()
	}
}

// Run expectations for jobs.
type runLevel int

const (
	// The job MUST run. This is the default.
	mustRun runLevel = iota

	// The job may or may not run. Use this for jobs where there's
	// ambiguity in whether they'll run depending on the scheduler's
	// behavior.
	mayRun

	// The job MUST NOT run.
	mustNotRun
)

type fakeJobConfig struct {
	Deps []*ScheduledJob

	// Error to return on job execution, if any.
	FailWith error

	// By default, all jobs MUST run. use this for more fine-grained
	// control on run expectations.
	Run runLevel
}

// NewJob constructs a new fake job.
func (c *fakeJobController) NewJob(cfg *fakeJobConfig) Job {
	j := &fakeJob{
		t:        c.t,
		name:     strconv.Itoa(len(c.jobs)),
		runLevel: cfg.Run,
		failWith: cfg.FailWith,
	}

	c.jobs = append(c.jobs, j)
	return Job{
		Run:          j.run,
		Dependencies: cfg.Deps,
	}
}

type fakeJob struct {
	t *testing.T

	name     string
	failWith error
	runLevel runLevel
	ran      atomic.Bool
}

func (j *fakeJob) verify() {
	if j.runLevel == mustRun && !j.ran.Load() {
		j.t.Errorf("job %q did not run, but was expected to run", j.name)
	}
}

func (j *fakeJob) run(ctx context.Context) error {
	if j.runLevel == mustNotRun {
		err := fmt.Errorf("job %q must not run", j.name)
		j.t.Error(err)
		return err
	}

	if j.ran.Swap(true) {
		err := fmt.Errorf("job %q ran multiple times", j.name)
		j.t.Error(err)
		return err
	}

	return j.failWith
}
