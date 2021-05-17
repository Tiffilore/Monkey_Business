package visualizer

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strings"

	"github.com/rwestlund/gotex"
)

/*
the functions in this file implement teX-specific functions
*/

func tex2pdf(document string, filename string, path string) error {

	var pdf, err = gotex.Render(document, gotex.Options{
		Command: path,
		Runs:    1})

	if err != nil {
		log.Println("render failed ", err)
		return err
	} else {
		err := ioutil.Writefile(filename, pdf, 0664)
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

func teXify(input string) (string, bool) {

	containsUnExpChars := false
	// replace all characters we provide no solution for

	re := regexp.MustCompile(`[^a-zA-Z0-9_\"=+\-!*//<>,;:(){}\[\]&%*$#~]`)
	replacement := re.ReplaceAllString(input, "\\textdagger") //

	//strings.Count("cheese", "e"))
	if input != replacement {
		containsUnExpChars = true
		input = replacement
	}
	// replace certain characters sequences by something teXy

	/// LaTeX special characters
	input = strings.ReplaceAll(input, "{", "\\{")
	input = strings.ReplaceAll(input, "}", "\\}")
	input = strings.ReplaceAll(input, "\\{", "\\{{}")
	input = strings.ReplaceAll(input, "\\}", "\\}{}")
	input = strings.ReplaceAll(input, "\\textdagger", "\\textdagger{}")
	input = strings.ReplaceAll(input, "&", "\\&{}")
	input = strings.ReplaceAll(input, "%", "\\%{}")
	input = strings.ReplaceAll(input, "$", "\\${}")
	input = strings.ReplaceAll(input, "#", "\\#{}")
	input = strings.ReplaceAll(input, "_", "\\_{}")
	input = strings.ReplaceAll(input, "[", "$[${}")
	input = strings.ReplaceAll(input, "]", "$]${}")
	input = strings.ReplaceAll(input, "~", "\\textasciitilde{}")
	input = strings.ReplaceAll(input, "^", "\\textasciicircum{}")
	input = strings.ReplaceAll(input, "\"", "\\textquotedbl{}")

	/// operators
	input = strings.ReplaceAll(input, "<", "$<${}")
	input = strings.ReplaceAll(input, ">", "$>${}")
	input = strings.ReplaceAll(input, "!", "$!${}")
	input = strings.ReplaceAll(input, "=", "$=${}")
	input = strings.ReplaceAll(input, "+", "$+${}")
	input = strings.ReplaceAll(input, "-", "$-${}")
	input = strings.ReplaceAll(input, "*", "$*${}")
	input = strings.ReplaceAll(input, "/", "$/${}")
	return input, containsUnExpChars
}

func roofify(str string) string {
	return fmt.Sprint("\\edge[roof];{", str, "}")
}

func texColorize(str, bcolor, tcolor string) string { //TODO: includes also \tt
	return "\\colorbox{" + bcolor + "}{\\textcolor{" + tcolor + "}{\\tt " + str + "}}"
}

var document_prefix = `
\documentclass[varwidth, border=0.2cm]{standalone} %varwidth for itemize
\usepackage[T1]{fontenc}
\usepackage[utf8]{inputenc}
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
