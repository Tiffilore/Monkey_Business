package visualizer

import (
	"fmt"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"testing"
)

/*
The tests here only test whether rendering works.
In addition they output pdf-files in directory test/ that may be checked
for accuracy
*/

const latexPath = "/usr/bin/pdflatex"

func Test_tex2pdf_Standalone(t *testing.T) {

	content := "hello"
	document := makeStandalone(content)

	err := tex2pdf(document, "tests/hello.pdf", latexPath)

	if err != nil {
		t.Errorf("Rendering did not succeed. Reason: %q", err)
	}
}

func Test_tex2pdf_StandaloneTikz(t *testing.T) {

	qtree := "\\Tree[.A [.B C D ] E ]"
	document := makeStandalone(makeTikz(qtree))

	err := tex2pdf(document, "tests/tree.pdf", latexPath)

	if err != nil {
		t.Errorf("Rendering did not succeed. Reason: %q", err)
	}
}

/*
	alphabet to be considered by Monkey fragment at end of chapter 3:
	--> all Tokens since there is the option inclToken
	IDENT may contains underscore
	--> in chapter 4 + String ==> all characters !!!!
*/

func Test_teXify(t *testing.T) {

	tests := []string{
		// Identifier
		"a12ab",
		"a12_c",
		// String
		"\"hello\"",
		"\"&%*$#~&\"",
		"\"a&a%a*a$a#a~a&a\"",
		// Not supported characters
		"°",
		"°°",
		"a°a",
		"+°+",
		"%°%",
		// Operators
		"=",  // ASSIGN
		"+",  // PLUS
		"-",  // MINUS
		"!",  // BANG
		"*",  // ASTERISK
		"/",  // SLASH
		"<",  // LT
		">",  // GT
		"==", // EQ
		"!=", // NOT_EQ
		"+-!*/<>==!=",
		// Delimiters
		",", // COMMA
		";", // SEMICOLON
		":", // COLON
		"(", // LPAREN
		")", // RPAREN
		"{", // LBRACE
		"}", // RBRACE
		"[", // LBRACKET
		"]", // RBRACKET
		",;:(){}[]",
	}

	content := "\\begin{itemize}"

	for _, input := range tests {
		translation, containsUnExpChars := teXify(input)
		content = content + "\n\\item " + translation

		if containsUnExpChars {
			fmt.Println(input, "\t", translation) // the go in the action's environment seems to dislike t.Log
		}
	}
	content = content + "\n\\end{itemize}"

	document := makeStandalone(content)

	err := tex2pdf(document, "tests/texify_all.pdf", latexPath)

	if err != nil {
		t.Errorf("Rendering did not succeed. Reason: %q", err)
		t.Error(document)

	}
}

func Test_visRunIndent(t *testing.T) {
	prefix := "|| "
	indent := " | "

	v := NewVisRun(
		prefix,
		indent,
		0,
		CONSOLE,
		PARSE,
		false,
		false,
		false,
	)
	content := []string{
		"1st line",
		"2nd line",
		"3rd line",
		"4th line",
	}

	v.testVisRunIndent(content)
	fmt.Println(v.out.String())

}

func (v *visRun) testVisRunIndent(content []string) {
	if len(content) == 0 {
		return
	}

	v.printInd("[ ")
	v.printW(content[0])
	content = content[1:]
	v.incrIndent()
	v.testVisRunIndent(content)
	v.decrIndent()
	v.printInd("] ")

}

func Test_envTableTeX(t *testing.T) {

	env := object.NewEnvironment()
	setup := []string{
		"let a = 1",
		"let b = true",
		"let a_b = fn(x){x}",
	}

	for _, input := range setup {
		l := lexer.New(input)
		p := parser.New(l)
		ast := p.ParseProgram()
		evaluator.Eval(ast, env)
	}

	content := texEnvTables(env, 0, true)
	document := makeStandalone(content)

	err := tex2pdf(document, "tests/env.pdf", latexPath)
	//t.Errorf(document)
	if err != nil {
		t.Errorf("Rendering did not succeed. Reason: %q", err)
	}

}

func Test_TeXParseTree(t *testing.T) {

	//at least one example for every nodetype, might be within a bigger example
	/*
		nil + isNil +

		Program

		LetStatement
		ReturnStatement
		ExpressionStatement
		BlockStatement

		Identifier
		Boolean
		IntegerLiteral
		PrefixExpression
		InfixExpression
		IfExpression
		FunctionLiteral
		CallExpression

		--- chap 4:
		StringLiteral
		ArrayLiteral
		IndexExpression
		HashLiteral

	*/
	tests := []struct {
		input string
		file  string
	}{
		{"@", "nil"},                         //nil
		{"if(1){}", "isNil"},                 // empty alternative - Nil value
		{"fn(){}", "function_with_0_params"}, // + empty block statement
		{"fn(x){x}", "function_with_1_params"},
		{"fn(x,y){x+y}", "function_with_2_params"},
		{"double(a)", "identifiers-call"},
		{"!true", "bang"},
		{"let a = if(a>2){}", "let-if"},
		{"return if(false){} else {}", "return-if"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		node := p.ParseProgram()
		file := fmt.Sprintf("tests/p_%v.pdf", tt.file)
		err := TeXParseTree(tt.input, node, 0, false, file, latexPath)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_TeXParseTree_inclToken_verb(t *testing.T) {

	tests := []struct {
		input string
		file  string
	}{
		{"@", "nil"},                         //nil
		{"if(1){}", "isNil"},                 // empty alternative - Nil value
		{"fn(){}", "function_with_0_params"}, // + empty block statement
		{"fn(x){x}", "function_with_1_params"},
		{"fn(x,y){x+y}", "function_with_2_params"},
		{"double(a)", "identifiers-call"},
		{"!true", "bang"},
		{"let a = if(a>2){}", "let-if"},
		{"return if(false){} else {}", "return-if"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)
		node := p.ParseProgram()
		for verbosity := 0; verbosity < 3; verbosity++ {
			file := fmt.Sprintf("tests/p_%v_tok_%v.pdf", tt.file, verbosity)
			err := TeXParseTree(tt.input, node, verbosity, true, file, latexPath)
			if err != nil {
				t.Error(err)
			}
		}
	}
}

/*
Object types:

	NULL_OBJ  = "NULL"
	ERROR_OBJ = "ERROR"

	INTEGER_OBJ = "INTEGER"
	BOOLEAN_OBJ = "BOOLEAN"
	STRING_OBJ  = "STRING"

	RETURN_VALUE_OBJ = "RETURN_VALUE"

	FUNCTION_OBJ = "FUNCTION"

	---
	chap4:
	BUILTIN_OBJ  = "BUILTIN"

	ARRAY_OBJ = "ARRAY"
	HASH_OBJ  = "HASH"


*/
func Test_TeXEvalTree_Objects(t *testing.T) {

	tests := []struct {
		setup string
		input string
		file  string
	}{
		{"", "if(1){}", "isNil"},                 // empty alternative - Nil value
		{"", "fn(){}", "function_with_0_params"}, // + empty block statement
		{"", "fn(x){x}", "function_with_1_params"},
		{"", "fn(x,y){x+y}", "function_with_2_params"},
		{"let double = fn(x){2*x} let a= 3", "double(a)", "identifiers-call"},
		{"", "!true", "bang-bool"},
		{"", "let a = if(3>2){1}", "let-if-nil"},
		{"", "return if(false){}", "return-if-null"},
		{"let v_int = 1; let v_bool = true; let v_null = if(false){}; let v_nil = if(true){}; let v_error = 1 + true",
			"v_int; v_bool; v_null; v_nil; v_error", "some_objects"},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		//setup
		l_setup := lexer.New(tt.setup)
		p_setup := parser.New(l_setup)
		node_setup := p_setup.ParseProgram()
		evaluator.EvalT(node_setup, env, true)
		//input
		l := lexer.New(tt.input)
		p := parser.New(l)
		node := p.ParseProgram()
		_, trace := evaluator.EvalT(node, env, true)

		file := fmt.Sprintf("tests/e_%v.pdf", tt.file)
		err := TeXEvalTree(tt.input, trace, 0, false, false, false, file, latexPath)
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_TeXEvalTree_goObjType_verbosity(t *testing.T) {

	tests := []struct {
		setup string
		input string
		file  string
	}{
		{"", "if(1){}", "isNil"},                 // empty alternative - Nil value
		{"", "fn(){}", "function_with_0_params"}, // + empty block statement
		{"", "fn(x){x}", "function_with_1_params"},
		{"", "fn(x,y){x+y}", "function_with_2_params"},
		{"let double = fn(x){2*x} let a= 3", "double(a)", "identifiers-call"},
		{"", "!true", "bang-bool"},
		{"", "let a = if(3>2){1}", "let-if-nil"},
		{"", "return if(false){}", "return-if-null"},
		{"let v_int = 1; let v_bool = true; let v_null = if(false){}; let v_nil = if(true){}; let v_error = 1 + true",
			"v_int; v_bool; v_null; v_nil; v_error", "some_objects"},
	}

	for _, tt := range tests {
		env := object.NewEnvironment()
		//setup
		l_setup := lexer.New(tt.setup)
		p_setup := parser.New(l_setup)
		node_setup := p_setup.ParseProgram()
		evaluator.EvalT(node_setup, env, true)
		//input
		l := lexer.New(tt.input)
		p := parser.New(l)
		node := p.ParseProgram()
		_, trace := evaluator.EvalT(node, env, true)
		for verbosity := 0; verbosity < 3; verbosity++ {
			file := fmt.Sprintf("tests/e_%v_got_%v.pdf", tt.file, verbosity)
			err := TeXEvalTree(tt.input, trace, verbosity, false, true, false, file, latexPath)
			if err != nil {
				t.Error(err)
			}
		}
	}
}
