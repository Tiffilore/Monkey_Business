package visualizer

import (
	"fmt"
	"monkey/ast"
	"monkey/object"
	"runtime"
	"strings"
)

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Cyan   = "\033[36m"
	White  = "\033[97m"
	// Purple = "\033[35m"
	// Gray = "\033[37m"
)

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Cyan = ""
		White = ""
		//	Purple = ""
		//	Gray = ""
	}
}

func consColorize(str, color string) string {
	return color + str + Reset
}

func consNode(node ast.Node, verbosity verbosity) string {
	// was	func (v *Visualizer) colorNode(str string, node ast.Node) string {

	nodeType := visNodeType(node, verbosity)
	return consColorNodeStr(nodeType, node)

}

func consColorNodeStr(nodeType string, node ast.Node) string {
	if _, ok := node.(ast.Expression); ok {
		return consColorize(nodeType, Cyan)
	} else if _, ok := node.(ast.Statement); ok {
		return consColorize(nodeType, Yellow)
	} else if _, ok := node.(*ast.Program); ok {
		return consColorize(nodeType, Blue)
	} else { //new nodes that fall under neither of these cases
		return consColorize(nodeType, Red)
	}
}

func VisObjectType(obj object.Object, verbosity int, goObjType bool) string {

	return visObjectType(obj, getVerbosity(verbosity), goObjType)
}

func visObjectType(obj object.Object, verbosity verbosity, goObjType bool) string {
	if obj == nil {
		return "<nil>"
	}

	if !goObjType {
		return string(obj.Type())
	}
	return goObjectType(obj, verbosity)
}

func goObjectType(obj object.Object, verbosity verbosity) string { // was (v *Visualizer) stringObjectType(obj object.Object) string

	str_objtype := fmt.Sprintf("%T", obj)
	if verbosity < VVV {
		str_objtype = strings.TrimLeft(str_objtype, "*object.")
	}
	if verbosity < VV {
		str_objtype = abbreviateGoObjectType(str_objtype)
	}

	return str_objtype
}

func visNodeType(node ast.Node, verbosity verbosity) string { // was (v *Visualizer) stringNodeType(node ast.Node) string {

	if node == nil {
		return "<nil>"
	}
	str_nodetype := fmt.Sprintf("%T", node)

	if verbosity < VVV {
		str_nodetype = strings.TrimLeft(str_nodetype, "*ast.")
	}
	if verbosity < VV {
		str_nodetype = abbreviateNodeType(str_nodetype)
	}

	return str_nodetype
}

func abbreviateGoObjectType(objtype string) string { // was abbreviateObjectType
	switch objtype {
	case "Integer":
		return "Int"
	case "Function":
		return "Fun"
	case "Error":
		return "Err"
	case "ReturnValue":
		return "RetV"
	default:
		if len(objtype) > 4 {
			return objtype[0:4]
		}
		return objtype
	}
}

func abbreviateNodeType(nodetype string) string {
	switch nodetype {
	case "Program": //Program
		return "Prog"
	case "LetStatement": //Statements
		return "LetS"
	case "ExpressionStatement":
		return "ExpS"
	case "BlockStatement":
		return "BlkS"
	case "ReturnStatement":
		return "RetS"
	case "IfExpression": //Expressions
		return "IfEx"
	case "InfixExpression":
		return "InfE"
	case "PrefixExpression":
		return "PreE"
	case "CallExpression":
		return "CalE"
	case "Identifier":
		return "Iden"
	case "Boolean":
		return "Bool"
	case "IntegerLiteral":
		return "IntL"
	case "FunctionLiteral":
		return "FctL"
	default:
		if len(nodetype) > 4 {
			return nodetype[0:4]
		}
		return nodetype
	}
}

func abbreviateFieldName(fieldname string) string {
	switch fieldname {
	case "Statements":
		return "Stmts"
	case "Name":
		return "Name"
	case "Value":
		return "Val"
	case "ReturnValue":
		return "RetV"
	case "Expression":
		return "Expr"
	case "Operator":
		return "Op"
	case "Right":
		return "Right"
	case "Left":
		return "Left"
	case "Condition":
		return "Cond"
	case "Consequence":
		return "Cons"
	case "Alternative":
		return "Altr"
	case "Parameters":
		return "Params"
	case "Body":
		return "Body"
	case "Function":
		return "Func"
	case "Arguments":
		return "Args"
	default:
		if len(fieldname) > 4 {
			return fieldname[0:4]
		}
		return fieldname
	}
}

func abbreviateObjectType(objtype string) string {
	switch objtype {
	case "Integer":
		return "Int"
	case "Function":
		return "Fun"
	case "Error":
		return "Err"
	case "ReturnValue":
		return "RetV"
	default:
		if len(objtype) > 4 {
			return objtype[0:4]
		}
		return objtype
	}
}
