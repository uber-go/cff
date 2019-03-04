// +build !cff

package example

import (
	"context"
	"strconv"
	"sync"

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
	err := func(ctx context.Context, scope tally.Scope,
		logger *zap.Logger, v1 *Request) (err error) {

		if ctx.Err() != nil {
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)
			scope.Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "HandleFoo"))
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
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)
			scope.Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "HandleFoo"))
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
			timer := scope.Timer("task.timing").Start()
			defer timer.Stop()

			var err1 error
			v5, err = h.users.List(v3)
			if err1 != nil {
				scope.Counter("task.error").Inc(1)
				scope.Counter("task.recovered").Inc(1)
				logger.Error("task error recovered",
					zap.String("name", "FormSendEmailRequest"),
					zap.Error(err),
				)

				v5, err = &ListUsersResponse{}, nil
			} else {
				scope.Counter("task.success").Inc(1)
				logger.Debug("task succeeded", zap.String("name", "FormSendEmailRequest"))
			}

		}()

		wg1.Wait()
		if err != nil {
			scope.Counter("taskflow.error").Inc(1)
			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v4
			_ = &v5
		)

		if ctx.Err() != nil {
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)
			scope.Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "HandleFoo"))
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
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped", zap.String("name", "FormSendEmailRequest"))
		}

		if ctx.Err() != nil {
			scope.Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "HandleFoo"))
			return ctx.Err()
		}

		var v7 []*SendEmailResponse
		v7, err = h.ses.BatchSendEmail(v6)
		if err != nil {

			scope.Counter("taskflow.error").Inc(1)
			return err
		}

		if ctx.Err() != nil {
			scope.Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "HandleFoo"))
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
			scope.Counter("taskflow.error").Inc(1)
		} else {
			scope.Counter("taskflow.success").Inc(1)
			logger.Debug("taskflow succeeded", zap.String("name", "HandleFoo"))
		}

		return err
	}(ctx, h.scope, h.logger, req)
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
