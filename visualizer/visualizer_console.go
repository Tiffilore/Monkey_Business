package visualizer

import (
	"monkey/ast"
	"reflect"
	"strings"
)

func representNodeType(n ast.Node, brevity int) string {

	//brevity 0: whole type
	//brevity 1: without *ast.
	//brevity 2: abbreviation

	if n == nil { //should not happen
		return Red + "nil" + Reset
	}

	nodetype := reflect.TypeOf(n)

	expr_interface := reflect.TypeOf((*ast.Expression)(nil)).Elem()
	stmt_interface := reflect.TypeOf((*ast.Statement)(nil)).Elem()
	node_interface := reflect.TypeOf((*ast.Node)(nil)).Elem()

	nodetype_str := nodetype.String()
	if brevity > 0 {
		nodetype_str = strings.TrimLeft(nodetype_str, "*ast.")
	}
	if brevity > 1 {
		nodetype_str = abbreviateNodeType(nodetype_str)

	}

	if nodetype.Implements(expr_interface) {
		return Yellow + nodetype_str + Reset
	}
	if nodetype.Implements(stmt_interface) {
		return Cyan + nodetype_str + Reset
	}
	if nodetype.Implements(node_interface) {
		return Blue + nodetype_str + Reset
	}

	return nodetype_str
}

// func representObjectType(o object.Object, brevity int) string {

// 	if o == nil { //should not happen
// 		return Red + "nil" + Reset
// 	}
// 	objtype := reflect.TypeOf(o)

// 	objtype_str := objtype.String()
// 	if brevity > 0 {
// 		objtype_str = strings.TrimLeft(objtype_str, "*object.")
// 	}
// 	if brevity > 1 {
// 		objtype_str = abbreviateObjectType(objtype_str)

// 	}

// 	return objtype_str
// }
