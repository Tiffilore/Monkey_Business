package vis2

import (
	"bytes"
	"monkey/ast"
	"monkey/evaluator"
	"monkey/object"
	"monkey/token"
	"reflect"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// preconditions:
//   non-circular
//   nodes are structs
//   Tokens have fields Type and Literal

func (v *Visualizer) VisualizeQTree(node ast.Node, //, t evaluator.Tracer
) string {
	v.process = PARSE
	v.display = TEX
	v.fillMaps(node)

	v.mode = WRITE
	visitedNodes = make(map[ast.Node]bool)
	//v.visualizeNode(node)

	v.visualizeFieldValue(node)

	return "\\Tree " + v.out.String()
}

func (v *Visualizer) VisualizeConsTree(node ast.Node) string {
	v.display = CONSOLE
	//v.visualizeNode(node)
	v.visualizeFieldValue(node)
	return v.prefix + v.out.String()
}

var visitedNodes map[ast.Node]bool
var visitedObjects map[object.Object]bool
var visitedEnvs map[*object.Environment]bool
var namesNodes map[string]map[ast.Node]string
var namesObjects map[string]map[object.Object]string
var namesEnvs map[*object.Environment]string
var envsOrdered []*object.Environment

func (v *Visualizer) setupEval(t *evaluator.Tracer) ast.Node {
	v.tracer = t
	v.process = EVAL
	node := t.GetRoot()
	v.fillMaps(node)
	v.mode = WRITE
	visitedNodes = make(map[ast.Node]bool)
	visitedObjects = make(map[object.Object]bool)
	visitedEnvs = make(map[*object.Environment]bool)
	v.out.Reset()
	return node
}

func (v *Visualizer) VisualizeEvalQTree(t *evaluator.Tracer, display Display) string {
	v.display = display
	node := v.setupEval(t)
	v.visualizeNode(node)
	if v.display == TEX {
		return "\\Tree " + v.out.String()
	}
	return v.out.String()
}

func (v *Visualizer) VisualizeEnvironments(t *evaluator.Tracer, display Display) string {
	v.display = display
	v.setupEval(t)

	for _, env := range envsOrdered {
		v.printInd()
		v.printInd()
		envName, _ := v.getEnvName(env)

		v.printInd(v.strEnvDep(env))
		v.printInd()
		envSnapsIntervals := v.getEnvSteps(env, t) //

		v.incrIndent()
		for _, e := range envSnapsIntervals {
			v.printInd(envName, " (", e.fromStep, " - ", e.toStep, "):")
			v.printInd()
			v.visualizeEnvTable(e.envSnap)
			v.printInd()

		}
		v.decrIndent()

	}

	return v.out.String()
}

type EnvSnapAt struct {
	envSnap *object.Environment
	step    int
}

type EnvSnapInterval struct {
	envSnap  *object.Environment
	fromStep int
	toStep   int
}

func (v *Visualizer) getEnvSteps(env *object.Environment, t *evaluator.Tracer) []*EnvSnapInterval {
	//TODO later maybe: also look at envs pointed to!!
	envSnaps := make([]EnvSnapAt, 0)
	envSnapsIntervals := make([]*EnvSnapInterval, 0)

	for i := 0; i < t.Steps(); i++ {
		if call, ok := t.Calls[i]; ok {
			if call.Env == env {
				envSnaps = append(envSnaps, EnvSnapAt{call.EnvSnap, call.No})
			}
		} else if exit, ok := t.Exits[i]; ok {
			if exit.Env == env {
				envSnaps = append(envSnaps, EnvSnapAt{exit.EnvSnap, exit.No})
			}
		}
	}

	if len(envSnaps) == 0 { //right now impossible
		return envSnapsIntervals
	}

	first := envSnaps[0]
	curInterval := &EnvSnapInterval{first.envSnap, first.step, first.step}
	envSnapsIntervals = append(envSnapsIntervals, curInterval)

	for _, e := range envSnaps {
		if reflect.DeepEqual(e.envSnap, curInterval.envSnap) {
			curInterval.toStep = e.step
		} else {
			curInterval = &EnvSnapInterval{e.envSnap, e.step, e.step}
			envSnapsIntervals = append(envSnapsIntervals, curInterval)
		}

	}
	return envSnapsIntervals
}

func (v *Visualizer) fillMaps(node ast.Node) {
	visitedNodes = make(map[ast.Node]bool)
	namesNodes = make(map[string]map[ast.Node]string)

	if v.process == EVAL {
		envsOrdered = make([]*object.Environment, 0)
		visitedObjects = make(map[object.Object]bool)
		visitedEnvs = make(map[*object.Environment]bool)
		namesObjects = make(map[string]map[object.Object]string)
		namesEnvs = make(map[*object.Environment]string)
	}
	v.mode = COLLECT
	v.visualizeNode(node)
}

func (v *Visualizer) visualizeFieldValue(i interface{}) { //visualize field

	//case nil
	if i == nil {
		v.visualizeNil() // fieldvalue
		return
	}

	// case slice
	if reflect.TypeOf(i).Kind() == reflect.Slice {

		values := reflect.Indirect(reflect.ValueOf(i))

		v.beginList(values.Len())

		for i := 0; i < values.Len(); i++ {
			if v.display == CONSOLE {
				v.printInd()
			}
			v.visualizeFieldValue(values.Index(i).Interface())
		}
		v.endList()
		return
	}

	switch i := i.(type) {

	case ast.Node:
		v.visualizeNode(i)
		return
	case token.Token:
		v.visualizeToken(i)
		return
	case *object.Environment:
		v.visualizeEnv(i)
		return
	case object.Object:
		v.visualizeObject(i)
		return
	default:
		v.visualizeLeaf(i, false)
		return

	}
}

func (v *Visualizer) visualizeNode(node ast.Node) {

	_, visited := visitedNodes[node]

	// case nil
	if node == nil {
		v.visualizeNil()
		return
	}

	//if reflect.TypeOf(node).Kind() == reflect.Ptr { // && !reflect.ValueOf(node).IsNil() { // to avoid repetitions and circles
	if _, ok := visitedNodes[node]; ok && v.mode == COLLECT { // we do not need to ask whether it is a pointer
		v.createNodeName(node)
	}

	// label node
	v.beginNode(node, visited) // case CONS: with objects!

	if reflect.ValueOf(node).IsNil() {
		v.visualizeNilValue()
		v.endNode(visited)
		return
	}

	if node, ok := node.(*ast.Identifier); ok && v.verbosity < VVV && v.display == TEX { // also if it has already been displayed!
		v.visualizeRoofed(node.String())
	}

	// children
	if _, ok := visitedNodes[node]; !ok { // we do not need to ask whether it is a pointer
		visitedNodes[node] = true

		if _, ok := node.(*ast.Identifier); ok && v.verbosity < VVV && v.display == TEX { // also if it has already been displayed!
			v.endNode(visited)
			return
		}
		nodeContVal := reflect.ValueOf(node).Elem()
		//if nodeContVal.Kind() != reflect.Struct {
		//	v.printW(" NO STRUCT VALUE") // TODO: might be an err ? für Erweiterungen
		//	return
		//}

		nodeContType := nodeContVal.Type()

		for i := 0; i < nodeContVal.NumField(); i++ {
			f := nodeContVal.Field(i)
			// label: fieldname
			fieldname := nodeContType.Field(i).Name
			if v.exclToken && fieldname == "Token" {
				continue
			}

			v.beginField(fieldname)

			// field value
			//fmt.Printf("%d: %s %s = %v\n", i,
			//	nodeContType.Field(i).Name, f.Type(), f.Interface())

			v.visualizeFieldValue(f.Interface())
			v.endField()

		}

		if v.process == EVAL && v.display == TEX {
			//add objects

			_, exits := v.getCallsAndExits(node)
			for _, exit := range exits {
				no := exit.No
				v.beginVal(no)              // label
				v.visualizeObject(exit.Val) // Name oder Darstellung
				v.endVal()                  // Klammer zu

			}
		}
	}

	v.endNode(visited)
	//TODO error: any node should be either a Statement an Expression or a Program
}

func (v *Visualizer) visualizeObject(obj object.Object) {

	// case nil
	if obj == nil {
		v.visualizeNil()
		return
	}

	if v.verbosity < VVV &&
		(obj == evaluator.FALSE ||
			obj == evaluator.TRUE ||
			obj == evaluator.NULL) {
		v.visualizeSimpleObj(obj)
		return

	}

	if reflect.ValueOf(obj).IsNil() { // so far never happens
		v.visualizeNilValue()
		v.endObject()
		return
	}

	if _, ok := visitedObjects[obj]; ok && v.mode == COLLECT { // we do not need to ask whether it is a pointer
		v.createObjectName(obj)
	} // TODO: evtl nach begin object verschieben?

	// label node
	v.beginObject(obj)

	if obj, ok := obj.(*object.Integer); ok && v.verbosity < VVV { // also if it has already been displayed!
		v.visualizeRoofed(obj.Inspect())

	}
	if obj, ok := obj.(*object.Error); ok && v.verbosity < VVV && v.display == CONSOLE { // also if it has already been displayed!
		v.visualizeRoofed(obj.Inspect())
	}

	// children --> Nilvalue
	if _, ok := visitedObjects[obj]; !ok { // we do not need to ask whether it is a pointer
		visitedObjects[obj] = true
		if _, ok := obj.(*object.Integer); ok && v.verbosity < VVV {
			v.endObject()
			return
		}
		if _, ok := obj.(*object.Error); ok && v.verbosity < VVV && v.display == CONSOLE {
			v.endObject()
			return
		}

		if obj, ok := obj.(*object.Error); ok && v.verbosity < VVV && v.display == TEX {
			v.visualizeErrorMsgShort(obj)
			v.endObject()
			return
		}

		if v.display == CONSOLE {
			v.printW(consColorize("{", Green))
		}

		objContVal := reflect.ValueOf(obj).Elem()
		//if nodeContVal.Kind() != reflect.Struct {
		//	v.printW(" NO STRUCT VALUE") // TODO: might be an err ? für Erweiterungen
		//	return
		//}

		nodeContType := objContVal.Type()

		for i := 0; i < objContVal.NumField(); i++ {
			f := objContVal.Field(i)
			// label: fieldname
			fieldname := nodeContType.Field(i).Name

			if fieldname == "Ennv" {
				continue
			}

			v.beginField(fieldname)
			//v.printW("%")
			v.visualizeFieldValue(f.Interface())
			_ = f
			v.endField()

		}

		if v.display == CONSOLE {
			v.printInd(consColorize("}", Green))
		}

	}

	v.endObject()
}

func (v *Visualizer) visualizeEnvTable(env *object.Environment) { // reusable by step-by-step-tracer?

	// case nil
	if env == nil {
		v.visualizeNil()
		return
	}

	v.printInd()

	table := getStoreTable(env)
	lines := strings.Split(table, "\n")
	for _, line := range lines {
		v.printInd(line)
	}

	v.printW("--> outer: ")

	v.incrIndent()
	v.visualizeEnvTable(env.Outer)
	v.decrIndent()
}

func getStoreTable(env *object.Environment) string {

	var temp_out bytes.Buffer
	store := env.Store

	//sort alphabetically
	keys := make([]string, 0, len(store))
	for k := range store {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	t := table.NewWriter()
	t.SetOutputMirror(&temp_out)
	t.AppendHeader(table.Row{"Identifier", "Type", "Value"})
	t.AppendSeparator()

	for _, key := range keys {
		object := store[key]
		nodetype := reflect.TypeOf(object)
		var value string
		if object == nil {
			value = "nil"
		} else {
			value = object.Inspect() //strings.ReplaceAll(object.Inspect(), "\n", "\n\t  ")
		}
		t.AppendRow([]interface{}{key, nodetype, value})
	}
	// //t.AppendFooter(table.Row{"", "", "Total", 10000})
	// //t.SetStyle(table.StyleColoredBright)
	t.Render()
	return temp_out.String()
}
