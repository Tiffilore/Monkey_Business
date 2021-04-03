package visualizer

import (
	"bytes"
	"fmt"
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
	//return "\\colorbox{black}{\\textcolor{white}{\\tt " + objtype_str + "}}"
	return objtype_str
}

func blacken(val string) string {
	return fmt.Sprint(
		"\\colorbox{black}{\\textcolor{white}{\\tt ",
		val,
		"}} ")
}
func fieldNameQTree(fieldname string, brevity int) string {

	if brevity > 1 {
		fieldname = abbreviateFieldName(fieldname)
	}

	return "{\\small " + fieldname + "}"
}

func fieldQTree(value string) string { // TODO: adapt ast as well; leave it for now
	return "{\\small " + value + "}"
}

func leafValueQTree(value interface{}, brevity int) string { // TODO: adapt ast as well; leave it for now
	var out bytes.Buffer

	if brevity > 1 {
		fmt.Fprintf(&out, "\\underline{\\it %+v}", value)

	} else {
		fmt.Fprintf(&out, "\\underline{\\it %T %+v}", value, value)
	}
	return strings.ReplaceAll(out.String(), "_", "\\_") // for functions
}

func leafQTree(value interface{}) string { // TODO: adapt ast as well; leave it for now
	//return "\\underline{\\it " + value + "}"
	var out bytes.Buffer
	fmt.Fprintf(&out, "\\underline{\\it %T %+v }", value, value)
	return strings.ReplaceAll(out.String(), "_", "\\_")
}
