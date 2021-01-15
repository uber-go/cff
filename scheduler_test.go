package cff

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAdaptSchedulerEmitter(t *testing.T) {
	t.Parallel()

	schedInfo := &SchedulerInfo{FlowInfo: &FlowInfo{}}

	t.Run("nil emitter", func(t *testing.T) {
		assert.Nil(t, adaptSchedulerEmitter(nil))
	})

	t.Run("no op emitter", func(t *testing.T) {
		assert.Nil(t, adaptSchedulerEmitter(NopEmitter().SchedulerInit(schedInfo)))
	})

	t.Run("live emitter", func(t *testing.T) {
		e := LogEmitter(nil).SchedulerInit(schedInfo)
		assert.NotNil(t, adaptSchedulerEmitter(e))
	})
}
