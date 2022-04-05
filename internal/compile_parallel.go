package internal

import (
	"errors"
	"go/ast"
	"go/types"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

type parallel struct {
	ast.Node

	Ctx         ast.Expr // initial ctx argument to cff.Parallel(...)
	Concurrency ast.Expr // argument to cff.Concurrency, if any.

	ContinueOnError ast.Expr // argument to cff.ContinueOnError.

	Emitters []ast.Expr // zero or more expressions of the type cff.Emitter.

	Tasks []*parallelTask

	SliceTasks []*sliceTask

	MapTasks []*mapTask

	Instrument *instrument

	PosInfo *PosInfo // Used to pass information to uniquely identify a task.
}

type parallelTask struct {
	Function *function
	// Serial is a unique serially incrementing number for each task.
	Serial int

	Instrument *instrument

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
		case "InstrumentFlow":
			c.errf(c.nodePosition(arg), "%q is an invalid cff.Parallel Option", f.Name())
			continue
		case "Task":
			if t := c.compileParallelTask(parallel, ce.Args[0], ce.Args[1:]); t != nil {
				parallel.Tasks = append(parallel.Tasks, t)
			}
		case "Tasks":
			parallel.Tasks = append(parallel.Tasks, c.compileParallelTasks(parallel, ce)...)
		case "Concurrency":
			parallel.Concurrency = ce.Args[0]
		case "ContinueOnError":
			parallel.ContinueOnError = ce.Args[0]
		case "InstrumentParallel":
			parallel.Instrument = c.compileInstrument(ce)
		case "Slice":
			if st := c.compileSlice(parallel, ce); st != nil {
				parallel.SliceTasks = append(parallel.SliceTasks, st)
			}
		case "WithEmitter":
			parallel.Emitters = append(parallel.Emitters, ce.Args[0])
		case "Map":
			if mt := c.compileMap(ce); mt != nil {
				parallel.MapTasks = append(parallel.MapTasks, mt)
			}
		}
	}
	c.validateParallelInstrument(parallel)

	return parallel
}

func (c *compiler) validateParallelInstrument(p *parallel) {
	// If the directive, or any task in the directive were instrumented, we require
	// at least one emitter to be provided.
	if len(p.Emitters) > 0 {
		return
	}

	if p.Instrument != nil {
		c.errf(c.nodePosition(p.Node), "cff.InstrumentParallel requires a cff.Emitter to be provided: use cff.WithEmitter")
	}

	for _, t := range p.Tasks {
		if t.Instrument != nil {
			c.errf(c.nodePosition(p.Node), "cff.Instrument requires a cff.Emitter to be provided: use cff.WithEmitter")
		}
	}

	for _, s := range p.SliceTasks {
		if s.SliceEndFn != nil && p.ContinueOnError != nil {
			c.errf(c.nodePosition(p.Node), `"cff.SliceEnd" is an invalid option when "ContinueOnError" is used`)
		}
	}
}

func (c *compiler) compileParallelTask(p *parallel, call ast.Expr, opts []ast.Expr) *parallelTask {
	t := c.compileParallelTaskFn(p, call)
	if t == nil {
		c.errf(c.nodePosition(call), "parallel task failed to compile")
		return nil
	}
	for _, opt := range opts {
		call, fn, err := c.identifyOption(opt)
		if err != nil {
			c.errf(c.nodePosition(opt), err.Error())
			continue
		}
		switch fn.Name() {
		case "Instrument":
			t.Instrument = c.compileInstrument(call)
		}
	}
	return t
}

func (c *compiler) compileParallelTasks(p *parallel, call *ast.CallExpr) []*parallelTask {
	var tasks []*parallelTask
	for _, arg := range call.Args {
		t := c.compileParallelTaskFn(p, arg)
		if t != nil {
			tasks = append(tasks, t)
		}
	}
	return tasks
}

func (c *compiler) compileParallelTaskFn(p *parallel, arg ast.Expr) *parallelTask {
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
		return errors.New("the only allowed argument is a single context.Context parameter")
	case len(fn.Outputs) != 0:
		return errors.New("the only allowed return value is an error")
	default:
		return nil
	}
}

type sliceTask struct {
	Function   *compiledFunc
	Slice      ast.Expr
	ElemType   types.Type
	SliceEndFn *compiledFunc

	// Serial is a unique serially incrementing number for each sliceTask.
	Serial int

	PosInfo *PosInfo // Used to pass information to uniquely identify a task.
}

func (c *compiler) applySliceOptions(t *sliceTask, opts []ast.Expr) {
	for _, opt := range opts {
		ce, fn, err := c.identifyOption(opt)
		if err != nil {
			c.errf(c.nodePosition(opt), err.Error())
			continue
		}

		switch fn.Name() {
		case "SliceEnd":
			sliceEndFn := c.compileSliceEnd(opt, ce)
			if sliceEndFn == nil {
				continue
			}
			if t.SliceEndFn != nil {
				c.errf(c.nodePosition(opt), "cff.Slice accepts at most one cff.SliceEnd option")
				continue
			}
			t.SliceEndFn = sliceEndFn
		}
	}
	return
}

func (c *compiler) compileSliceEnd(opt ast.Expr, ce *ast.CallExpr) *compiledFunc {
	fn := c.compileFunction(ce.Args[0])
	switch {
	case fn == nil:
		c.errf(c.nodePosition(opt), "SliceEnd function failed to compile")
		return nil
	case len(fn.Inputs) != 0:
		c.errf(c.nodePosition(opt), "the only allowed argument is a single context.Context parameter")
		return nil
	case len(fn.Outputs) != 0:
		c.errf(c.nodePosition(opt), "the only allowed return value is an error")
		return nil
	default:
		return fn
	}
}

func (c *compiler) compileSlice(p *parallel, ce *ast.CallExpr) *sliceTask {
	sliceFn, slce := ce.Args[0], ce.Args[1]
	fn := c.compileFunction(sliceFn)
	if fn == nil {
		c.errf(c.nodePosition(sliceFn), "slice function failed to compile")
		return nil
	}

	if len(fn.Outputs) != 0 {
		c.errf(c.nodePosition(sliceFn), "the only allowed return value is an error")
		return nil
	}

	if len(fn.Inputs) != 2 {
		c.errf(c.nodePosition(slce), "slice function expects two non-context arguments: slice index and slice element")
		return nil
	}

	if t, ok := fn.Inputs[0].(*types.Basic); !ok || t.Kind() != types.Int {
		c.errf(c.nodePosition(slce), "the first non-context argument of the slice function must be an int, got %v", fn.Inputs[0])
		return nil
	}

	typ := c.info.TypeOf(slce)
	if typ == nil {
		c.errf(c.nodePosition(slce), "type of the slice argument is not found")
		return nil
	}

	slc, ok := typ.(*types.Slice)
	if !ok {
		c.errf(c.nodePosition(slce), "the second argument to cff.Slice must be a slice, got %v", typ)
		return nil
	}

	if !types.AssignableTo(fn.Inputs[1], slc.Elem()) {
		c.errf(c.nodePosition(slce), "slice element of type %v cannot be passed as a parameter to function expecting %v", fn.Inputs[1], slc.Elem())
		return nil
	}

	s := &sliceTask{
		Function: fn,
		Slice:    slce,
		ElemType: slc.Elem(),
		Serial:   c.taskSerial,
		PosInfo:  c.getPosInfo(ce),
	}
	c.taskSerial++
	c.applySliceOptions(s, ce.Args[2:])
	return s
}

type mapTask struct {
	Function *compiledFunc
	Map      ast.Expr
	KeyType  types.Type
	ElemType types.Type

	// Serial is a unique serially incrementing number for each mapTask.
	Serial int

	PosInfo *PosInfo // Used to pass information to uniquely identify a task.
}

func (c *compiler) compileMap(ce *ast.CallExpr) *mapTask {
	mapFun, mmap := ce.Args[0], ce.Args[1]
	fn := c.compileFunction(mapFun)
	if fn == nil {
		c.errf(c.nodePosition(mapFun), "map function failed to compile")
		return nil
	}

	if len(fn.Outputs) != 0 {
		c.errf(c.nodePosition(mapFun), "the only allowed return value is an error")
		return nil
	}

	if len(fn.Inputs) != 2 {
		c.errf(c.nodePosition(mmap), "map function expects two non-context arguments: key and value elements from a map")
		return nil
	}

	typ := c.info.TypeOf(mmap)
	mtype, ok := typ.(*types.Map)
	if !ok {
		c.errf(c.nodePosition(mmap), "the second argument to cff.Map must be a map, got %v", typ)
		return nil
	}

	if !types.AssignableTo(fn.Inputs[0], mtype.Key()) {
		c.errf(c.nodePosition(mmap), "key element of type %v cannot be passed as a parameter to function expecting %v", mtype.Key(), fn.Inputs[0])
		return nil
	}

	if !types.AssignableTo(fn.Inputs[1], mtype.Elem()) {
		c.errf(c.nodePosition(mmap), "value element of type %v cannot be passed as a parameter to function expecting %v", mtype.Elem(), fn.Inputs[1])
		return nil
	}

	m := &mapTask{
		Function: fn,
		Map:      mmap,
		KeyType:  mtype.Key(),
		ElemType: mtype.Elem(),
		Serial:   c.taskSerial,
		PosInfo:  c.getPosInfo(ce),
	}
	c.taskSerial++
	return m
}
