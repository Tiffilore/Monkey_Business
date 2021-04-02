package evaluator

import (
	"monkey/ast"
	"monkey/object"
)

type Call struct {
	No    int
	Id    int
	Depth int
	Node  ast.Node
	Env   object.Environment
}

type Exit struct {
	No    int
	Id    int
	Depth int
	Node  ast.Node
	Env   object.Environment
	Val   object.Object
}

type Tracer struct {
	Calls   map[int]Call
	Exits   map[int]Exit
	active  bool
	counter int
	id      int
	depth   int
}

func NewTracer() *Tracer {
	calls := make(map[int]Call)
	exits := make(map[int]Exit)

	return &Tracer{
		Calls:   calls,
		Exits:   exits,
		active:  false,
		counter: 0,
		id:      0,
		depth:   0,
	}
}

func (t Tracer) Steps() int {
	return t.counter
}

func (t Tracer) Clear() {
	t.Calls = make(map[int]Call)
	t.Exits = make(map[int]Exit)
	t.counter = 0
	t.depth = 0
	t.id = 0
}

func StartTracer() {
	T = NewTracer()
	T.active = true
}

func StopTracer() {
	T.active = false
}

func traceCall(node ast.Node, env *object.Environment) int {
	if !T.active {
		return 0
	}
	var call Call
	call.No = T.counter
	no := call.No
	T.counter++
	call.Depth = T.depth
	T.depth++
	call.Id = T.id
	T.id++
	call.Node = node
	//call.Env =
	T.Calls[no] = call
	return call.Id
}

func traceExit(id int, node ast.Node, env *object.Environment, val object.Object) {
	if !T.active {
		return
	}
	var exit Exit
	exit.No = T.counter
	no := exit.No
	T.counter++
	T.depth--
	exit.Depth = T.depth
	exit.Id = id
	exit.Node = node
	//exit.Env =
	exit.Val = val
	T.Exits[no] = exit
}
