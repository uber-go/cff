// +build cff

package example

import (
	"context"
	"strconv"

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
	err := cff.Flow(ctx,
		cff.Params(req),
		cff.Results(&res),
		cff.Scope(h.scope),
		cff.Logger(h.logger),
		cff.InstrumentFlow("HandleFoo"),

		cff.Tasks(
			func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
				return &GetManagerRequest{
						LDAPGroup: req.LDAPGroup,
					}, &ListUsersRequest{
						LDAPGroup: req.LDAPGroup,
					}
			},

			h.mgr.Get,
			h.ses.BatchSendEmail,

			func(responses []*SendEmailResponse) *Response {
				var r Response
				for _, res := range responses {
					r.MessageIDs = append(r.MessageIDs, res.MessageID)
				}
				return &r
			},
		),
		cff.Task(
			h.users.List,
			cff.RecoverWith(&ListUsersResponse{}),
			cff.Instrument("FormSendEmailRequest"),
		),
		cff.Task(
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
