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

// RequestV2 TODO
type RequestV2 struct {
	LDAPGroup string
}

// ResponseV2 TODO
type ResponseV2 struct {
	MessageIDs []string
}

type fooHandlerV2 struct {
	mgr    *ManagerRepository
	users  *UserRepository
	ses    *SESClient
	scope  tally.Scope
	logger *zap.Logger
}

func (h *fooHandlerV2) HandleFoo(ctx context.Context, req *Request) (*Response, error) {
	var res *Response
	err := _cffFlowmagicv2_33_9(ctx,
		_cffParamsmagicv2_34_3(req),
		_cffResultsmagicv2_35_3(&res),
		_cffConcurrencymagicv2_36_3(8),
		_cffWithEmittermagicv2_37_3(cff.TallyEmitter(h.scope)),
		_cffWithEmittermagicv2_38_3(cff.LogEmitter(h.logger)),
		_cffInstrumentFlowmagicv2_39_3("HandleFoo"),

		_cffTaskmagicv2_41_3(
			func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
				return &GetManagerRequest{
						LDAPGroup: req.LDAPGroup,
					}, &ListUsersRequest{
						LDAPGroup: req.LDAPGroup,
					}
			}),
		_cffTaskmagicv2_49_3(
			h.mgr.Get),
		_cffTaskmagicv2_51_3(h.ses.BatchSendEmail),
		_cffTaskmagicv2_52_3(
			func(responses []*SendEmailResponse) *Response {
				var r Response
				for _, res := range responses {
					r.MessageIDs = append(r.MessageIDs, res.MessageID)
				}
				return &r
			},
		),
		_cffTaskmagicv2_61_3(
			h.users.List,
			cff.Predicate(func(req *GetManagerRequest) bool {
				return req.LDAPGroup != "everyone"
			}),
			cff.FallbackWith(&ListUsersResponse{}),
			cff.Instrument("FormSendEmailRequest"),
		),
		_cffTaskmagicv2_69_3(
			func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest {
				var reqs []*SendEmailRequest
				for _, u := range users.Emails {
					reqs = append(reqs, &SendEmailRequest{Address: u})
				}
				return reqs
			},
			cff.Predicate(func(req *GetManagerRequest) bool {
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
				return SendMessage()
			},
			SendMessage,
		),
		cff.Task(
			func() error {
				return SendMessage()
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
			cff.SliceEnd(func(context.Context) error {
				return nil
			}),
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
			cff.MapEnd(func(context.Context) {
				_ = fmt.Sprint("}")
			}),
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
func (*ManagerRepositoryV2) Get(req *GetManagerRequest) (*GetManagerResponse, error) {
	return &GetManagerResponse{Email: "boss@example.com"}, nil
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
func (*UserRepositoryV2) List(req *ListUsersRequest) (*ListUsersResponse, error) {
	return &ListUsersResponse{
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
func (*SESClientV2) BatchSendEmail(req []*SendEmailRequest) ([]*SendEmailResponse, error) {
	res := make([]*SendEmailResponse, len(req))
	for i := range req {
		res[i] = &SendEmailResponse{MessageID: strconv.Itoa(i)}
	}
	return res, nil
}

// SendMessageV2 returns nil error.
func SendMessageV2() error {
	return nil
}
func _cffFlowmagicv2_33_9(
	ctx context.Context,
	mmagicv234_3 func() *Request,
	mmagicv235_3 func() **Response,
	mmagicv236_3 func() int,
	mmagicv237_3 func() cff.Emitter,
	mmagicv238_3 func() cff.Emitter,
	mmagicv239_3 func() string,
	mmagicv241_3 func() func(req *Request) (*GetManagerRequest, *ListUsersRequest),
	mmagicv249_3 func() func(req *GetManagerRequest) (*GetManagerResponse, error),
	mmagicv251_3 func() func(req []*SendEmailRequest) ([]*SendEmailResponse, error),
	mmagicv252_3 func() func(responses []*SendEmailResponse) *Response,
	mmagicv261_3 func() (func(req *ListUsersRequest) (*ListUsersResponse, error), cff.TaskOption, cff.TaskOption, cff.TaskOption),
	_ cff.Option,
	mmagicv269_3 func() (func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest, cff.TaskOption, cff.TaskOption),
	_ cff.Option,
) error {
	_34_14 := mmagicv234_3()
	_ = _34_14 // possibly unused.
	_35_15 := mmagicv235_3()
	_ = _35_15 // possibly unused.
	_36_19 := mmagicv236_3()
	_ = _36_19 // possibly unused.
	_37_19 := mmagicv237_3()
	_ = _37_19 // possibly unused.
	_38_19 := mmagicv238_3()
	_ = _38_19 // possibly unused.
	_39_22 := mmagicv239_3()
	_ = _39_22 // possibly unused.
	_42_4 := mmagicv241_3()
	_ = _42_4 // possibly unused.
	_50_4 := mmagicv249_3()
	_ = _50_4 // possibly unused.
	_51_12 := mmagicv251_3()
	_ = _51_12 // possibly unused.
	_53_4 := mmagicv252_3()
	_ = _53_4 // possibly unused.
	_62_4, _63_4, _66_4, _67_4 := mmagicv261_3()
	_, _, _, _ = _62_4, _63_4, _66_4, _67_4 // possibly unused.
	_70_4, _77_4, _80_4 := mmagicv269_3()
	_, _, _ = _70_4, _77_4, _80_4 // possibly unused.

	var v1 *Request = _34_14
	emitter := cff.EmitterStack(_37_19, _38_19)

	var (
		flowInfo = &cff.FlowInfo{
			Name:   _39_22,
			File:   "go.uber.org/cff/examples/magic_v2.go",
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

	// go.uber.org/cff/examples/magic_v2.go:42:4
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

	task0.run = func(ctx context.Context) (err error) {
		defer func() {
			recovered := recover()
			if recovered != nil {
				err = fmt.Errorf("task panic: %v", recovered)
			}
		}()

		v2, v3 = _42_4(v1)
		return
	}

	task0.job = sched.Enqueue(ctx, cff.Job{
		Run: task0.run,
	})

	tasks = append(tasks, task0)

	// go.uber.org/cff/examples/magic_v2.go:50:4
	var (
		v4 *GetManagerResponse
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

		v4, err = _50_4(v2)
		return
	}

	task1.job = sched.Enqueue(ctx, cff.Job{
		Run: task1.run,
		Dependencies: []*cff.ScheduledJob{
			task0.job,
		},
	})

	tasks = append(tasks, task1)

	// go.uber.org/cff/examples/magic_v2.go:62:4
	var (
		v5 *ListUsersResponse
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

		v5, err = _62_4(v3)
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

	// go.uber.org/cff/examples/magic_v2.go:70:4
	var (
		v6 []*SendEmailRequest
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

		v6 = _70_4(v4, v5)
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

	// go.uber.org/cff/examples/magic_v2.go:51:12
	var (
		v7 []*SendEmailResponse
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

		v7, err = _51_12(v6)
		return
	}

	task2.job = sched.Enqueue(ctx, cff.Job{
		Run: task2.run,
		Dependencies: []*cff.ScheduledJob{
			task5.job,
		},
	})

	tasks = append(tasks, task2)

	// go.uber.org/cff/examples/magic_v2.go:53:4
	var (
		v8 *Response
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

		v8 = _53_4(v7)
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
}

func _cffParamsmagicv2_34_3(mmagicv234_14 *Request) func() *Request {
	return func() *Request { return mmagicv234_14 }
}

func _cffResultsmagicv2_35_3(mmagicv235_15 **Response) func() **Response {
	return func() **Response { return mmagicv235_15 }
}

func _cffConcurrencymagicv2_36_3(c int) func() int {
	return func() int { return c }
}

func _cffWithEmittermagicv2_37_3(mmagicv237_19 cff.Emitter) func() cff.Emitter {
	return func() cff.Emitter { return mmagicv237_19 }
}

func _cffWithEmittermagicv2_38_3(mmagicv238_19 cff.Emitter) func() cff.Emitter {
	return func() cff.Emitter { return mmagicv238_19 }
}

func _cffInstrumentFlowmagicv2_39_3(mmagicv239_22 string) func() string {
	return func() string { return mmagicv239_22 }
}

func _cffTaskmagicv2_41_3(mmagicv242_4 func(req *Request) (*GetManagerRequest, *ListUsersRequest)) func() func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
	return func() func(req *Request) (*GetManagerRequest, *ListUsersRequest) { return mmagicv242_4 }
}

func _cffTaskmagicv2_49_3(mmagicv250_4 func(req *GetManagerRequest) (*GetManagerResponse, error)) func() func(req *GetManagerRequest) (*GetManagerResponse, error) {
	return func() func(req *GetManagerRequest) (*GetManagerResponse, error) { return mmagicv250_4 }
}

func _cffTaskmagicv2_51_3(mmagicv251_12 func(req []*SendEmailRequest) ([]*SendEmailResponse, error)) func() func(req []*SendEmailRequest) ([]*SendEmailResponse, error) {
	return func() func(req []*SendEmailRequest) ([]*SendEmailResponse, error) { return mmagicv251_12 }
}

func _cffTaskmagicv2_52_3(mmagicv253_4 func(responses []*SendEmailResponse) *Response) func() func(responses []*SendEmailResponse) *Response {
	return func() func(responses []*SendEmailResponse) *Response { return mmagicv253_4 }
}

func _cffTaskmagicv2_61_3(mmagicv262_4 func(req *ListUsersRequest) (*ListUsersResponse, error), mmagicv263_4 cff.TaskOption, mmagicv266_4 cff.TaskOption, mmagicv267_4 cff.TaskOption) func() (func(req *ListUsersRequest) (*ListUsersResponse, error), cff.TaskOption, cff.TaskOption, cff.TaskOption) {
	return func() (func(req *ListUsersRequest) (*ListUsersResponse, error), cff.TaskOption, cff.TaskOption, cff.TaskOption) {
		return mmagicv262_4, mmagicv263_4, mmagicv266_4, mmagicv267_4
	}
}

func _cffTaskmagicv2_69_3(mmagicv270_4 func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest, mmagicv277_4 cff.TaskOption, mmagicv280_4 cff.TaskOption) func() (func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest, cff.TaskOption, cff.TaskOption) {
	return func() (func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest, cff.TaskOption, cff.TaskOption) {
		return mmagicv270_4, mmagicv277_4, mmagicv280_4
	}
}
