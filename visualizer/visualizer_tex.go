package visualizer

import (
	"io/ioutil"
	"log"

	"github.com/rwestlund/gotex"
)

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

func makeStandalone(content string) string {
	return document_prefix + content + document_suffix
}

func makeTikz(qtree string) string {
	return tikz_prefix + qtree + tikz_suffix
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
