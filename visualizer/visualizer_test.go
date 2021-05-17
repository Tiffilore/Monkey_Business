package visualizer

import (
	"fmt"
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
