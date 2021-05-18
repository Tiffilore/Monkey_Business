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

	inputs := []string{
		"1",
		"true",
		"fn(x){x}",
		"@",
	}

	for i, input := range inputs {
		l := lexer.New(input)
		p := parser.New(l)
		node := p.ParseProgram()
		file := fmt.Sprintf("tests/p_%v.pdf", i)
		err := TeXParseTree(node, 0, false, file, latexPath)
		//t.Error("aha")
		if err != nil {
			t.Error(err)
		}
	}
}

func Test_TeXEvalTree(t *testing.T) {

	inputs := []string{
		"1",
		"true",
		"fn(x){x}",
		// "@",
	}

	for i, input := range inputs {
		l := lexer.New(input)
		p := parser.New(l)
		node := p.ParseProgram()
		env := object.NewEnvironment()
		_, trace := evaluator.EvalT(node, env, true)
		file := fmt.Sprintf("tests/e_%v.pdf", i)
		err := TeXEvalTree(trace, 0, false, false, false, file, latexPath)
		//t.Error("aha")
		if err != nil {
			t.Error(err)
		}
	}
}