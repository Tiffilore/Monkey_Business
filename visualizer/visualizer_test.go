package visualizer

import "testing"

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
