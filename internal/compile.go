package internal

import (
	"container/list"
	"errors"
	"fmt"
	"go/ast"
	"go/constant"
	"go/token"
	"go/types"
	"path/filepath"
	"strconv"
	"strings"

	"code.uber.internal/go/importer"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

// funcIndex marks cff.Results outputs.
type funcIndex int

const (
	cffImportPath   = "go.uber.org/cff"
	funcIndexResult = funcIndex(-1)
)

type compiler struct {
	pkg        *types.Package
	fset       *token.FileSet
	info       *types.Info
	taskSerial int
	errors     []error

	instrumentAllTasks bool
}

type compilerOpts struct {
	Fset               *token.FileSet
	Info               *types.Info
	Package            *types.Package
	InstrumentAllTasks bool
}

func newCompiler(opts compilerOpts) *compiler {
	return &compiler{
		fset:               opts.Fset,
		info:               opts.Info,
		pkg:                opts.Package,
		instrumentAllTasks: opts.InstrumentAllTasks,
	}
}

func (c *compiler) errf(pos token.Position, msg string, args ...interface{}) {
	formattedMsg := fmt.Sprintf("%v: ", pos) + fmt.Sprintf(msg, args...)

	c.errors = append(c.errors, errors.New(formattedMsg))
}

func (c *compiler) position(pos token.Pos) token.Position {
	return c.fset.Position(pos)
}

func (c *compiler) nodePosition(n ast.Node) token.Position {
	return c.position(n.Pos())
}

type file struct {
	AST     *ast.File
	Package *importer.Package

	// Map from import path to local names of the import. If the import is
	// unnamed, it will be recorded as the package name.
	Imports map[string][]string
	// Slice because you can have the same import path multiple times with
	// different names.

	// Packages that were imported unnamed.
	UnnamedImports map[string]struct{}

	Filepath   string
	Flows      []*flow
	Parallels  []*parallel
	Generators []directiveGenerator
}

func (c *compiler) CompileFile(file *ast.File, pkg *importer.Package) (*file, error) {
	f := c.compileFile(file, pkg)
	return f, multierr.Combine(c.errors...)
}

func (c *compiler) compileFile(astFile *ast.File, pkg *importer.Package) *file {
	file := file{
		AST:            astFile,
		Package:        pkg,
		Filepath:       c.fset.File(astFile.Pos()).Name(),
		Imports:        make(map[string][]string),
		UnnamedImports: make(map[string]struct{}),
	}

	astWalk(astFile, func(n ast.Node) bool {
		switch n := n.(type) {
		case *ast.ImportSpec:
			// If the user defines a name for an import, we would like to track their name if we use it in the generated
			// code
			importPath, _ := strconv.Unquote(n.Path.Value)
			var (
				importName string
			)
			if n.Name != nil {
				importName = n.Name.String()
			} else {
				// Not using a named import. Ask go/types.Info for the
				// implicit name.
				obj := c.info.Implicits[n]

				p := obj.(*types.PkgName)
				importName = p.Name()
			}
			file.Imports[importPath] = append(file.Imports[importPath], importName)
			file.UnnamedImports[importPath] = struct{}{}
			return false
		case *ast.CallExpr:
			// We're looking for a call in the form "cff.Flow" or
			// "cff.Parallel". It will be a SelectorExpr where the "X" is a
			// reference to the "cff" package.
			sel, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return true // keep looking
			}

			fn, ok := c.info.Uses[sel.Sel]
			if !ok {
				return true // invalid code b/c identifier doesn't exist. keep looking
			}

			if !isPackagePathEquivalent(fn.Pkg(), cffImportPath) {
				return true // keep looking
			}

			// Inside a +cff file, code generation directives
			// should only appear inside a cff.Flow or cff.Parallel.
			//
			// Code generation directives by themselves at the
			// top-level are not allowed.

			switch {
			case fn.Name() == "Flow":
				flow := c.compileFlow(astFile, n)
				file.Flows = append(file.Flows, flow)
				file.Generators = append(
					file.Generators,
					flowGenerator{
						flow: flow,
					},
				)

			case fn.Name() == "Parallel":
				parallel := c.compileParallel(astFile, n)
				file.Parallels = append(file.Parallels, parallel)
				file.Generators = append(
					file.Generators,
					parallelGenerator{
						parallel: parallel,
					},
				)

			case IsCodegenDirective(fn.Name()):
				c.errf(c.nodePosition(n), "unexpected code generation directive %q: "+"only cff.Flow or cff.Parallel may be called at the top-level", fn.Name())
			default:
				// Calls to functions that are not code
				// generation directives (LogEmitter, etc.)
				// are allowed.
			}

			return false
		default:
			return true // keep looking
		}
	})

	return &file
}

type flow struct {
	ast.Node

	Ctx         ast.Expr // initial ctx argument to cff.Flow(...)
	Concurrency ast.Expr // argument to cff.Concurrency, if any.

	Emitters []ast.Expr // zero or more expressions of the type cff.Emitter.

	Inputs  []*input
	Outputs []*output
	Tasks   []*task

	Predicates []*predicate
	Funcs      []*function

	// Topologically ordered list of functions.
	// For all i < j, TopoFuncs[i] cannot depend on TopoFuncs[j].
	TopoFuncs []*function

	Instrument *instrument

	providers *typeutil.Map // map[types.Type]int (index in Tasks)
	receivers *typeutil.Map // map[types.Type][]funcIndex tracks types needed to detect unused inputs

	invokeTypeCnt int       // input to make unique invokeType sentinels
	invokeTypes   []*output // tracks tasks of with either no results or a single error

	predicateTypeCnt int                // input to make unique predicateType sentinels.
	predicateTypes   []*predicateOutput // tracks cff.Predicate sentinel types.

	PosInfo *PosInfo // Used to pass information to uniquely identify a task.
}

// addNoOutput adds a unique invokeType sentinel type to the invokeTypes list.
// The breadth-first searching algorithm visits based on result types, but
// functions with no return values would not be visited since no function can
// depend on its outputs and therefore would not be visited or included in the graph.
// addNoOutput creates a unique sentinel type so that we can pretend that this
// function is needed to provide the sentinel type for scheduling purposes.
func (f *flow) addNoOutput() *noOutput {
	f.invokeTypeCnt++
	name := strconv.Itoa(f.invokeTypeCnt)
	field := types.NewVar(0, nil, name, &types.Struct{})
	no := types.NewStruct([]*types.Var{field}, nil)
	f.invokeTypes = append(f.invokeTypes, &output{Type: no})
	return no
}

// mustSetNoOutputProvider sets the provider for the no-output, panicking if the no-output sentinel type was already
// present.
func (f *flow) mustSetNoOutputProvider(key *function, value int) {
	prev := f.providers.Set(key.Task.invokeType, value)
	if prev != nil {
		panic(fmt.Sprintf("cff assertion error: invokeType sentinel types should be unique, found %T for %dth task (defined at %v), expected to be nil", prev, value, key.Node))
	}

	if typ, ok := f.receivers.At(key.Task.invokeType).([]funcIndex); ok {
		typ := append(typ, funcIndex(value))
		f.receivers.Set(key.Task.invokeType, typ)
	} else {
		f.receivers.Set(key.Task.invokeType, []funcIndex{funcIndex(value)})
	}
}

// addPreciateOutput creates a unique predicate sentinel type. Since CFF's
// dependency resolution logic only allows for one of each type to be provided
// to the dependency graph and cff.Predicate functions return booleans,
// these sentinel types distinguish the outputs of cff.Predicates.
func (f *flow) addPredicateOutput() *predicateOutput {
	f.predicateTypeCnt++
	name := strconv.Itoa(f.predicateTypeCnt)
	field := types.NewVar(0, nil, name, &types.Basic{})
	np := types.NewStruct([]*types.Var{field}, nil)
	f.predicateTypes = append(f.predicateTypes, np)
	return np
}

func (c *compiler) compileFlow(file *ast.File, call *ast.CallExpr) *flow {
	if len(call.Args) == 1 {
		c.errf(c.nodePosition(call), "cff.Flow expects at least one function")
		return nil
	}

	flow := flow{
		Ctx:       call.Args[0],
		Node:      call,
		PosInfo:   c.getPosInfo(call),
		providers: new(typeutil.Map),
		receivers: new(typeutil.Map),
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
		case "ContinueOnError", "Slice", "Map":
			c.errf(c.nodePosition(arg), "%q is an invalid cff.Flow Option", f.Name())
			continue
		case "Params":
			provided := new(typeutil.Map) // *type => *input
			for _, i := range ce.Args {
				in := c.compileInput(i)
				if other, _ := provided.At(in.Type).(*input); other != nil {
					c.errf(c.nodePosition(i), "type %v already provided to cff.Params at %v", other.Type, c.nodePosition(other.Node))
					continue
				}
				flow.Inputs = append(flow.Inputs, in)
				provided.Set(in.Type, in)
			}
		case "Results":
			for _, o := range ce.Args {
				if output := c.compileOutput(o); output != nil {
					flow.Outputs = append(flow.Outputs, output)
					// receivers is used to look up Task for compilation checks. Since we have
					// Results, we dont need to find an associated task.
					// We don't care about other values, cff.Results should be the only receiver.
					flow.receivers.Set(output.Type, []funcIndex{funcIndexResult})
				}
			}
		case "InstrumentFlow":
			flow.Instrument = c.compileInstrument(ce)
		case "WithEmitter":
			flow.Emitters = append(flow.Emitters, ce.Args[0])
		case "Concurrency":
			flow.Concurrency = ce.Args[0]
		case "Task":
			if task := c.compileTask(&flow, ce.Args[0], ce.Args[1:]); task != nil {
				flow.Tasks = append(flow.Tasks, task)
				flow.Funcs = append(flow.Funcs, task.Function)
				if task.Predicate != nil {
					flow.Funcs = append(flow.Funcs, task.Predicate.Function)
					flow.Predicates = append(flow.Predicates, task.Predicate)
				}
			}
		}
	}
	// At this point, c.errors may be non-empty but we are continuing with more checks to catch all
	// possible errors prior to scheduling attempt and return them at once.
	c.validateInstrument(&flow)

	for i, fn := range flow.Funcs {
		for _, in := range fn.Dependencies {
			if typ, ok := flow.receivers.At(in).([]funcIndex); ok {
				typ := append(typ, funcIndex(i))
				flow.receivers.Set(in, typ)
			} else {
				flow.receivers.Set(in, []funcIndex{funcIndex(i)})
			}
		}
		for _, o := range fn.outputs() {
			prev := flow.providers.Set(o, i)
			if prev != nil {
				pIdx := prev.(int)
				p := flow.Funcs[pIdx]
				c.errf(c.nodePosition(fn), "type %v already provided at %v", o, c.nodePosition(p))
				continue
			}
		}
		if fn.Task != nil && fn.Task.invokeType != nil {
			flow.mustSetNoOutputProvider(fn, i)
		}
		// Unlike invokeTypes, predicate sentinel types are already registered
		// as part of a function's Dependencies.
	}

	c.validateNoUnusedOutputTypes(&flow)
	c.validateFuncs(&flow)
	// At this point we may have already found some errors in c.errors.
	if err := validateFlowCycles(&flow, c.fset); err != nil {
		c.errors = append(c.errors, err)
		return nil
	}
	if len(c.errors) > 0 {
		return nil
	}

	c.scheduleFlowAndToposort(&flow)

	return &flow
}

type validateVisitedType struct {
	Type types.Type

	// Node is the place in the code (either a task or a flow output) that we needed the type
	Node ast.Node
}

// validateNoUnusedOutputTypes ensures that every output type is consumed by either a cff.Results or another task
// or a predicate of another task.
func (c *compiler) validateNoUnusedOutputTypes(f *flow) {
	for _, t := range f.Funcs {
		for _, o := range t.outputs() {
			if f.receivers.At(o) == nil {
				c.errf(c.nodePosition(t.Node), "unused output type %v", o)
			}
		}
	}
}

// validateFuncs walks the graph from the bottom of the graph (the outputs) to validate that
// all outputs are provided by some function. we also walk up the graph in case cff.Results
// is not the root, and check if there are any tasks with output past cff.Results.
func (c *compiler) validateFuncs(f *flow) {
	var (
		queue      = list.New() // []validateVisitedType
		visited    typeutil.Map // map[types.Type]struct{}
		flowInputs typeutil.Map // map[types.Type]*input
	)

	for _, i := range f.Inputs {
		flowInputs.Set(i.Type, i)
	}

	for _, o := range f.Outputs {
		queue.PushBack(validateVisitedType{Type: o.Type, Node: o.Node})
	}
	for _, o := range f.invokeTypes {
		// We do not need to walk forward for invokeType tasks since they aren't expected to return
		// non-error results. The case when they don't return anything and don't use cff.Invoke will
		// be handled after interpreting task options.
		queue.PushBack(validateVisitedType{Type: o.Type, Node: o.Node})
	}

	// A list of predicateOutput types are not pushed into the queue as they
	// are an internal book-keeping type not declared as an input or output
	// by CFF2's public APIs.

	for queue.Len() > 0 {
		t := queue.Remove(queue.Front()).(validateVisitedType)

		if visited.At(t.Type) != nil {
			// Two tasks can depend on the same input, and that is OK, but
			// we cannot allow cycles. Skip processing of a task that has
			// already been processed and handle cycle detection at a separate stage.
			continue
		}
		visited.Set(t.Type, struct{}{})

		if funcIdx, ok := f.providers.At(t.Type).(int); ok {
			fn := f.Funcs[funcIdx]
			for _, i := range fn.Dependencies {
				queue.PushBack(validateVisitedType{Type: i, Node: fn.Node})
			}

			continue
		}

		if flowInputs.Delete(t.Type) {
			continue
		}

		c.errf(c.nodePosition(t.Node), "no provider found for %v", t.Type)
	}

	if flowInputs.Len() > 0 {
		inputs := flowInputs.Keys()
		for _, inputType := range inputs {
			inputUntyped := flowInputs.At(inputType)
			input := inputUntyped.(*input)
			c.errf(c.nodePosition(input.Node), "unused input type %v", input.Type)
		}
	}
}

func (c *compiler) validateInstrument(f *flow) {
	instrumented := f.Instrument != nil
	if !instrumented {
		for _, t := range f.Tasks {
			if t.Instrument != nil {
				instrumented = true
				break
			}
		}
	}

	// If the flow, or any task in the flow were instrumented, we require
	// at least one emitter to be provided.
	if !instrumented {
		return
	}

	if len(f.Emitters) == 0 {
		c.errf(c.nodePosition(f.Node), "cff.Instrument requires a cff.Emitter to be provided: use cff.WithEmitter")
	}
}

func (c *compiler) scheduleFlowAndToposort(f *flow) {
	g := graph{
		Count: len(f.Funcs),
		Dependencies: func(funcIdx int) []int {
			var deps []int
			for _, typ := range f.Funcs[funcIdx].Dependencies {
				if i, ok := f.providers.At(typ).(int); ok {
					// For non-ok case, if we do not find a dependency amongst providers, then it
					// was passed in from Params annotation.
					deps = append(deps, i)
				}
			}
			return deps
		},
	}

	for idx, fn := range f.Funcs {
		for _, depIdx := range g.Dependencies(idx) {
			fn.DependsOn = append(fn.DependsOn, f.Funcs[depIdx])
		}
	}

	var topo []*function
	for _, idx := range toposort(g) {
		topo = append(topo, f.Funcs[idx])
	}

	f.TopoFuncs = topo
}

// PosInfo contains positional information about a Flow or Task. This may be
// used to uniquely identify a flow or a task.
type PosInfo struct {
	File         string // the file in which it's defined
	Line, Column int    // line and column in the file where the flow was defined
}

type task struct {
	ast.Node

	// Function is the object of a task's execution, a task must have
	// a function.
	Function *function

	// Serial is a unique serially incrementing number for each task.
	Serial int

	Inputs  []types.Type // non ctx params
	Outputs []types.Type // non error results

	// A task has at most one predicate.
	Predicate  *predicate  // non-nil if Predicate was provided
	Instrument *instrument // non-nil if instrumentation was enabled

	FallbackWith        bool       // whether we should ignore errors from this function
	FallbackWithResults []ast.Expr // expressions that return a value for each return type of this function

	invokeType *noOutput // non-nil if there are no non-error results

	PosInfo *PosInfo // Used to pass information to uniquely identify a task.
}

// invokeType is a sentinel return type for tasks that have no non-error results.
// It can not be custom defined type, otherwise it won't work with typeutil.Map.
type noOutput = types.Struct

// predicateOutput is a sentinel return type for cff.Predicates that return
// a boolean results that would otherwise conflict other cff.Predicate return
// values in the dependency graph.
type predicateOutput = types.Struct

func (c *compiler) compileTask(flow *flow, expr ast.Expr, opts []ast.Expr) *task {
	compiledFunc := c.compileFunction(expr)
	if compiledFunc == nil {
		return nil
	}

	taskFunc := &function{
		Node:         compiledFunc.Node,
		Sig:          compiledFunc.Sig,
		WantCtx:      compiledFunc.WantCtx,
		HasError:     compiledFunc.HasError,
		Dependencies: compiledFunc.Inputs,
		PosInfo:      compiledFunc.PosInfo,
	}

	t := task{
		Node:     expr,
		Function: taskFunc,
		Serial:   c.taskSerial,
		Inputs:   compiledFunc.Inputs,
		Outputs:  compiledFunc.Outputs,
		PosInfo:  c.getPosInfo(expr),
	}

	taskFunc.Task = &t
	c.taskSerial++

	c.interpretTaskOptions(flow, &t, opts)
	if t.Predicate != nil {
		t.Function.Dependencies = append(t.Function.Dependencies, t.Predicate.SentinelOutput)
	}

	// Check if we return nothing and we don't have an Invoke call.
	if len(t.Outputs) == 0 && t.invokeType == nil {
		c.errf(c.nodePosition(expr), "task must return at least one non-error value but currently produces zero."+" Did you intend to use cff.Invoke(true)?")
	}
	if len(t.Outputs) > 0 && t.invokeType != nil {
		c.errf(c.nodePosition(expr), "cff.Invoke cannot be provided on a Task that produces values besides errors")
	}

	// Create an implied Instrument(...) annotation for all tasks if the
	// flow is instrumented and the --instrument-all-tasks flag was
	// passed.
	if flow.Instrument != nil && c.instrumentAllTasks && t.Instrument == nil {
		taskPos := c.nodePosition(t)
		name := fmt.Sprintf("%s.%d", filepath.Base(taskPos.Filename), taskPos.Line)
		t.Instrument = c.compileInstrumentName(name)
	}

	return &t
}

// function is the smallest unit of execution. Higher level functionality like
// cff.Tasks, cff.Predicates, and cff.Invokes are mapped to one or more
// functions.
type function struct {
	ast.Node

	Sig *types.Signature

	// Whether the first argument to this task is a context.Context.
	WantCtx bool

	// Whether the last result is an error.
	HasError bool

	// Dependencies are types.Type dependencies of this function.
	Dependencies []types.Type

	// DependsOn are function dependencies of this function.
	DependsOn []*function

	Task      *task      // non-nil if function executes a task
	Predicate *predicate // non-nil if function executes a predicate

	PosInfo *PosInfo // Used to pass information to uniquely identify a function.
}

// Inputs returns the types consumed by this function.
func (f *function) inputs() []types.Type {
	if f.Predicate != nil {
		return f.Predicate.Inputs
	}
	return f.Task.Inputs
}

// Outputs returns the types produced by this function.
func (f *function) outputs() []types.Type {
	if f.Predicate != nil {
		return []types.Type{f.Predicate.SentinelOutput}
	}
	return f.Task.Outputs
}

// compiledFunc is a compiled function expression.
type compiledFunc struct {
	ast.Node

	Sig *types.Signature

	// Whether the first argument to this task is a context.Context.
	WantCtx bool

	// Whether the last result is an error.
	HasError bool

	Inputs  []types.Type // non ctx params
	Outputs []types.Type // non error results

	PosInfo *PosInfo // Used to pass information to uniquely identify a function.
}

func (c *compiler) compileFunction(expr ast.Expr) *compiledFunc {
	typ := c.info.TypeOf(expr)
	sig, ok := typ.(*types.Signature)
	if !ok {
		c.errf(c.nodePosition(expr), "expected function, got %v", typ)
		return nil
	}

	if sig.Variadic() {
		c.errf(c.nodePosition(expr), "variadic functions are not yet supported")
		return nil
	}

	f := compiledFunc{
		Node:    expr,
		Sig:     sig,
		PosInfo: c.getPosInfo(expr),
	}

	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		ptype := param.Type()
		if !isContext(ptype) {
			f.Inputs = append(f.Inputs, ptype)
			continue
		}

		if i != 0 {
			c.errf(c.position(param.Pos()), "only the first argument may be context.Context")
			return nil
		}
		f.WantCtx = true
	}

	results := sig.Results()
	for i := 0; i < results.Len(); i++ {
		result := results.At(i)
		rtype := result.Type()
		if !isError(rtype) {
			f.Outputs = append(f.Outputs, rtype)
			continue
		}
		// Error case.
		if i != results.Len()-1 {
			c.errf(c.position(result.Pos()), "only the last result may be an error")
			return nil
		}
		f.HasError = true
	}

	return &f
}

func (c *compiler) getPosInfo(n ast.Node) *PosInfo {
	pos := c.nodePosition(n)
	posInfo := &PosInfo{
		File:   filepath.Join(c.pkg.Path(), filepath.Base(pos.Filename)),
		Line:   pos.Line,
		Column: pos.Column,
	}

	return posInfo
}

func (c *compiler) interpretTaskOptions(flow *flow, t *task, opts []ast.Expr) {
	for _, opt := range opts {
		call, fn, err := c.identifyOption(opt)
		if err != nil {
			c.errf(c.nodePosition(opt), err.Error())
			continue
		}

		switch fn.Name() {
		case "FallbackWith":
			errResults := call.Args
			if len(errResults) != len(t.Outputs) {
				c.errf(c.nodePosition(opt), "cff.FallbackWith must produce the same number of results as the task: "+"expected %v, got %v", len(t.Outputs), len(errResults))
				continue
			}
			// Verify that Task returns an error for FallbackWith to be used.
			var hasError = false
			results := t.Function.Sig.Results()
			for i := 0; i < results.Len(); i++ {
				result := results.At(i)
				rtype := result.Type()
				if isError(rtype) {
					// Found error.
					hasError = true
				}
			}
			if !hasError {
				c.errf(c.nodePosition(opt), "Task must return an error for FallbackWith to be used")
				continue
			}
			for i, er := range errResults {
				give := c.info.TypeOf(er)
				want := t.Outputs[i]
				if !types.AssignableTo(give, want) {
					c.errf(
						c.nodePosition(er),
						"cff.FallbackWith result at position %v of type %v cannot be used as %v",
						i+1, give, want)
				}
			}

			t.FallbackWith = true
			t.FallbackWithResults = call.Args
		case "Predicate":
			t.Predicate = c.compilePredicate(flow, t, call)
		case "Instrument":
			t.Instrument = c.compileInstrument(call)
		case "Invoke":
			t.invokeType = c.compileInvoke(flow, call)
		}
	}
}

func (c *compiler) identifyOption(opt ast.Expr) (*ast.CallExpr, types.Object, error) {
	// All options are function calls right now.
	call, ok := opt.(*ast.CallExpr)
	if !ok {
		return nil, nil, fmt.Errorf("expected a function call, got %v", astutil.NodeDescription(opt))
	}

	sel, ok := call.Fun.(*ast.SelectorExpr)
	if !ok {
		return nil, nil, errors.New("only cff functions can be passed as task options")
	}

	fn, ok := c.info.Uses[sel.Sel]
	if !ok {
		return nil, nil, fmt.Errorf("unresolvable reference in call: %v", sel.Sel.String())
	}

	if !isPackagePathEquivalent(fn.Pkg(), cffImportPath) {
		return nil, nil, fmt.Errorf("only cff functions may be passed as task options: "+"found package %q", fn.Pkg().Path())
	}
	return call, fn, nil
}

type predicate struct {
	ast.Node

	Function *function

	Inputs []types.Type // non ctx params

	// Output is the parsed return type from the cff.Predicate invocation.
	// SentinelOutput should be used when there is a need to uniquely
	// identify the output of the predicate.
	Output types.Type

	// SentinelOutput is the sentinel value which represents the boolean
	// output of the predicate and is used to distinguish the results of the
	// predicate in the CFF graph.
	SentinelOutput *predicateOutput

	// Serial is a unique serially incrementing number for each predicate.
	Serial int

	// Task that predicate stops.
	Task *task

	PosInfo *PosInfo // Used to pass information to uniquely identify a predicate.
}

func (c *compiler) compilePredicate(f *flow, t *task, call *ast.CallExpr) *predicate {
	fn := call.Args[0]
	fnType := c.info.TypeOf(fn)

	sig, ok := fnType.(*types.Signature)
	if !ok {
		c.errf(c.nodePosition(fn), "cff.Predicate expected a function but received %v", fnType)
		return nil
	}

	if sig.Variadic() {
		c.errf(c.nodePosition(fn), "variadic functions are not yet supported")
		return nil
	}

	results := sig.Results()
	if results.Len() != 1 {
		c.errf(c.nodePosition(fn), "the function must return a single boolean result")
		return nil
	}

	if rtype, ok := results.At(0).Type().(*types.Basic); !ok || rtype.Kind() != types.Bool {
		c.errf(c.nodePosition(fn), "the function must return a single boolean result")
		return nil
	}

	compiledFunc := c.compileFunction(fn)
	if compiledFunc == nil {
		return nil
	}

	predFunc := &function{
		Node:         compiledFunc.Node,
		Sig:          compiledFunc.Sig,
		WantCtx:      compiledFunc.WantCtx,
		HasError:     compiledFunc.HasError,
		Dependencies: compiledFunc.Inputs,
		PosInfo:      compiledFunc.PosInfo,
	}

	p := &predicate{
		Node:           predFunc.Node,
		PosInfo:        c.getPosInfo(call),
		Function:       predFunc,
		Inputs:         compiledFunc.Inputs,
		Output:         compiledFunc.Outputs[0], // Predicates must have one output.
		SentinelOutput: f.addPredicateOutput(),
		Serial:         f.predicateTypeCnt,
		Task:           t,
	}
	predFunc.Predicate = p
	return p
}

type instrument struct {
	Name ast.Expr // name to use in metrics for this task
}

func (c *compiler) compileInstrument(call *ast.CallExpr) *instrument {
	name := call.Args[0]
	return &instrument{Name: name}
}

func (c *compiler) compileInstrumentName(name string) *instrument {
	return &instrument{
		Name: &ast.BasicLit{
			Kind:  token.STRING,
			Value: strconv.Quote(name),
		},
	}
}

func (c *compiler) compileInvoke(flow *flow, o *ast.CallExpr) *noOutput {
	// Bool type checking is satisfied by cff.Invoke interface.
	if len(o.Args) != 1 {
		c.errf(c.nodePosition(o.Fun), "invoke expects exactly one argument")
	}
	val, ok := c.info.Types[o.Args[0]]
	if !ok {
		c.errf(c.nodePosition(o), "expected to find a bool, found %v instead", astutil.NodeDescription(o.Args[0]))
		return nil
	}
	if constant.BoolVal(val.Value) {
		return flow.addNoOutput()
	}
	return nil
}

type input struct {
	// Reference to the parameter.
	Node ast.Expr

	// Type of the value.
	Type types.Type
}

func (c *compiler) compileInput(i ast.Expr) *input {
	return &input{
		Node: i,
		Type: c.info.TypeOf(i),
	}
}

type output struct {
	// Reference to the &foo.
	Node ast.Expr

	// Type of the target value, not the pointer.
	Type types.Type
}

func (c *compiler) compileOutput(o ast.Expr) *output {
	t := c.info.TypeOf(o)
	p, ok := t.(*types.Pointer)
	if !ok {
		c.errf(c.nodePosition(o), "invalid parameter to cff.Results: "+"expected pointer, got %v", t)
		return nil
	}

	return &output{
		Node: o,
		Type: p.Elem(),
	}
}

// isPackagePathEquivalent returns whether the path of the package is exactly equal to the path given or is equivalent due to vendoring.
//
// The package path when used in an external repo as a vendored dependency will have a different
// import path; when used in package a it will be a/vendor/b, it may even be a/vendor/b/vendor/c
// See https://github.com/golang/go/issues/12739
func isPackagePathEquivalent(pkg *types.Package, path string) bool {
	if pkg == nil {
		// pkg will be nil when the package is part of the language builtins
		return false
	}

	if pkg.Path() == path {
		return true
	}

	if strings.HasSuffix(pkg.Path(), "/vendor/"+path) {
		return true
	}

	return pkg.Path() == "vendor/"+path
}
