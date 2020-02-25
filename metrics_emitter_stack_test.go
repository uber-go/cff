package cff_test

import (
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
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			m.flow1.EXPECT().FlowSuccess()
			m.flow2.EXPECT().FlowSuccess()
			m.stack.FlowInit("foo").FlowSuccess()
		})
		t.Run("FlowError", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			m.flow1.EXPECT().FlowError()
			m.flow2.EXPECT().FlowError()
			m.stack.FlowInit("foo").FlowError()
		})
		t.Run("FlowSkipped", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			m.flow1.EXPECT().FlowSkipped()
			m.flow2.EXPECT().FlowSkipped()
			m.stack.FlowInit("foo").FlowSkipped()
		})
		t.Run("FlowDone", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			m.flow1.EXPECT().FlowDone(time.Duration(1))
			m.flow2.EXPECT().FlowDone(time.Duration(1))
			m.stack.FlowInit("foo").FlowDone(time.Duration(1))
		})
		t.Run("FlowFailedTask", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit("foo").Return(m.flow1)
			m.emitter2.EXPECT().FlowInit("foo").Return(m.flow2)

			newFlow1 := cff.NewMockFlowEmitter(m.ctrl)
			newFlow2 := cff.NewMockFlowEmitter(m.ctrl)

			m.flow1.EXPECT().FlowFailedTask("foobar").Return(newFlow1)
			m.flow2.EXPECT().FlowFailedTask("foobar").Return(newFlow2)

			newEmitter := m.stack.FlowInit("foo").FlowFailedTask("foobar")

			// Asserts that the subsequent requests should go to the return-value from FlowFailedTask, not m.flow1, m.flow2

			newFlow1.EXPECT().FlowSuccess()
			newFlow2.EXPECT().FlowSuccess()

			newEmitter.FlowSuccess()
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
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			m.task1.EXPECT().TaskSuccess()
			m.task2.EXPECT().TaskSuccess()
			m.stack.TaskInit("foo").TaskSuccess()
		})
		t.Run("TaskError", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			m.task1.EXPECT().TaskError()
			m.task2.EXPECT().TaskError()
			m.stack.TaskInit("foo").TaskError()
		})
		t.Run("TaskSkipped", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			m.task1.EXPECT().TaskSkipped()
			m.task2.EXPECT().TaskSkipped()
			m.stack.TaskInit("foo").TaskSkipped()
		})
		t.Run("TaskPanic", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			m.task1.EXPECT().TaskPanic()
			m.task2.EXPECT().TaskPanic()
			m.stack.TaskInit("foo").TaskPanic()
		})
		t.Run("TaskRecovered", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			m.task1.EXPECT().TaskRecovered()
			m.task2.EXPECT().TaskRecovered()
			m.stack.TaskInit("foo").TaskRecovered()
		})
		t.Run("TaskDone", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit("foo").Return(m.task1)
			m.emitter2.EXPECT().TaskInit("foo").Return(m.task2)

			m.task1.EXPECT().TaskDone(time.Duration(1))
			m.task2.EXPECT().TaskDone(time.Duration(1))
			m.stack.TaskInit("foo").TaskDone(time.Duration(1))
		})
	})
}
