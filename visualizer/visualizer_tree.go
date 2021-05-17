package visualizer

import (
	"fmt"
	"monkey/ast"
	"monkey/evaluator"
)

func ConsParseTree(node ast.Node, verbosity int, inclToken bool) string {
	strT := ""
	if inclToken {
		strT = ", with Tokens"
	}
	return fmt.Sprintf("representation of parsetree of %v with verbosity %v%v", node, verbosity, strT)
}

func TeXParseTree(node ast.Node, verbosity int, inclToken bool, file string) {
	strT := ""
	if inclToken {
		strT = ", with Tokens"
	}
	fmt.Printf("representation of parsetree of %v with verbosity %v%v in file %v\n", node, verbosity, strT, file)
}

func ConsEvalTree(
	trace *evaluator.Trace,
	verbosity int,
	inclToken bool,
	goObjType bool,
	inclEnv bool) string {
	strT := ""
	if inclToken {
		strT = ", with Tokens"
	}
	strE := ""
	if inclEnv {
		strE = ", with representation of environments"
	}
	strO := "using Monkey object types"
	if goObjType {
		strO = "using Go object types"
	}
	return fmt.Sprintf("representation of evaltree of %v with verbosity %v %v%v%v", trace.GetRoot(), verbosity, strO, strT, strE)
}

func TeXEvalTree(
	trace *evaluator.Trace,
	verbosity int,
	inclToken bool,
	goObjType bool,
	inclEnv bool,
	file string) {
	strT := ""
	if inclToken {
		strT = ", with Tokens"
	}
	strE := ""
	if inclEnv {
		strE = ", with representation of environments"
	}
	strO := "using Monkey object types"
	if goObjType {
		strO = "using Go object types"
	}
	fmt.Printf("representation of evaltree of %v with verbosity %v %v%v%v in file %v\n", trace.GetRoot(), verbosity, strO, strT, strE, file)

}
