package internal

import (
	"errors"
	"go/ast"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

type parallel struct {
	ast.Node

	Ctx         ast.Expr // initial ctx argument to cff.Parallel(...)
	Concurrency ast.Expr // argument to cff.Concurrency, if any.

	Tasks []*parallelTask

	PosInfo *PosInfo // Used to pass information to uniquely identify a task.
}

type parallelTask struct {
	Function *function
	// Serial is a unique serially incrementing number for each task.
	Serial int

	PosInfo *PosInfo // Used to pass information to uniquely identify a task.
}

func (c *compiler) compileParallel(file *ast.File, call *ast.CallExpr) *parallel {
	if len(call.Args) == 1 {
		c.errf(c.nodePosition(call), "cff.Parallel expects at least one function")
		return nil
	}

	parallel := &parallel{
		Ctx:     call.Args[0],
		Node:    call,
		PosInfo: c.getPosInfo(call),
	}
	for _, arg := range call.Args[1:] {
		arg := astutil.Unparen(arg)

		ce, ok := arg.(*ast.CallExpr)
		if !ok {
			c.errf(c.nodePosition(arg), "expected a function call, got %v", astutil.NodeDescription(arg))
			continue
		}

		f := typeutil.StaticCallee(c.info, ce)
		if f == nil || !isPackagePathEquivalent(f.Pkg(), cffImportPath) {
			c.errf(c.nodePosition(arg), "expected cff call but got %v", typeutil.Callee(c.info, ce))
			continue
		}

		switch f.Name() {
		case "Tasks":
			parallel.Tasks = append(parallel.Tasks, c.compileParallelTasks(parallel, ce)...)
		case "Concurrency":
			parallel.Concurrency = ce.Args[0]
		}
		// TODO(GO-84): WithEmitter, Instrument, Task, ContinueOnError,
		// Map, Slice.
	}
	return parallel
}

func (c *compiler) compileParallelTasks(p *parallel, call *ast.CallExpr) []*parallelTask {
	var tasks []*parallelTask
	for _, arg := range call.Args {
		t := c.compileParallelTask(p, arg)
		if t != nil {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func (c *compiler) compileParallelTask(p *parallel, arg ast.Expr) *parallelTask {
	taskF := c.compileFunction(arg)
	if taskF == nil {
		c.errf(c.nodePosition(arg), "parallel tasks function failed to compile")
		return nil
	}
	if err := checkParallelTask(taskF); err != nil {
		c.errf(c.nodePosition(arg), "parallel tasks function is invalid: %v", err)
		return nil
	}
	fn := &function{
		Node:     taskF.Node,
		Sig:      taskF.Sig,
		WantCtx:  taskF.WantCtx,
		HasError: taskF.HasError,
		PosInfo:  taskF.PosInfo,
	}
	t := &parallelTask{
		Function: fn,
		Serial:   c.taskSerial,
		PosInfo:  taskF.PosInfo,
	}
	c.taskSerial++
	return t
}

func checkParallelTask(fn *compiledFunc) error {
	switch {
	case len(fn.Inputs) != 0:
		return errors.New("the only allowed argument is a single context parameter")
	case len(fn.Outputs) != 0:
		return errors.New("the only allowed return value is a single error")
	default:
		return nil
	}
}
