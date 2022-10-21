//go:build !cff && v2
// +build !cff,v2

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

// RequestV2 TODO
type RequestV2 struct {
	LDAPGroup string
}

// ResponseV2 TODO
type ResponseV2 struct {
	MessageIDs []string
}

type fooHandlerV2 struct {
	mgr    *ManagerRepositoryV2
	users  *UserRepositoryV2
	ses    *SESClientV2
	scope  tally.Scope
	logger *zap.Logger
}

func (h *fooHandlerV2) HandleFoo(ctx context.Context, req *RequestV2) (*ResponseV2, error) {
	var res *ResponseV2
	err := _cffFlowmagicv2_36_9(ctx,
		_cffParamsmagicv2_37_3(req),
		_cffResultsmagicv2_38_3(&res),
		_cffConcurrencymagicv2_39_3(8),
		_cffWithEmittermagicv2_40_3(cff.TallyEmitter(h.scope)),
		_cffWithEmittermagicv2_41_3(cff.LogEmitter(h.logger)),
		_cffInstrumentFlowmagicv2_42_3("HandleFoo"),

		_cffTaskmagicv2_44_3(
			func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2) {
				return &GetManagerRequestV2{
						LDAPGroup: req.LDAPGroup,
					}, &ListUsersRequestV2{
						LDAPGroup: req.LDAPGroup,
					}
			}),
		_cffTaskmagicv2_52_3(
			h.mgr.Get),
		_cffTaskmagicv2_54_3(h.ses.BatchSendEmail),
		_cffTaskmagicv2_55_3(
			func(responses []*SendEmailResponseV2) *ResponseV2 {
				var r ResponseV2
				for _, res := range responses {
					r.MessageIDs = append(r.MessageIDs, res.MessageID)
				}
				return &r
			},
		),
		_cffTaskmagicv2_64_3(
			h.users.List,
			cff.Predicate(func(req *GetManagerRequestV2) bool {
				return req.LDAPGroup != "everyone"
			}),
			cff.FallbackWith(&ListUsersResponseV2{}),
			cff.Instrument("FormSendEmailRequest"),
		),
		_cffTaskmagicv2_72_3(
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
			cff.Instrument("FormSendEmailRequest"),
		),
	)

	err = cff.Parallel(
		ctx,
		cff.Concurrency(2),
		cff.WithEmitter(cff.TallyEmitter(h.scope)),
		cff.WithEmitter(cff.LogEmitter(h.logger)),
		cff.InstrumentParallel("SendParallel"),
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
			cff.Instrument("SendMsg"),
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
func _cffFlowmagicv2_36_9(
	ctx context.Context,
	mmagicv237_3 func() *RequestV2,
	mmagicv238_3 func() **ResponseV2,
	mmagicv239_3 func() int,
	mmagicv240_3 func() cff.Emitter,
	mmagicv241_3 func() cff.Emitter,
	mmagicv242_3 func() string,
	mmagicv244_3 func() func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2),
	mmagicv252_3 func() func(req *GetManagerRequestV2) (*GetManagerResponseV2, error),
	mmagicv254_3 func() func(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error),
	mmagicv255_3 func() func(responses []*SendEmailResponseV2) *ResponseV2,
	mmagicv264_3 func() (func(req *ListUsersRequestV2) (*ListUsersResponseV2, error), cff.TaskOption, cff.TaskOption, cff.TaskOption),
	_ cff.Option,
	mmagicv272_3 func() (func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2, cff.TaskOption, cff.TaskOption),
	_ cff.Option,
) error {
	_37_14 := mmagicv237_3()
	_ = _37_14 // possibly unused.
	_38_15 := mmagicv238_3()
	_ = _38_15 // possibly unused.
	_39_19 := mmagicv239_3()
	_ = _39_19 // possibly unused.
	_40_19 := mmagicv240_3()
	_ = _40_19 // possibly unused.
	_41_19 := mmagicv241_3()
	_ = _41_19 // possibly unused.
	_42_22 := mmagicv242_3()
	_ = _42_22 // possibly unused.
	_45_4 := mmagicv244_3()
	_ = _45_4 // possibly unused.
	_53_4 := mmagicv252_3()
	_ = _53_4 // possibly unused.
	_54_12 := mmagicv254_3()
	_ = _54_12 // possibly unused.
	_56_4 := mmagicv255_3()
	_ = _56_4 // possibly unused.
	_65_4, _66_4, _69_4, _70_4 := mmagicv264_3()
	_, _, _, _ = _65_4, _66_4, _69_4, _70_4 // possibly unused.
	_73_4, _80_4, _83_4 := mmagicv272_3()
	_, _, _ = _73_4, _80_4, _83_4 // possibly unused.

	var v1 *RequestV2 = _37_14
	emitter := cff.EmitterStack(_40_19, _41_19)

	var (
		flowInfo = &cff.FlowInfo{
			Name:   _42_22,
			File:   "go.uber.org/cff/examples/magic_v2.go",
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

	// go.uber.org/cff/examples/magic_v2.go:45:4
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
				err = fmt.Errorf("task panic: %v", recovered)
			}
		}()

		v2, v3 = _45_4(v1)
		return
	}

	task0.job = sched.Enqueue(ctx, cff.Job{
		Run: task0.run,
	})

	tasks = append(tasks, task0)

	// go.uber.org/cff/examples/magic_v2.go:53:4
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
				err = fmt.Errorf("task panic: %v", recovered)
			}
		}()

		v4, err = _53_4(v2)
		return
	}

	task1.job = sched.Enqueue(ctx, cff.Job{
		Run: task1.run,
		Dependencies: []*cff.ScheduledJob{
			task0.job,
		},
	})

	tasks = append(tasks, task1)

	// go.uber.org/cff/examples/magic_v2.go:65:4
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
				err = fmt.Errorf("task panic: %v", recovered)
			}
		}()

		v5, err = _65_4(v3)
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

	// go.uber.org/cff/examples/magic_v2.go:73:4
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
				err = fmt.Errorf("task panic: %v", recovered)
			}
		}()

		v6 = _73_4(v4, v5)
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

	// go.uber.org/cff/examples/magic_v2.go:54:12
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
				err = fmt.Errorf("task panic: %v", recovered)
			}
		}()

		v7, err = _54_12(v6)
		return
	}

	task2.job = sched.Enqueue(ctx, cff.Job{
		Run: task2.run,
		Dependencies: []*cff.ScheduledJob{
			task5.job,
		},
	})

	tasks = append(tasks, task2)

	// go.uber.org/cff/examples/magic_v2.go:56:4
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
				err = fmt.Errorf("task panic: %v", recovered)
			}
		}()

		v8 = _56_4(v7)
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

	*(_38_15) = v8 // *go.uber.org/cff/examples.ResponseV2

	flowEmitter.FlowSuccess(ctx)
	return nil
}

func _cffParamsmagicv2_37_3(mmagicv237_14 *RequestV2) func() *RequestV2 {
	return func() *RequestV2 { return mmagicv237_14 }
}

func _cffResultsmagicv2_38_3(mmagicv238_15 **ResponseV2) func() **ResponseV2 {
	return func() **ResponseV2 { return mmagicv238_15 }
}

func _cffConcurrencymagicv2_39_3(c int) func() int {
	return func() int { return c }
}

func _cffWithEmittermagicv2_40_3(mmagicv240_19 cff.Emitter) func() cff.Emitter {
	return func() cff.Emitter { return mmagicv240_19 }
}

func _cffWithEmittermagicv2_41_3(mmagicv241_19 cff.Emitter) func() cff.Emitter {
	return func() cff.Emitter { return mmagicv241_19 }
}

func _cffInstrumentFlowmagicv2_42_3(mmagicv242_22 string) func() string {
	return func() string { return mmagicv242_22 }
}

func _cffTaskmagicv2_44_3(mmagicv245_4 func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2)) func() func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2) {
	return func() func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2) { return mmagicv245_4 }
}

func _cffTaskmagicv2_52_3(mmagicv253_4 func(req *GetManagerRequestV2) (*GetManagerResponseV2, error)) func() func(req *GetManagerRequestV2) (*GetManagerResponseV2, error) {
	return func() func(req *GetManagerRequestV2) (*GetManagerResponseV2, error) { return mmagicv253_4 }
}

func _cffTaskmagicv2_54_3(mmagicv254_12 func(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error)) func() func(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error) {
	return func() func(req []*SendEmailRequestV2) ([]*SendEmailResponseV2, error) { return mmagicv254_12 }
}

func _cffTaskmagicv2_55_3(mmagicv256_4 func(responses []*SendEmailResponseV2) *ResponseV2) func() func(responses []*SendEmailResponseV2) *ResponseV2 {
	return func() func(responses []*SendEmailResponseV2) *ResponseV2 { return mmagicv256_4 }
}

func _cffTaskmagicv2_64_3(mmagicv265_4 func(req *ListUsersRequestV2) (*ListUsersResponseV2, error), mmagicv266_4 cff.TaskOption, mmagicv269_4 cff.TaskOption, mmagicv270_4 cff.TaskOption) func() (func(req *ListUsersRequestV2) (*ListUsersResponseV2, error), cff.TaskOption, cff.TaskOption, cff.TaskOption) {
	return func() (func(req *ListUsersRequestV2) (*ListUsersResponseV2, error), cff.TaskOption, cff.TaskOption, cff.TaskOption) {
		return mmagicv265_4, mmagicv266_4, mmagicv269_4, mmagicv270_4
	}
}

func _cffTaskmagicv2_72_3(mmagicv273_4 func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2, mmagicv280_4 cff.TaskOption, mmagicv283_4 cff.TaskOption) func() (func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2, cff.TaskOption, cff.TaskOption) {
	return func() (func(mgr *GetManagerResponseV2, users *ListUsersResponseV2) []*SendEmailRequestV2, cff.TaskOption, cff.TaskOption) {
		return mmagicv273_4, mmagicv280_4, mmagicv283_4
	}
}
