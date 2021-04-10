package vis2

import (
	"monkey/ast"
	"monkey/token"
	"reflect"
)

// preconditions:
//   non-circular
//   nodes are structs
//   Tokens have fields Type and Literal
func (v *Visualizer) VisualizeQTree(node ast.Node) string {

	v.fillNodeMap(node)

	visited = make(map[ast.Node]bool)
	v.display = TEX
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

func (v *Visualizer) fillNodeMap(node ast.Node) {
	names = make(map[ast.Node]string)
	visited = make(map[ast.Node]bool)
	v.display = NONE
	//v.visualizeNode(node)
	v.visualizeFieldValue(node)
}

var visited map[ast.Node]bool = make(map[ast.Node]bool)
var names map[ast.Node]string = make(map[ast.Node]string)

func (v *Visualizer) visualizeFieldValue(i interface{}) { //visualize field

	//case nil
	if i == nil {
		v.visualizeNil() // fieldvalue
		return
	}

	// case slice
	if reflect.TypeOf(i).Kind() == reflect.Slice {

		v.beginList()
		values := reflect.Indirect(reflect.ValueOf(i))

		for i := 0; i < values.Len(); i++ {
			if i > 0 || v.display == CONSOLE {
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
	default:
		v.visualizeLeaf(i, false)
		return

	}
}

func (v *Visualizer) visualizeNode(node ast.Node) {

	// case nil
	if node == nil { // unnötig, wenn wir mit visualizeFieldValue starten!
		v.visualizeNil()
		return
	}

	//if reflect.TypeOf(node).Kind() == reflect.Ptr { // && !reflect.ValueOf(node).IsNil() { // to avoid repetitions and circles
	if _, ok := visited[node]; ok && v.display == NONE { // we do not need to ask whether it is a pointer
		v.createName(node)
	}

	// label node
	v.beginNode(node)

	// children
	if _, ok := visited[node]; !ok { // we do not need to ask whether it is a pointer
		visited[node] = true
		nodeContVal := reflect.ValueOf(node).Elem()
		//if nodeContVal.Kind() != reflect.Struct {
		//	v.printW(" NO STRUCT VALUE") // TODO: might be an err ? für Erweiterungen
		//	return
		//}
		if reflect.ValueOf(node).IsNil() { // später hoch, jetzt zum testen
			v.visualizeNilValue()

		} else {
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
		}
	}

	v.endNode()
	//TODO error: any node should be either a Statement an Expression or a Program
}
