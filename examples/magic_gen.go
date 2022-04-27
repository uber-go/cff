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
		_33_18 := ctx
		_34_14 := req
		_35_15 := &res
		_36_19 := 8
		_37_19 := cff.TallyEmitter(h.scope)
		_38_19 := cff.LogEmitter(h.logger)
		_39_22 := "HandleFoo"
		_42_4 := func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
			return &GetManagerRequest{
					LDAPGroup: req.LDAPGroup,
				}, &ListUsersRequest{
					LDAPGroup: req.LDAPGroup,
				}
		}
		_50_4 := h.mgr.Get
		_51_12 := h.ses.BatchSendEmail
		_53_4 := func(responses []*SendEmailResponse) *Response {
			var r Response
			for _, res := range responses {
				r.MessageIDs = append(r.MessageIDs, res.MessageID)
			}
			return &r
		}
		_62_4 := h.users.List
		_63_21 := &ListUsersResponse{}
		_64_19 := "FormSendEmailRequest"
		_67_4 := func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest {
			var reqs []*SendEmailRequest
			for _, u := range users.Emails {
				reqs = append(reqs, &SendEmailRequest{Address: u})
			}
			return reqs
		}
		_75_18 := func(req *GetManagerRequest) bool {
			return req.LDAPGroup != "everyone"
		}
		_79_19 := "FormSendEmailRequest"
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
				Name:   _64_19,
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
				if recovered != nil {
					taskEmitter.TaskPanicRecovered(ctx, recovered)
					v5, err = _63_21, nil
				}
			}()

			defer task4.ran.Store(true)

			v5, err = _62_4(v3)

			if err != nil {
				taskEmitter.TaskErrorRecovered(ctx, err)
				v5, err = _63_21, nil
			} else {
				taskEmitter.TaskSuccess(ctx)
			}

			return
		}

		task4.job = sched.Enqueue(ctx, cff.Job{
			Run: task4.run,
			Dependencies: []*cff.ScheduledJob{
				task0.job,
			},
		})
		tasks = append(tasks, task4)

		// go.uber.org/cff/examples/magic.go:75:4
		var p0 bool
		pred1 := new(struct {
			ran cff.AtomicBool
			run func(context.Context) error
			job *cff.ScheduledJob
		})
		pred1.run = func(ctx context.Context) (err error) {
			p0 = _75_18(v2)
			return nil
		}

		pred1.job = sched.Enqueue(ctx, cff.Job{
			Run: pred1.run,
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
		task5.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _79_19,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   67,
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
				if recovered != nil {
					taskEmitter.TaskPanic(ctx, recovered)
					err = fmt.Errorf("task panic: %v", recovered)
				}
			}()

			if !p0 {
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
				pred1.job,
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
		return nil
	}()

	err = func() (err error) {
		_84_3 := ctx
		_85_19 := 2
		_86_19 := cff.TallyEmitter(h.scope)
		_87_19 := cff.LogEmitter(h.logger)
		_88_26 := "SendParallel"
		_89_23 := true
		_91_4 := func(_ context.Context) error {
			return SendMessage()
		}
		_94_4 := SendMessage
		_97_4 := func() error {
			return SendMessage()
		}
		_100_19 := "SendMsg"
		_103_4 := func(ctx context.Context, idx int, s string) error {
			_ = fmt.Sprintf("%d and %q", idx, s)
			_, _ = ctx.Deadline()
			return nil
		}
		_108_4 := []string{"message", "to", "send"}
		_111_4 := func(ctx context.Context, idx int, s string) error {
			_ = fmt.Sprintf("%d and %q", idx, s)
			ctx.Deadline()
			return nil
		}
		_116_4 := []string{"more", "messages", "sent"}
		_117_17 := func(context.Context) error {
			return nil
		}
		_122_4 := func(ctx context.Context, key string, value string) error {
			_ = fmt.Sprintf("%q : %q", key, value)
			_, _ = ctx.Deadline()
			return nil
		}
		_127_4 := map[string]string{"key": "value"}
		ctx := _84_3
		emitter := cff.EmitterStack(_86_19, _87_19)

		var (
			parallelInfo = &cff.ParallelInfo{
				Name:   _88_26,
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
				Concurrency: _85_19, Emitter: schedEmitter,
				ContinueOnError: _89_23,
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

		// go.uber.org/cff/examples/magic.go:91:4
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

			err = _91_4(ctx)

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

		// go.uber.org/cff/examples/magic.go:94:4
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

			err = _94_4()

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

		// go.uber.org/cff/examples/magic.go:97:4
		task8 := new(struct {
			emitter cff.TaskEmitter
			fn      func(context.Context) error
			ran     cff.AtomicBool
		})
		task8.emitter = emitter.TaskInit(
			&cff.TaskInfo{
				Name:   _100_19,
				File:   "go.uber.org/cff/examples/magic.go",
				Line:   97,
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

			err = _97_4()

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

		// go.uber.org/cff/examples/magic.go:102:3
		sliceTask9Slice := _108_4
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

				err = _103_4(ctx, idx, val)
				return
			}
			sched.Enqueue(ctx, cff.Job{
				Run: sliceTask9.fn,
			})
		}

		// go.uber.org/cff/examples/magic.go:110:3
		sliceTask10Slice := _116_4
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

				err = _111_4(ctx, idx, val)
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

				err = _117_17(ctx)
				return
			},
		})

		// go.uber.org/cff/examples/magic.go:121:3
		for key, val := range _127_4 {
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

				err = _122_4(ctx, key, val)
				return
			}

			sched.Enqueue(ctx, cff.Job{
				Run: mapTask11.fn,
			})
		}

		if err := sched.Wait(ctx); err != nil {
			parallelEmitter.ParallelError(ctx, err)
			return err
		}
		parallelEmitter.ParallelSuccess(ctx)
		return nil
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
