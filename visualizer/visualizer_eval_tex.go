package visualizer

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"monkey/evaluator"
	"monkey/object"
)

// use indent from visualizer_ast_tex.go

func QTreeEval(t *evaluator.Trace) string {
	rootNode := t.Calls[0].Node
	return "\\Tree" + evalNodeQtree(rootNode, "", t) + "\n"
}

func evalNodeQtree(node ast.Node, thisIndent string, t *evaluator.Trace) string {

	left := ""
	right := ""
	for _, call := range t.Calls {
		if call.Node == node {
			left = left + fmt.Sprint(call.No) + "$\\downarrow$ "
		}
	}
	for _, exit := range t.Exits {
		if exit.Node == node {
			right = right + " $\\uparrow$" + fmt.Sprint(exit.No)
		}
	}
	typestr := nodeTypeQTree(node, 2)

	var out bytes.Buffer
	fmt.Fprint(&out, thisIndent, "[.{", left, typestr, right, "}")

	if node == nil {
		fmt.Fprint(&out, thisIndent, " ]")
		return out.String()
	}

	// add children

	if hasNilValue(node) {
		fmt.Fprint(&out, " $\\emptyset$\n", thisIndent, " ]")
		return out.String()
	}

	nextIndent := thisIndent + indent
	fmt.Fprint(&out, evalChildrenQTree(node, nextIndent, t))

	// add return values

	for _, exit := range t.Exits {
		if exit.Node == node {
			fmt.Fprintf(&out, "\n%v\\edge node[auto=%v]{\\tiny %v};  ",
				nextIndent,
				"left",
				fmt.Sprintf("Val %v", exit.No),
			)
			//val := "nil"
			//if exit.Val != nil {
			//	val = strings.ReplaceAll(exit.Val.Inspect(), "\n", " ")
			//}
			//fmt.Fprintf(&out, "\n%v%v", nextIndent, val)
			fmt.Fprint(&out, evalObjQTree(exit.Val, nextIndent, t))
		}
	}

	fmt.Fprint(&out, "\n", thisIndent, "]")

	return out.String()

}

func evalChildrenQTree(node ast.Node, thisIndent string, t *evaluator.Trace) string {
	var out bytes.Buffer
	nextIndent := thisIndent + indent
	switch node := node.(type) {
	case *ast.Program:
		// Statements []Statement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Statements"))
		for _, stmt := range node.Statements {
			fmt.Fprint(&out, "\n", evalNodeQtree(stmt, nextIndent, t))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.LetStatement:
		// Name  *Identifier
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Name"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Name, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Value Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Value"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Value, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.ReturnStatement:
		// ReturnValue Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("ReturnValue"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.ReturnValue, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.ExpressionStatement:
		// Expression Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Expression"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Expression, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.BlockStatement:
		// Statements []Statement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Statements"))
		for _, stmt := range node.Statements {
			fmt.Fprint(&out, "\n", evalNodeQtree(stmt, nextIndent, t))
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
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Right, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.InfixExpression:
		// Left    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Left"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Left, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Operator string
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Operator"))
		fmt.Fprintf(&out, "\n%v%v", thisIndent, leafQTree(lateXify(node.Operator)))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Right    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Right"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Right, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.IfExpression:
		// Condition    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Condition"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Condition, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Consequence *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Consequence"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Consequence, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Alternative *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Alternative"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Alternative, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.FunctionLiteral:
		// Parameters []*Identifier
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Parameters"))
		for _, id := range node.Parameters {
			fmt.Fprint(&out, "\n", evalNodeQtree(id, nextIndent, t))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")

		// Body       *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Body"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Body, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.CallExpression:
		// Function  Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Function"))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Function, nextIndent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Arguments []Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Arguments"))
		for _, arg := range node.Arguments {
			fmt.Fprint(&out, "\n", evalNodeQtree(arg, nextIndent, t))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
	default:
		fmt.Fprint(&out, "\n", thisIndent, " TODO")
	}
	return out.String()
}

func evalObjQTree(obj object.Object, thisIndent string, t *evaluator.Trace) string {
	typestr := objTypeQTree(obj, 2)

	var out bytes.Buffer
	fmt.Fprint(&out, thisIndent, "[.{", typestr, "}")

	if obj == nil {
		fmt.Fprint(&out, thisIndent, " ]")
		return out.String()
	}
	if hasNilValue(obj) { //can that ever happen?
		fmt.Fprint(&out, " $\\emptyset$\n", thisIndent, "]")
		return out.String()
	}

	switch obj := obj.(type) {
	case *object.Integer, *object.Boolean, *object.Null, *object.ReturnValue:
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Value"))
		fmt.Fprint(&out, "\n", thisIndent, obj.Inspect())
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *object.Error:
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Message"))
		fmt.Fprint(&out, "\n", thisIndent, "{", obj.Message, "}") //TODO: in ein Kasterl oder so
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *object.Function:
		//Parameters []*ast.Identifier
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Parameters"))
		for _, id := range obj.Parameters {
			fmt.Fprint(&out, "\n", evalNodeQtree(id, thisIndent+indent, t))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")

		//Body       *ast.BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Body"))
		fmt.Fprint(&out, "\n", evalNodeQtree(obj.Body, thisIndent+indent, t))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		//Env        *Environment
	//	fmt.Fprint(&out, "\n", thisIndent, "[.", fieldQTree("Environment"))
	//	fmt.Fprint(&out, "\n", "--")
	//	fmt.Fprint(&out, "\n", thisIndent, "]")
	default:
		fmt.Fprint(&out, "\n", thisIndent, " TODO")

	}

	fmt.Fprint(&out, "\n", thisIndent, "]")
	return out.String()
}
