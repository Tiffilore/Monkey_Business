package evaluator

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

func TestRuntimeErrorsWithNil1(t *testing.T) {

	tests := []string{
		"let nil=if(true){}; !nil", // works already
		"let nil=if(true){}; -nil",
		"let nil=if(true){}; nil +0",
		"let nil=if(true){}; 0+nil",
		"let nil=if(true){}; nil-0",
		"let nil=if(true){}; 0-nil",
		"let nil=if(true){}; nil*0",
		"let nil=if(true){}; 0*nil",
		"let nil=if(true){}; nil/0",
		"let nil=if(true){}; 0/nil",
		"let nil=if(true){}; nil<0",
		"let nil=if(true){}; 0<nil",
		"let nil=if(true){}; nil>0",
		"let nil=if(true){}; 0>nil",
		"let nil=if(true){}; nil==0",
		"let nil=if(true){}; 0==nil",
		"let nil=if(true){}; nil!=0",
		"let nil=if(true){}; 0!=nil",
		"let nil=if(true){}; nil()",
	}
	for _, input := range tests {

		l := lexer.New(input)
		p := parser.New(l)
		ast := p.ParseProgram()

		// we check specifically for a runtime error caused by the evaluation
		checkRuntimeError(input, ast, len(p.Errors()) > 0, t)
	}
}

func TestRuntimeErrorsWithNil2(t *testing.T) {

	tests := []string{
		"let nil=fn(){}(); !nil", // works already
		"let nil=fn(){}(); -nil",
		"let nil=fn(){}(); nil +0",
		"let nil=fn(){}(); 0+nil",
		"let nil=fn(){}(); nil-0",
		"let nil=fn(){}(); 0-nil",
		"let nil=fn(){}(); nil*0",
		"let nil=fn(){}(); 0*nil",
		"let nil=fn(){}(); nil/0",
		"let nil=fn(){}(); 0/nil",
		"let nil=fn(){}(); nil<0",
		"let nil=fn(){}(); 0<nil",
		"let nil=fn(){}(); nil>0",
		"let nil=fn(){}(); 0>nil",
		"let nil=fn(){}(); nil==0",
		"let nil=fn(){}(); 0==nil",
		"let nil=fn(){}(); nil!=0",
		"let nil=fn(){}(); 0!=nil",
		"let nil=fn(){}(); nil()",
	}
	for _, input := range tests {

		l := lexer.New(input)
		p := parser.New(l)
		ast := p.ParseProgram()

		// we check specifically for a runtime error caused by the evaluation
		checkRuntimeError(input, ast, len(p.Errors()) > 0, t)
	}
}

func TestRuntimeErrorsWithNull(t *testing.T) {
	// suceeds already; no runtime errors
	// although looking at the results might be interesting!

	tests := []string{
		"let null=if(false){}; !null",
		"let null=if(false){}; -null",
		"let null=if(false){}; null +0",
		"let null=if(false){}; 0+null",
		"let null=if(false){}; null-0",
		"let null=if(false){}; 0-null",
		"let null=if(false){}; null*0",
		"let null=if(false){}; 0*null",
		"let null=if(false){}; null/0",
		"let null=if(false){}; 0/null",
		"let null=if(false){}; null<0",
		"let null=if(false){}; 0<null",
		"let null=if(false){}; null>0",
		"let null=if(false){}; 0>null",
		"let null=if(false){}; null==0",
		"let null=if(false){}; 0==null",
		"let null=if(false){}; null!=0",
		"let null=if(false){}; 0!=null",
		"let null= if(false){}; null()",
	}
	for _, input := range tests {

		l := lexer.New(input)
		p := parser.New(l)
		ast := p.ParseProgram()

		// we check specifically for a runtime error caused by the evaluation
		checkRuntimeError(input, ast, len(p.Errors()) > 0, t)
	}
}

func TestRuntimeErrorsWithInvalidPrograms(t *testing.T) {
	//whether that is really a problem could be discussed, since usually only valid problems are evaluated.

	// nil vs isNil !
	tests := []string{
		"let",
		"@",
		"@ let",
		"let;@; ",
	}
	for _, input := range tests {

		l := lexer.New(input)
		p := parser.New(l)
		ast := p.ParseProgram()

		// we check specifically for a runtime error caused by the evaluation
		checkRuntimeError(input, ast, len(p.Errors()) > 0, t)
	}
}

func checkRuntimeError(input string, ast *ast.Program, hasParseErrors bool, t *testing.T) {
	env := object.NewEnvironment()
	defer func() { // idea from https://golang.org/doc/effective_go#recover
		if err := recover(); err != nil {
			if hasParseErrors {
				t.Errorf("Runtime error after parse errors for %v: %q", input, err)
			} else {
				t.Errorf("Runtime error though no parse errors for %v: %q", input, err)
			}
		}
	}()
	//value :=
	Eval(ast, env)
	//t.Error(value.Inspect())
}
