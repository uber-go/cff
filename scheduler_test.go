package cff_test

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	. "go.uber.org/cff"
	"go.uber.org/cff/internal/emittertest"
)

func TestAdaptSchedulerEmitter(t *testing.T) {
	t.Parallel()

	schedInfo := &SchedulerInfo{}

	t.Run("nil emitter", func(t *testing.T) {
		assert.Nil(t, AdaptSchedulerEmitter(nil))
	})

	t.Run("no op emitter", func(t *testing.T) {
		assert.Nil(t, AdaptSchedulerEmitter(NopEmitter().SchedulerInit(schedInfo)))
	})

	t.Run("live emitter", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		e := emittertest.NewMockSchedulerEmitter(ctrl)
		assert.NotNil(t, AdaptSchedulerEmitter(e))
	})
}
