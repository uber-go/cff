// +build cff

package externalpackage

import (
	"context"

	"go.uber.org/cff"
	"go.uber.org/cff/internal/tests/externalpackage/external"
	"github.com/gofrs/uuid"
)

// NestedType is a flow that has a dep on an external struct.
func NestedType(ctx context.Context, driverUUID uuid.UUID) error {
	var ok bool
	return cff.Flow(ctx,
		cff.Params(driverUUID),
		cff.Results(&ok),

		cff.Task(func(c uuid.UUID) bool {
			return true
		}),
	)
}

// ImplicitType is a flow that has a dep on an external struct where the package is not imported in this file
func ImplicitType(ctx context.Context) error {
	var ok bool
	return cff.Flow(ctx,
		cff.Results(&ok),

		cff.Task(external.ProvidesUUID),
		cff.Task(external.NeedsUUID),
	)
}
