package visualizer

import (
	"monkey/ast"
	"monkey/object"
	"reflect"
	"strings"
)

const indent = "    "

func lateXify(input string) string {
	input = strings.ReplaceAll(input, "{", "\\{")
	input = strings.ReplaceAll(input, "}", "\\}")
	input = strings.ReplaceAll(input, "<", "$<$")
	input = strings.ReplaceAll(input, ">", "$>$")
	input = strings.ReplaceAll(input, "!=", "$=$")
	input = strings.ReplaceAll(input, "==", "$==$")
	input = strings.ReplaceAll(input, "!", "$!$")
	input = strings.ReplaceAll(input, "=", "$=$")
	input = strings.ReplaceAll(input, "+", "$+$")
	input = strings.ReplaceAll(input, "-", "$-$")
	input = strings.ReplaceAll(input, "*", "$*$")
	input = strings.ReplaceAll(input, "/", "$/$")
	input = strings.ReplaceAll(input, "_", "\\_")

	return input
}

func nodeTypeQTree(n ast.Node, brevity int) string {
	var bcolor string
	var tcolor string

	if n == nil { //should not happen
		return "nil"
	}

	nodetype := reflect.TypeOf(n)

	nodetype_str := nodetype.String()
	if brevity > 0 {
		nodetype_str = strings.TrimLeft(nodetype_str, "*ast.")
	}
	if brevity > 1 {
		nodetype_str = abbreviateNodeType(nodetype_str)

	}
	expr_interface := reflect.TypeOf((*ast.Expression)(nil)).Elem()
	stmt_interface := reflect.TypeOf((*ast.Statement)(nil)).Elem()
	node_interface := reflect.TypeOf((*ast.Node)(nil)).Elem()

	if nodetype.Implements(expr_interface) {
		bcolor = "bluish"
		tcolor = "black"
	} else if nodetype.Implements(stmt_interface) {
		bcolor = "yellish"
		tcolor = "black"
	} else if nodetype.Implements(node_interface) { //program
		bcolor = "dbluish"
		tcolor = "white"
	} else { // should never happen
		bcolor = "red"
		tcolor = "black"
	}

	return "\\colorbox{" + bcolor + "}{\\textcolor{" + tcolor + "}{\\tt " + nodetype_str + "}}"
}

func objTypeQTree(o object.Object, brevity int) string {

	var objtype_str string

	if o == nil {
		objtype_str = "nil"
	} else {
		objtype_str = reflect.TypeOf(o).String()

	}
	if brevity > 0 {
		objtype_str = strings.TrimLeft(objtype_str, "*object.")
	}
	if brevity > 1 {
		objtype_str = abbreviateObjectType(objtype_str)

	}
	return "\\colorbox{" + "black" + "}{\\textcolor{" + "white" + "}{\\tt " + objtype_str + "}}"

}

func fieldQTree(value string) string {
	return "{\\small " + value + "}"
}
