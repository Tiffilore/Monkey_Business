package visualizer

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"reflect"
	"strings"
)

const indent = "    "

func QTree(node ast.Node, exclToken bool) string {
	return "\\Tree" + nodeQtree(node, "", exclToken) + "\n"
}

func nodeQtree(node ast.Node, thisIndent string, exclToken bool) string {

	if node == nil {
		return thisIndent + "[.nil ]"
	}

	typestr := typeQTree(node)
	var out bytes.Buffer
	fmt.Fprint(&out, thisIndent, "[.", typestr)
	if ast.HasNilValue(node) {
		fmt.Fprint(&out, " $\\emptyset$\n", thisIndent, "]")
		return out.String()
	}
	nextIndent := thisIndent + indent

	//Token
	if !exclToken {
		switch node := node.(type) {
		case *ast.Program:
		default:
			fmt.Fprintf(&out, "\n%v\\edge node[auto=%v]{\\tiny %v};  ",
				nextIndent,
				"right",
				"Token",
			)
			fmt.Fprintf(&out, "\n%v%v", nextIndent, tokenleafQTree(node.TokenLiteral()))
		}
	}
	switch node := node.(type) {
	case *ast.Program:
		// Statements []Statement
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Statements"))
		for _, stmt := range node.Statements {
			fmt.Fprint(&out, "\n", nodeQtree(stmt, nextIndent+indent, exclToken))
		}
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.LetStatement:
		// Name  *Identifier
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Name"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Name, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
		// Value Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Value"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Value, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.ReturnStatement:
		// ReturnValue Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("ReturnValue"))
		fmt.Fprint(&out, "\n", nodeQtree(node.ReturnValue, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.ExpressionStatement:
		// Expression Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Expression"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Expression, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.BlockStatement:
		// Statements []Statement
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Statements"))
		for _, stmt := range node.Statements {
			fmt.Fprint(&out, "\n", nodeQtree(stmt, nextIndent+indent, exclToken))
		}
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.Identifier:
		// Value string
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Value"))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafQTree(node.Value))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.Boolean:
		// Value bool
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Value"))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafQTree(node.Value))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.IntegerLiteral:
		// Value int64
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Value"))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafQTree(node.Value))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.PrefixExpression:
		// Operator string
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Operator"))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafQTree(lateXify(node.Operator)))
		fmt.Fprint(&out, "\n", nextIndent, "]")
		// Right    Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Right"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Right, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.InfixExpression:
		// Left    Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Left"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Left, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
		// Operator string
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Operator"))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafQTree(lateXify(node.Operator)))
		fmt.Fprint(&out, "\n", nextIndent, "]")
		// Right    Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Right"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Right, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.IfExpression:
		// Condition    Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Condition"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Condition, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
		// Consequence *BlockStatement
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Consequence"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Consequence, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
		// Alternative *BlockStatement
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Alternative"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Alternative, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.FunctionLiteral:
		// Parameters []*Identifier
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Parameters"))
		for _, id := range node.Parameters {
			fmt.Fprint(&out, "\n", nodeQtree(id, nextIndent+indent, exclToken))
		}
		fmt.Fprint(&out, "\n", nextIndent, "]")

		// Body       *BlockStatement
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Body"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Body, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
	case *ast.CallExpression:
		// Function  Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Function"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Function, nextIndent+indent, exclToken))
		fmt.Fprint(&out, "\n", nextIndent, "]")
		// Arguments []Expression
		fmt.Fprint(&out, "\n", nextIndent, "[.", fieldQTree("Arguments"))
		for _, arg := range node.Arguments {
			fmt.Fprint(&out, "\n", nodeQtree(arg, nextIndent+indent, exclToken))
		}
		fmt.Fprint(&out, "\n", nextIndent, "]")
	default:
		fmt.Fprint(&out, "\n", nextIndent, " TODO")

	}

	//

	fmt.Fprint(&out, "\n", thisIndent, "]")
	return out.String()

}

func tokenleafQTree(value string) string {
	return "{\\bf " + lateXify(value) + "}"
}

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

func leafQTree(value interface{}) string {
	//return "\\underline{\\it " + value + "}"
	var out bytes.Buffer
	fmt.Fprintf(&out, "\\underline{\\it %T %+v }", value, value)
	return strings.ReplaceAll(out.String(), "_", "\\_")
}

func fieldQTree(value string) string {
	return "{\\small " + value + "}"
}

func typeQTree(n ast.Node) string {
	var bcolor string
	var tcolor string

	if n == nil { //should not happen
		return "\\{nil\\}"
	}

	nodetype := reflect.TypeOf(n)

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

	nodetype_str := strings.TrimLeft(nodetype.String(), "*ast.")

	return "\\colorbox{" + bcolor + "}{\\textcolor{" + tcolor + "}{\\tt " + nodetype_str + "}}"

}
