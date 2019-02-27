// +build !cff

package example

import (
	"context"
	"strconv"
	"sync"

	"github.com/uber-go/tally"
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
	mgr   *ManagerRepository
	users *UserReposistory
	ses   *SESClient
	scope tally.Scope
}

func (h *fooHandler) HandleFoo(ctx context.Context, req *Request) (*Response, error) {
	var res *Response
	err := func(ctx context.Context, v1 *Request) (err error) {

		if ctx.Err() != nil {
			h.scope.Counter("task.skipped").Inc(1)
			h.scope.Counter("task.skipped").Inc(1)
			h.scope.Counter("taskflow.skipped").Inc(1)
			return ctx.Err()
		}

		var (
			v2 *GetManagerRequest
			v3 *ListUsersRequest
		)
		v2, v3 = func(req *Request) (*GetManagerRequest, *ListUsersRequest) {
			return &GetManagerRequest{
					LDAPGroup: req.LDAPGroup,
				}, &ListUsersRequest{
					LDAPGroup: req.LDAPGroup,
				}
		}(v1)

		if ctx.Err() != nil {
			h.scope.Counter("task.skipped").Inc(1)
			h.scope.Counter("task.skipped").Inc(1)
			h.scope.Counter("taskflow.skipped").Inc(1)
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(2)

		var v4 *GetManagerResponse
		go func() {
			defer wg1.Done()

			var err1 error
			v4, err = h.mgr.Get(v2)
			if err1 != nil {

				once1.Do(func() {
					err = err1
				})
			}

		}()

		var v5 *ListUsersResponse
		go func() {
			defer wg1.Done()
			timer := h.scope.Timer("task.timing").Start()
			defer timer.Stop()

			var err1 error
			v5, err = h.users.List(v3)
			if err1 != nil {
				h.scope.Counter("task.error").Inc(1)
				h.scope.Counter("task.recovered").Inc(1)

				v5, err = &ListUsersResponse{}, nil
			} else {
				h.scope.Counter("task.success").Inc(1)
			}

		}()

		wg1.Wait()
		if err != nil {
			h.scope.Counter("taskflow.error").Inc(1)
			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v4
			_ = &v5
		)

		if ctx.Err() != nil {
			h.scope.Counter("task.skipped").Inc(1)
			h.scope.Counter("taskflow.skipped").Inc(1)
			return ctx.Err()
		}

		var v6 []*SendEmailRequest
		if func(req *GetManagerRequest) bool {
			return req.LDAPGroup != "everyone"
		}(v2) {
			v6 = func(mgr *GetManagerResponse, users *ListUsersResponse) []*SendEmailRequest {
				var reqs []*SendEmailRequest
				for _, u := range users.Emails {
					reqs = append(reqs, &SendEmailRequest{Address: u})
				}
				return reqs
			}(v4, v5)

		} else {
			h.scope.Counter("task.skipped").Inc(1)
		}

		if ctx.Err() != nil {
			h.scope.Counter("taskflow.skipped").Inc(1)
			return ctx.Err()
		}

		var v7 []*SendEmailResponse
		v7, err = h.ses.BatchSendEmail(v6)
		if err != nil {

			h.scope.Counter("taskflow.error").Inc(1)
			return err
		}

		if ctx.Err() != nil {
			h.scope.Counter("taskflow.skipped").Inc(1)
			return ctx.Err()
		}

		var v8 *Response
		v8 = func(responses []*SendEmailResponse) *Response {
			var r Response
			for _, res := range responses {
				r.MessageIDs = append(r.MessageIDs, res.MessageID)
			}
			return &r
		}(v7)

		*(&res) = v8

		if err != nil {
			h.scope.Counter("taskflow.error").Inc(1)
		} else {
			h.scope.Counter("taskflow.success").Inc(1)
		}

		return err
	}(ctx, req)
	return res, err
}

type ManagerRepository struct{}

type GetManagerRequest struct {
	LDAPGroup string
}

type GetManagerResponse struct {
	Email string
}

func (*ManagerRepository) Get(req *GetManagerRequest) (*GetManagerResponse, error) {
	return &GetManagerResponse{Email: "boss@example.com"}, nil
}

type UserReposistory struct{}

type ListUsersRequest struct {
	LDAPGroup string
}

type ListUsersResponse struct {
	Emails []string
}

func (*UserReposistory) List(req *ListUsersRequest) (*ListUsersResponse, error) {
	return &ListUsersResponse{
		Emails: []string{"a@example.com", "b@example.com"},
	}, nil
}

type SESClient struct{}

type SendEmailRequest struct {
	Address string
}

type SendEmailResponse struct {
	MessageID string
}

func (*SESClient) BatchSendEmail(req []*SendEmailRequest) ([]*SendEmailResponse, error) {
	res := make([]*SendEmailResponse, len(req))
	for i := range req {
		res[i] = &SendEmailResponse{MessageID: strconv.Itoa(i)}
	}
	return res, nil
}
