//line magic.go:1
//go:build !cff
// +build !cff

package example

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/uber-go/tally"
	"go.uber.org/cff"
	"go.uber.org/zap"
)

// Request TODO
type Request struct {
	LDAPGroup string
}

// Response TODO
type Response struct {
	MessageIDs []string
}

type fooHandler struct {
	mgr    *ManagerRepository
	users  *UserRepository
	ses    *SESClient
	scope  tally.Scope
	logger *zap.Logger
}

func (h *fooHandler) HandleFoo(ctx context.Context, req *Request) (*Response, error) {
	var res *Response
	err := func() (err error) {
		/*line magic.go:36:18*/
		_36_18 := ctx
		/*line magic.go:37:14*/
		_37_14 := req
		/*line magic.go:38:15*/
		_38_15 := &res
		/*line magic.go:39:19*/
		_39_19 := 8
		/*line magic.go:40:19*/
		_40_19 := cff.TallyEmitter(h.scope)
		/*line magic.go:41:19*/
		_41_19 := cff.LogEmitter(h.logger)
		/*line magic.go:42:22*/
		_42_22 := "HandleFoo"
		/*line magic.go:45:4*/
		_45_4 := func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
			return &GetManagerRequest{
					LDAPGroup: req.LDAPGroup,
				}, &ListUsersRequest{
					LDAPGroup: req.LDAPGroup,
				}
		}
		/*line magic.go:53:4*/
		_53_4 := h.mgr.Get
		/*line magic.go:54:12*/
		_54_12 := h.ses.BatchSendEmail
		/*line magic.go:56:4*/
		_56_4 := func(responses []*SendEmailResponse) *Response {
			var r Response
			for _, res := range responses {
				r.MessageIDs = append(r.MessageIDs, res.MessageID)
			}
			return &r
		}
		/*line magic.go:65:4*/
		_65_4 := h.users.List
		/*line magic.go:66:18*/
		_66_18 := func(req *GetManagerRequest) bool {
			return req.LDAPGroup != "everyone"
		}
		/*line magic.go:69:21*/
		_69_21 := &ListUsersResponse{}
		/*line magic.go:70:19*/
		_70_19 := "FormSendEmailRequest"
		/*line magic.go:73:4*/
		_73_4 := func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest {
			var reqs []*SendEmailRequest
			for _, u := range users.Emails {
				reqs = append(reqs, &SendEmailRequest{Address: u})
			}
			return reqs
		}
		/*line magic.go:80:18*/
		_80_18 := func(req *GetManagerRequest) bool {
			return req.LDAPGroup != "everyone"
		}
		/*line magic.go:83:19*/
		_83_19 := "FormSendEmailRequest"

		/*line magic_gen.go:99*/
		ctx := _36_18
		var v1 *Request = _37_14
		emitter := cff.EmitterStack(_40_19, _41_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _42_22,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   36,
				Column: 9,
			}
			flowEmitter = emitter.FlowInit(flowInfo)

			schedInfo = &cff.SchedulerInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			}

			// possibly unused
			_ = flowInfo
		)

		startTime := time.Now()
		defer func() { flowEmitter.FlowDone(ctx, time.Since(startTime)) }()

		schedEmitter := emitter.SchedulerInit(schedInfo)

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Concurrency: _39_19, Emitter: schedEmitter,
			},
		)

		var tasks []*struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		}
		defer func() {
			for _, t := range tasks {
				if !t.ran.Load() {
					t.emitter.TaskSkipped(ctx, err)
				}
			}
		}()

		// go.uber.org/cff/examples/magic.go:45:4
		var (
			v2 *GetManagerRequest
			v3 *ListUsersRequest
		)
		task0 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task0.emitter = cff.NopTaskEmitter()
		task0.run = func(ctx context.Context) (err error) {
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task0.ran.Store(true)

			v2, v3 = _45_4(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cff.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/examples/magic.go:53:4
		var (
			v4 *GetManagerResponse
		)
		task1 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task1.emitter = cff.NopTaskEmitter()
		task1.run = func(ctx context.Context) (err error) {
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task1.ran.Store(true)

			v4, err = _53_4(v2)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task1.job = sched.Enqueue(ctx, cff.Job{
			Run: task1.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})
		tasks = append(tasks, task1)

		// go.uber.org/cff/examples/magic.go:66:4
		var p0 bool
		var p0PanicRecover interface{}
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
				}
			}()
			p0 = _66_18(v2)
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})

		// go.uber.org/cff/examples/magic.go:65:4
		var (
			v5 *ListUsersResponse
		)
		task4 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task4.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _70_19,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   65,
				Column: 4,
			},
			&cff.DirectiveInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			},
		)
		task4.run = func(ctx context.Context) (err error) {
			taskEmitter := task4.emitter
			startTime := time.Now()
			defer func() {
				if task4.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil && p0PanicRecover != nil {
					recovered = p0PanicRecover
				}
				if recovered != nil {
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v5, err = _69_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task4.ran.Store(true)

			v5, err = _65_4(v3)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v5, err = _69_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task4.job = sched.Enqueue(ctx, cff.Job{
			Run: task4.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
				pred1.job,
			},
		})
		tasks = append(tasks, task4)

		// go.uber.org/cff/examples/magic.go:80:4
		var p1 bool
		var p1PanicRecover interface{}
		pred2 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred2.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p1PanicRecover = recovered
				}
			}()
			p1 = _80_18(v2)
			return nil
		}

		pred2.job = sched.Enqueue(ctx, cff.Job{
			Run: pred2.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})

		// go.uber.org/cff/examples/magic.go:73:4
		var (
			v6 []*SendEmailRequest
		)
		task5 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task5.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _83_19,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   73,
				Column: 4,
			},
			&cff.DirectiveInfo{
				Name:      flowInfo.Name,
				Directive: cff.FlowDirective,
				File:      flowInfo.File,
				Line:      flowInfo.Line,
				Column:    flowInfo.Column,
			},
		)
		task5.run = func(ctx context.Context) (err error) {
			taskEmitter := task5.emitter
			startTime := time.Now()
			defer func() {
				if task5.ran.Load() {
					taskEmitter.TaskDone(ctx, time.Since(startTime))
				}
			}()

			defer func() {
				recovered := recover()
				if recovered == nil && p1PanicRecover != nil {
					recovered = p1PanicRecover
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !p1 {
				return nil
			}

			defer task5.ran.Store(true)

			v6 = _73_4(v4, v5)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task5.job = sched.Enqueue(ctx, cff.Job{
			Run: task5.run,
			Dependencies: []*cff.ScheduledJob{
				task1.job,
				task4.job,
				pred2.job,
			},
		})
		tasks = append(tasks, task5)

		// go.uber.org/cff/examples/magic.go:54:12
		var (
			v7 []*SendEmailResponse
		)
		task2 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task2.emitter = cff.NopTaskEmitter()
		task2.run = func(ctx context.Context) (err error) {
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task2.ran.Store(true)

			v7, err = _54_12(v6)

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return err
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task2.job = sched.Enqueue(ctx, cff.Job{
			Run: task2.run,
			Dependencies: []*cff.ScheduledJob{
				task5.job,
			},
		})
		tasks = append(tasks, task2)

		// go.uber.org/cff/examples/magic.go:56:4
		var (
			v8 *Response
		)
		task3 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task3.emitter = cff.NopTaskEmitter()
		task3.run = func(ctx context.Context) (err error) {
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
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			defer task3.ran.Store(true)

			v8 = _56_4(v7)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task3.job = sched.Enqueue(ctx, cff.Job{
			Run: task3.run,
			Dependencies: []*cff.ScheduledJob{
				task2.job,
			},
		})
		tasks = append(tasks, task3)

		if err := sched.Wait(ctx); err != nil {
			flowEmitter.FlowError(ctx, err)
			return err
		}

		*(_38_15) = v8 // *go.uber.org/cff/examples.Response

		flowEmitter.FlowSuccess(ctx)
		return nil /*line magic.go:84*/
	}()

	err = func() (err error) {
		/*line magic.go:88:3*/
		_88_3 := ctx
		/*line magic.go:89:19*/
		_89_19 := 2
		/*line magic.go:90:19*/
		_90_19 := cff.TallyEmitter(h.scope)
		/*line magic.go:91:19*/
		_91_19 := cff.LogEmitter(h.logger)
		/*line magic.go:92:26*/
		_92_26 := "SendParallel"
		/*line magic.go:93:23*/
		_93_23 := true
		/*line magic.go:95:4*/
		_95_4 := func(_ context.Context) error {
			return SendMessage()
		}
		/*line magic.go:98:4*/
		_98_4 := SendMessage
		/*line magic.go:101:4*/
		_101_4 := func() error {
			return SendMessage()
		}
		/*line magic.go:104:19*/
		_104_19 := "SendMsg"
		/*line magic.go:107:4*/
		_107_4 := func(ctx context.Context, idx int, s string) error {
			_ = fmt.Sprintf("%d and %q", idx, s)
			_, _ = ctx.Deadline()
			return nil
		}
		/*line magic.go:112:4*/
		_112_4 := []string{"message", "to", "send"}
		/*line magic.go:115:4*/
		_115_4 := func(ctx context.Context, s string) error {
			_ = fmt.Sprintf("%q", s)
			_, _ = ctx.Deadline()
			return nil
		}
		/*line magic.go:120:4*/
		_120_4 := []string{"message", "to", "send"}
		/*line magic.go:123:4*/
		_123_4 := func(ctx context.Context, idx int, s string) error {
			_ = fmt.Sprintf("%d and %q", idx, s)
			ctx.Deadline()
			return nil
		}
		/*line magic.go:128:4*/
		_128_4 := []string{"more", "messages", "sent"}
		/*line magic.go:129:17*/
		_129_17 := func(context.Context) error {
			return nil
		}
		/*line magic.go:134:4*/
		_134_4 := func(ctx context.Context, key string, value string) error {
			_ = fmt.Sprintf("%q : %q", key, value)
			_, _ = ctx.Deadline()
			return nil
		}
		/*line magic.go:139:4*/
		_139_4 := map[string]string{"key": "value"}
		/*line magic.go:142:4*/
		_142_4 := func(ctx context.Context, key string, value int) error {
			_ = fmt.Sprintf("%q: %v", key, value)
			return nil
		}
		/*line magic.go:146:4*/
		_146_4 := map[string]int{"a": 1, "b": 2, "c": 3}
		/*line magic.go:147:15*/
		_147_15 := func(context.Context) {
			_ = fmt.Sprint("}")
		}

		/*line magic_gen.go:612*/
		ctx := _88_3
		emitter := cff.EmitterStack(_90_19, _91_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _92_26,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   87,
				Column: 8,
			}
			directiveInfo = &cff.DirectiveInfo{
				Name:      parallelInfo.Name,
				Directive: cff.ParallelDirective,
				File:      parallelInfo.File,
				Line:      parallelInfo.Line,
				Column:    parallelInfo.Column,
			}
			parallelEmitter = emitter.ParallelInit(parallelInfo)

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

		sched := cff.BeginFlow(
			cff.SchedulerParams{
				Concurrency: _89_19, Emitter: schedEmitter,
				ContinueOnError: _93_23,
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

		// go.uber.org/cff/examples/magic.go:95:4
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

			err = _95_4(ctx)

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

		// go.uber.org/cff/examples/magic.go:98:4
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

			err = _98_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task7.fn,
		})
		tasks = append(tasks, task7)

		// go.uber.org/cff/examples/magic.go:101:4
		task8 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task8.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _104_19,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   101,
				Column: 4,
			},
			directiveInfo,
		)
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

			err = _101_4()

			if err != nil {
				taskEmitter.TaskError(ctx, err)
				return
			}
			taskEmitter.TaskSuccess(ctx)
			return
		}

		sched.Enqueue(ctx, cff.Job{
			Run: task8.fn,
		})
		tasks = append(tasks, task8)

		// go.uber.org/cff/examples/magic.go:106:3
		sliceTask9Slice := _112_4
		for idx, val := range sliceTask9Slice {
			idx := idx
			val := val
			sliceTask9 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask9.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _107_4(ctx, idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask9.fn,
			})
		}

		// go.uber.org/cff/examples/magic.go:114:3
		sliceTask10Slice := _120_4
		for _, val := range sliceTask10Slice {

			val := val
			sliceTask10 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask10.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _115_4(ctx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask10.fn,
			})
		}

		// go.uber.org/cff/examples/magic.go:122:3
		sliceTask11Slice := _128_4
		sliceTask11Jobs := make([]*cff.ScheduledJob, len(sliceTask11Slice))
		for idx, val := range sliceTask11Slice {
			idx := idx
			val := val
			sliceTask11 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			sliceTask11.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()
				err = _123_4(ctx, idx, val)
				return
			}
			sliceTask11Jobs[idx] = sched.Enqueue(ctx, cff.Job{
				Run: sliceTask11.fn,
			})
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: sliceTask11Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _129_17(ctx)
				return
			},
		})

		// go.uber.org/cff/examples/magic.go:133:3
		for key, val := range _139_4 {
			key := key
			val := val
			mapTask12 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			mapTask12.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _134_4(ctx, key, val)
				return
			}

			sched.Enqueue(ctx, cff.Job{
				Run: mapTask12.fn,
			})
		}

		mapTask13Jobs := make([]*cff.ScheduledJob, 0, len(_146_4))
		// go.uber.org/cff/examples/magic.go:141:3
		for key, val := range _146_4 {
			key := key
			val := val
			mapTask13 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			mapTask13.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _142_4(ctx, key, val)
				return
			}

			mapTask13Jobs = append(mapTask13Jobs, sched.Enqueue(ctx, cff.Job{
				Run: mapTask13.fn,
			}))
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: mapTask13Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					if recovered := recover(); recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_147_15(ctx)
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line magic.go:150*/
	}()
	return res, err
}

// ManagerRepository TODO
type ManagerRepository struct{}

// GetManagerRequest TODO
type GetManagerRequest struct {
	LDAPGroup string
}

// GetManagerResponse TODO
type GetManagerResponse struct {
	Email string
}

// Get TODO
func (*ManagerRepository) Get(req *GetManagerRequest) (*GetManagerResponse, error) {
	return &GetManagerResponse{Email: "boss@example.com"}, nil
}

// UserRepository TODO
type UserRepository struct{}

// ListUsersRequest TODO
type ListUsersRequest struct {
	LDAPGroup string
}

// ListUsersResponse TODO
type ListUsersResponse struct {
	Emails []string
}

// List TODO
func (*UserRepository) List(req *ListUsersRequest) (*ListUsersResponse, error) {
	return &ListUsersResponse{
		Emails: []string{"a@example.com", "b@example.com"},
	}, nil
}

// SESClient TODO
type SESClient struct{}

// SendEmailRequest TODO
type SendEmailRequest struct {
	Address string
}

// SendEmailResponse TODO
type SendEmailResponse struct {
	MessageID string
}

// BatchSendEmail TODO
func (*SESClient) BatchSendEmail(req []*SendEmailRequest) ([]*SendEmailResponse, error) {
	res := make([]*SendEmailResponse, len(req))
	for i := range req {
		res[i] = &SendEmailResponse{MessageID: strconv.Itoa(i)}
	}
	return res, nil
}

// SendMessage returns nil error.
func SendMessage() error {
	return nil
}
