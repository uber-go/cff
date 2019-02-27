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
	err = func(ctx context.Context, v1 string) (err error) {

		if ctx.Err() != nil {
			h.scope.Counter("task.skipped").Inc(1)
			h.scope.Counter("taskflow.skipped").Inc(1)
			return ctx.Err()
		}

		var v2 int
		v2, err = strconv.Atoi(v1)
		if err != nil {
			h.scope.Counter("task.error").Inc(1)
			h.scope.Counter("taskflow.error").Inc(1)
			return err
		} else {
			h.scope.Counter("task.success").Inc(1)
		}

		*(&res) = v2

		if err != nil {
			h.scope.Counter("taskflow.error").Inc(1)
		} else {
			h.scope.Counter("taskflow.success").Inc(1)
		}

		return err
	}(ctx, req)
	return
}
