//go:build !cff && v2
// +build !cff,v2

package example

import (
	"context"
	"fmt"
	"runtime/debug"
	"strconv"
	"time"

	"go.uber.org/cff"
)

// RequestV2 TODO
type RequestV2 struct {
	LDAPGroup string
}

// ResponseV2 TODO
type ResponseV2 struct {
	MessageIDs []string
}

type fooHandlerV2 struct {
	mgr   *ManagerRepositoryV2
	users *UserRepositoryV2
	ses   *SESClientV2
}

func (h *fooHandlerV2) HandleFoo(ctx context.Context, req *RequestV2) (*ResponseV2, error) {
	var res *ResponseV2
	err := _cffFlowmagicv2_32_9(ctx,
		_cffParamsmagicv2_33_3(req),
		_cffResultsmagicv2_34_3(&res),
		_cffConcurrencymagicv2_35_3(8),

		_cffTaskmagicv2_37_3(
			func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2) {
				return &GetManagerRequestV2{
						LDAPGroup: req.LDAPGroup,
					}, &ListUsersRequestV2{
						LDAPGroup: req.LDAPGroup,
					}
			}),
		_cffTaskmagicv2_45_3(
			h.mgr.Get),
		_cffTaskmagicv2_47_3(h.ses.BatchSendEmail),
		_cffTaskmagicv2_48_3(
			func(responses []*SendEmailResponseV2) *ResponseV2 {
				var r ResponseV2
				for _, res := range responses {
					r.MessageIDs = append(r.MessageIDs, res.MessageID)
				}
				return &r
			},
		),
		_cffTaskmagicv2_57_3(
			h.users.List,
			cff.Predicate(func(req *GetManagerRequestV2) bool {
				return req.LDAPGroup != "everyone"
			}),
			cff.FallbackWith(&ListUsersResponseV2{}),
		),
		_cffTaskmagicv2_64_3(
			func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2 {
				var reqs []*SendEmailRequestV2
				for _, u := range users.Emails {
					reqs = append(reqs, &SendEmailRequestV2{Address: u})
				}
				return reqs
			},
			cff.Predicate(func(req *GetManagerRequestV2) bool {
				return req.LDAPGroup != "everyone"
			}),
		),
	)

	err = cff.Parallel(
		ctx,
		cff.Concurrency(2),
		cff.ContinueOnError(true),
		cff.Tasks(
			func(_ context.Context) error {
				return SendMessageV2()
			},
			SendMessageV2,
		),
		cff.Task(
			func() error {
				return SendMessageV2()
			},
		),
		cff.Slice(
			func(ctx context.Context, idx int, s string) error {
				_ = fmt.Sprintf("%d and %q", idx, s)
				_, _ = ctx.Deadline()
				return nil
			},
			[]string{"message", "to", "send"},
		),
		cff.Slice(
			func(ctx context.Context, s string) error {
				_ = fmt.Sprintf("%q", s)
				_, _ = ctx.Deadline()
				return nil
			},
			[]string{"message", "to", "send"},
		),
		cff.Slice(
			func(ctx context.Context, idx int, s string) error {
				_ = fmt.Sprintf("%d and %q", idx, s)
				ctx.Deadline()
				return nil
			},
			[]string{"more", "messages", "sent"},
		),
		cff.Map(
			func(ctx context.Context, key string, value string) error {
				_ = fmt.Sprintf("%q : %q", key, value)
				_, _ = ctx.Deadline()
				return nil
			},
			map[string]string{"key": "value"},
		),
		cff.Map(
			func(ctx context.Context, key string, value int) error {
				_ = fmt.Sprintf("%q: %v", key, value)
				return nil
			},
			map[string]int{"a": 1, "b": 2, "c": 3},
		),
	)
	return res, err
}

// ManagerRepositoryV2 TODO
type ManagerRepositoryV2 struct{}

// GetManagerRequestV2 TODO
type GetManagerRequestV2 struct {
	LDAPGroup string
}

// GetManagerResponseV2 TODO
type GetManagerResponseV2 struct {
	Email string
}

// Get TODO
func (*ManagerRepositoryV2) Get(req *GetManagerRequestV2) (*GetManagerResponseV2, error) {
	return &GetManagerResponseV2{Email: "boss@example.com"}, nil
}

// UserRepositoryV2 TODO
type UserRepositoryV2 struct{}

// ListUsersRequestV2 TODO
type ListUsersRequestV2 struct {
	LDAPGroup string
}

// ListUsersResponseV2 TODO
type ListUsersResponseV2 struct {
	Emails []string
}

// List TODO
func (*UserRepositoryV2) List(req *ListUsersRequestV2) (*ListUsersResponseV2, error) {
	return &ListUsersResponseV2{
		Emails: []string{"a@example.com", "b@example.com"},
	}, nil
}

// SESClientV2 TODO
type SESClientV2 struct{}

// SendEmailRequestV2 TODO
type SendEmailRequestV2 struct {
	Address string
}

// SendEmailResponseV2 TODO
type SendEmailResponseV2 struct {
	MessageID string
}

// BatchSendEmail TODO
func (*SESClientV2) BatchSendEmail(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error) {
	res := make([]*SendEmailResponseV2, len(req))
	for i := range req {
		res[i] = &SendEmailResponseV2{MessageID: strconv.Itoa(i)}
	}
	return res, nil
}

// SendMessageV2 returns nil error.
func SendMessageV2() error {
	return nil
}
func _cffFlowmagicv2_32_9(
	ctx context.Context,
	mmagicv233_3 func() *RequestV2,
	mmagicv234_3 func() **ResponseV2,
	mmagicv235_3 func() int,
	mmagicv237_3 func() func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2),
	mmagicv245_3 func() func(req *GetManagerRequestV2) (*GetManagerResponseV2, error),
	mmagicv247_3 func() func(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error),
	mmagicv248_3 func() func(responses []*SendEmailResponseV2) *ResponseV2,
	mmagicv257_3 func() (func(req *ListUsersRequestV2) (*ListUsersResponseV2, error), cff.TaskOption, cff.TaskOption),
	_ cff.Option,
	mmagicv264_3 func() (func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2, cff.TaskOption),
	_ cff.Option,
) error {
	_33_14 := mmagicv233_3()
	_ = _33_14 // possibly unused.
	_34_15 := mmagicv234_3()
	_ = _34_15 // possibly unused.
	_35_19 := mmagicv235_3()
	_ = _35_19 // possibly unused.
	_38_4 := mmagicv237_3()
	_ = _38_4 // possibly unused.
	_46_4 := mmagicv245_3()
	_ = _46_4 // possibly unused.
	_47_12 := mmagicv247_3()
	_ = _47_12 // possibly unused.
	_49_4 := mmagicv248_3()
	_ = _49_4 // possibly unused.
	_58_4, _59_4, _62_4 := mmagicv257_3()
	_, _, _ = _58_4, _59_4, _62_4 // possibly unused.
	_65_4, _72_4 := mmagicv264_3()
	_, _ = _65_4, _72_4 // possibly unused.

	var v1 *RequestV2 = _33_14
	emitter := cff.NopEmitter()

	var (
		flowInfo = &cff.FlowInfo{
			File:   "go.uber.org/cff/examples/magic_v2.go",
			Line:   32,
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
			Concurrency: _35_19, Emitter: schedEmitter,
		},
	)

	var tasks []*struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	}

	// go.uber.org/cff/examples/magic_v2.go:38:4
	var (
		v2 *GetManagerRequestV2
		v3 *ListUsersRequestV2
	)
	task0 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task0.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v2, v3 = _38_4(v1)
		return
	}

	task0.job = sched.Enqueue(ctx, cff.Job{
		Run: task0.run,
	})

	tasks = append(tasks, task0)

	// go.uber.org/cff/examples/magic_v2.go:46:4
	var (
		v4 *GetManagerResponseV2
	)
	task1 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task1.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v4, err = _46_4(v2)
		return
	}

	task1.job = sched.Enqueue(ctx, cff.Job{
		Run: task1.run,
		Dependencies: []*cff.ScheduledJob{
			task0.job,
		},
	})

	tasks = append(tasks, task1)

	// go.uber.org/cff/examples/magic_v2.go:58:4
	var (
		v5 *ListUsersResponseV2
	)
	task4 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task4.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v5, err = _58_4(v3)
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

	// go.uber.org/cff/examples/magic_v2.go:65:4
	var (
		v6 []*SendEmailRequestV2
	)
	task5 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task5.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v6 = _65_4(v4, v5)
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

	// go.uber.org/cff/examples/magic_v2.go:47:12
	var (
		v7 []*SendEmailResponseV2
	)
	task2 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task2.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v7, err = _47_12(v6)
		return
	}

	task2.job = sched.Enqueue(ctx, cff.Job{
		Run: task2.run,
		Dependencies: []*cff.ScheduledJob{
			task5.job,
		},
	})

	tasks = append(tasks, task2)

	// go.uber.org/cff/examples/magic_v2.go:49:4
	var (
		v8 *ResponseV2
	)
	task3 := new(struct {
		emitter cff.TaskEmitter
		ran     cff.AtomicBool
		run     func(context.Context) error
		job     *cff.ScheduledJob
	})

	task3.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = &cff.PanicError{
					Value:      recovered,
					Stacktrace: debug.Stack(),
				}
			}
		}()

		v8 = _49_4(v7)
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

	*(_34_15) = v8 // *go.uber.org/cff/examples.ResponseV2

	flowEmitter.FlowSuccess(ctx)
	return nil
}

func _cffParamsmagicv2_33_3(mmagicv233_14 *RequestV2) func() *RequestV2 {
	return func() *RequestV2 { return mmagicv233_14 }
}

func _cffResultsmagicv2_34_3(mmagicv234_15 **ResponseV2) func() **ResponseV2 {
	return func() **ResponseV2 { return mmagicv234_15 }
}

func _cffConcurrencymagicv2_35_3(c int) func() int {
	return func() int { return c }
}

func _cffTaskmagicv2_37_3(mmagicv238_4 func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2)) func() func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2) {
	return func() func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2) { return mmagicv238_4 }
}

func _cffTaskmagicv2_45_3(mmagicv246_4 func(req *GetManagerRequestV2) (*GetManagerResponseV2, error)) func() func(req *GetManagerRequestV2) (*GetManagerResponseV2, error) {
	return func() func(req *GetManagerRequestV2) (*GetManagerResponseV2, error) { return mmagicv246_4 }
}

func _cffTaskmagicv2_47_3(mmagicv247_12 func(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error)) func() func(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error) {
	return func() func(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error) { return mmagicv247_12 }
}

func _cffTaskmagicv2_48_3(mmagicv249_4 func(responses []*SendEmailResponseV2) *ResponseV2) func() func(responses []*SendEmailResponseV2) *ResponseV2 {
	return func() func(responses []*SendEmailResponseV2) *ResponseV2 { return mmagicv249_4 }
}

func _cffTaskmagicv2_57_3(mmagicv258_4 func(req *ListUsersRequestV2) (*ListUsersResponseV2, error), mmagicv259_4 cff.TaskOption, mmagicv262_4 cff.TaskOption) func() (func(req *ListUsersRequestV2) (*ListUsersResponseV2, error), cff.TaskOption, cff.TaskOption) {
	return func() (func(req *ListUsersRequestV2) (*ListUsersResponseV2, error), cff.TaskOption, cff.TaskOption) {
		return mmagicv258_4, mmagicv259_4, mmagicv262_4
	}
}

func _cffTaskmagicv2_64_3(mmagicv265_4 func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2, mmagicv272_4 cff.TaskOption) func() (func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2, cff.TaskOption) {
	return func() (func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2, cff.TaskOption) {
		return mmagicv265_4, mmagicv272_4
	}
}
