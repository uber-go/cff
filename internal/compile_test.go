package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/types/typeutil"
)

func TestNoOutputTypes(t *testing.T) {
	defer func() {
		err := recover()
		if err == nil {
			t.Fatalf("expected mustSetNoOutputProvider to panic")
		}
	}()
	providers := &typeutil.Map{}
	f := &flow{
		providers: providers,
	}
	noOutput := f.addNoOutput()
	task := &task{
		invokeType: noOutput,
	}
	f.mustSetNoOutputProvider(task, 0)
	f.mustSetNoOutputProvider(task, 0)
	assert.Equal(t, task.invokeType, f.providers.At(task.invokeType))
}
