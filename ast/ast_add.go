package ast

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
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

func RepresentNodeConsoleTree(node Node, indent string, exclToken bool) string {
	return representNode(node, "", indent, exclToken)
}

func representStmtList(nodes []Statement, thisIndent string, indent string, exclToken bool) string {
	nextIndent := thisIndent + indent
	var out bytes.Buffer
	out.WriteString("[")

	for _, stmt := range nodes {
		out.WriteString("\n" + nextIndent)
		out.WriteString(representNode(stmt, nextIndent, indent, exclToken))
	}
	out.WriteString("\n" + thisIndent + "]")

	return out.String()
}

func representExprList(nodes []Expression, thisIndent string, indent string, exclToken bool) string {
	nextIndent := thisIndent + indent
	var out bytes.Buffer
	out.WriteString("[")

	for _, expr := range nodes {
		out.WriteString("\n" + nextIndent)
		out.WriteString(representNode(expr, nextIndent, indent, exclToken))
	}
	out.WriteString("\n" + thisIndent + "]")

	return out.String()
}

func representIdentifierList(identifiers []*Identifier, thisIndent string, indent string, exclToken bool) string {
	nextIndent := thisIndent + indent
	var out bytes.Buffer
	out.WriteString("[")

	for _, id := range identifiers {
		out.WriteString("\n" + nextIndent)
		fmt.Fprintf(&out, "%T %+v", id, id)
	}
	out.WriteString("\n" + thisIndent + "]")

	return out.String()
}

func representNode(node Node, thisIndent string, indent string, exclToken bool) string {

	if node == nil {
		return White + "nil" + Reset
	}

	//type
	typestr := representType(node)

	var out bytes.Buffer
	out.WriteString(typestr)
	out.WriteString(" {")

	if HasNilValue(node) {
		out.WriteString("}")
		return out.String()
	}

	nextIndent := thisIndent + indent
	switch node := node.(type) {
	case *Program:
		// Statements []Statement
		out.WriteString("\n" + nextIndent + "Statements: ")
		out.WriteString(representStmtList(node.Statements, nextIndent, indent, exclToken))

	case *LetStatement:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Name  *Identifier
		out.WriteString("\n" + nextIndent + "Name: ")
		out.WriteString(representNode(node.Name, nextIndent, indent, exclToken))
		// Value Expression
		out.WriteString("\n" + nextIndent + "Value: ")
		out.WriteString(representNode(node.Value, nextIndent, indent, exclToken))
	case *ReturnStatement:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// ReturnValue Expression
		out.WriteString("\n" + nextIndent + "ReturnValue: ")
		out.WriteString(representNode(node.ReturnValue, nextIndent, indent, exclToken))
	case *ExpressionStatement:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Expression Expression
		out.WriteString("\n" + nextIndent + "Expression: ")
		out.WriteString(representNode(node.Expression, nextIndent, indent, exclToken))
	case *BlockStatement:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Statements []Statement
		out.WriteString("\n" + nextIndent + "Statements: ")
		out.WriteString(representStmtList(node.Statements, nextIndent, indent, exclToken))
	case *Identifier:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Value string
		out.WriteString("\n" + nextIndent + "Value: ")
		fmt.Fprintf(&out, "%T %+v", node.Value, node.Value)
	case *Boolean:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Value bool
		out.WriteString("\n" + nextIndent + "Value: ")
		fmt.Fprintf(&out, "%T %+v", node.Value, node.Value)
	case *IntegerLiteral:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Value int64
		out.WriteString("\n" + nextIndent + "Value: ")
		fmt.Fprintf(&out, "%T %+v", node.Value, node.Value)
	case *PrefixExpression:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Operator string
		out.WriteString("\n" + nextIndent + "Operator: ")
		fmt.Fprintf(&out, "%T %+v", node.Operator, node.Operator)
		// Right    Expression
		out.WriteString("\n" + nextIndent + "Right: ")
		out.WriteString(representNode(node.Right, nextIndent, indent, exclToken))
	case *InfixExpression:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Left    Expression
		out.WriteString("\n" + nextIndent + "Left: ")
		out.WriteString(representNode(node.Left, nextIndent, indent, exclToken))
		// Operator string
		out.WriteString("\n" + nextIndent + "Operator: ")
		fmt.Fprintf(&out, "%T %+v", node.Operator, node.Operator)
		// Right    Expression
		out.WriteString("\n" + nextIndent + "Right: ")
		out.WriteString(representNode(node.Right, nextIndent, indent, exclToken))
	case *IfExpression:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Condition    Expression
		out.WriteString("\n" + nextIndent + "Condition: ")
		out.WriteString(representNode(node.Condition, nextIndent, indent, exclToken))
		// Consequence *BlockStatement
		out.WriteString("\n" + nextIndent + "Consequence: ")
		out.WriteString(representNode(node.Consequence, nextIndent, indent, exclToken))
		// Alternative *BlockStatement
		out.WriteString("\n" + nextIndent + "Alternative: ")
		out.WriteString(representNode(node.Alternative, nextIndent, indent, exclToken))
	case *FunctionLiteral:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Parameters []*Identifier
		out.WriteString("\n" + nextIndent + "Parameters: ")
		out.WriteString(representIdentifierList(node.Parameters, nextIndent, indent, exclToken))
		// Body       *BlockStatement
		out.WriteString("\n" + nextIndent + "Body: ")
		out.WriteString(representNode(node.Body, nextIndent, indent, exclToken))
	case *CallExpression:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Function  Expression
		out.WriteString("\n" + nextIndent + "Function: ")
		out.WriteString(representNode(node.Function, nextIndent, indent, exclToken))
		// Arguments []Expression
		out.WriteString("\n" + nextIndent + "Arguments: ")
		out.WriteString(representExprList(node.Arguments, nextIndent, indent, exclToken))

	default:
		out.WriteString("\n" + thisIndent + indent + "Fields not yet implemented")
	}

	out.WriteString("\n" + thisIndent + "}")
	return out.String()

}

func representType(n Node) string {

	if n == nil { //should not happen
		return Red + "nil" + Reset
	}

	nodetype := reflect.TypeOf(n)

	expr_interface := reflect.TypeOf((*Expression)(nil)).Elem()
	stmt_interface := reflect.TypeOf((*Statement)(nil)).Elem()
	node_interface := reflect.TypeOf((*Node)(nil)).Elem()

	if nodetype.Implements(expr_interface) {
		return Yellow + nodetype.String() + Reset
	}
	if nodetype.Implements(stmt_interface) {
		return Cyan + nodetype.String() + Reset
	}
	if nodetype.Implements(node_interface) {
		return Blue + nodetype.String() + Reset
	}

	return nodetype.String()
}

func IsLiterallyNil(i interface{}) bool {
	return i == nil
}

func HasNilValue(i interface{}) bool {
	if IsLiterallyNil(i) {
		return false
	}
	return reflect.ValueOf(i).IsNil()
}

func HasNilValue2(i Node) bool {
	if IsLiterallyNil(i) {
		return false
	}
	t := i.(Node)
	return reflect.ValueOf(t).IsNil()
}

func WhatAmI(i interface{}) string {

	if i == nil {
		return "no type"
	}

	nodetype := reflect.TypeOf(i)
	return nodetype.String()
}

func WhatNodeInterfaceAmI(n Node) string {

	if n == nil {
		return "no interface"
	}

	nodetype := reflect.TypeOf(n)

	expr_interface := reflect.TypeOf((*Expression)(nil)).Elem()
	stmt_interface := reflect.TypeOf((*Statement)(nil)).Elem()
	node_interface := reflect.TypeOf((*Node)(nil)).Elem()

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
	node_interface := reflect.TypeOf((*Node)(nil)).Elem()
	return nodetype.Implements(node_interface)
}

func RepresentAsJson(i interface{}, indent string) string {
	json, err := json.MarshalIndent(i, "", indent)
	if err == nil {
		return string(json)
	}
	return ""
}
