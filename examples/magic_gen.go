//line magic.go:1
//go:build !cff
// +build !cff

package example

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"go.uber.org/cff"
)

// Request TODO
type Request struct {
	LDAPGroup string
}

// Response TODO
type Response struct {
	MessageIDs []string
}

// FooHandler does things.
type FooHandler struct {
	mgr   *ManagerRepository
	users *UserRepository
	ses   *SESClient
}

// HandleFoo handles foo requests.
func (h *FooHandler) HandleFoo(ctx context.Context, req *Request) (*Response, error) {
	var res *Response
	err := func() (err error) {
		/*line magic.go:34:18*/
		_34_18 := ctx
		/*line magic.go:35:14*/
		_35_14 := req
		/*line magic.go:36:15*/
		_36_15 := &res
		/*line magic.go:37:19*/
		_37_19 := 8
		/*line magic.go:40:4*/
		_40_4 := func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
			return &GetManagerRequest{
					LDAPGroup: req.LDAPGroup,
				}, &ListUsersRequest{
					LDAPGroup: req.LDAPGroup,
				}
		}
		/*line magic.go:48:4*/
		_48_4 := h.mgr.Get
		/*line magic.go:49:12*/
		_49_12 := h.ses.BatchSendEmail
		/*line magic.go:51:4*/
		_51_4 := func(responses []*SendEmailResponse) *Response {
			var r Response
			for _, res := range responses {
				r.MessageIDs = append(r.MessageIDs, res.MessageID)
			}
			return &r
		}
		/*line magic.go:60:4*/
		_60_4 := h.users.List
		/*line magic.go:61:18*/
		_61_18 := func(req *GetManagerRequest) bool {
			return req.LDAPGroup != "everyone"
		}
		/*line magic.go:64:21*/
		_64_21 := &ListUsersResponse{}
		/*line magic.go:67:4*/
		_67_4 := func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest {
			var reqs []*SendEmailRequest
			for _, u := range users.Emails {
				reqs = append(reqs, &SendEmailRequest{Address: u})
			}
			return reqs
		}
		/*line magic.go:74:18*/
		_74_18 := func(req *GetManagerRequest) bool {
			return req.LDAPGroup != "everyone"
		}

		/*line magic_gen.go:88*/
		ctx := _34_18
		var v1 *Request = _35_14
		emitter := cff.NopEmitter()

		var (
			flowInfo = &cff.FlowInfo{
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   34,
				Column: 9,
			}
			flowEmitter = cff.NopFlowEmitter()

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

		sched := cff.NewScheduler(
			cff.SchedulerParams{
				Concurrency: _37_19, Emitter: schedEmitter,
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

		// go.uber.org/cff/examples/magic.go:40:4
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
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task0.ran.Store(true)

			v2, v3 = _40_4(v1)

			taskEmitter.TaskSuccess(ctx)

			return
		}

		task0.job = sched.Enqueue(ctx, cff.Job{
			Run: task0.run,
		})
		tasks = append(tasks, task0)

		// go.uber.org/cff/examples/magic.go:48:4
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
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task1.ran.Store(true)

			v4, err = _48_4(v2)

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

		// go.uber.org/cff/examples/magic.go:61:4
		var p0 bool
		var p0PanicRecover interface{}
		var p0PanicStacktrace string
		_ = p0PanicStacktrace // possibly unused.
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p0PanicRecover = recovered
					p0PanicStacktrace = string(debug.Stack())
				}
			}()
			p0 = _61_18(v2)
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})

		// go.uber.org/cff/examples/magic.go:60:4
		var (
			v5 *ListUsersResponse
		)
		task4 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task4.emitter = cff.NopTaskEmitter()
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
					v5, err = _64_21, nil
				}
			}()

			if !p0 {
				return nil
			}

			defer task4.ran.Store(true)

			v5, err = _60_4(v3)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v5, err = _64_21, nil
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

		// go.uber.org/cff/examples/magic.go:74:4
		var p1 bool
		var p1PanicRecover interface{}
		var p1PanicStacktrace string
		_ = p1PanicStacktrace // possibly unused.
		pred2 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred2.run = func(ctx context.Context) (err error) {
			defer func() {
				if recovered := recover(); recovered != nil {
					p1PanicRecover = recovered
					p1PanicStacktrace = string(debug.Stack())
				}
			}()
			p1 = _74_18(v2)
			return nil
		}

		pred2.job = sched.Enqueue(ctx, cff.Job{
			Run: pred2.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})

		// go.uber.org/cff/examples/magic.go:67:4
		var (
			v6 []*SendEmailRequest
		)
		task5 := new(struct {
			emitter cff.TaskEmitter
			ran     cff.AtomicBool
			run     func(context.Context) error
			job     *cff.ScheduledJob
		})
		task5.emitter = cff.NopTaskEmitter()
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
				var stacktrace string
				if recovered != nil {
					stacktrace = string(debug.Stack())
				}
				if recovered == nil && p1PanicRecover != nil {
					recovered = p1PanicRecover
					stacktrace = p1PanicStacktrace
				}
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: stacktrace,
					}
				}
			}()

			if !p1 {
				return nil
			}

			defer task5.ran.Store(true)

			v6 = _67_4(v4, v5)

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

		// go.uber.org/cff/examples/magic.go:49:12
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
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task2.ran.Store(true)

			v7, err = _49_12(v6)

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

		// go.uber.org/cff/examples/magic.go:51:4
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
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}
			}()

			defer task3.ran.Store(true)

			v8 = _51_4(v7)

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

		*(_36_15) = v8 // *go.uber.org/cff/examples.Response

		flowEmitter.FlowSuccess(ctx)
		return nil /*line magic.go:77*/
	}()
	if err != nil {
		return nil, err
	}

	err = func() (err error) {
		/*line magic.go:84:3*/
		_84_3 := ctx
		/*line magic.go:85:19*/
		_85_19 := 2
		/*line magic.go:86:23*/
		_86_23 := true
		/*line magic.go:88:4*/
		_88_4 := func(_ context.Context) error {
			return SendMessage()
		}
		/*line magic.go:91:4*/
		_91_4 := SendMessage
		/*line magic.go:94:4*/
		_94_4 := func() error {
			return SendMessage()
		}
		/*line magic.go:99:4*/
		_99_4 := func(ctx context.Context, idx int, s string) error {
			_ = fmt.Sprintf("%d and %q", idx, s)
			_, _ = ctx.Deadline()
			return nil
		}
		/*line magic.go:104:4*/
		_104_4 := []string{"message", "to", "send"}
		/*line magic.go:107:4*/
		_107_4 := func(ctx context.Context, s string) error {
			_ = fmt.Sprintf("%q", s)
			_, _ = ctx.Deadline()
			return nil
		}
		/*line magic.go:112:4*/
		_112_4 := []string{"message", "to", "send"}
		/*line magic.go:115:4*/
		_115_4 := func(ctx context.Context, idx int, s string) error {
			_ = fmt.Sprintf("%d and %q", idx, s)
			ctx.Deadline()
			return nil
		}
		/*line magic.go:120:4*/
		_120_4 := []string{"more", "messages", "sent"}
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

		/*line magic_gen.go:587*/
		ctx := _84_3
		emitter := cff.NopEmitter()

		var (
			parallelInfo = &cff.ParallelInfo{
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   83,
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
				Concurrency: _85_19, Emitter: schedEmitter,
				ContinueOnError: _86_23,
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

		// go.uber.org/cff/examples/magic.go:88:4
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
				if recovered == nil {
					return
				}
				taskEmitter.TaskPanic(ctx, recovered)
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: string(debug.Stack()),
				}
			}()

			defer task6.ran.Store(true)

			err = _88_4(ctx)

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

		// go.uber.org/cff/examples/magic.go:91:4
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
				if recovered == nil {
					return
				}
				taskEmitter.TaskPanic(ctx, recovered)
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: string(debug.Stack()),
				}
			}()

			defer task7.ran.Store(true)

			err = _91_4()

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

		// go.uber.org/cff/examples/magic.go:94:4
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
				if recovered == nil {
					return
				}
				taskEmitter.TaskPanic(ctx, recovered)
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: string(debug.Stack()),
				}
			}()

			defer task8.ran.Store(true)

			err = _94_4()

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

		// go.uber.org/cff/examples/magic.go:98:3
		sliceTask9Slice := _104_4
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
					if recovered == nil {
						return
					}
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}()
				err = _99_4(ctx, idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask9.fn,
			})
		}

		// go.uber.org/cff/examples/magic.go:106:3
		sliceTask10Slice := _112_4
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
					if recovered == nil {
						return
					}
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}()
				err = _107_4(ctx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask10.fn,
			})
		}

		// go.uber.org/cff/examples/magic.go:114:3
		sliceTask11Slice := _120_4
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
					if recovered == nil {
						return
					}
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}()
				err = _115_4(ctx, idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask11.fn,
			})
		}

		// go.uber.org/cff/examples/magic.go:122:3
		for key, val := range _128_4 {
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
					if recovered == nil {
						return
					}
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}()

				err = _123_4(ctx, key, val)
				return
			}

			sched.Enqueue(ctx, cff.Job{
				Run: mapTask12.fn,
			})
		}

		// go.uber.org/cff/examples/magic.go:130:3
		for key, val := range _135_4 {
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
					if recovered == nil {
						return
					}
					err = &cff.PanicError{
						Value:      recovered,
						Stacktrace: string(debug.Stack()),
					}
				}()

				err = _131_4(ctx, key, val)
				return
			}

			sched.Enqueue(ctx, cff.Job{
				Run: mapTask13.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil /*line magic.go:136*/
	}()
	if err != nil {
		return nil, err
	}

	return res, nil
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
