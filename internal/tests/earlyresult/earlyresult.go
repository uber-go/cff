package earlyresult

import (
	"context"

	"go.uber.org/cff"
)

type foo struct{}
type bar struct{}
type baz struct{}
type qux struct{}

// EarlyResult makes sure ordering for an early cff.Results doesn't cause compiler to error.
func EarlyResult(ctx context.Context) error {
	request := int(2)
	var out *bar
	var out2 *foo
	return cff.Flow(
		ctx,
		cff.Params(request),
		cff.Results(&out, &out2),
		cff.Task(
			func(*foo) *bar {
				return &bar{}
			}),
		cff.Task(
			func(*foo) *baz {
				return &baz{}
			}),
		cff.Task(
			func(*bar, *baz) *qux {
				return &qux{}
			}),
		cff.Task(
			func(int) *foo {
				return &foo{}
			}),
		cff.Task(
			func(*qux) error {
				return nil
			},
			cff.Invoke(true),
		),
	)
}

// ConsumesResult makes sure that we can have an early cff.Results and run post-processing tasks.
func ConsumesResult() error {
	// t1 -> genService.Status_GetStatus_Args
	type t1 struct{}
	// t2 -> statusValidator.Request
	type t2 struct{}
	// t3 -> genService.StatusResponse
	type t3 struct{}
	// t4 -> statusValidator.Response
	type t4 struct{}
	// t5 -> node.NodeContext
	type t5 struct{}
	// t6 -> statusPostProcessor.Request
	type t6 struct{}
	// t7 -> statusPostProcessor.Response
	type t7 struct{}

	var v1 *t3
	var request *t1

	return cff.Flow(context.Background(),
		cff.Results(&v1),
		cff.Params(request),
		cff.Task(
			// bindRequestToGetStatusValidatorRequest
			func(*t1) *t2 { return &t2{} },
		),
		cff.Task(
			// bindRequestToGetStatusNodeContext
			func(*t4) *t5 { return &t5{} },
		),
		cff.Task(
			// e.Tasks.AppsRiderStatusValidator.Execute,
			func(*t2) (*t4, error) { return nil, nil },
		),

		cff.Task(
			// e.Tasks.AppsRiderStatusNodes.Execute,
			func(*t5) (*t3, error) { return nil, nil },
		),

		cff.Task(
			// bindResponseToGetStatusPostProcessorRequest
			func(*t3) *t6 { return &t6{} },
		),
		cff.Task(
			// e.Tasks.AppsRiderStatusPostprocessor.Execute,
			func(*t6) (*t7, error) { return nil, nil },
		),
		cff.Task(
			func(*t7) error {
				return nil
			},
			cff.Invoke(true)),
	)
}
