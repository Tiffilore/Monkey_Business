package visualizer

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"monkey/evaluator"
	"monkey/object"
	"strings"
)

var objects []object.Object

func QTreeEval(t *evaluator.Tracer, brevity int) string {
	rootNode := t.Calls[0].Node
	if brevity > 1 {
		objects = make([]object.Object, 0)
	}
	qtree := "\\Tree" + evalNodeQtree(rootNode, "", t, brevity) + "\n"
	if brevity > 1 {
		objects = nil
	}
	return qtree
}

func evalNodeQtree(node ast.Node, thisIndent string, t *evaluator.Tracer, brevity int) string {

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
	typestr := nodeTypeQTree(node, brevity)

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
	fmt.Fprint(&out, evalChildrenQTree(node, nextIndent, t, brevity))

	// add return values

	for _, exit := range t.Exits {
		if exit.Node == node {
			fmt.Fprintf(&out, "\n%v\\edge node[auto=%v]{\\tiny %v};  ",
				nextIndent,
				"left",
				fmt.Sprintf("Val %v", exit.No),
			)
			fmt.Fprint(&out, evalObjQTree(exit.Val, nextIndent, t, brevity))
		}
	}

	fmt.Fprint(&out, "\n", thisIndent, "]")

	return out.String()
}

func evalChildrenQTree(node ast.Node, thisIndent string, t *evaluator.Tracer, brevity int) string {
	var out bytes.Buffer
	nextIndent := thisIndent + indent
	switch node := node.(type) {
	case *ast.Program:
		// Statements []Statement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Statements", brevity))
		for _, stmt := range node.Statements {
			fmt.Fprint(&out, "\n", evalNodeQtree(stmt, nextIndent, t, brevity))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.LetStatement:
		// Name  *Identifier
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Name", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Name, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Value Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Value", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Value, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.ReturnStatement:
		// ReturnValue Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("ReturnValue", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.ReturnValue, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.ExpressionStatement:
		// Expression Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Expression", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Expression, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.BlockStatement:
		// Statements []Statement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Statements", brevity))
		for _, stmt := range node.Statements {
			fmt.Fprint(&out, "\n", evalNodeQtree(stmt, nextIndent, t, brevity))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.Identifier:
		// Value string
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Value", brevity))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafValueQTree(node.Value, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.Boolean:
		// Value bool
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Value", brevity))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafValueQTree(node.Value, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.IntegerLiteral:
		// Value int64
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Value", brevity))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafValueQTree(node.Value, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.PrefixExpression:
		// Operator string
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Operator", brevity))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafValueQTree(lateXify(node.Operator), brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Right    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Right", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Right, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.InfixExpression:
		// Left    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Left", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Left, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Operator string
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Operator", brevity))
		fmt.Fprintf(&out, "\n%v%v", nextIndent, leafValueQTree(lateXify(node.Operator), brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Right    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Right", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Right, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.IfExpression:
		// Condition    Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Condition", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Condition, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Consequence *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Consequence", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Consequence, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Alternative *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Alternative", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Alternative, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.FunctionLiteral:
		// Parameters []*Identifier
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Parameters", brevity))
		for _, id := range node.Parameters {
			fmt.Fprint(&out, "\n", evalNodeQtree(id, nextIndent, t, brevity))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")

		// Body       *BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Body", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Body, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *ast.CallExpression:
		// Function  Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Function", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(node.Function, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		// Arguments []Expression
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Arguments", brevity))
		for _, arg := range node.Arguments {
			fmt.Fprint(&out, "\n", evalNodeQtree(arg, nextIndent, t, brevity))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
	default:
		fmt.Fprint(&out, "\n", thisIndent, " TODO")
	}
	return out.String()
}

func evalObjQTree(obj object.Object, thisIndent string, t *evaluator.Tracer, brevity int) string {
	typestr := objTypeQTree(obj, brevity)

	var out bytes.Buffer

	if obj == nil {
		fmt.Fprint(&out, thisIndent, "[.{", blacken(typestr), "}")
		fmt.Fprint(&out, thisIndent, " ]")
		return out.String()
	}
	if hasNilValue(obj) { //can that ever happen?
		fmt.Fprint(&out, thisIndent, "[.{", blacken(typestr), "}")
		fmt.Fprint(&out, " $\\emptyset$\n", thisIndent, "]")
		return out.String()
	}

	if brevity > 1 { //special treatment for certain multiply used objects
		// exclude booleans + null
		switch obj := obj.(type) {
		case *object.Boolean, *object.Null: // nothing
		default:

			// check whether multiply used
			used := 0
			for _, exit := range t.Exits {
				if exit.Val == obj {
					used++
				}
			}
			if used > 1 {
				// check whether already used
				for id, o := range objects {
					if o == obj {
						// just print out id
						return fmt.Sprint(thisIndent, blacken(typestr+fmt.Sprintf("$_%v$", id)))
					}
				}
				// not already used
				id := len(objects)
				typestr = typestr + fmt.Sprintf("$_%v$", id)
				objects = append(objects, obj)
			}
		}
	}
	if brevity > 0 { //special treatment for Boolean, Null and Error
		if obj, ok := obj.(*object.Boolean); ok {
			valstr := "FALSE"
			if obj.Value {
				valstr = "TRUE"
			}
			return fmt.Sprint(thisIndent, blacken(valstr))
		}
		if _, ok := obj.(*object.Null); ok {
			return fmt.Sprint(thisIndent, blacken("NULL"))
		}
		if obj, ok := obj.(*object.Error); ok {
			fmt.Fprint(&out, thisIndent, "[.{", blacken(typestr), "} \\edge[roof];{\\small")
			message := obj.Message
			message = strings.ReplaceAll(message, ":", "\\\\")
			message = strings.ReplaceAll(message, "INTEGER", "INT")
			message = strings.ReplaceAll(message, "BOOLEAN", "BOOL")
			fmt.Fprint(&out, "\n", thisIndent, message, "} ]")
			return out.String()
		}

	}

	fmt.Fprint(&out, thisIndent, "[.{", blacken(typestr), "}")

	if obj == nil {
		fmt.Fprint(&out, thisIndent, " ]")
		return out.String()
	}
	if hasNilValue(obj) { //can that ever happen?
		fmt.Fprint(&out, " $\\emptyset$\n", thisIndent, "]")
		return out.String()
	}

	nextIndent := thisIndent + indent
	switch obj := obj.(type) {
	case *object.Integer, *object.Boolean, *object.Null:
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Value", brevity))
		fmt.Fprint(&out, "\n", nextIndent, obj.Inspect())
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *object.ReturnValue:
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("ReturnValue", brevity))
		fmt.Fprint(&out, "\n", nextIndent, evalObjQTree(obj.Value, thisIndent+indent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *object.Error:
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Message", brevity))
		fmt.Fprint(&out, "\n", nextIndent, "{", obj.Message, "}") //TODO: in ein Kasterl oder so
		fmt.Fprint(&out, "\n", thisIndent, "]")
	case *object.Function:
		//Parameters []*ast.Identifier
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Parameters", brevity))
		for _, id := range obj.Parameters {
			fmt.Fprint(&out, "\n", evalNodeQtree(id, nextIndent, t, brevity))
		}
		fmt.Fprint(&out, "\n", thisIndent, "]")
		//Body       *ast.BlockStatement
		fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Body", brevity))
		fmt.Fprint(&out, "\n", evalNodeQtree(obj.Body, nextIndent, t, brevity))
		fmt.Fprint(&out, "\n", thisIndent, "]")
		//Env        *Environment
	//	fmt.Fprint(&out, "\n", thisIndent, "[.", fieldNameQTree("Environment"))
	//	fmt.Fprint(&out, "\n", "--")
	//	fmt.Fprint(&out, "\n", thisIndent, "]")
	default:
		fmt.Fprint(&out, "\n", thisIndent, " TODO")

	}
	fmt.Fprint(&out, "\n", thisIndent, "]")

	return out.String()
}
