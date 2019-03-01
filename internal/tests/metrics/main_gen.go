// +build !cff

package main

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/uber-go/tally"
)

func main() {
	scope := tally.NoopScope
	h := &h{scope: scope}
	ctx := context.Background()
	res, err := h.run(ctx, os.Args[1])
	if err != nil {
		panic(err)
	}
	fmt.Printf("%d\n", res)
}

type h struct {
	scope tally.Scope
}

func (h *h) run(ctx context.Context, req string) (res int, err error) {
	err = func(ctx context.Context, scope tally.Scope, v1 string) (err error) {

		if ctx.Err() != nil {
			scope.Counter("task.skipped").Inc(1)
			scope.Counter("taskflow.skipped").Inc(1)
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
		}

		*(&res) = v2

		if err != nil {
			scope.Counter("taskflow.error").Inc(1)
		} else {
			scope.Counter("taskflow.success").Inc(1)
		}

		return err
	}(ctx, h.scope, req)
	return
}
