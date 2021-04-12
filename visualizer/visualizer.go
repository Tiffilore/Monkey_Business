package visualizer

import "reflect"

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

func abbreviateObjectType(objecttype string) string {
	switch objecttype {
	case "Integer":
		return "Intg"
	case "Boolean":
		return "Bool"
	case "Null":
		return "Null"
	case "ReturnValue":
		return "RetV"
	case "Error":
		return "Error"
	case "Function":
		return "Func"
	default:
		return objecttype
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

func isLiterallyNil(i interface{}) bool {
	return i == nil
}

func hasNilValue(i interface{}) bool {
	if isLiterallyNil(i) {
		return false
	}
	return reflect.ValueOf(i).IsNil()
}

/* fooling around:

func hasNilValue2(i ast.Node) bool {
	if isLiterallyNil(i) {
		return false
	}
	t := i.(ast.Node)
	return reflect.ValueOf(t).IsNil()
}

func WhatAmI(i interface{}) string {

	if i == nil {
		return "no type"
	}

	nodetype := reflect.TypeOf(i)
	return nodetype.String()
}

func WhatNodeInterfaceAmI(n ast.Node) string {

	if n == nil {
		return "no interface"
	}

	nodetype := reflect.TypeOf(n)

	expr_interface := reflect.TypeOf((*ast.Expression)(nil)).Elem()
	stmt_interface := reflect.TypeOf((*ast.Statement)(nil)).Elem()
	node_interface := reflect.TypeOf((*ast.Node)(nil)).Elem()

	if nodetype.Implements(expr_interface) {
		return expr_interface.String()
	}
	if nodetype.Implements(stmt_interface) {
		return stmt_interface.String()
	}
	if nodetype.Implements(node_interface) {
		return node_interface.String()
	}

	return nodetype.String()
}

func IsNode(i interface{}) bool {
	nodetype := reflect.TypeOf(i)
	node_interface := reflect.TypeOf((*ast.Node)(nil)).Elem()
	return nodetype.Implements(node_interface)
}

func RepresentAsJson(i interface{}, indent string) string {
	json, err := json.MarshalIndent(i, "", indent)
	if err == nil {
		return string(json)
	}
	return ""
}

*/