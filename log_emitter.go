package cff

import (
	"context"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LogEmitterOption customizes a LogEmitter.
type LogEmitterOption interface {
	applyLogEmitterOption(*logEmitter)
}

// LogErrors determines the log level for logging *unrecoverable* errors.
//
// Defaults to Debug because unrecoverable errors are surfaced back to caller
// of cff.Flow, which the caller may log at their chosen level.
func LogErrors(lvl zapcore.Level) LogEmitterOption {
	return logErrors(lvl)
}

type logErrors zapcore.Level

func (lvl logErrors) applyLogEmitterOption(e *logEmitter) {
	e.errLevel = zapcore.Level(lvl)
}

// LogPanics determines the log level for logging panics.
//
// Defaults to Error.
func LogPanics(lvl zapcore.Level) LogEmitterOption {
	return logPanics(lvl)
}

type logPanics zapcore.Level

func (lvl logPanics) applyLogEmitterOption(e *logEmitter) {
	e.panicLevel = zapcore.Level(lvl)
}

// LogRecovers determines the log level for logging recovered errors and
// panics.
//
// Defaults to Error.
func LogRecovers(lvl zapcore.Level) LogEmitterOption {
	return logRecovers(lvl)
}

type logRecovers zapcore.Level

func (lvl logRecovers) applyLogEmitterOption(e *logEmitter) {
	e.recoverLevel = zapcore.Level(lvl)
}

// logEmitter is an Emitter that writes to a Zap logger.
type logEmitter struct {
	logger *zap.Logger

	errLevel     zapcore.Level
	panicLevel   zapcore.Level
	recoverLevel zapcore.Level
}

// LogEmitter builds a CFF emitter which writes logs to the provided Zap
// logger.
func LogEmitter(log *zap.Logger, opts ...LogEmitterOption) Emitter {
	e := logEmitter{
		logger:       log,
		errLevel:     zapcore.DebugLevel,
		panicLevel:   zapcore.ErrorLevel,
		recoverLevel: zapcore.ErrorLevel,
	}
	for _, opt := range opts {
		opt.applyLogEmitterOption(&e)
	}
	return &e
}

type logFlowEmitter struct {
	// Field holding the flow name.
	flow zap.Field

	logger   *zap.Logger
	errLevel zapcore.Level
}

func (logFlowEmitter) flowEmitter() {}

func (e *logEmitter) FlowInit(info *FlowInfo) FlowEmitter {
	return &logFlowEmitter{
		flow:     zap.String("flow", info.Name),
		logger:   e.logger,
		errLevel: e.errLevel,
	}
}

func (e *logFlowEmitter) FlowSuccess(context.Context) {
	e.logger.Debug("flow success", e.flow)
}

func (e *logFlowEmitter) FlowError(ctx context.Context, err error) {
	if ce := e.logger.Check(e.errLevel, "flow error"); ce != nil {
		ce.Write(e.flow, zap.Error(err))
	}
}

func (e *logFlowEmitter) FlowDone(ctx context.Context, d time.Duration) {
	e.logger.Debug("flow done", e.flow)
}

type logParallelEmitter struct {
	// Field holding the parallel name.
	parallel zap.Field

	logger   *zap.Logger
	errLevel zapcore.Level
}

func (logParallelEmitter) parallelEmitter() {}

func (e *logEmitter) ParallelInit(info *ParallelInfo) ParallelEmitter {
	return &logParallelEmitter{
		parallel: zap.String("parallel", info.Name),
		logger:   e.logger,
		errLevel: e.errLevel,
	}
}

func (e *logParallelEmitter) ParallelSuccess(context.Context) {
	e.logger.Debug("parallel success", e.parallel)
}

func (e *logParallelEmitter) ParallelError(ctx context.Context, err error) {
	if ce := e.logger.Check(e.errLevel, "parallel error"); ce != nil {
		ce.Write(e.parallel, zap.Error(err))
	}
}

func (e *logParallelEmitter) ParallelDone(ctx context.Context, d time.Duration) {
	e.logger.Debug("parallel done", e.parallel)
}

type logTaskEmitter struct {
	// Fields holding the directive name.
	directive, task zap.Field

	logger       *zap.Logger
	errLevel     zapcore.Level
	panicLevel   zapcore.Level
	recoverLevel zapcore.Level
}

func (logTaskEmitter) taskEmitter() {}

func (e *logEmitter) TaskInit(task *TaskInfo, d *DirectiveInfo) TaskEmitter {
	return &logTaskEmitter{
		directive:    zap.String(d.Directive.String(), d.Name),
		task:         zap.String("task", task.Name),
		logger:       e.logger,
		errLevel:     e.errLevel,
		panicLevel:   e.panicLevel,
		recoverLevel: e.recoverLevel,
	}
}

func (e *logTaskEmitter) TaskSuccess(context.Context) {
	e.logger.Debug("task success", e.directive, e.task)
}

func (e *logTaskEmitter) TaskError(ctx context.Context, err error) {
	if ce := e.logger.Check(e.errLevel, "task error"); ce != nil {
		ce.Write(e.directive, e.task, zap.Error(err))
	}
}

func (e *logTaskEmitter) TaskErrorRecovered(ctx context.Context, err error) {
	if ce := e.logger.Check(e.recoverLevel, "task error recovered"); ce != nil {
		ce.Write(e.directive, e.task, zap.Error(err))
	}
}

func (e *logTaskEmitter) TaskSkipped(ctx context.Context, err error) {
	e.logger.Debug("task skipped", e.directive, e.task, zap.Error(err))
}

func (e *logTaskEmitter) TaskPanic(ctx context.Context, pv interface{}) {
	if ce := e.logger.Check(e.panicLevel, "task panic"); ce != nil {
		ce.Write(
			e.directive,
			e.task,
			zap.Stack("stack"),
			zap.Any("panic-value", pv),
			maybeErrorField(pv),
		)
	}
}

func (e *logTaskEmitter) TaskPanicRecovered(ctx context.Context, pv interface{}) {
	if ce := e.logger.Check(e.recoverLevel, "task panic recovered"); ce != nil {
		ce.Write(
			e.directive,
			e.task,
			zap.Stack("stack"),
			zap.Any("panic-value", pv),
			maybeErrorField(pv),
		)
	}
}

func (e *logTaskEmitter) TaskDone(ctx context.Context, _ time.Duration) {
	e.logger.Debug("task done", e.directive, e.task)
}

// SchedulerInit constructs a logging scheduler emitter.
func (e *logEmitter) SchedulerInit(info *SchedulerInfo) SchedulerEmitter {
	directive := zap.Skip()
	if info.Name != "" && info.Directive != UnknownDirective {
		directive = zap.String(info.Directive.String(), info.Name)
	}
	return &logSchedulerEmitter{
		directive: directive,
		logger:    e.logger,
	}
}

type logSchedulerEmitter struct {
	directive zap.Field
	logger    *zap.Logger
}

func (e *logSchedulerEmitter) EmitScheduler(s SchedulerState) {
	e.logger.Debug(
		"scheduler state",
		e.directive,
		zap.Int("pending", s.Pending),
		zap.Int("ready", s.Ready),
		zap.Int("waiting", s.Waiting),
		zap.Int("idle_workers", s.IdleWorkers),
		zap.Int("concurrency", s.Concurrency),
	)
}

func (logSchedulerEmitter) schedulerEmitter() {}

func maybeErrorField(pv interface{}) zap.Field {
	if err, ok := pv.(error); ok {
		return zap.Error(err)
	}
	return zap.Skip()
}
