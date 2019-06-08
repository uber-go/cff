// +build cff

package externalpackage

import (
	"context"

	"github.com/gofrs/uuid"

	"go.uber.org/cff"
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
