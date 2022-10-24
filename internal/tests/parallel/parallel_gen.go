//go:build !cff
// +build !cff

package parallel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/cff"
)

// TasksAndTask runs cff.Tasks and cff.Task in parallel to populate the
// provided map.
func TasksAndTask(m *sync.Map) error {
	return func() (err error) {

		_19_3 := context.Background()

		_20_19 := 2

		_22_4 := func() {
			m.Store("foo", "bar")
		}

		_25_4 := func(_ context.Context) {
			m.Store("fiz", "buzz")
		}

		_30_4 := func(_ context.Context) {
			m.Store("go", "lang")
		}
		ctx := _19_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   18,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _20_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:22:4
		task0 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task0.emitter = cff.NopTaskEmitter()
		task0.fn = func(ctx context.Context) (err error) {
			taskEmitter := task0.emitter
			startTime := time.Now()
			defer func() {
				if task0.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task0.ran.Store(true)

			_22_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task0.fn,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:25:4
		task1 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task1.emitter = cff.NopTaskEmitter()
		task1.fn = func(ctx context.Context) (err error) {
			taskEmitter := task1.emitter
			startTime := time.Now()
			defer func() {
				if task1.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task1.ran.Store(true)

			_25_4(ctx)

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task1.fn,
		})
		tasks = append(tasks, task1)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:30:4
		task2 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task2.emitter = cff.NopTaskEmitter()
		task2.fn = func(ctx context.Context) (err error) {
			taskEmitter := task2.emitter
			startTime := time.Now()
			defer func() {
				if task2.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task2.ran.Store(true)

			_30_4(ctx)

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task2.fn,
		})
		tasks = append(tasks, task2)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:33*/
	}()
}

// TasksWithError runs a parallel cff.Tasks that errors.
func TasksWithError() error {
	return func() (err error) {

		_40_3 := context.Background()

		_41_19 := 2

		_43_4 := func() error {
			return errors.New("sad times")
		}
		ctx := _40_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   39,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _41_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:43:4
		task3 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task3.emitter = cff.NopTaskEmitter()
		task3.fn = func(ctx context.Context) (err error) {
			taskEmitter := task3.emitter
			startTime := time.Now()
			defer func() {
				if task3.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task3.ran.Store(true)

			err = _43_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task3.fn,
		})
		tasks = append(tasks, task3)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:46*/
	}()
}

// TasksWithPanic runs a parallel cff.Tasks that panics.
func TasksWithPanic() error {
	return func() (err error) {

		_53_3 := context.Background()

		_54_19 := 2

		_56_4 := func() {
			panic("sad times")
		}
		ctx := _53_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   52,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _54_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:56:4
		task4 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task4.emitter = cff.NopTaskEmitter()
		task4.fn = func(ctx context.Context) (err error) {
			taskEmitter := task4.emitter
			startTime := time.Now()
			defer func() {
				if task4.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task4.ran.Store(true)

			_56_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task4.fn,
		})
		tasks = append(tasks, task4)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:59*/
	}()
}

// MultipleTasks runs multiple cff.Tasks in parallel to populate the provided
// channel.
func MultipleTasks(c chan<- string) error {
	send := func(_ context.Context) error {
		c <- "send"
		return nil
	}
	return func() (err error) {

		_71_3 := context.Background()

		_72_19 := 2

		_74_4 := func() {
			c <- "multiple"
		}

		_79_4 := send
		ctx := _71_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   70,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _72_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:74:4
		task5 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task5.emitter = cff.NopTaskEmitter()
		task5.fn = func(ctx context.Context) (err error) {
			taskEmitter := task5.emitter
			startTime := time.Now()
			defer func() {
				if task5.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task5.ran.Store(true)

			_74_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task5.fn,
		})
		tasks = append(tasks, task5)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:79:4
		task6 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task6.emitter = cff.NopTaskEmitter()
		task6.fn = func(ctx context.Context) (err error) {
			taskEmitter := task6.emitter
			startTime := time.Now()
			defer func() {
				if task6.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task6.ran.Store(true)

			err = _79_4(ctx)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task6.fn,
		})
		tasks = append(tasks, task6)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:80*/
	}()
}

// ContextErrorBefore runs a cff.Tasks function to test that the function is
// not run if the context errors before scheduler execution.
func ContextErrorBefore(ctx context.Context, src, target []int) error {
	return func() (err error) {

		_88_3 := ctx

		_89_19 := 2

		_91_4 := func() {
			target[0] = src[0]
		}
		ctx := _88_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   87,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _89_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:91:4
		task7 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task7.emitter = cff.NopTaskEmitter()
		task7.fn = func(ctx context.Context) (err error) {
			taskEmitter := task7.emitter
			startTime := time.Now()
			defer func() {
				if task7.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task7.ran.Store(true)

			_91_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task7.fn,
		})
		tasks = append(tasks, task7)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:94*/
	}()
}

// ContextErrorInFlight runs a cff.Tasks function to test that the function is
// not run if the context errors during scheduler execution.
func ContextErrorInFlight(ctx context.Context, cancel func(), src, target []int) error {
	blocker := make(chan struct{})
	return func() (err error) {

		_103_3 := ctx

		_104_19 := 2

		_108_4 := func() {
			cancel()
			close(blocker)
		}

		_112_4 := func() {
			<-blocker
		}

		_115_4 := func() {
			target[0] = src[0]
		}
		ctx := _103_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   102,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _104_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:108:4
		task8 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task8.emitter = cff.NopTaskEmitter()
		task8.fn = func(ctx context.Context) (err error) {
			taskEmitter := task8.emitter
			startTime := time.Now()
			defer func() {
				if task8.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task8.ran.Store(true)

			_108_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task8.fn,
		})
		tasks = append(tasks, task8)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:112:4
		task9 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task9.emitter = cff.NopTaskEmitter()
		task9.fn = func(ctx context.Context) (err error) {
			taskEmitter := task9.emitter
			startTime := time.Now()
			defer func() {
				if task9.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task9.ran.Store(true)

			_112_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task9.fn,
		})
		tasks = append(tasks, task9)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:115:4
		task10 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task10.emitter = cff.NopTaskEmitter()
		task10.fn = func(ctx context.Context) (err error) {
			taskEmitter := task10.emitter
			startTime := time.Now()
			defer func() {
				if task10.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task10.ran.Store(true)

			_115_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task10.fn,
		})
		tasks = append(tasks, task10)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:118*/
	}()
}

// TaskWithError runs a parallel cff.Task that errors.
func TaskWithError() error {
	return func() (err error) {

		_125_3 := context.Background()

		_126_19 := 2

		_128_4 := func() error {
			return errors.New("sad times")
		}
		ctx := _125_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   124,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _126_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:128:4
		task11 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task11.emitter = cff.NopTaskEmitter()
		task11.fn = func(ctx context.Context) (err error) {
			taskEmitter := task11.emitter
			startTime := time.Now()
			defer func() {
				if task11.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task11.ran.Store(true)

			err = _128_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task11.fn,
		})
		tasks = append(tasks, task11)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:131*/
	}()
}

// TaskWithPanic runs a parallel cff.task that panics.
func TaskWithPanic() error {
	return func() (err error) {

		_138_3 := context.Background()

		_139_19 := 2

		_141_4 := func() {
			panic("sad times")
		}
		ctx := _138_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   137,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _139_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:141:4
		task12 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task12.emitter = cff.NopTaskEmitter()
		task12.fn = func(ctx context.Context) (err error) {
			taskEmitter := task12.emitter
			startTime := time.Now()
			defer func() {
				if task12.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task12.ran.Store(true)

			_141_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task12.fn,
		})
		tasks = append(tasks, task12)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:144*/
	}()
}

// MultipleTask runs multiple cff.Task in parallel to populate the provided
// slice.
func MultipleTask(src, target []int) error {
	send := func(ctx context.Context) error {
		_, _ = ctx.Deadline()
		target[1] = src[1]
		return nil
	}
	return func() (err error) {

		_157_3 := context.Background()

		_158_19 := 2

		_160_4 := func() {
			target[0] = src[0]
		}

		_165_4 := send
		ctx := _157_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   156,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _158_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:160:4
		task13 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task13.emitter = cff.NopTaskEmitter()
		task13.fn = func(ctx context.Context) (err error) {
			taskEmitter := task13.emitter
			startTime := time.Now()
			defer func() {
				if task13.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task13.ran.Store(true)

			_160_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task13.fn,
		})
		tasks = append(tasks, task13)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:165:4
		task14 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task14.emitter = cff.NopTaskEmitter()
		task14.fn = func(ctx context.Context) (err error) {
			taskEmitter := task14.emitter
			startTime := time.Now()
			defer func() {
				if task14.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task14.ran.Store(true)

			err = _165_4(ctx)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task14.fn,
		})
		tasks = append(tasks, task14)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:166*/
	}()
}

// ContinueOnError runs multiple parallel tasks through errors in the
// scheduler. Despite tasks that error, the contents of the src slice should
// be copied into the target slice.
func ContinueOnError(src []int, target []int) error {
	blockerA := make(chan struct{})
	blockerB := make(chan struct{})
	return func() (err error) {

		_177_3 := context.Background()

		_178_19 := 2

		_179_23 := true

		_181_4 := func(_ context.Context) error {
			close(blockerA)
			return errors.New("sad times")
		}

		_185_4 := func() {

			<-blockerA
			target[0] = src[0]
		}

		_193_4 := func() {

			<-blockerA
			close(blockerB)
			panic("sadder times")
		}

		_201_4 := func() {

			<-blockerB
			target[1] = src[1]
		}
		ctx := _177_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   176,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _178_19, Emitter: schedEmitter,
				ContinueOnError: _179_23,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:181:4
		task15 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task15.emitter = cff.NopTaskEmitter()
		task15.fn = func(ctx context.Context) (err error) {
			taskEmitter := task15.emitter
			startTime := time.Now()
			defer func() {
				if task15.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task15.ran.Store(true)

			err = _181_4(ctx)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task15.fn,
		})
		tasks = append(tasks, task15)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:185:4
		task16 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task16.emitter = cff.NopTaskEmitter()
		task16.fn = func(ctx context.Context) (err error) {
			taskEmitter := task16.emitter
			startTime := time.Now()
			defer func() {
				if task16.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task16.ran.Store(true)

			_185_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task16.fn,
		})
		tasks = append(tasks, task16)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:193:4
		task17 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task17.emitter = cff.NopTaskEmitter()
		task17.fn = func(ctx context.Context) (err error) {
			taskEmitter := task17.emitter
			startTime := time.Now()
			defer func() {
				if task17.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task17.ran.Store(true)

			_193_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task17.fn,
		})
		tasks = append(tasks, task17)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:201:4
		task18 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task18.emitter = cff.NopTaskEmitter()
		task18.fn = func(ctx context.Context) (err error) {
			taskEmitter := task18.emitter
			startTime := time.Now()
			defer func() {
				if task18.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task18.ran.Store(true)

			_201_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task18.fn,
		})
		tasks = append(tasks, task18)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:207*/
	}()
}

// ContinueOnErrorBoolExpr parameterizes the ContinueOnError Option
// with a boolean expression. Despite tasks that error, the contents of
// the src slice should be copied into the target slice. The error returned
// by this function is checked to verify that it's underlying type is
// unchanged.
func ContinueOnErrorBoolExpr(src, target []int, fn func() bool) error {
	blockerA := make(chan struct{})
	blockerB := make(chan struct{})
	return func() (err error) {

		_220_3 := context.Background()

		_221_19 := 2

		_222_23 := fn()

		_224_4 := func(_ context.Context) error {
			close(blockerA)

			_, err := os.Open("non-existing")
			return err
		}

		_231_4 := func() {

			<-blockerA
			target[0] = src[0]
			close(blockerB)
		}

		_240_4 := func() {

			<-blockerB
			target[1] = src[1]
		}
		ctx := _220_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   219,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _221_19, Emitter: schedEmitter,
				ContinueOnError: _222_23,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:224:4
		task19 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task19.emitter = cff.NopTaskEmitter()
		task19.fn = func(ctx context.Context) (err error) {
			taskEmitter := task19.emitter
			startTime := time.Now()
			defer func() {
				if task19.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task19.ran.Store(true)

			err = _224_4(ctx)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task19.fn,
		})
		tasks = append(tasks, task19)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:231:4
		task20 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task20.emitter = cff.NopTaskEmitter()
		task20.fn = func(ctx context.Context) (err error) {
			taskEmitter := task20.emitter
			startTime := time.Now()
			defer func() {
				if task20.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task20.ran.Store(true)

			_231_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task20.fn,
		})
		tasks = append(tasks, task20)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:240:4
		task21 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task21.emitter = cff.NopTaskEmitter()
		task21.fn = func(ctx context.Context) (err error) {
			taskEmitter := task21.emitter
			startTime := time.Now()
			defer func() {
				if task21.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task21.ran.Store(true)

			_240_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task21.fn,
		})
		tasks = append(tasks, task21)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:247*/
	}()
}

// ContinueOnErrorCancelled runs a Parallel with a cancelled context.
func ContinueOnErrorCancelled(ctx context.Context, src []int, target []int) error {
	return func() (err error) {

		_254_3 := ctx

		_255_19 := 2

		_256_23 := true

		_258_4 := func() {
			target[0] = src[0]
		}
		ctx := _254_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   253,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _255_19, Emitter: schedEmitter,
				ContinueOnError: _256_23,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:258:4
		task22 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task22.emitter = cff.NopTaskEmitter()
		task22.fn = func(ctx context.Context) (err error) {
			taskEmitter := task22.emitter
			startTime := time.Now()
			defer func() {
				if task22.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task22.ran.Store(true)

			_258_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task22.fn,
		})
		tasks = append(tasks, task22)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:261*/
	}()
}

// ContinueOnErrorCancelledDuring runs a parallel with a context that is
// cancelled during scheduler operation.
func ContinueOnErrorCancelledDuring(ctx context.Context, cancel func(), src []int, target []int) error {
	blocker := make(chan struct{})
	return func() (err error) {

		_270_3 := ctx

		_271_19 := 2

		_272_23 := true

		_276_4 := func() {
			cancel()
			close(blocker)
		}

		_280_4 := func() {
			<-blocker
		}

		_285_4 := func() {
			target[0] = src[0]
		}
		ctx := _270_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   269,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _271_19, Emitter: schedEmitter,
				ContinueOnError: _272_23,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:276:4
		task23 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task23.emitter = cff.NopTaskEmitter()
		task23.fn = func(ctx context.Context) (err error) {
			taskEmitter := task23.emitter
			startTime := time.Now()
			defer func() {
				if task23.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task23.ran.Store(true)

			_276_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task23.fn,
		})
		tasks = append(tasks, task23)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:280:4
		task24 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task24.emitter = cff.NopTaskEmitter()
		task24.fn = func(ctx context.Context) (err error) {
			taskEmitter := task24.emitter
			startTime := time.Now()
			defer func() {
				if task24.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task24.ran.Store(true)

			_280_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task24.fn,
		})
		tasks = append(tasks, task24)

		// go.uber.org/cff/internal/tests/parallel/parallel.go:285:4
		task25 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task25.emitter = cff.NopTaskEmitter()
		task25.fn = func(ctx context.Context) (err error) {
			taskEmitter := task25.emitter
			startTime := time.Now()
			defer func() {
				if task25.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("panic: %v", recovered)
				}
			}()

			defer task25.ran.Store(true)

			_285_4()

			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task25.fn,
		})
		tasks = append(tasks, task25)

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:288*/
	}()
}

// SliceMultiple runs multiple cff.Slices to populate the provided slice.
// Slice will panic if len(target) < len(src).
func SliceMultiple(srcA, srcB, targetA, targetB []int) error {
	return func() (err error) {

		_296_3 := context.Background()

		_297_19 := 2

		_299_4 := func(idx int, val int) error {
			targetA[idx] = val
			return nil
		}

		_303_4 := srcA

		_306_4 := func(_ context.Context, idx int, val int) {
			targetB[idx] = val
		}

		_309_4 := srcB
		ctx := _296_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   295,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _297_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:298:3
		sliceTask26Slice := _303_4
		for idx, val := range sliceTask26Slice {
			idx := idx
			val := val
			sliceTask26 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask26.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _299_4(idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask26.fn,
			})
		}

		// go.uber.org/cff/internal/tests/parallel/parallel.go:305:3
		sliceTask27Slice := _309_4
		for idx, val := range sliceTask27Slice {
			idx := idx
			val := val
			sliceTask27 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask27.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				_306_4(ctx, idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask27.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:310*/
	}()
}

// SliceNoIndex runs multiple cff.Slices to populate the provided slice
func SliceNoIndex(srcA, srcB, targetA, targetB []int) error {
	return func() (err error) {

		_317_3 := context.Background()

		_318_19 := 2

		_320_4 := func(val int) error {
			targetA[val] = val
			return nil
		}

		_324_4 := srcA

		_327_4 := func(_ context.Context, val int) {
			targetB[val] = val
		}

		_330_4 := srcB
		ctx := _317_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   316,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _318_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:319:3
		sliceTask28Slice := _324_4
		for _, val := range sliceTask28Slice {

			val := val
			sliceTask28 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask28.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _320_4(val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask28.fn,
			})
		}

		// go.uber.org/cff/internal/tests/parallel/parallel.go:326:3
		sliceTask29Slice := _330_4
		for _, val := range sliceTask29Slice {

			val := val
			sliceTask29 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask29.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				_327_4(ctx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask29.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:331*/
	}()
}

type manyInts []int

// SliceWrapped runs cff.Slice with a wrapped type for a slice.
func SliceWrapped(src, target manyInts) error {
	return func() (err error) {

		_340_3 := context.Background()

		_341_19 := 2

		_343_4 := func(idx int, val int) error {
			target[idx] = val
			return nil
		}

		_347_4 := src
		ctx := _340_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   339,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _341_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:342:3
		sliceTask30Slice := _347_4
		for idx, val := range sliceTask30Slice {
			idx := idx
			val := val
			sliceTask30 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask30.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _343_4(idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask30.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:348*/
	}()
}

// AssignSliceItems runs cff.Slice in parallel to populate the provided slices.
// AssignSliceItems expects len(target) >= len(src).
func AssignSliceItems(src, target []string, keepgoing bool) error {
	return func() (err error) {

		_356_3 := context.Background()

		_357_19 := 2

		_358_23 := keepgoing

		_360_4 := func(idx int, val string) error {
			target[idx] = val
			switch val {
			case "error":
				return errors.New("sad times")
			case "panic":
				panic("sadder times")
			default:
				return nil
			}
		}

		_371_4 := src
		ctx := _356_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   355,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _357_19, Emitter: schedEmitter,
				ContinueOnError: _358_23,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:359:3
		sliceTask31Slice := _371_4
		for idx, val := range sliceTask31Slice {
			idx := idx
			val := val
			sliceTask31 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask31.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _360_4(idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask31.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:372*/
	}()
}

// SliceEnd runs cff.Slice in parallel and calls sliceEndFn after all items
// in the slice have finished.
func SliceEnd(src []int, sliceFn func(idx, val int) error, sliceEndFn func()) (err error) {
	err = func() (err error) {

		_380_3 := context.Background()

		_381_19 := 2

		_383_4 := sliceFn

		_384_4 := src

		_385_17 := sliceEndFn
		ctx := _380_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   379,
				Column: 8,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _381_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:382:3
		sliceTask32Slice := _384_4
		sliceTask32Jobs := make([]*cff.ScheduledJob, len(sliceTask32Slice))
		for idx, val := range sliceTask32Slice {
			idx := idx
			val := val
			sliceTask32 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask32.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _383_4(idx, val)
				return
			}
			sliceTask32Jobs[idx] = sched.Enqueue(ctx, cff.Job{
				Run: sliceTask32.fn,
			})
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: sliceTask32Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_385_17()
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:386*/
	}()
	return err
}

// SliceEndWithErr runs cff.Slice in parallel and calls an erroring sliceEndFn
// after all items in the slice have finished.
func SliceEndWithErr(src []int, sliceFn func(idx, val int) error, sliceEndFn func() error) (err error) {
	err = func() (err error) {

		_395_3 := context.Background()

		_396_19 := 2

		_398_4 := sliceFn

		_399_4 := src

		_400_17 := sliceEndFn
		ctx := _395_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   394,
				Column: 8,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _396_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:397:3
		sliceTask33Slice := _399_4
		sliceTask33Jobs := make([]*cff.ScheduledJob, len(sliceTask33Slice))
		for idx, val := range sliceTask33Slice {
			idx := idx
			val := val
			sliceTask33 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask33.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _398_4(idx, val)
				return
			}
			sliceTask33Jobs[idx] = sched.Enqueue(ctx, cff.Job{
				Run: sliceTask33.fn,
			})
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: sliceTask33Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _400_17()
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:401*/
	}()
	return err
}

// SliceEndWithCtx runs cff.Slice in parallel and calls sliceEndFn after all items
// in the slice have finished.
func SliceEndWithCtx(src []int, sliceFn func(idx, val int) error, sliceEndFn func(context.Context)) (err error) {
	err = func() (err error) {

		_410_3 := context.Background()

		_411_19 := 2

		_413_4 := sliceFn

		_414_4 := src

		_415_17 := sliceEndFn
		ctx := _410_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   409,
				Column: 8,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _411_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:412:3
		sliceTask34Slice := _414_4
		sliceTask34Jobs := make([]*cff.ScheduledJob, len(sliceTask34Slice))
		for idx, val := range sliceTask34Slice {
			idx := idx
			val := val
			sliceTask34 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask34.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _413_4(idx, val)
				return
			}
			sliceTask34Jobs[idx] = sched.Enqueue(ctx, cff.Job{
				Run: sliceTask34.fn,
			})
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: sliceTask34Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_415_17(ctx)
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:416*/
	}()
	return err
}

// SliceEndWithCtxAndErr runs cff.Slice in parallel and calls an erroring sliceEndFn
// after all items in the slice have finished.
func SliceEndWithCtxAndErr(src []int, sliceFn func(idx, val int) error, sliceEndFn func(context.Context) error) (err error) {
	err = func() (err error) {

		_425_3 := context.Background()

		_426_19 := 2

		_428_4 := sliceFn

		_429_4 := src

		_430_17 := sliceEndFn
		ctx := _425_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   424,
				Column: 8,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _426_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:427:3
		sliceTask35Slice := _429_4
		sliceTask35Jobs := make([]*cff.ScheduledJob, len(sliceTask35Slice))
		for idx, val := range sliceTask35Slice {
			idx := idx
			val := val
			sliceTask35 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask35.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _428_4(idx, val)
				return
			}
			sliceTask35Jobs[idx] = sched.Enqueue(ctx, cff.Job{
				Run: sliceTask35.fn,
			})
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: sliceTask35Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _430_17(ctx)
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:431*/
	}()
	return err
}

// AssignMapItems runs cff.Map in parallel to populate the provided slices.
func AssignMapItems(src map[string]int, keys []string, values []int, keepgoing bool) error {
	return func() (err error) {

		_439_3 := context.Background()

		_440_19 := 2

		_441_23 := keepgoing

		_443_4 := func(key string, val int) error {
			switch key {
			case "error":
				return errors.New("sad times")
			case "panic":
				panic("sadder times")
			default:
				keys[val] = key
				values[val] = val
				return nil
			}
		}

		_455_4 := src
		ctx := _439_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   438,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _440_19, Emitter: schedEmitter,
				ContinueOnError: _441_23,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/internal/tests/parallel/parallel.go:442:3
		for key, val := range _455_4 {
			key := key
			val := val
			mapTask36 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			mapTask36.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _443_4(key, val)
				return
			}

			sched.Enqueue(ctx, cff.Job{
				Run: mapTask36.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:456*/
	}()
}

// ForEachMapItem runs the given function on all items of the given map
// concurrently, and then the provided function when all items have been
// processed.
func ForEachMapItem[K comparable, V any](
	src map[K]V,
	fn func(K, V),
	after func(),
) error {
	return func() (err error) {

		_469_3 := context.Background()

		_470_19 := 2

		_471_11 := fn

		_471_15 := src

		_471_31 := after
		ctx := _469_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   468,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _470_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		mapTask37Jobs := make([]*cff.ScheduledJob, 0, len(_471_15))
		// go.uber.org/cff/internal/tests/parallel/parallel.go:471:3
		for key, val := range _471_15 {
			key := key
			val := val
			mapTask37 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			mapTask37.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_471_11(key, val)
				return
			}

			mapTask37Jobs = append(mapTask37Jobs, sched.Enqueue(ctx, cff.Job{
				Run: mapTask37.fn,
			}))
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: mapTask37Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					if recovered := recover(); recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_471_31()
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:471*/
	}()
}

// ForEachMapItemError is a variant of ForEachMapItem where all functions
// can return an error.
func ForEachMapItemError[K comparable, V any](
	src map[K]V,
	fn func(K, V) error,
	after func() error,
) error {
	return func() (err error) {

		_483_3 := context.Background()

		_484_19 := 2

		_485_11 := fn

		_485_15 := src

		_485_31 := after
		ctx := _483_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   482,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _484_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		mapTask38Jobs := make([]*cff.ScheduledJob, 0, len(_485_15))
		// go.uber.org/cff/internal/tests/parallel/parallel.go:485:3
		for key, val := range _485_15 {
			key := key
			val := val
			mapTask38 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			mapTask38.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _485_11(key, val)
				return
			}

			mapTask38Jobs = append(mapTask38Jobs, sched.Enqueue(ctx, cff.Job{
				Run: mapTask38.fn,
			}))
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: mapTask38Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					if recovered := recover(); recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _485_31()
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:485*/
	}()
}

// ForEachMapItemContext is a variant of ForEachMapItem where all functions
// accept a context.
func ForEachMapItemContext[K comparable, V any](
	ctx context.Context,
	src map[K]V,
	fn func(context.Context, K, V),
	after func(context.Context),
) error {
	return func() (err error) {

		_498_3 := ctx

		_499_19 := 2

		_500_11 := fn

		_500_15 := src

		_500_31 := after
		ctx := _498_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/internal/tests/parallel/parallel.go",
				Line:   497,
				Column: 9,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = cff.NopParallelEmitter()

			schedInfo = &cff.SchedulerInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}

			// possibly unused
			_ = parallelInfo
			_ = directiveInfo
		)

		startTime := time.Now()
		defer func() { parallelEmitter.ParallelDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _499_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		mapTask39Jobs := make([]*cff.ScheduledJob, 0, len(_500_15))
		// go.uber.org/cff/internal/tests/parallel/parallel.go:500:3
		for key, val := range _500_15 {
			key := key
			val := val
			mapTask39 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			mapTask39.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_500_11(ctx, key, val)
				return
			}

			mapTask39Jobs = append(mapTask39Jobs, sched.Enqueue(ctx, cff.Job{
				Run: mapTask39.fn,
			}))
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: mapTask39Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					if recovered := recover(); recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_500_31(ctx)
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line parallel.go:500*/
	}()
}
