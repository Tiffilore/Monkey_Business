package vis2

import (
	"fmt"
	"io/ioutil"
	"log"
	"monkey/ast"
	"monkey/evaluator"
	"monkey/object"

	"github.com/rwestlund/gotex"
)

func Ast2pdf(node ast.Node, filename string, path string, verb Verbosity, exclToken bool) error {
	vis := NewVisualizer("", "  ", verb, exclToken)
	qtreenode := vis.VisualizeQTree(node)
	//fmt.Println(qtreenode)
	document := makeTeX(qtreenode)
	return tex2pdf(document, filename, path)
}

func Eval2pdf(node ast.Node, filename string, path string, verb Verbosity, exclToken bool) error {
	vis := NewVisualizer("", "  ", verb, exclToken)
	env := object.NewEnvironment()
	evaluator.StartTracer()
	evaluator.Eval(node, env)
	evaluator.StopTracer()

	qtreenode := vis.VisualizeEvalQTree(evaluator.T)

	fmt.Println("hier:\n", qtreenode)
	document := makeTeX(qtreenode)
	return tex2pdf(document, filename, path)
}

// func EvalTree2pdf(t *evaluator.Tracer, filename string, path string, brevity int) {

// 	evalqtreenode := QTreeEval(t, brevity)
// 	document := makeTeX(evalqtreenode)
// 	tex2pdf(document, filename, path)
// }

func tex2pdf(document string, filename string, path string) error {

	var pdf, err = gotex.Render(document, gotex.Options{
		Command: path,
		Runs:    1})

	if err != nil {
		log.Println("render failed ", err)
		return err
	} else {
		err := ioutil.WriteFile(filename, pdf, 0664)
		if err != nil {
			log.Println(" failed ", err)
			return err
		}
	}
	return nil
}

var document_prefix = `
\documentclass[border=0.2cm]{standalone}
\usepackage{xcolor}
%\usepackage{qtree}
\usepackage{tikz}
\usepackage{tikz-qtree}

\definecolor{bluish}{HTML}{E0EBF5}
\definecolor{yellish}{HTML}{FFFFA8}
\definecolor{dbluish}{HTML}{375EAB}
\definecolor{brownish}{HTML}{BC8C64}
\definecolor{darkish}{HTML}{6F6B69}
\definecolor{darkish2}{HTML}{848475}

%047C9C dark go-blue


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
