package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

type Call struct {
	No      int
	Id      int
	Depth   int
	Node    ast.Node
	Env     *object.Environment
	EnvSnap *object.Environment
}

type Exit struct {
	No      int
	Id      int
	Depth   int
	Node    ast.Node
	Env     *object.Environment
	EnvSnap *object.Environment
	Val     object.Object
}

type Trace struct {
	Calls        map[int]Call
	Exits        map[int]Exit
	Environments []*object.Environment
	counter      int
}

type tracer struct {
	active       bool
	counter      int
	id           int
	depth        int
	calls        map[int]Call
	exits        map[int]Exit
	environments []*object.Environment
}

func newTracer() *tracer {
	calls := make(map[int]Call)
	exits := make(map[int]Exit)
	environments := make([]*object.Environment, 0)

	return &tracer{
		calls:        calls,
		exits:        exits,
		active:       false,
		counter:      0,
		id:           0,
		depth:        0,
		environments: environments,
	}
}

func startTracer() {
	t = newTracer()
	t.active = true
}

func stopTracer() {
	t.active = false
}

func traceCall(node ast.Node, env *object.Environment) int {
	if !t.active {
		return 0
	}
	var call Call
	call.No = t.counter
	no := call.No
	t.counter++
	call.Depth = t.depth
	t.depth++
	call.Id = t.id
	t.id++
	call.Node = node
	call.Env = env
	call.EnvSnap = copyEnv(env)
	t.calls[no] = call
	new := true
	for _, e := range t.environments {
		if env == e {
			new = false
		}
	}
	if new {
		t.environments = append(t.environments, env)
	}
	return call.Id
}

func traceExit(id int, node ast.Node, env *object.Environment, val object.Object) {
	if !t.active {
		return
	}
	var exit Exit
	exit.No = t.counter
	no := exit.No
	t.counter++
	t.depth--
	exit.Depth = t.depth
	exit.Id = id
	exit.Node = node
	exit.Env = env
	exit.EnvSnap = copyEnv(env)
	exit.Val = val
	t.exits[no] = exit
}

func copyEnv(env *object.Environment) *object.Environment {
	newEnv := object.NewEnvironment()
	for name, val := range env.Store {
		newEnv.Set(name, val)
	}
	if env.Outer != nil {
		newEnv.Outer = copyEnv(env.Outer)
	}
	return newEnv
}

func (t tracer) getTrace() *Trace {
	return &Trace{
		Calls:        t.calls,
		Exits:        t.exits,
		Environments: t.environments,
		counter:      t.counter,
	}
}

func (t Trace) GetRoot() ast.Node {
	if len(t.Calls) == 0 {
		return nil
	}
	return t.Calls[0].Node
}

func (t Trace) Steps() int {
	return t.counter
}
