package visualizer

import (
	"monkey/ast"
	"reflect"
	"runtime"
	"strings"
)

var Reset = "\033[0m"
var Red = "\033[31m"
var Green = "\033[32m"
var Yellow = "\033[33m"
var Blue = "\033[34m"
var Purple = "\033[35m"
var Cyan = "\033[36m"
var Gray = "\033[37m"
var White = "\033[97m"

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
		Gray = ""
		White = ""
	}
}

func RepresentType(n ast.Node, brevity int) string {

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
		nodetype_str = abbreviate(nodetype_str)

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
