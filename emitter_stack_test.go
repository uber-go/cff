package cff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"go.uber.org/cff"
	"github.com/golang/mock/gomock"
	"github.com/uber-go/tally"
	"go.uber.org/zap"
)

func TestEmitterStackConstruction(t *testing.T) {
	tests := []struct {
		desc string
		give []cff.Emitter
	}{
		{desc: "empty"},
		{
			desc: "single",
			give: []cff.Emitter{
				cff.LogEmitter(zap.NewNop()),
			},
		},
		{
			desc: "multiple",
			give: []cff.Emitter{
				cff.LogEmitter(zap.NewNop()),
				cff.NopEmitter(),
				cff.TallyEmitter(tally.NoopScope),
			},
		},
		{
			desc: "nested",
			give: []cff.Emitter{
				cff.LogEmitter(zap.NewNop()),
				cff.EmitterStack(
					cff.NopEmitter(),
					cff.TallyEmitter(tally.NoopScope),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx := context.Background()

			e := cff.EmitterStack(tt.give...)
			e.FlowInit(&cff.FlowInfo{Flow: "foo"}).
				FlowFailedTask(ctx, "bar", errors.New("baz")).
				FlowDone(ctx, time.Second)
		})
	}
}

type testStructs struct {
	ctrl     *gomock.Controller
	task1    *cff.MockTaskEmitter
	flow1    *cff.MockFlowEmitter
	emitter1 *cff.MockEmitter
	task2    *cff.MockTaskEmitter
	flow2    *cff.MockFlowEmitter
	emitter2 *cff.MockEmitter
	stack    cff.Emitter
}

func mocks(t *testing.T) testStructs {
	m := testStructs{}
	m.ctrl = gomock.NewController(t)
	m.task1 = cff.NewMockTaskEmitter(m.ctrl)
	m.flow1 = cff.NewMockFlowEmitter(m.ctrl)
	m.emitter1 = cff.NewMockEmitter(m.ctrl)
	m.task2 = cff.NewMockTaskEmitter(m.ctrl)
	m.flow2 = cff.NewMockFlowEmitter(m.ctrl)
	m.emitter2 = cff.NewMockEmitter(m.ctrl)
	m.stack = cff.EmitterStack(
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

			m.emitter1.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Times(1)
			m.emitter2.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Times(1)
			m.stack.FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0})
		})

		t.Run("FlowSuccess", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow1)
			m.emitter2.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow2)

			m.flow1.EXPECT().FlowSuccess(ctx)
			m.flow2.EXPECT().FlowSuccess(ctx)
			m.stack.FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).FlowSuccess(ctx)
		})
		t.Run("FlowError", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow1)
			m.emitter2.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow2)

			err := errors.New("foobar")
			m.flow1.EXPECT().FlowError(ctx, err)
			m.flow2.EXPECT().FlowError(ctx, err)
			m.stack.FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).FlowError(ctx, err)
		})
		t.Run("FlowSkipped", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow1)
			m.emitter2.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow2)

			err := errors.New("foobar")
			m.flow1.EXPECT().FlowSkipped(ctx, err)
			m.flow2.EXPECT().FlowSkipped(ctx, err)
			m.stack.FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).FlowSkipped(ctx, err)
		})
		t.Run("FlowDone", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow1)
			m.emitter2.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow2)

			m.flow1.EXPECT().FlowDone(ctx, time.Duration(1))
			m.flow2.EXPECT().FlowDone(ctx, time.Duration(1))
			m.stack.FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).FlowDone(ctx, time.Duration(1))
		})
		t.Run("FlowFailedTask", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow1)
			m.emitter2.EXPECT().FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).Return(m.flow2)

			newFlow1 := cff.NewMockFlowEmitter(m.ctrl)
			newFlow2 := cff.NewMockFlowEmitter(m.ctrl)

			err := errors.New("foobar")

			m.flow1.EXPECT().FlowFailedTask(ctx, "foobar", err).Return(newFlow1)
			m.flow2.EXPECT().FlowFailedTask(ctx, "foobar", err).Return(newFlow2)

			newEmitter := m.stack.FlowInit(&cff.FlowInfo{"foo", "foo.go", 0, 0}).FlowFailedTask(ctx, "foobar", err)

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

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16},
				&cff.FlowInfo{"fooFlow", "foo.go", 10, 12},
			).Times(1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16},
				&cff.FlowInfo{"fooFlow", "foo.go", 10, 12},
			).Times(1)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16},
				&cff.FlowInfo{"fooFlow", "foo.go", 10, 12},
			)
		})

		t.Run("TaskSuccess", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			m.task1.EXPECT().TaskSuccess(ctx)
			m.task2.EXPECT().TaskSuccess(ctx)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskSuccess(ctx)
		})
		t.Run("TaskError", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("foobar")

			m.task1.EXPECT().TaskError(ctx, err)
			m.task2.EXPECT().TaskError(ctx, err)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskError(ctx, err)
		})
		t.Run("TaskErrorRecovered", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("great sadness")

			m.task1.EXPECT().TaskErrorRecovered(ctx, err)
			m.task2.EXPECT().TaskErrorRecovered(ctx, err)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskErrorRecovered(ctx, err)
		})
		t.Run("TaskSkipped", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("foobar")

			m.task1.EXPECT().TaskSkipped(ctx, err)
			m.task2.EXPECT().TaskSkipped(ctx, err)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskSkipped(ctx, err)
		})
		t.Run("TaskPanic", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			pv := int(1)

			m.task1.EXPECT().TaskPanic(ctx, pv)
			m.task2.EXPECT().TaskPanic(ctx, pv)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskPanic(ctx, pv)
		})
		t.Run("TaskPanicRecovered", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			pv := int(1)

			m.task1.EXPECT().TaskPanicRecovered(ctx, pv)
			m.task2.EXPECT().TaskPanicRecovered(ctx, pv)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskPanicRecovered(ctx, pv)
		})
		t.Run("TaskDone", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).Return(m.task2)

			m.task1.EXPECT().TaskDone(ctx, time.Duration(1))
			m.task2.EXPECT().TaskDone(ctx, time.Duration(1))
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.FlowInfo{"fooFlow", "foo.go", 10, 12}).TaskDone(ctx, time.Duration(1))
		})
	})
}
