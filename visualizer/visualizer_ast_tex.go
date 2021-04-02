package visualizer

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"strings"
)

func QTree(node ast.Node, exclToken bool) string {
	return "\\Tree" + nodeQtree(node, "", exclToken) + "\n"
}

func nodeQtree(node ast.Node, thisIndent string, exclToken bool) string {

	if node == nil {
		return thisIndent + "[.nil ]"
	}

	typestr := nodeTypeQTree(node, 1)
	var out bytes.Buffer
	fmt.Fprint(&out, thisIndent, "[.", typestr)
	if hasNilValue(node) {
		fmt.Fprint(&out, " $\\emptyset$\n", thisIndent, "]")
		return out.String()
	}

	fmt.Fprint(&out, childrenQTree(node, thisIndent+indent, exclToken))

	fmt.Fprint(&out, "\n", thisIndent, "]")
	return out.String()
}

func childrenQTree(node ast.Node, thisIndent string, exclToken bool) string {
	var out bytes.Buffer
	//Token
	if !exclToken {
		switch node := node.(type) {
		case *ast.Program:
		default:
			fmt.Fprintf(&out, "\n%v\\edge node[auto=%v]{\\tiny %v};  ",
				thisIndent,
				"right",
				"Token",
			)
			fmt.Fprintf(&out, "\n%v%v", thisIndent, tokenleafQTree(node.TokenLiteral()))
		}
	}

	nextIndent := thisIndent + indent
	switch node := node.(type) {
	case *ast.Program:
		// Statements []Statement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Statements"))
		for _, stmt := range node.Statements {
			fmt.Fprint(&out, "\n", nodeQtree(stmt, nextIndent, exclToken))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.LetStatement:
		// Name  *Identifier
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Name"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Name, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Value Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Value"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Value, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.ReturnStatement:
		// ReturnValue Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("ReturnValue"))
		fmt.Fprint(&out, "\n", nodeQtree(node.ReturnValue, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.ExpressionStatement:
		// Expression Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Expression"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Expression, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.BlockStatement:
		// Statements []Statement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Statements"))
		for _, stmt := range node.Statements {
			fmt.Fprint(&out, "\n", nodeQtree(stmt, nextIndent, exclToken))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.Identifier:
		// Value string
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Value"))
		fmt.Fprintf(&out, "\n%v%v", thisIndent, leafQTree(node.Value))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.Boolean:
		// Value bool
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Value"))
		fmt.Fprintf(&out, "\n%v%v", thisIndent, leafQTree(node.Value))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.IntegerLiteral:
		// Value int64
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Value"))
		fmt.Fprintf(&out, "\n%v%v", thisIndent, leafQTree(node.Value))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.PrefixExpression:
		// Operator string
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Operator"))
		fmt.Fprintf(&out, "\n%v%v", thisIndent, leafQTree(lateXify(node.Operator)))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Right    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Right"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Right, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.InfixExpression:
		// Left    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Left"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Left, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Operator string
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Operator"))
		fmt.Fprintf(&out, "\n%v%v", thisIndent, leafQTree(lateXify(node.Operator)))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Right    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Right"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Right, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.IfExpression:
		// Condition    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Condition"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Condition, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Consequence *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Consequence"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Consequence, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Alternative *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Alternative"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Alternative, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.FunctionLiteral:
		// Parameters []*Identifier
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Parameters"))
		for _, id := range node.Parameters {
			fmt.Fprint(&out, "\n", nodeQtree(id, nextIndent, exclToken))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")

		// Body       *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Body"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Body, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.CallExpression:
		// Function  Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Function"))
		fmt.Fprint(&out, "\n", nodeQtree(node.Function, nextIndent, exclToken))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Arguments []Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Arguments"))
		for _, arg := range node.Arguments {
			fmt.Fprint(&out, "\n", nodeQtree(arg, nextIndent, exclToken))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
	default:
		fmt.Fprint(&out, "\n", thisIndent, " TODO")
	}
	return out.String()
}

func tokenleafQTree(value string) string {
	return "{\\bf " + lateXify(value) + "}"
}

func leafQTree(value interface{}) string {
	//return "\\underline{\\it " + value + "}"
	var out bytes.Buffer
	fmt.Fprintf(&out, "\\underline{\\it %T %+v }", value, value)
	return strings.ReplaceAll(out.String(), "_", "\\_")
}
