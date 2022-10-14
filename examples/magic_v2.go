//go:build cff
// +build cff

package example

import (
	"context"
	"fmt"
	"strconv"

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
	mgr    *ManagerRepositoryV2
	users  *UserRepositoryV2
	ses    *SESClientV2
	scope  tally.Scope
	logger *zap.Logger
}

func (h *fooHandlerV2) HandleFoo(ctx context.Context, req *RequestV2) (*ResponseV2, error) {
	var res *ResponseV2
	err := cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.Concurrency(8),
		cff.WithEmitter(cff.TallyEmitter(h.scope)),
		cff.WithEmitter(cff.LogEmitter(h.logger)),
		cff.InstrumentFlow("HandleFoo"),

		cff.Task(
			func(req *RequestV2) (*GetManagerRequestV2, *ListUsersRequestV2) {
				return &GetManagerRequestV2{
						LDAPGroup: req.LDAPGroup,
					}, &ListUsersRequestV2{
						LDAPGroup: req.LDAPGroup,
					}
			}),
		cff.Task(
			h.mgr.Get),
		cff.Task(h.ses.BatchSendEmail),
		cff.Task(
			func(responses []*SendEmailResponseV2) *ResponseV2 {
				var r ResponseV2
				for _, res := range responses {
					r.MessageIDs = append(r.MessageIDs, res.MessageID)
				}
				return &r
			},
		),
		cff.Task(
			h.users.List,
			cff.Predicate(func(req *GetManagerRequestV2) bool {
				return req.LDAPGroup != "everyone"
			}),
			cff.FallbackWith(&ListUsersResponseV2{}),
			cff.Instrument("FormSendEmailRequest"),
		),
		cff.Task(
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
