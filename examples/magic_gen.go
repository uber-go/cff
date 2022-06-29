//line magic.go:1
package example

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"go.uber.org/cff"
	"github.com/uber-go/tally"
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
		/*line magic.go:33:18*/
		_33_18 := ctx
		/*line magic.go:34:14*/
		_34_14 := req
		/*line magic.go:35:15*/
		_35_15 := &res
		/*line magic.go:36:19*/
		_36_19 := 8
		/*line magic.go:37:19*/
		_37_19 := cff.TallyEmitter(h.scope)
		/*line magic.go:38:19*/
		_38_19 := cff.LogEmitter(h.logger)
		/*line magic.go:39:22*/
		_39_22 := "HandleFoo"
		/*line magic.go:42:4*/
		_42_4 := func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
			return &GetManagerRequest{
					LDAPGroup: req.LDAPGroup,
				}, &ListUsersRequest{
					LDAPGroup: req.LDAPGroup,
				}
		}
		/*line magic.go:50:4*/
		_50_4 := h.mgr.Get
		/*line magic.go:51:12*/
		_51_12 := h.ses.BatchSendEmail
		/*line magic.go:53:4*/
		_53_4 := func(responses []*SendEmailResponse) *Response {
			var r Response
			for _, res := range responses {
				r.MessageIDs = append(r.MessageIDs, res.MessageID)
			}
			return &r
		}
		/*line magic.go:62:4*/
		_62_4 := h.users.List
		/*line magic.go:63:18*/
		_63_18 := func(req *GetManagerRequest) bool {
			return req.LDAPGroup != "everyone"
		}
		/*line magic.go:66:21*/
		_66_21 := &ListUsersResponse{}
		/*line magic.go:67:19*/
		_67_19 := "FormSendEmailRequest"
		/*line magic.go:70:4*/
		_70_4 := func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest {
			var reqs []*SendEmailRequest
			for _, u := range users.Emails {
				reqs = append(reqs, &SendEmailRequest{Address: u})
			}
			return reqs
		}
		/*line magic.go:77:18*/
		_77_18 := func(req *GetManagerRequest) bool {
			return req.LDAPGroup != "everyone"
		}
		/*line magic.go:80:19*/
		_80_19 := "FormSendEmailRequest"

		/*line magic_gen.go:96*/
		ctx := _33_18
		var v1 *Request = _34_14
		emitter := cff.EmitterStack(_37_19, _38_19)

		var (
			flowInfo = &cff.FlowInfo{
				Name:   _39_22,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   33,
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
				Concurrency: _36_19, Emitter: schedEmitter,
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

		// go.uber.org/cff/examples/magic.go:42:4
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

			v2, v3 = _42_4(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cff.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/examples/magic.go:50:4
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

			v4, err = _50_4(v2)

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

		// go.uber.org/cff/examples/magic.go:63:4
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
			p0 = _63_18(v2)
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})

		// go.uber.org/cff/examples/magic.go:62:4
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
				Name:   _67_19,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   62,
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
					v5, err = _66_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task4.ran.Store(true)

			v5, err = _62_4(v3)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v5, err = _66_21, nil
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

		// go.uber.org/cff/examples/magic.go:77:4
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
			p1 = _77_18(v2)
			return nil
		}

		pred2.job = sched.Enqueue(ctx, cff.Job{
			Run: pred2.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})

		// go.uber.org/cff/examples/magic.go:70:4
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
				Name:   _80_19,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   70,
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

			v6 = _70_4(v4, v5)

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

		// go.uber.org/cff/examples/magic.go:51:12
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

			v7, err = _51_12(v6)

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

		// go.uber.org/cff/examples/magic.go:53:4
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

			v8 = _53_4(v7)

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

		*(_35_15) = v8 // *go.uber.org/cff/examples.Response

		flowEmitter.FlowSuccess(ctx)
		return nil /*line magic.go:81*/
	}()

	err = func() (err error) {
		/*line magic.go:85:3*/
		_85_3 := ctx
		/*line magic.go:86:19*/
		_86_19 := 2
		/*line magic.go:87:19*/
		_87_19 := cff.TallyEmitter(h.scope)
		/*line magic.go:88:19*/
		_88_19 := cff.LogEmitter(h.logger)
		/*line magic.go:89:26*/
		_89_26 := "SendParallel"
		/*line magic.go:90:23*/
		_90_23 := true
		/*line magic.go:92:4*/
		_92_4 := func(_ context.Context) error {
			return SendMessage()
		}
		/*line magic.go:95:4*/
		_95_4 := SendMessage
		/*line magic.go:98:4*/
		_98_4 := func() error {
			return SendMessage()
		}
		/*line magic.go:101:19*/
		_101_19 := "SendMsg"
		/*line magic.go:104:4*/
		_104_4 := func(ctx context.Context, idx int, s string) error {
			_ = fmt.Sprintf("%d and %q", idx, s)
			_, _ = ctx.Deadline()
			return nil
		}
		/*line magic.go:109:4*/
		_109_4 := []string{"message", "to", "send"}
		/*line magic.go:112:4*/
		_112_4 := func(ctx context.Context, idx int, s string) error {
			_ = fmt.Sprintf("%d and %q", idx, s)
			ctx.Deadline()
			return nil
		}
		/*line magic.go:117:4*/
		_117_4 := []string{"more", "messages", "sent"}
		/*line magic.go:118:17*/
		_118_17 := func(context.Context) error {
			return nil
		}
		/*line magic.go:123:4*/
		_123_4 := func(ctx context.Context, key string, value string) error {
			_ = fmt.Sprintf("%q : %q", key, value)
			_, _ = ctx.Deadline()
			return nil
		}
		/*line magic.go:128:4*/
		_128_4 := map[string]string{"key": "value"}
		/*line magic.go:131:4*/
		_131_4 := func(ctx context.Context, key string, value int) error {
			_ = fmt.Sprintf("%q: %v", key, value)
			return nil
		}
		/*line magic.go:135:4*/
		_135_4 := map[string]int{"a": 1, "b": 2, "c": 3}
		/*line magic.go:136:15*/
		_136_15 := func(context.Context) {
			_ = fmt.Sprint("}")
		}

		/*line magic_gen.go:601*/
		ctx := _85_3
		emitter := cff.EmitterStack(_87_19, _88_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _89_26,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   84,
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
				Concurrency: _86_19, Emitter: schedEmitter,
				ContinueOnError: _90_23,
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

		// go.uber.org/cff/examples/magic.go:92:4
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

			err = _92_4(ctx)

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

		// go.uber.org/cff/examples/magic.go:95:4
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

			err = _95_4()

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

		// go.uber.org/cff/examples/magic.go:98:4
		task8 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task8.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _101_19,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   98,
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

			err = _98_4()

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

		// go.uber.org/cff/examples/magic.go:103:3
		sliceTask9Slice := _109_4
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

				err = _104_4(ctx, idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask9.fn,
			})
		}

		// go.uber.org/cff/examples/magic.go:111:3
		sliceTask10Slice := _117_4
		sliceTask10Jobs := make([]*cff.ScheduledJob, len(sliceTask10Slice))
		for idx, val := range sliceTask10Slice {
			idx := idx
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

				err = _112_4(ctx, idx, val)
				return
			}
			sliceTask10Jobs[idx] = sched.Enqueue(ctx, cff.Job{
				Run: sliceTask10.fn,
			})
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: sliceTask10Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _118_17(ctx)
				return
			},
		})

		// go.uber.org/cff/examples/magic.go:122:3
		for key, val := range _128_4 {
			key := key
			val := val
			mapTask11 := new(struct {
				emitter cff.TaskEmitter
				fn      func(context.Context) error
				ran     cff.AtomicBool
			})
			mapTask11.fn = func(ctx context.Context) (err error) {
				defer func() {
					recovered := recover()
					if recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				err = _123_4(ctx, key, val)
				return
			}

			sched.Enqueue(ctx, cff.Job{
				Run: mapTask11.fn,
			})
		}

		mapTask12Jobs := make([]*cff.ScheduledJob, 0, len(_135_4))
		// go.uber.org/cff/examples/magic.go:130:3
		for key, val := range _135_4 {
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

				err = _131_4(ctx, key, val)
				return
			}

			mapTask12Jobs = append(mapTask12Jobs, sched.Enqueue(ctx, cff.Job{
				Run: mapTask12.fn,
			}))
		}

		sched.Enqueue(ctx, cff.Job{
			Dependencies: mapTask12Jobs,
			Run: func(ctx context.Context) (err error) {
				defer func() {
					if recovered := recover(); recovered != nil {
						err = fmt.Errorf("panic: %v", recovered)
					}
				}()

				_136_15(ctx)
				return
			},
		})

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line magic.go:139*/
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
