package archive_tests

import (
	"context"

	errors "example.import/archivedata"
)

func Test(ctx context.Context) {
	_ = errors.New("hello")
}
