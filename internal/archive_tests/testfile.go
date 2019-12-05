// +build ignore

package archivetests

import (
	"context"

	errors "example.import/archivedata"
)

// Test function
func Test(ctx context.Context) {
	_ = errors.New("hello")
}
