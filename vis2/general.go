package vis2

import (
	"fmt"
	"runtime"
	"strings"
)

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
		//return fieldname[0:4] //<-- make sure that length is sufficient
		return fieldname
	}
}

func abbreviateObjectType(objtype string) string {
	switch objtype {
	case "Integer":
		return "Int"
	case "Function":
		return "Fun"

	default:
		return objtype
	}
}

func roofify(str string) string {

	return fmt.Sprint("\\edge[roof];{\\small", str, "}")

}
func teXify(input string) string {
	input = strings.ReplaceAll(input, "&", "\\& ")
	input = strings.ReplaceAll(input, "{", "\\{ ")
	input = strings.ReplaceAll(input, "}", "\\} ")
	input = strings.ReplaceAll(input, "<", "$<$ ")
	input = strings.ReplaceAll(input, ">", "$>$ ")
	input = strings.ReplaceAll(input, "!=", "$!=$ ")
	input = strings.ReplaceAll(input, "==", "$==$ ")
	input = strings.ReplaceAll(input, "!", "$!$ ")
	input = strings.ReplaceAll(input, "=", "$=$ ")
	input = strings.ReplaceAll(input, "+", "$+$ ")
	input = strings.ReplaceAll(input, "-", "$-$ ")
	input = strings.ReplaceAll(input, "*", "$*$ ")
	input = strings.ReplaceAll(input, "/", "$/$ ")
	input = strings.ReplaceAll(input, "_", "\\_ ")

	return strings.Trim(input, " ")
}

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

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

// no method!
func consColorize(str, color string) string {
	return color + str + Reset
}

// no method!
func texColorize(str, bcolor, tcolor string) string { //TODO: includes also \tt
	return "\\colorbox{" + bcolor + "}{\\textcolor{" + tcolor + "}{\\tt " + str + "}}"
}
