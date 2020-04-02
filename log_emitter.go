package cff

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
)

// logEmitter is an Emitter that writes to a Zap logger.
type logEmitter struct {
	logger *zap.Logger
}

// LogEmitter builds a CFF2 emitter which writes logs to the provided Zap
// logger.
func LogEmitter(log *zap.Logger) Emitter {
	return &logEmitter{logger: log}
}

type logFlowEmitter struct {
	// Field holding the flow name.
	flow zap.Field

	logger *zap.Logger
}

func (e *logEmitter) FlowInit(info *FlowInfo) FlowEmitter {
	return &logFlowEmitter{
		flow:   zap.String("flow", info.Flow),
		logger: e.logger,
	}
}

func (e *logFlowEmitter) FlowSuccess(context.Context) {
	e.logger.Debug("taskflow succeeded", e.flow)
}

func (e *logFlowEmitter) FlowError(ctx context.Context, err error) {
}

func (e *logFlowEmitter) FlowSkipped(ctx context.Context, err error) {
	e.logger.Debug("taskflow skipped", e.flow, zap.Error(err))
}

func (e *logFlowEmitter) FlowDone(ctx context.Context, d time.Duration) {
}

func (e *logFlowEmitter) FlowFailedTask(ctx context.Context, task string, err error) FlowEmitter {
	return e
}

type logTaskEmitter struct {
	// Fields holding the flow and task name.
	flow, task zap.Field

	logger *zap.Logger
}

func (e *logEmitter) TaskInit(task *TaskInfo, flow *FlowInfo) TaskEmitter {
	return &logTaskEmitter{
		flow:   zap.String("flow", flow.Flow),
		task:   zap.String("task", task.Task),
		logger: e.logger,
	}
}

func (e *logTaskEmitter) TaskSuccess(context.Context) {
	e.logger.Debug("task succeeded", e.flow, e.task)
}

func (e *logTaskEmitter) TaskError(ctx context.Context, err error) {
}

func (e *logTaskEmitter) TaskErrorRecovered(ctx context.Context, err error) {
	e.logger.Error("task error recovered", e.flow, e.task, zap.Error(err))
}

func (e *logTaskEmitter) TaskSkipped(ctx context.Context, err error) {
	e.logger.Debug("task skipped", e.flow, e.task, zap.Error(err))
}

func (e *logTaskEmitter) TaskPanic(ctx context.Context, x interface{}) {
	err, ok := x.(error)
	if !ok {
		err = fmt.Errorf("task panic: %v", x)
	}

	e.logger.Error("task panic", e.flow, e.task, zap.Stack("stack"), zap.Error(err))
}

func (e *logTaskEmitter) TaskPanicRecovered(ctx context.Context, x interface{}) {
	err, ok := x.(error)
	if !ok {
		err = fmt.Errorf("task panic: %v", x)
	}

	e.logger.Error("task panic recovered", e.flow, e.task, zap.Stack("stack"), zap.Error(err))
}

func (e *logTaskEmitter) TaskDone(ctx context.Context, _ time.Duration) {
}
