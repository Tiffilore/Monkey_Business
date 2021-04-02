package visualizer

import (
	"bytes"
	"fmt"
	"monkey/ast"
)

func RepresentNodeConsoleTree(node ast.Node, indent string, exclToken bool) string {
	return representNode(node, "", indent, exclToken)
}

func representStmtList(nodes []ast.Statement, thisIndent string, indent string, exclToken bool) string {
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

func representExprList(nodes []ast.Expression, thisIndent string, indent string, exclToken bool) string {
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

func representIdentifierList(identifiers []*ast.Identifier, thisIndent string, indent string, exclToken bool) string {
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

func representNode(node ast.Node, thisIndent string, indent string, exclToken bool) string {

	if node == nil {
		return White + "nil" + Reset
	}

	//type
	typestr := representNodeType(node, 1)

	var out bytes.Buffer
	out.WriteString(typestr)
	out.WriteString(" {")

	if hasNilValue(node) {
		out.WriteString("}")
		return out.String()
	}

	nextIndent := thisIndent + indent
	switch node := node.(type) {
	case *ast.Program:
		// Statements []Statement
		out.WriteString("\n" + nextIndent + "Statements: ")
		out.WriteString(representStmtList(node.Statements, nextIndent, indent, exclToken))
	case *ast.LetStatement:
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
	case *ast.ReturnStatement:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// ReturnValue Expression
		out.WriteString("\n" + nextIndent + "ReturnValue: ")
		out.WriteString(representNode(node.ReturnValue, nextIndent, indent, exclToken))
	case *ast.ExpressionStatement:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Expression Expression
		out.WriteString("\n" + nextIndent + "Expression: ")
		out.WriteString(representNode(node.Expression, nextIndent, indent, exclToken))
	case *ast.BlockStatement:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Statements []Statement
		out.WriteString("\n" + nextIndent + "Statements: ")
		out.WriteString(representStmtList(node.Statements, nextIndent, indent, exclToken))
	case *ast.Identifier:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Value string
		out.WriteString("\n" + nextIndent + "Value: ")
		fmt.Fprintf(&out, "%T %+v", node.Value, node.Value)
	case *ast.Boolean:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Value bool
		out.WriteString("\n" + nextIndent + "Value: ")
		fmt.Fprintf(&out, "%T %+v", node.Value, node.Value)
	case *ast.IntegerLiteral:
		if !exclToken { //Token
			out.WriteString("\n" + nextIndent + "Token: ")
			fmt.Fprintf(&out, "%T %+v", node.Token, node.Token)
		}
		// Value int64
		out.WriteString("\n" + nextIndent + "Value: ")
		fmt.Fprintf(&out, "%T %+v", node.Value, node.Value)
	case *ast.PrefixExpression:
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
	case *ast.InfixExpression:
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
	case *ast.IfExpression:
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
	case *ast.FunctionLiteral:
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
	case *ast.CallExpression:
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
