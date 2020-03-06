package cff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/cff"
	"github.com/golang/mock/gomock"
)

type testStructs struct {
	ctrl     *gomock.Controller
	task1    *cff.MockTaskEmitter
	flow1    *cff.MockFlowEmitter
	emitter1 *cff.MockMetricsEmitter
	task2    *cff.MockTaskEmitter
	flow2    *cff.MockFlowEmitter
	emitter2 *cff.MockMetricsEmitter
	stack    cff.MetricsEmitter
}

func mocks(t *testing.T) testStructs {
	m := testStructs{}
	m.ctrl = gomock.NewController(t)
	m.task1 = cff.NewMockTaskEmitter(m.ctrl)
	m.flow1 = cff.NewMockFlowEmitter(m.ctrl)
	m.emitter1 = cff.NewMockMetricsEmitter(m.ctrl)
	m.task2 = cff.NewMockTaskEmitter(m.ctrl)
	m.flow2 = cff.NewMockFlowEmitter(m.ctrl)
	m.emitter2 = cff.NewMockMetricsEmitter(m.ctrl)
	m.stack = cff.MetricsEmitterStack([]cff.MetricsEmitter{
		m.emitter1,
		m.emitter2,
	})

	return m
}

func TestMetricsEmitterStack(t *testing.T) {
	t.Run("Flow", func(t *testing.T) {
		t.Run("Init", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Times(1)
			m.emitter2.EXPECT().FlowInit("foo").Times(1)
			m.stack.FlowInit("foo")
		})

		t.Run("FlowSuccess", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			m.flow1.EXPECT().FlowSuccess(ctx)
			m.flow2.EXPECT().FlowSuccess(ctx)
			m.stack.FlowInit("foo").FlowSuccess(ctx)
		})
		t.Run("FlowError", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			err := errors.New("foobar")
			m.flow1.EXPECT().FlowError(ctx, err)
			m.flow2.EXPECT().FlowError(ctx, err)
			m.stack.FlowInit("foo").FlowError(ctx, err)
		})
		t.Run("FlowSkipped", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			err := errors.New("foobar")
			m.flow1.EXPECT().FlowSkipped(ctx, err)
			m.flow2.EXPECT().FlowSkipped(ctx, err)
			m.stack.FlowInit("foo").FlowSkipped(ctx, err)
		})
		t.Run("FlowDone", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			m.flow1.EXPECT().FlowDone(ctx, time.Duration(1))
			m.flow2.EXPECT().FlowDone(ctx, time.Duration(1))
			m.stack.FlowInit("foo").FlowDone(ctx, time.Duration(1))
		})
		t.Run("FlowFailedTask", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			newFlow1 := cff.NewMockFlowEmitter(m.ctrl)
			newFlow2 := cff.NewMockFlowEmitter(m.ctrl)

			err := errors.New("foobar")

			m.flow1.EXPECT().FlowFailedTask(ctx, "foobar", err).Return(newFlow1)
			m.flow2.EXPECT().FlowFailedTask(ctx, "foobar", err).Return(newFlow2)

			newEmitter := m.stack.FlowInit("foo").FlowFailedTask(ctx, "foobar", err)

			// Asserts that the subsequent requests should go to the return-value from FlowFailedTask, not m.flow1, m.flow2

			newFlow1.EXPECT().FlowSuccess(ctx)
			newFlow2.EXPECT().FlowSuccess(ctx)

			newEmitter.FlowSuccess(ctx)
		})
	})

	t.Run("Task", func(t *testing.T) {
		t.Run("Init", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Times(1)
			m.emitter2.EXPECT().TaskInit("foo").Times(1)
			m.stack.TaskInit("foo")
		})

		t.Run("TaskSuccess", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			m.task1.EXPECT().TaskSuccess(ctx)
			m.task2.EXPECT().TaskSuccess(ctx)
			m.stack.TaskInit("foo").TaskSuccess(ctx)
		})
		t.Run("TaskError", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			err := errors.New("foobar")

			m.task1.EXPECT().TaskError(ctx, err)
			m.task2.EXPECT().TaskError(ctx, err)
			m.stack.TaskInit("foo").TaskError(ctx, err)
		})
		t.Run("TaskSkipped", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			err := errors.New("foobar")

			m.task1.EXPECT().TaskSkipped(ctx, err)
			m.task2.EXPECT().TaskSkipped(ctx, err)
			m.stack.TaskInit("foo").TaskSkipped(ctx, err)
		})
		t.Run("TaskPanic", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			pv := int(1)

			m.task1.EXPECT().TaskPanic(ctx, pv)
			m.task2.EXPECT().TaskPanic(ctx, pv)
			m.stack.TaskInit("foo").TaskPanic(ctx, pv)
		})
		t.Run("TaskRecovered", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			pv := int(1)

			m.task1.EXPECT().TaskRecovered(ctx, pv)
			m.task2.EXPECT().TaskRecovered(ctx, pv)
			m.stack.TaskInit("foo").TaskRecovered(ctx, pv)
		})
		t.Run("TaskDone", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			m.task1.EXPECT().TaskDone(ctx, time.Duration(1))
			m.task2.EXPECT().TaskDone(ctx, time.Duration(1))
			m.stack.TaskInit("foo").TaskDone(ctx, time.Duration(1))
		})
	})
}
