// +build !cff

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"sync"

	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

func main() {
	scope := tally.NoopScope
	logger := zap.NewNop()
	h := &h{
		scope:  scope,
		logger: logger,
	}
	ctx := context.Background()
	res, err := h.run(ctx, os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", res)
}

type h struct {
	scope  tally.Scope
	logger *zap.Logger
}

func (h *h) run(ctx context.Context, req string) (res uint8, err error) {
	err = func(ctx context.Context, scope tally.Scope,
		logger *zap.Logger, v1 string) (err error) {
		flowTags := map[string]string{"name": "AtoiRun"}
		if ctx.Err() != nil {
			s0t0Tags := map[string]string{"name": "Atoi"}
			scope.Tagged(s0t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "Atoi"),
				zap.Error(ctx.Err()),
			)

			s1t0Tags := map[string]string{"name": "uint8"}
			scope.Tagged(s1t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "uint8"),
				zap.Error(ctx.Err()),
			)
			scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "AtoiRun"))
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v2 int
		var err0 error
		go func() {
			defer wg0.Done()
			tags := map[string]string{"name": "Atoi"}
			timer := scope.Tagged(tags).Timer("task.timing").Start()
			defer timer.Stop()
			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)
						scope.Tagged(map[string]string{"name": "Atoi"}).Counter("task.panic").Inc(1)
						logger.Error("task panic",
							zap.String("name", "Atoi"),
							zap.Stack("stack"),
							zap.Error(recoveredErr))
						err = recoveredErr
					})
				}
			}()

			v2, err0 = strconv.Atoi(v1)
			if err0 != nil {
				scope.Tagged(tags).Counter("task.error").Inc(1)
				once0.Do(func() {
					err = err0
				})
			} else {
				scope.Tagged(tags).Counter("task.success").Inc(1)
				logger.Debug("task succeeded", zap.String("name", "Atoi"))
			}

		}()

		wg0.Wait()
		if err != nil {
			scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v2
		)

		if ctx.Err() != nil {
			s1t0Tags := map[string]string{"name": "uint8"}
			scope.Tagged(s1t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "uint8"),
				zap.Error(ctx.Err()),
			)
			scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "AtoiRun"))
			return ctx.Err()
		}
		var (
			wg1   sync.WaitGroup
			once1 sync.Once
		)

		wg1.Add(1)
		var v3 uint8
		var err1 error
		go func() {
			defer wg1.Done()
			tags := map[string]string{"name": "uint8"}
			timer := scope.Tagged(tags).Timer("task.timing").Start()
			defer timer.Stop()
			defer func() {
				recovered := recover()
				if recovered != nil {
					once1.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)
						scope.Tagged(map[string]string{"name": "uint8"}).Counter("task.panic").Inc(1)
						logger.Error("task panic",
							zap.String("name", "uint8"),
							zap.Stack("stack"),
							zap.Error(recoveredErr))
						err = recoveredErr
					})
				}
			}()

			v3, err1 = func(i int) (uint8, error) {
				if i > -1 && i < 256 {
					return uint8(i), nil
				}
				return 0, errors.New("int can not fit into 8 bits")
			}(v2)
			if err1 != nil {
				scope.Tagged(tags).Counter("task.error").Inc(1)
				scope.Tagged(tags).Counter("task.recovered").Inc(1)
				logger.Error("task error recovered",
					zap.String("name", "uint8"),
					zap.Error(err1),
				)

				v3, err1 = uint8(0), nil
			} else {
				scope.Tagged(tags).Counter("task.success").Inc(1)
				logger.Debug("task succeeded", zap.String("name", "uint8"))
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
			_ = &v3
		)

		*(&res) = v3

		if err != nil {
			scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
		} else {
			scope.Tagged(flowTags).Counter("taskflow.success").Inc(1)
			logger.Debug("taskflow succeeded", zap.String("name", "AtoiRun"))
		}

		return err
	}(ctx, h.scope, h.logger, req)
	return
}

func (h *h) do(ctx context.Context, req string) (res int, err error) {
	err = func(ctx context.Context, scope tally.Scope,
		logger *zap.Logger, v1 string) (err error) {
		flowTags := map[string]string{"name": "AtoiDo"}
		if ctx.Err() != nil {
			s0t0Tags := map[string]string{"name": "Atoi"}
			scope.Tagged(s0t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "Atoi"),
				zap.Error(ctx.Err()),
			)
			scope.Tagged(flowTags).Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "AtoiDo"))
			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v2 int
		var err2 error
		go func() {
			defer wg0.Done()
			tags := map[string]string{"name": "Atoi"}
			timer := scope.Tagged(tags).Timer("task.timing").Start()
			defer timer.Stop()
			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)
						scope.Tagged(map[string]string{"name": "Atoi"}).Counter("task.panic").Inc(1)
						logger.Error("task panic",
							zap.String("name", "Atoi"),
							zap.Stack("stack"),
							zap.Error(recoveredErr))
						err = recoveredErr
					})
				}
			}()

			v2, err2 = strconv.Atoi(v1)
			if err2 != nil {
				scope.Tagged(tags).Counter("task.error").Inc(1)
				once0.Do(func() {
					err = err2
				})
			} else {
				scope.Tagged(tags).Counter("task.success").Inc(1)
				logger.Debug("task succeeded", zap.String("name", "Atoi"))
			}

		}()

		wg0.Wait()
		if err != nil {
			scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v2
		)

		*(&res) = v2

		if err != nil {
			scope.Tagged(flowTags).Counter("taskflow.error").Inc(1)
		} else {
			scope.Tagged(flowTags).Counter("taskflow.success").Inc(1)
			logger.Debug("taskflow succeeded", zap.String("name", "AtoiDo"))
		}

		return err
	}(ctx, h.scope, h.logger, req)
	return
}

func (h *h) work(ctx context.Context, req string) (res int, err error) {
	err = func(ctx context.Context, scope tally.Scope,
		logger *zap.Logger, v1 string) (err error) {
		if ctx.Err() != nil {
			s0t0Tags := map[string]string{"name": "Atoi"}
			scope.Tagged(s0t0Tags).Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "Atoi"),
				zap.Error(ctx.Err()),
			)

			return ctx.Err()
		}
		var (
			wg0   sync.WaitGroup
			once0 sync.Once
		)

		wg0.Add(1)
		var v2 int
		var err3 error
		go func() {
			defer wg0.Done()
			tags := map[string]string{"name": "Atoi"}
			timer := scope.Tagged(tags).Timer("task.timing").Start()
			defer timer.Stop()
			defer func() {
				recovered := recover()
				if recovered != nil {
					once0.Do(func() {
						recoveredErr := fmt.Errorf("task panic: %v", recovered)
						scope.Tagged(map[string]string{"name": "Atoi"}).Counter("task.panic").Inc(1)
						logger.Error("task panic",
							zap.String("name", "Atoi"),
							zap.Stack("stack"),
							zap.Error(recoveredErr))
						err = recoveredErr
					})
				}
			}()

			v2, err3 = strconv.Atoi(v1)
			if err3 != nil {
				scope.Tagged(tags).Counter("task.error").Inc(1)
				once0.Do(func() {
					err = err3
				})
			} else {
				scope.Tagged(tags).Counter("task.success").Inc(1)
				logger.Debug("task succeeded", zap.String("name", "Atoi"))
			}

		}()

		wg0.Wait()
		if err != nil {

			return err
		}

		// Prevent variable unused errors.
		var (
			_ = &once0
			_ = &v2
		)

		*(&res) = v2

		return err
	}(ctx, h.scope, h.logger, req)
	return
}
