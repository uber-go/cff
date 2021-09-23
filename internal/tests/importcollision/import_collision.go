package importcollision

import (
	"context"

	"go.uber.org/cff/internal/tests/importcollision/template"

	cff2 "go.uber.org/cff"
)

// Flow tests a flow that requires code generation to resolve multiple imports with
// the same base path.
func Flow() (string, error) {
	var result string
	err := cff2.Flow(
		context.Background(),
		cff2.Results(&result),
		cff2.Task(GetHTMLTemplate),
		cff2.Task(GetTextTemplate),
		cff2.Task(GetFoo),
		cff2.Task(GetResult),
		cff2.Task(template.GetError, cff2.Invoke(true)),
	)

	return result, err
}
