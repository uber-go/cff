package internal

import (
	"container/list"
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strconv"
	"strings"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/types/typeutil"
)

const (
	cffImportPath = "go.uber.org/cff"
)

type compiler struct {
	pkg        *types.Package
	fset       *token.FileSet
	info       *types.Info
	taskSerial int
	errors     []error
}

func newCompiler(fset *token.FileSet, info *types.Info, pkg *types.Package) *compiler {
	return &compiler{
		fset: fset,
		info: info,
		pkg:  pkg,
	}
}

func (c *compiler) errf(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}
	c.errors = append(c.errors, errors.New(msg))
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

				if p, ok := obj.(*types.PkgName); ok {
					importName = p.Name()
				} else {
					// This will usually not happen if the code compiles but
					// we can fall back to the base name.
					importName = filepath.Base(importPath)
				}
			}
			file.Imports[importPath] = append(file.Imports[importPath], importName)
			file.UnnamedImports[importPath] = struct{}{}
			return false
		case *ast.CallExpr:
			// We're looking for a call in the form "ctf.Flow". It will be a
			// SelectorExpr where the "X" is a reference to the "ctf" package.
			sel, ok := n.Fun.(*ast.SelectorExpr)
			if !ok {
				return true // keep looking
			}

			fn := c.info.Uses[sel.Sel]

			if fn == nil || !isPackagePathEquivalent(fn.Pkg(), cffImportPath) {
				return true // keep looking
			}

			if fn.Name() != "Flow" {
				c.errf("%v: unknown top-level cff function %q: "+
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
}

func (c *compiler) compileFlow(file *ast.File, call *ast.CallExpr) *flow {
	if len(call.Args) == 0 {
		c.errf("%v: ctf.Flow expects at least one argument", c.nodePosition(call))
		return nil
	}

	flow := flow{
		Node:      call,
		Ctx:       call.Args[0],
		providers: new(typeutil.Map),
	}
	for _, arg := range call.Args[1:] {
		arg := astutil.Unparen(arg)
		ce, ok := arg.(*ast.CallExpr)
		if !ok {
			c.errf("%v: expected a function call, got %v",
				c.nodePosition(arg), astutil.NodeDescription(arg))
			continue
		}

		f := typeutil.StaticCallee(c.info, ce)
		if f == nil || !isPackagePathEquivalent(f.Pkg(), cffImportPath) {
			c.errf("%v: expected cff call but got %v", c.nodePosition(arg), typeutil.Callee(c.info, ce))
			continue
		}

		switch f.Name() {
		case "Params":
			for _, i := range ce.Args {
				flow.Inputs = append(flow.Inputs, c.compileInput(i))
			}
		case "Results":
			for _, o := range ce.Args {
				flow.Outputs = append(flow.Outputs, c.compileOutput(o))
			}
		case "Metrics":
			flow.Metrics = c.compileMetrics(&flow, ce)
		case "Logger":
			flow.Logger = c.compileLogger(&flow, ce)
		case "InstrumentFlow":
			flow.Instrument = c.compileInstrument(&flow, ce)
		case "Tasks":
			for _, f := range ce.Args {
				flow.Tasks = append(flow.Tasks, c.compileTask(&flow, f, nil /* options */))
			}
		case "Task":
			if len(ce.Args) == 0 {
				c.errf("%v: cff.Task requires at least one argument", c.nodePosition(ce))
				continue
			}
			flow.Tasks = append(flow.Tasks, c.compileTask(&flow, ce.Args[0], ce.Args[1:]))
		default:
			c.errf("%v: undefined cff function %q", c.nodePosition(ce), f.Name())
		}
	}

	for i, t := range flow.Tasks {
		for _, o := range t.Outputs {
			prev := flow.providers.Set(o, i)
			if prev != nil {
				pIdx := prev.(int)
				p := flow.Tasks[pIdx]
				c.errf("%v: type %v already provided at %v",
					c.nodePosition(t), o, c.nodePosition(p))
				continue
			}
		}
	}

	c.validateTasks(&flow)
	if len(c.errors) > 0 {
		// Can't proceed with the remaining checks if there was an error.
		return nil
		// TODO(abg): This is ugly. Fix.
	}

	// FIXME(abg): Re-enable after fixing.
	if err := validateFlowCycles(&flow); err != nil {
		c.errors = append(c.errors, err)
		return nil
	}

	c.scheduleFlow(&flow)
	return &flow
}

// validateTasks walks the graph from the bottom of the graph (the outputs) to validate that
// all outputs are provided by some function.
func (c *compiler) validateTasks(f *flow) {
	var (
		queue      = list.New() // []type.Type
		visited    typeutil.Map // map[types.Type]struct{}
		flowInputs typeutil.Map // map[types.Type]struct{}
	)

	for _, i := range f.Inputs {
		flowInputs.Set(i.Type, struct{}{})
	}

	for _, o := range f.Outputs {
		queue.PushBack(o.Type)
	}

	for queue.Len() > 0 {
		t := queue.Remove(queue.Front()).(types.Type)

		if visited.At(t) != nil {
			// Two tasks can depend on the same input, and that is OK, but
			// we cannot allow cycles. Skip processing of a task that has
			// already been processed and handle cycle detection at a separate stage.
			continue
		}
		visited.Set(t, struct{}{})

		if taskIdx, ok := f.providers.At(t).(int); ok {
			task := f.Tasks[taskIdx]
			for _, i := range task.Dependencies {
				queue.PushBack(i)
			}

			continue
		}

		if flowInputs.Delete(t) {
			continue
		}

		// TODO: include path
		c.errf("no provider found for %v", t)
	}

	if flowInputs.Len() > 0 {
		// TODO: ordering is unspecified so sort
		c.errf("unused inputs: %v", flowInputs.KeysString())
	}
}

func (c *compiler) scheduleFlow(f *flow) {
	g := graph{
		Count: len(f.Tasks),
		Dependencies: func(taskIdx int) []int {
			var deps []int
			for _, typ := range f.Tasks[taskIdx].Dependencies {
				if i, ok := f.providers.At(typ).(int); ok {
					// TODO: Handle non-ok
					deps = append(deps, i)
				}
			}
			return deps
		},
	}

	for _, o := range f.Outputs {
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

	Predicate   *predicate  // non-nil if Predicate was provided
	Instrument  *instrument // non-nil if Metrics and Logger were provided
	RecoverWith []ast.Expr
}

func (c *compiler) compileTask(flow *flow, f ast.Expr, opts []ast.Expr) *task {
	typ := c.info.TypeOf(f)
	sig, ok := typ.(*types.Signature)
	if !ok {
		c.errf("%v: expected function, got %v", c.nodePosition(f), typ)
		return nil
	}

	if sig.Variadic() {
		c.errf("%v: variadic functions are not yet supported", c.nodePosition(f))
		return nil
	}

	t := task{Sig: sig, Node: f, Serial: c.taskSerial}
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
			c.errf("%v: only the first argument may be context.Context", c.position(param.Pos()))
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

		if i != results.Len()-1 {
			c.errf("%v: only the last result may be an error", c.position(result.Pos()))
			return nil
		}
		t.HasError = true
	}

	c.interpretTaskOptions(flow, &t, opts)
	return &t
}

func (c *compiler) interpretTaskOptions(flow *flow, t *task, opts []ast.Expr) {
	for _, opt := range opts {
		// All options are function calls right now.
		call, ok := opt.(*ast.CallExpr)
		if !ok {
			c.errf("%v: expected a function call, got %v", c.nodePosition(opt), astutil.NodeDescription(opt))
			continue
		}

		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			c.errf("%v: only cff functions may be passed as task options", c.nodePosition(opt))
			continue
		}

		fn := c.info.Uses[sel.Sel]
		if fn == nil || !isPackagePathEquivalent(fn.Pkg(), cffImportPath) {
			c.errf("%v: only cff functions may be passed as task options: "+
				"found package %q", c.nodePosition(opt), fn.Pkg().Path())
			continue
		}

		switch fn.Name() {
		case "RecoverWith":
			errResults := call.Args
			if len(errResults) != len(t.Outputs) {
				c.errf("%v: cff.RecoverWith must produce the same number of results as the task: "+
					"expected %v, got %v", c.nodePosition(opt), len(t.Outputs), len(errResults))
				continue
			}

			for i, er := range errResults {
				give := c.info.TypeOf(er)
				want := t.Outputs[i]
				if !types.AssignableTo(give, want) {
					c.errf("%v: cff.RecoverWith result at position %v of type %v cannot be used as %v",
						c.nodePosition(er), i+1, give, want)
				}
			}

			t.RecoverWith = call.Args
		case "Predicate":
			t.Predicate = c.compilePredicate(t, call)
		case "Instrument":
			t.Instrument = c.compileInstrument(flow, call)
		default:
			c.errf("%v: unknown task option %q", c.nodePosition(opt), fn.Name())
		}
	}
}

type predicate struct {
	WantCtx bool // TODO(abg): support this

	Node   ast.Expr
	Inputs []types.Type
}

func (c *compiler) compilePredicate(t *task, call *ast.CallExpr) *predicate {
	if len(call.Args) != 1 {
		c.errf("%v: cff.Predicate accepts exactly one argument: received %v", c.nodePosition(call), len(call.Args))
		return nil
	}

	fn := call.Args[0]
	fnType := c.info.TypeOf(fn)

	sig, ok := fnType.(*types.Signature)
	if !ok {
		c.errf("%v: cff.Predicate expected a function but received %v", c.nodePosition(fn), fnType)
		return nil
	}

	if sig.Variadic() {
		c.errf("%v: variadic functions are not yet supported", c.nodePosition(fn))
		return nil
	}

	results := sig.Results()
	if results.Len() != 1 {
		c.errf("%v: the function must return a single boolean result", c.nodePosition(fn))
		return nil
	}

	if rtype, ok := results.At(0).Type().(*types.Basic); !ok || rtype.Kind() != types.Bool {
		c.errf("%v: the function must return a single boolean result", c.nodePosition(fn))
		return nil
	}

	params := sig.Params()
	var inputs []types.Type
	for i := 0; i < params.Len(); i++ {
		param := params.At(i)
		ptype := param.Type()
		if isContext(ptype) {
			// TODO(abg): We can support this pretty easily.
			c.errf("%v: cff.Predicate may not depend on the context", c.position(param.Pos()))
		} else {
			inputs = append(inputs, ptype)
		}
	}

	t.Dependencies = append(t.Dependencies, inputs...)
	return &predicate{
		Node:   fn,
		Inputs: inputs,
	}
}

type instrument struct {
	Name ast.Expr // name to use in metrics for this task
}

func (c *compiler) compileInstrument(flow *flow, call *ast.CallExpr) *instrument {
	// TODO(jacobg): Accept additional tags
	if len(call.Args) != 1 {
		c.errf("%v: cff.Instrument accepts exactly one argument: received %v", c.nodePosition(call), len(call.Args))
		return nil
	}

	if flow.Metrics == nil || flow.Logger == nil {
		c.errf("%v: cff.Instrument requires a tally.Scope and *zap.Logger to be provided: use cff.Metrics and cff.Logger", c.nodePosition(call))
		return nil
	}
	flow.ObservabilityEnabled = true

	name := call.Args[0]
	nameType := c.info.TypeOf(name)
	if nt, ok := nameType.(*types.Basic); !ok || nt.Kind() != types.String {
		c.errf("%v: cff.Instrument accepts a single string argument, got %v", c.nodePosition(name), nameType)
		return nil
	}

	return &instrument{Name: name}
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
		c.errf("%v: invalid parameter to cff.Results: "+
			"expected pointer, got %v", c.nodePosition(o), t)
		return nil
	}

	return &output{
		Node: o,
		Type: p.Elem(),
	}
}

func (c *compiler) compileMetrics(flow *flow, call *ast.CallExpr) ast.Expr {
	if len(call.Args) != 1 {
		c.errf("%v: cff.Metrics accepts exactly one argument: received %v", c.nodePosition(call), len(call.Args))
		return nil
	}

	return call.Args[0]
}

func (c *compiler) compileLogger(flow *flow, call *ast.CallExpr) ast.Expr {
	if len(call.Args) != 1 {
		c.errf("%v: cff.Logger accepts exactly one argument: received %v", c.nodePosition(call), len(call.Args))
		return nil
	}

	return call.Args[0]
}

// isPackagePathEquivalent returns whether the path of the package is exactly equal to the path given or is equivalent due to vendoring.
//
// The package path when used in an external repo as a vendored dependency will have a different
// import path; when used in package a it will be a/vendor/b, it may even be a/vendor/b/vendor/c
// See https://github.com/golang/go/issues/12739
func isPackagePathEquivalent(pkg *types.Package, path string) bool {
	return pkg.Path() == path || strings.HasSuffix(pkg.Path(), "/vendor/"+path)
}
