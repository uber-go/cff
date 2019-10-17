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

	"go.uber.org/multierr"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

// taskIndex marks cff.Results outputs.
type taskIndex int

const (
	cffImportPath   = "go.uber.org/cff"
	taskIndexRESULT = taskIndex(-1)
)

type compiler struct {
	pkg          *types.Package
	fset         *token.FileSet
	info         *types.Info
	compilerOpts CompilerOpts
	taskSerial   int
	errors       []error
}

func newCompiler(fset *token.FileSet, info *types.Info, pkg *types.Package, compilerOpts CompilerOpts) *compiler {
	return &compiler{
		fset:         fset,
		info:         info,
		pkg:          pkg,
		compilerOpts: compilerOpts,
	}
}

// CompilerOpts is a set of options to pass to the compiler that control the output of the generated code.
type CompilerOpts struct {
	InstrumentAllTasks bool
}

func (c *compiler) errf(msg string, pos token.Position, args ...interface{}) {
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
	AST *ast.File

	// Map from import path to local names of the import. If the import is
	// unnamed, it will be recorded as the package name.
	Imports map[string][]string
	// Slice because you can have the same import path multiple times with
	// different names.

	// Packages that were imported unnamed.
	UnnamedImports map[string]struct{}

	Filepath string
	Flows    []*flow
}

func (c *compiler) CompileFile(file *ast.File) (*file, error) {
	f := c.compileFile(file)
	return f, multierr.Combine(c.errors...)
}

func (c *compiler) compileFile(astFile *ast.File) *file {
	file := file{
		AST:            astFile,
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
			// We're looking for a call in the form "cff.Flow". It will be a
			// SelectorExpr where the "X" is a reference to the "cff" package.
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

			if fn.Name() != "Flow" {
				c.errf("unknown top-level cff function %q: "+
					"only cff.Flow may be called at the top-level", c.nodePosition(n), fn.Name())
			} else {
				file.Flows = append(file.Flows, c.compileFlow(astFile, n))
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

	Ctx     ast.Expr // the expression that is a local variable of type context.Context
	Metrics ast.Expr // the expression that is a local variable of type tally.Scope
	Logger  ast.Expr // the expression that is a local variable of type *zap.Logger

	Inputs  []*input
	Outputs []*output
	Tasks   []*task

	// Partition of all tasks defining a schedule in which the tasks must be
	// executed. All tasks in one of the subsets can be executed in parallel,
	// and they must all have finished executing before the next subset of
	// tasks is called.
	//
	// So, for all i, j where i < j, all tasks in Schedule[i] may be executed
	// in parallel, and they must all finish before tasks in Schedule[j] are
	// executed.
	Schedule [][]*task

	Instrument           *instrument
	ObservabilityEnabled bool

	providers *typeutil.Map // map[types.Type]int (index in Tasks)
	receivers *typeutil.Map // map[types.Type][]taskIndex tracks types needed to detect unused inputs

	invokeTypeCnt int       // input to make unique invokeType sentinels
	invokeTypes   []*output // tracks tasks of with either no results or a single error
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

func (f *flow) addInstrument(name ast.Expr) {
	f.ObservabilityEnabled = true
	f.Instrument = &instrument{Name: name}
}

// mustSetNoOutputProvider sets the provider for the no-output, panicking if the no-output sentinel type was already
// present.
func (f *flow) mustSetNoOutputProvider(key *task, value int) {
	prev := f.providers.Set(key.invokeType, value)
	if prev != nil {
		panic(fmt.Sprintf("cff assertion error: invokeType sentinel types should be unique, found %T for %dth task (defined at %v), expected to be nil", prev, value, key.Node))
	}

	if typ, ok := f.receivers.At(key.invokeType).([]taskIndex); ok {
		typ := append(typ, taskIndex(value))
		f.receivers.Set(key.invokeType, typ)
	} else {
		f.receivers.Set(key.invokeType, []taskIndex{taskIndex(value)})
	}

}

func (c *compiler) compileFlow(file *ast.File, call *ast.CallExpr) *flow {
	if len(call.Args) == 1 {
		c.errf("cff.Flow expects at least one function", c.nodePosition(call))
		return nil
	}

	flow := flow{
		Node:      call,
		Ctx:       call.Args[0],
		providers: new(typeutil.Map),
		receivers: new(typeutil.Map),
	}
	for _, arg := range call.Args[1:] {
		arg := astutil.Unparen(arg)

		ce, ok := arg.(*ast.CallExpr)
		if !ok {
			c.errf("expected a function call, got %v",
				c.nodePosition(arg), astutil.NodeDescription(arg))
			continue
		}

		f := typeutil.StaticCallee(c.info, ce)
		if f == nil || !isPackagePathEquivalent(f.Pkg(), cffImportPath) {
			c.errf("expected cff call but got %v", c.nodePosition(arg), typeutil.Callee(c.info, ce))
			continue
		}

		switch f.Name() {
		case "Params":
			for _, i := range ce.Args {
				flow.Inputs = append(flow.Inputs, c.compileInput(i))
			}
		case "Results":
			for _, o := range ce.Args {
				if output := c.compileOutput(o); output != nil {
					flow.Outputs = append(flow.Outputs, output)
					// receivers is used to look up Task for compilation checks. Since we have
					// Results, we dont need to find an associated task.
					// We don't care about other values, cff.Results should be the only receiver.
					flow.receivers.Set(output.Type, []taskIndex{taskIndexRESULT})
				}
			}
		case "Metrics":
			flow.Metrics = c.compileMetrics(&flow, ce)
		case "Logger":
			flow.Logger = c.compileLogger(&flow, ce)
		case "InstrumentFlow":
			flow.addInstrument(ce.Args[0])
		case "Task":
			if task := c.compileTask(&flow, ce.Args[0], ce.Args[1:]); task != nil {
				flow.Tasks = append(flow.Tasks, task)
			}
		}
	}
	// At this point, c.errors may be non-empty but we are continuing with more checks to catch all
	// possible errors prior to scheduling attempt and return them at once.
	c.validateInstrument(&flow)

	for i, t := range flow.Tasks {
		for _, in := range t.Dependencies {
			if typ, ok := flow.receivers.At(in).([]taskIndex); ok {
				typ := append(typ, taskIndex(i))
				flow.receivers.Set(in, typ)
			} else {
				flow.receivers.Set(in, []taskIndex{taskIndex(i)})
			}
		}

		for _, o := range t.Outputs {
			prev := flow.providers.Set(o, i)
			if prev != nil {
				pIdx := prev.(int)
				p := flow.Tasks[pIdx]
				c.errf("type %v already provided at %v",
					c.nodePosition(t), o, c.nodePosition(p))
				continue
			}
		}

		if t.invokeType != nil {
			flow.mustSetNoOutputProvider(t, i)
		}
	}

	c.validateNoUnusedOutputTypes(&flow)
	c.validateTasks(&flow)
	// At this point we may have already found some errors in c.errors.
	if err := validateFlowCycles(&flow, c.fset); err != nil {
		c.errors = append(c.errors, err)
		return nil
	}
	if len(c.errors) > 0 {
		return nil
	}
	c.scheduleFlow(&flow)
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
	for _, t := range f.Tasks {
		for _, o := range t.Outputs {
			if f.receivers.At(o) == nil {
				c.errf("unused output type %v", c.nodePosition(t.Node), o)
			}
		}
	}
}

// validateTasks walks the graph from the bottom of the graph (the outputs) to validate that
// all outputs are provided by some function. we also walk up the graph in case cff.Results
// is not the root, and check if there are any tasks with output past cff.Results.
func (c *compiler) validateTasks(f *flow) {
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

	for queue.Len() > 0 {
		t := queue.Remove(queue.Front()).(validateVisitedType)

		if visited.At(t.Type) != nil {
			// Two tasks can depend on the same input, and that is OK, but
			// we cannot allow cycles. Skip processing of a task that has
			// already been processed and handle cycle detection at a separate stage.
			continue
		}
		visited.Set(t.Type, struct{}{})

		if taskIdx, ok := f.providers.At(t.Type).(int); ok {
			task := f.Tasks[taskIdx]
			for _, i := range task.Dependencies {
				queue.PushBack(validateVisitedType{Type: i, Node: task.Node})
			}

			continue
		}

		if flowInputs.Delete(t.Type) {
			continue
		}

		c.errf("no provider found for %v", c.nodePosition(t.Node), t.Type)
	}

	if flowInputs.Len() > 0 {
		inputs := flowInputs.Keys()
		for _, inputType := range inputs {
			inputUntyped := flowInputs.At(inputType)
			input := inputUntyped.(*input)
			c.errf("unused input type %v", c.nodePosition(input.Node), input.Type)
		}
	}
}

func (c *compiler) validateInstrument(f *flow) {
	if f.ObservabilityEnabled {
		if f.Metrics == nil || f.Logger == nil {
			c.errf("cff.Instrument requires a tally.Scope and *zap.Logger to be provided: use cff.Metrics and cff.Logger", c.nodePosition(f.Node))
		}
	}
}

func (c *compiler) scheduleFlow(f *flow) {
	g := graph{
		Count: len(f.Tasks),
		Dependencies: func(taskIdx int) []int {
			var deps []int
			for _, typ := range f.Tasks[taskIdx].Dependencies {
				if i, ok := f.providers.At(typ).(int); ok {
					// For non-ok case, if we do not find a dependency amongst providers, then it
					// was passed in from Params annotation.
					deps = append(deps, i)
				}
			}
			return deps
		},
	}

	for _, o := range f.Outputs {
		g.Roots = append(g.Roots, f.providers.At(o.Type).(int))
	}
	for _, o := range f.invokeTypes {
		g.Roots = append(g.Roots, f.providers.At(o.Type).(int))
	}

	var schedule [][]*task
	for _, idxSet := range scheduleGraph(g) {
		var tasks []*task
		for _, idx := range idxSet {
			tasks = append(tasks, f.Tasks[idx])
		}
		schedule = append(schedule, tasks)
	}
	f.Schedule = schedule
}

type task struct {
	ast.Node

	Sig *types.Signature

	// Whether the first argument to this task is a context.Context.
	WantCtx bool

	// Whether the last result is an error.
	HasError bool
	// Serial is a unique serially incrementing number for each task.
	Serial int

	// Dependencies are the types required for the task, including inputs and
	// predicate inputs.
	Dependencies []types.Type

	Inputs  []types.Type // non ctx params
	Outputs []types.Type // non error results

	Predicate  *predicate  // non-nil if Predicate was provided
	Instrument *instrument // non-nil if Scope and Logger were provided

	FallbackWith        bool       // whether we should ignore errors from this function
	FallbackWithResults []ast.Expr // expressions that return a value for each return type of this function

	invokeType *noOutput // non-nil if there are no non-error results
}

// invokeType is a sentinel return type for tasks that have no non-error results.
// It can not be custom defined type, otherwise it won't work with typeutil.Map.
type noOutput = types.Struct

func (c *compiler) compileTask(flow *flow, expr ast.Expr, opts []ast.Expr) *task {
	typ := c.info.TypeOf(expr)

	sig, ok := typ.(*types.Signature)

	if !ok {
		c.errf("expected function, got %v", c.nodePosition(expr), typ)
		return nil
	}

	if sig.Variadic() {
		c.errf("variadic functions are not yet supported", c.nodePosition(expr))
		return nil
	}

	t := task{Sig: sig, Node: expr, Serial: c.taskSerial}
	c.taskSerial++

	params := sig.Params()
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		ptype := param.Type()
		if !isContext(ptype) {
			t.Inputs = append(t.Inputs, ptype)
			continue
		}

		if i != 0 {
			c.errf("only the first argument may be context.Context", c.position(param.Pos()))
			return nil
		}
		t.WantCtx = true
	}
	t.Dependencies = append(t.Dependencies, t.Inputs...)

	results := sig.Results()
	for i := 0; i < results.Len(); i++ {
		result := results.At(i)
		rtype := result.Type()
		if !isError(rtype) {
			t.Outputs = append(t.Outputs, rtype)
			continue
		}
		// Error case.
		if i != results.Len()-1 {
			c.errf("only the last result may be an error", c.position(result.Pos()))
			return nil
		}
		t.HasError = true
	}

	c.interpretTaskOptions(flow, &t, opts)

	// Check if we return nothing and we don't have an Invoke call.
	if len(t.Outputs) == 0 && t.invokeType == nil {
		c.errf("task must return at least one non-error value but currently produces zero."+
			"Did you intend to use cff.Invoke(true)?", c.nodePosition(expr))
	}
	if len(t.Outputs) > 0 && t.invokeType != nil {
		c.errf("cff.Invoke cannot be provided on a Task that produces values besides errors",
			c.nodePosition(expr))
	}
	// Create an implied Instrument(...) annotation.
	if flow.ObservabilityEnabled && c.compilerOpts.InstrumentAllTasks && t.Instrument == nil {
		taskPos := c.nodePosition(t)
		literalImpliedName := fmt.Sprintf("%s.%d", filepath.Base(taskPos.Filename), taskPos.Line)
		impliedNameQuoted := strconv.Quote(literalImpliedName)
		t.Instrument = &instrument{Name: &ast.BasicLit{
			Kind:  token.STRING,
			Value: impliedNameQuoted,
		}}
	}

	return &t
}

func (c *compiler) interpretTaskOptions(flow *flow, t *task, opts []ast.Expr) {
	for _, opt := range opts {
		// All options are function calls right now.
		call, ok := opt.(*ast.CallExpr)
		if !ok {
			c.errf("expected a function call, got %v", c.nodePosition(opt), astutil.NodeDescription(opt))
			continue
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			c.errf("only cff functions may be passed as task options", c.nodePosition(opt))
			continue
		}

		fn, ok := c.info.Uses[sel.Sel]
		if !ok {
			c.errf("unresolvable reference in call: %v", c.nodePosition(opt), sel.Sel.String())
			continue
		}

		if !isPackagePathEquivalent(fn.Pkg(), cffImportPath) {
			c.errf("only cff functions may be passed as task options: "+
				"found package %q", c.nodePosition(opt), fn.Pkg().Path())
			continue
		}

		switch fn.Name() {
		case "FallbackWith":
			errResults := call.Args
			if len(errResults) != len(t.Outputs) {
				c.errf("cff.FallbackWith must produce the same number of results as the task: "+
					"expected %v, got %v", c.nodePosition(opt), len(t.Outputs), len(errResults))
				continue
			}
			// Verify that Task returns an error for FallbackWith to be used.
			var hasError = false
			results := t.Sig.Results()
			for i := 0; i < results.Len(); i++ {
				result := results.At(i)
				rtype := result.Type()
				if isError(rtype) {
					// Found error.
					hasError = true
				}
			}
			if !hasError {
				c.errf("Task must return an error for FallbackWith to be used", c.nodePosition(opt))
				continue
			}
			for i, er := range errResults {
				give := c.info.TypeOf(er)
				want := t.Outputs[i]
				if !types.AssignableTo(give, want) {
					c.errf("cff.FallbackWith result at position %v of type %v cannot be used as %v",
						c.nodePosition(er), i+1, give, want)
				}
			}

			t.FallbackWith = true
			t.FallbackWithResults = call.Args
		case "Predicate":
			t.Predicate = c.compilePredicate(t, call)
		case "Instrument":
			t.Instrument = c.compileInstrument(flow, call)
		case "Invoke":
			t.invokeType = c.compileInvoke(flow, call)
		}
	}
}

type predicate struct {
	WantCtx bool

	Node   ast.Expr
	Inputs []types.Type
}

func (c *compiler) compilePredicate(t *task, call *ast.CallExpr) *predicate {
	fn := call.Args[0]
	fnType := c.info.TypeOf(fn)

	sig, ok := fnType.(*types.Signature)
	if !ok {
		c.errf("cff.Predicate expected a function but received %v", c.nodePosition(fn), fnType)
		return nil
	}

	if sig.Variadic() {
		c.errf("variadic functions are not yet supported", c.nodePosition(fn))
		return nil
	}

	results := sig.Results()
	if results.Len() != 1 {
		c.errf("the function must return a single boolean result", c.nodePosition(fn))
		return nil
	}

	if rtype, ok := results.At(0).Type().(*types.Basic); !ok || rtype.Kind() != types.Bool {
		c.errf("the function must return a single boolean result", c.nodePosition(fn))
		return nil
	}
	var wantCtx = false
	params := sig.Params()
	var inputs []types.Type
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		ptype := param.Type()
		if !isContext(ptype) {
			inputs = append(inputs, ptype)
			continue
		}
		if i != 0 {
			c.errf("only the first argument may be context.Context", c.position(param.Pos()))
			return nil
		}
		wantCtx = true
	}

	t.Dependencies = append(t.Dependencies, inputs...)
	return &predicate{
		Node:    fn,
		Inputs:  inputs,
		WantCtx: wantCtx,
	}
}

type instrument struct {
	Name ast.Expr // name to use in metrics for this task
}

func (c *compiler) compileInstrument(flow *flow, call *ast.CallExpr) *instrument {
	name := call.Args[0]

	// It's possible to enable observability for a single task without enabling it for the flow.
	flow.ObservabilityEnabled = true

	return &instrument{Name: name}
}

func (c *compiler) compileInvoke(flow *flow, o *ast.CallExpr) *noOutput {
	// Bool type checking is satisfied by cff.Invoke interface.
	if len(o.Args) != 1 {
		c.errf("invoke expects exactly one argument", c.nodePosition(o.Fun))
	}
	val, ok := c.info.Types[o.Args[0]]
	if !ok {
		c.errf("expected to find a bool, found %v instead",
			c.nodePosition(o), astutil.NodeDescription(o.Args[0]))
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
		c.errf("invalid parameter to cff.Results: "+
			"expected pointer, got %v", c.nodePosition(o), t)
		return nil
	}

	return &output{
		Node: o,
		Type: p.Elem(),
	}
}

func (c *compiler) compileMetrics(flow *flow, call *ast.CallExpr) ast.Expr {
	return call.Args[0]
}

func (c *compiler) compileLogger(flow *flow, call *ast.CallExpr) ast.Expr {
	return call.Args[0]
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
