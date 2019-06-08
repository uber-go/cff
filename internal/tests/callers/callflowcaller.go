package callers

import (
	"go.uber.org/cff/internal/tests/sandwich"
)

// PackageCall exports a function that calls a "sandwich" CFF flow and is used by a unit test.
func PackageCall() (string, string) {
	return sandwich.CallFlow()
}
