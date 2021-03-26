package evaluator

import (
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/token"

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

func TestEvalToBoolConsistency(t *testing.T) {

	// prep: create environment and build the function object (fn(){})
	env := object.NewEnvironment()
	params := []*ast.Identifier{}
	body := &ast.BlockStatement{
		Token: token.Token{Type: token.LBRACE, Literal: "{"}}
	body.Statements = []ast.Statement{}
	functionObj := &object.Function{Parameters: params, Env: env, Body: body}

	tests := []struct {
		object      object.Object
		description string
	}{
		{TRUE, "Boolean with value true"},
		{FALSE, "Boolean with value false"},
		{&object.Integer{Value: -1},
			"Integer with negative value"},
		{&object.Integer{Value: 0},
			"Integer with zero value"},
		{&object.Integer{Value: 1},
			"Integer with positive value"},
		{NULL, "Null"},
		{&object.Error{Message: ""}, "Error"},
		{functionObj, "Function"},
		{nil, "nil"},
	}

	for _, tt := range tests {
		env.Set("a", tt.object)

		expr1 := "if(a){true} else {false}"
		expr2 := "!!a"

		if evaluate(expr1, env, t) != evaluate(expr2, env, t) {
			t.Errorf("inconsistent evaluation to bool for " + tt.description)
		}
	}
}

func TestEvalToBoolCorrectness(t *testing.T) {

	// prep: create environment and build the function object (fn(){})
	env := object.NewEnvironment()
	params := []*ast.Identifier{}
	body := &ast.BlockStatement{
		Token: token.Token{Type: token.LBRACE, Literal: "{"}}
	body.Statements = []ast.Statement{}
	functionObj := &object.Function{Parameters: params, Env: env, Body: body}

	tests := []struct {
		object      object.Object
		description string
		expected    string
	}{
		{TRUE, "Boolean with value true", "true"},
		{FALSE, "Boolean with value false", "false"},
		{&object.Integer{Value: -1},
			"Integer with negative value", "error"},
		{&object.Integer{Value: 0},
			"Integer with zero value", "error"},
		{&object.Integer{Value: 1},
			"Integer with positive value", "error"},
		{NULL, "Null", "error"},
		{&object.Error{Message: ""}, "Error", "error"},
		{functionObj, "Function", "error"},
		{nil, "nil", "true"},
	}

	for _, tt := range tests {
		env.Set("a", tt.object)

		result := evaluate("!!a", env, t)

		switch tt.expected {
		case "true":
			if result != TRUE {
				t.Errorf(tt.description + " does not evaluate to true")
			}
		case "false":
			if result != FALSE {
				t.Errorf(tt.description + " does not evaluate to false")
			}
		case "error":
			if result.Type() != "ERROR" {
				t.Errorf(tt.description + " does not evaluate to an error")
			}
		default:
			fmt.Println("ha")
		}
	}
}

func TestDivisionByZero(t *testing.T) {

	tests := []string{
		"3/0",
		"-3/(1-1)",
	}
	for _, input := range tests {
		testDivisionByZero(input, t)
	}

}

func testDivisionByZero(input string, t *testing.T) {
	l := lexer.New(input)
	p := parser.New(l)
	ast := p.ParseProgram()

	defer func() {
		if err := recover(); err != nil {
			t.Errorf("Runtime error for " + input)
		}
	}()

	env := object.NewEnvironment()
	result := Eval(ast, env)
	if result.Type() != "ERROR" {
		t.Errorf("division by zero does not evaluate to an error")
	}

}
func evaluate(input string, env *object.Environment, t *testing.T) object.Object {

	l := lexer.New(input)
	p := parser.New(l)
	ast := p.ParseProgram()
	return Eval(ast, env)
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
