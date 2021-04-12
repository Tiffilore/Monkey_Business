package vis2

import (
	"monkey/ast"
	"monkey/evaluator"
	"monkey/object"
	"monkey/token"
	"reflect"
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

func (v *Visualizer) VisualizeEvalQTree(t *evaluator.Tracer) string {
	v.tracer = t
	v.process = EVAL
	v.display = TEX
	node := t.GetRoot()

	v.fillMaps(node)
	v.mode = WRITE
	visitedNodes = make(map[ast.Node]bool)
	visitedObjects = make(map[object.Object]bool)
	visitedEnvs = make(map[*object.Environment]bool)

	v.visualizeNode(node)

	return "\\Tree " + v.out.String()
}

func (v *Visualizer) fillMaps(node ast.Node) {
	visitedNodes = make(map[ast.Node]bool)
	namesNodes = make(map[string]map[ast.Node]string)

	if v.process == EVAL {
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
	v.beginNode(node)

	if reflect.ValueOf(node).IsNil() {
		v.visualizeNilValue()
		v.endNode()
		return
	}

	// children
	if _, ok := visitedNodes[node]; !ok { // we do not need to ask whether it is a pointer
		visitedNodes[node] = true
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

		if v.process == EVAL {
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

	v.endNode()
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

	if reflect.ValueOf(obj).IsNil() {
		v.visualizeNilValue()
		v.endObject()
		return
	}

	if _, ok := visitedObjects[obj]; ok && v.mode == COLLECT { // we do not need to ask whether it is a pointer
		v.createObjectName(obj)
	} // TODO: evtl nach begin object verschieben?

	// label node
	v.beginObject(obj)

	// children --> Nilvalue
	if _, ok := visitedObjects[obj]; !ok { // we do not need to ask whether it is a pointer
		visitedObjects[obj] = true

		if obj, ok := obj.(*object.Error); ok && v.verbosity < VVV {
			v.visualizeErrorMsgShort(obj)
			v.endObject()
			return
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

	}

	v.endObject()
}

// TODO:

// STUB: abbreviateObjectType
