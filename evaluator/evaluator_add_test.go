package evaluator

import (
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"

	"testing"
)

func TestRuntimeErrorsNotEnoughArguments(t *testing.T) {

	tests := []string{
		"let zero = fn(x){0}; zero()",
		"let id = fn(x){x}; id()",
		"let add = fn(x,y){x+y}; add(2)",
	}
	for _, input := range tests {

		l := lexer.New(input)
		p := parser.New(l)
		ast := p.ParseProgram()

		// we check specifically for a runtime error caused by the evaluation
		checkRuntimeError(input, ast, len(p.Errors()) > 0, t)
	}
}

func TestArityCallExpressions(t *testing.T) {

	//after fixing the runtime error problem, we need to specify how to treat function calls with not enough arguments

	tests := []struct {
		input  string
		expErr bool
	}{
		{"let id = fn(x){x}; id()", true},    // not enough arguments
		{"let id = fn(x){x}; id(1)", false},  // just right
		{"let id = fn(x){x}; id(1,2)", true}, // too many arguments // that is disputable!
	}

	for _, tt := range tests {
		testArityCallExpressions(tt.input, tt.expErr, t)
	}

}

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
				t.Errorf("Runtime error after parse errors " + input)

			} else {
				t.Errorf("Runtime error though no parse errors for " + input)
			}
		}
	}()
	//value :=
	Eval(ast, env)
	//t.Errorf(value.Inspect())
}

func testArityCallExpressions(input string, expErr bool, t *testing.T) {

	l := lexer.New(input)
	p := parser.New(l)
	ast := p.ParseProgram()
	env := object.NewEnvironment()

	defer func() {
		if err := recover(); err != nil {
			t.Errorf("Runtime error for " + input)
		}
	}()

	value := Eval(ast, env)
	_, hasErr := value.(*object.Error)

	if expErr && !hasErr {
		t.Errorf("Error message missing for " + input)
		return
	}
	if !expErr && hasErr {
		t.Errorf("Explainworthy error message for " + input)
	}
}
