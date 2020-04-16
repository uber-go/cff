package cff

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

func TestEmitterStackConstruction(t *testing.T) {
	tests := []struct {
		desc string
		give []Emitter
	}{
		{desc: "empty"},
		{
			desc: "single",
			give: []Emitter{
				LogEmitter(zap.NewNop()),
			},
		},
		{
			desc: "multiple",
			give: []Emitter{
				LogEmitter(zap.NewNop()),
				NopEmitter(),
				TallyEmitter(tally.NoopScope),
			},
		},
		{
			desc: "nested",
			give: []Emitter{
				LogEmitter(zap.NewNop()),
				EmitterStack(
					NopEmitter(),
					TallyEmitter(tally.NoopScope),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx := context.Background()

			e := EmitterStack(tt.give...)
			e.FlowInit(&FlowInfo{Flow: "foo"}).
				FlowDone(ctx, time.Second)
		})
	}
}

type testStructs struct {
	ctrl     *gomock.Controller
	task1    *MockTaskEmitter
	flow1    *MockFlowEmitter
	emitter1 *MockEmitter
	task2    *MockTaskEmitter
	flow2    *MockFlowEmitter
	emitter2 *MockEmitter
	stack    Emitter
}

func mocks(t *testing.T) testStructs {
	m := testStructs{}
	m.ctrl = gomock.NewController(t)
	m.task1 = NewMockTaskEmitter(m.ctrl)
	m.flow1 = NewMockFlowEmitter(m.ctrl)
	m.emitter1 = NewMockEmitter(m.ctrl)
	m.task2 = NewMockTaskEmitter(m.ctrl)
	m.flow2 = NewMockFlowEmitter(m.ctrl)
	m.emitter2 = NewMockEmitter(m.ctrl)
	m.stack = EmitterStack(
		m.emitter1,
		m.emitter2,
	)

	return m
}

func TestEmitterStack(t *testing.T) {
	t.Run("Flow", func(t *testing.T) {
		t.Run("Init", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).Times(1)
			m.emitter2.EXPECT().FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).Times(1)
			m.stack.FlowInit(&FlowInfo{"foo", "foo.go", 0, 0})
		})

		t.Run("FlowSuccess", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow1)
			m.emitter2.EXPECT().FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow2)

			m.flow1.EXPECT().FlowSuccess(ctx)
			m.flow2.EXPECT().FlowSuccess(ctx)
			m.stack.FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).FlowSuccess(ctx)
		})
		t.Run("FlowError", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow1)
			m.emitter2.EXPECT().FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow2)

			err := errors.New("foobar")
			m.flow1.EXPECT().FlowError(ctx, err)
			m.flow2.EXPECT().FlowError(ctx, err)
			m.stack.FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).FlowError(ctx, err)
		})
		t.Run("FlowDone", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow1)
			m.emitter2.EXPECT().FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow2)

			m.flow1.EXPECT().FlowDone(ctx, time.Duration(1))
			m.flow2.EXPECT().FlowDone(ctx, time.Duration(1))
			m.stack.FlowInit(&FlowInfo{"foo", "foo.go", 0, 0}).FlowDone(ctx, time.Duration(1))
		})
	})

	t.Run("Task", func(t *testing.T) {
		t.Run("Init", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16},
				&FlowInfo{"fooFlow", "foo.go", 10, 12},
			).Times(1)
			m.emitter2.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16},
				&FlowInfo{"fooFlow", "foo.go", 10, 12},
			).Times(1)
			m.stack.TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16},
				&FlowInfo{"fooFlow", "foo.go", 10, 12},
			)
		})

		t.Run("TaskSuccess", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			m.task1.EXPECT().TaskSuccess(ctx)
			m.task2.EXPECT().TaskSuccess(ctx)
			m.stack.TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskSuccess(ctx)
		})
		t.Run("TaskError", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("foobar")

			m.task1.EXPECT().TaskError(ctx, err)
			m.task2.EXPECT().TaskError(ctx, err)
			m.stack.TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskError(ctx, err)
		})
		t.Run("TaskErrorRecovered", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("great sadness")

			m.task1.EXPECT().TaskErrorRecovered(ctx, err)
			m.task2.EXPECT().TaskErrorRecovered(ctx, err)
			m.stack.TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskErrorRecovered(ctx, err)
		})
		t.Run("TaskSkipped", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("foobar")

			m.task1.EXPECT().TaskSkipped(ctx, err)
			m.task2.EXPECT().TaskSkipped(ctx, err)
			m.stack.TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskSkipped(ctx, err)
		})
		t.Run("TaskPanic", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			pv := int(1)

			m.task1.EXPECT().TaskPanic(ctx, pv)
			m.task2.EXPECT().TaskPanic(ctx, pv)
			m.stack.TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskPanic(ctx, pv)
		})
		t.Run("TaskPanicRecovered", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			pv := int(1)

			m.task1.EXPECT().TaskPanicRecovered(ctx, pv)
			m.task2.EXPECT().TaskPanicRecovered(ctx, pv)
			m.stack.TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskPanicRecovered(ctx, pv)
		})
		t.Run("TaskDone", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			m.task1.EXPECT().TaskDone(ctx, time.Duration(1))
			m.task2.EXPECT().TaskDone(ctx, time.Duration(1))
			m.stack.TaskInit(
				&TaskInfo{"foo", "foo.go", 14, 16}, &FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskDone(ctx, time.Duration(1))
		})
	})
}
