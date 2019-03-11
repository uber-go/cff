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
		flowTags := map[string]string{"name": "HandleFoo"}
		if ctx.Err() != nil {
			s1t1Tags := map[string]string{"name": "FormSendEmailRequest"}
			scope.Tagged(s1t1Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)

			s2t0Tags := map[string]string{"name": "FormSendEmailRequest"}
			scope.Tagged(s2t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)
			scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
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
			s1t1Tags := map[string]string{"name": "FormSendEmailRequest"}
			scope.Tagged(s1t1Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)

			s2t0Tags := map[string]string{"name": "FormSendEmailRequest"}
			scope.Tagged(s2t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)
			scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "HandleFoo"))
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(2)
		var v4 *GetManagerResponse
		var err1 error
		go func() {
			defer wg1.Done()

			v4, err1 = h.mgr.Get(v2)
			if err1 != nil {

				once1.Do(func() {
					err = err1
				})
			}

		}()
		var v5 *ListUsersResponse
		var err4 error
		go func() {
			defer wg1.Done()
			tags := map[string]string{"name": "FormSendEmailRequest"}
			timer := scope.Tagged(tags).Timer("task.timing").Start()
			defer timer.Stop()

			v5, err4 = h.users.List(v3)
			if err4 != nil {
				scope.Tagged(tags).Counter("task.error").Inc(1)
				scope.Tagged(tags).Counter("task.recovered").Inc(1)
				logger.Error("task error recovered",
					zap.String("name", "FormSendEmailRequest"),
					zap.Error(err4),
				)

				v5, err4 = &ListUsersResponse{}, nil
			} else {
				scope.Tagged(tags).Counter("task.success").Inc(1)
				logger.Debug("task succeeded", zap.String("name", "FormSendEmailRequest"))
			}

		}()

		wg1.Wait()
		if err != nil {
			scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once1
			_ = &v4
			_ = &v5
		)

		if ctx.Err() != nil {
			s2t0Tags := map[string]string{"name": "FormSendEmailRequest"}
			scope.Tagged(s2t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "FormSendEmailRequest"),
				zap.Error(ctx.Err()),
			)
			scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "HandleFoo"))
			return ctx.Err()
		}

		s2t0Tags := map[string]string{"name": "FormSendEmailRequest"}
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
			scope.Tagged(s2t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped", zap.String("name", "FormSendEmailRequest"))
		}

		if ctx.Err() != nil {
			scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "HandleFoo"))
			return ctx.Err()
		}

		var v7 []*SendEmailResponse
		var err2 error
		v7, err2 = h.ses.BatchSendEmail(v6)
		if err2 != nil {

			scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
			return err2
		}

		if ctx.Err() != nil {
			scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
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
			scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
		} else {
			scope.Tagged(flowTags).Counter("taskflow.success").Inc(1)
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
