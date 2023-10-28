package cff_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"go.uber.org/cff"
	"go.uber.org/cff/internal/emittertest"
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
				cff.NopEmitter(),
			},
		},
		{
			desc: "multiple",
			give: []cff.Emitter{
				cff.NopEmitter(),
				cff.NopEmitter(),
				cff.NopEmitter(),
			},
		},
		{
			desc: "nested",
			give: []cff.Emitter{
				cff.NopEmitter(),
				cff.EmitterStack(
					cff.NopEmitter(),
					cff.NopEmitter(),
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			ctx := context.Background()

			e := cff.EmitterStack(tt.give...)
			e.FlowInit(&cff.FlowInfo{Name: "foo"}).
				FlowDone(ctx, time.Second)
			e.ParallelInit(&cff.ParallelInfo{Name: "bar"}).
				ParallelDone(ctx, time.Second)
		})
	}
}

type testStructs struct {
	ctrl      *gomock.Controller
	task1     *emittertest.MockTaskEmitter
	flow1     *emittertest.MockFlowEmitter
	emitter1  *emittertest.MockEmitter
	task2     *emittertest.MockTaskEmitter
	flow2     *emittertest.MockFlowEmitter
	emitter2  *emittertest.MockEmitter
	parallel1 *emittertest.MockParallelEmitter
	parallel2 *emittertest.MockParallelEmitter
	stack     cff.Emitter
}

func mocks(t *testing.T) testStructs {
	m := testStructs{}
	m.ctrl = gomock.NewController(t)
	m.task1 = emittertest.NewMockTaskEmitter(m.ctrl)
	m.flow1 = emittertest.NewMockFlowEmitter(m.ctrl)
	m.emitter1 = emittertest.NewMockEmitter(m.ctrl)
	m.task2 = emittertest.NewMockTaskEmitter(m.ctrl)
	m.flow2 = emittertest.NewMockFlowEmitter(m.ctrl)
	m.emitter2 = emittertest.NewMockEmitter(m.ctrl)
	m.parallel1 = emittertest.NewMockParallelEmitter(m.ctrl)
	m.parallel2 = emittertest.NewMockParallelEmitter(m.ctrl)
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
	})

	t.Run("Parallel", func(t *testing.T) {
		t.Run("Init", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).Times(1)
			m.emitter2.EXPECT().ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).Times(1)
			m.stack.ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0})
		})

		t.Run("ParallelSuccess", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).Return(m.parallel1)
			m.emitter2.EXPECT().ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).Return(m.parallel2)

			m.parallel1.EXPECT().ParallelSuccess(ctx)
			m.parallel2.EXPECT().ParallelSuccess(ctx)
			m.stack.ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).ParallelSuccess(ctx)
		})
		t.Run("ParallelError", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).Return(m.parallel1)
			m.emitter2.EXPECT().ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).Return(m.parallel2)

			err := errors.New("foobar")
			m.parallel1.EXPECT().ParallelError(ctx, err)
			m.parallel2.EXPECT().ParallelError(ctx, err)
			m.stack.ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).ParallelError(ctx, err)
		})
		t.Run("ParallelDone", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).Return(m.parallel1)
			m.emitter2.EXPECT().ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).Return(m.parallel2)

			m.parallel1.EXPECT().ParallelDone(ctx, time.Duration(1))
			m.parallel2.EXPECT().ParallelDone(ctx, time.Duration(1))
			m.stack.ParallelInit(&cff.ParallelInfo{"foo", "foo.go", 0, 0}).ParallelDone(ctx, time.Duration(1))
		})
	})

	t.Run("Task", func(t *testing.T) {
		t.Run("Init", func(t *testing.T) {
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16},
				&cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12},
			).Times(1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16},
				&cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12},
			).Times(1)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16},
				&cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12},
			)
		})

		t.Run("TaskSuccess", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task2)

			m.task1.EXPECT().TaskSuccess(ctx)
			m.task2.EXPECT().TaskSuccess(ctx)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).TaskSuccess(ctx)
		})
		t.Run("TaskError", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("foobar")

			m.task1.EXPECT().TaskError(ctx, err)
			m.task2.EXPECT().TaskError(ctx, err)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).TaskError(ctx, err)
		})
		t.Run("TaskErrorRecovered", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("great sadness")

			m.task1.EXPECT().TaskErrorRecovered(ctx, err)
			m.task2.EXPECT().TaskErrorRecovered(ctx, err)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).TaskErrorRecovered(ctx, err)
		})
		t.Run("TaskSkipped", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task2)

			err := errors.New("foobar")

			m.task1.EXPECT().TaskSkipped(ctx, err)
			m.task2.EXPECT().TaskSkipped(ctx, err)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).TaskSkipped(ctx, err)
		})
		t.Run("TaskPanic", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task2)

			pv := int(1)

			m.task1.EXPECT().TaskPanic(ctx, pv)
			m.task2.EXPECT().TaskPanic(ctx, pv)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).TaskPanic(ctx, pv)
		})
		t.Run("TaskPanicRecovered", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task2)

			pv := int(1)

			m.task1.EXPECT().TaskPanicRecovered(ctx, pv)
			m.task2.EXPECT().TaskPanicRecovered(ctx, pv)
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).TaskPanicRecovered(ctx, pv)
		})
		t.Run("TaskDone", func(t *testing.T) {
			ctx := context.Background()
			m := mocks(t)
			defer m.ctrl.Finish()

			m.emitter1.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task1)
			m.emitter2.EXPECT().TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).Return(m.task2)

			m.task1.EXPECT().TaskDone(ctx, time.Duration(1))
			m.task2.EXPECT().TaskDone(ctx, time.Duration(1))
			m.stack.TaskInit(
				&cff.TaskInfo{"foo", "foo.go", 14, 16}, &cff.DirectiveInfo{"fooFlow", cff.FlowDirective, "foo.go", 10, 12}).TaskDone(ctx, time.Duration(1))
		})
	})
}
