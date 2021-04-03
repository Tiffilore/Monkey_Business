package visualizer

import (
	"io/ioutil"
	"log"
	"monkey/ast"
	"monkey/evaluator"

	"github.com/rwestlund/gotex"
)

func Ast2pdf(node ast.Node, excltoken bool, filename string, path string) {

	qtreenode := QTree(node, excltoken)
	document := makeTeX(qtreenode)
	tex2pdf(document, filename, path)
}

func EvalTree2pdf(t *evaluator.Tracer, filename string, path string, brevity int) {

	evalqtreenode := QTreeEval(t, brevity)
	document := makeTeX(evalqtreenode)
	tex2pdf(document, filename, path)
}

func tex2pdf(document string, filename string, path string) {

	var pdf, err = gotex.Render(document, gotex.Options{
		Command: path,
		Runs:    1})

	if err != nil {
		log.Println("render failed ", err)
	} else {
		err := ioutil.WriteFile(filename, pdf, 0664)
		if err != nil {
			log.Println(" failed ", err)
		}
	}
}

var document_prefix = `
\documentclass[border=0.2cm]{standalone}
\usepackage{xcolor}
%\usepackage{qtree}
\usepackage{tikz}
\usepackage{tikz-qtree}

\definecolor{yellish}{HTML}{E0EBF5}
\definecolor{bluish}{HTML}{FFFFA8}
\definecolor{dbluish}{HTML}{375EAB}

\begin{document}

`

var document_suffix = `
\end{document}
`

var tikz_prefix = `
\begin{tikzpicture}[
   every tree node/.style={anchor=north},
   every node/.append style={align=left}  
]

`

var tikz_suffix = `
\end{tikzpicture}
`

func makeTeX(qtreenode string) string {
	return document_prefix + tikz_prefix + qtreenode + tikz_suffix + document_suffix
}
