// +build !cff

package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"

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

		if ctx.Err() != nil {
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "Atoi"),
				zap.Error(ctx.Err()),
			)
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "uint8"),
				zap.Error(ctx.Err()),
			)
			scope.Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "AtoiRun"))
			return ctx.Err()
		}

		var v2 int
		v2, err = strconv.Atoi(v1)
		if err != nil {
			scope.Counter("task.error").Inc(1)
			scope.Counter("taskflow.error").Inc(1)
			return err
		} else {
			scope.Counter("task.success").Inc(1)
			logger.Debug("task succeeded", zap.String("name", "Atoi"))
		}

		if ctx.Err() != nil {
			scope.Counter("task.skipped").Inc(1)
			logger.Debug("task skipped",
				zap.String("name", "uint8"),
				zap.Error(ctx.Err()),
			)
			scope.Counter("taskflow.skipped").Inc(1)
			logger.Debug("taskflow skipped", zap.String("name", "AtoiRun"))
			return ctx.Err()
		}

		var v3 uint8
		v3, err = func(i int) (uint8, error) {
			if i > -1 && i < 256 {
				return uint8(i), nil
			}
			return 0, errors.New("int can not fit into 8 bits")
		}(v2)
		if err != nil {
			scope.Counter("task.error").Inc(1)
			scope.Counter("task.recovered").Inc(1)
			logger.Error("task error recovered",
				zap.String("name", "uint8"),
				zap.Error(err),
			)
			v3, err = uint8(0), nil
		} else {
			scope.Counter("task.success").Inc(1)
			logger.Debug("task succeeded", zap.String("name", "uint8"))
		}

		*(&res) = v3

		if err != nil {
			scope.Counter("taskflow.error").Inc(1)
		} else {
			scope.Counter("taskflow.success").Inc(1)
			logger.Debug("taskflow succeeded", zap.String("name", "AtoiRun"))
		}

		return err
	}(ctx, h.scope, h.logger, req)
	return
}
