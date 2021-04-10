package vis2

import (
	"encoding/json"
	"fmt"
	"monkey/ast"
	"monkey/lexer"
	"monkey/parser"
	"testing"
)

// teste ob tex compiliert
// teste ob alles bisher implementiert
//

var inputs = []string{
	"if(1){} if(2){} if(3){}", // 2 nil-Alternatives
	"",                        //empty program
	"@",                       //nil
	"let",                     //nil value
	"1",
	"1 2",
	"let a = 5",
	"return true",
	"{}", //empty BlkStmt (level s)
	"if(true){}",
	"if(1){} if(2){}", // 2 nil-Alternatives
	"if(true){a}",
	"!a",
	"if(1){2+3}",
	"fn(){}",
	"fn(x,y){}",
	"id()",
	"fn(+,-,!,/,*,==,=,!=,<,>){}", // test operator display as long as it is not fixed
}

func _Test_cons(t *testing.T) {

	exclToken := false

	for index, input := range inputs {
		if index > 0 {
			continue
		}
		levels := []string{"p", "s", "e"}

		for _, level := range levels {
			if level != "p" {
				continue
			}
			l := lexer.New(input)
			p := parser.New(l)
			node := parse_level(p, level)
			// for _, err := range p.Errors() {
			// 	t.Log(err)
			// }
			verbs := []Verbosity{V, VV, VVV}

			for _, verb := range verbs {
				if verb != V {
					continue
				}
				vis := NewVisualizer("", "|   ", verb, exclToken)
				fmt.Println(vis.VisualizeConsTree(node))
				//vis.VisualizeConsTree(node)

			}
		}
	}
}
func Test_tex(t *testing.T) {

	exclToken := true

	for index, input := range inputs {
		if index > 0 {
			continue
		}
		levels := []string{"p", "s", "e"}

		for _, level := range levels {
			l := lexer.New(input)
			p := parser.New(l)
			node := parse_level(p, level)
			for _, err := range p.Errors() {
				t.Error(err)
			}
			verbs := []Verbosity{V, VV, VVV}

			for _, verb := range verbs {
				//t.Error(prettyPrint(node))

				path := "/usr/bin/pdflatex"
				file := "test/test" + fmt.Sprint(index) + level + fmt.Sprint(verb) + ".pdf"
				err := Ast2pdf(node, file, path, verb, exclToken)
				if err != nil {
					t.Errorf("Problem for %v level %v", input, level)
				}
				//	vis := NewVisualizer("", "   ", verb, exclToken)
				//	fmt.Println(vis.VisualizeQTree(node))
				//	time.Sleep(3 * 1000 * 1000 * 1000)
			}
		}
	}
}
func _Test_vis2(t *testing.T) {

	inputs := []string{
		" ",
		"@",
		"let",
		"1",
		"1 2",
		"let a = 5",
		"return true",
		"{}",
		"if(true){}",
		"if(true){a}",
		"!a",
		"if(1){2+3}",
		"fn(){}",
		"fn(1,2){}",
		"id()",
	}
	_ = inputs
	inputs2 := []string{
		"1+1",
	}

	verb := VV
	for _, input := range inputs2 {
		test_level(input, t, "e", verb)
		//test_level(input, t, "s", verb)
		//test_level(input, t, "p", verb)
	}

}
func prettyPrint(i interface{}) string {
	s, _ := json.MarshalIndent(i, " ", " ")
	return string(s)
}
func test_level(input string, t *testing.T, level string, verb Verbosity) {

	l := lexer.New(input)
	p := parser.New(l)

	node := parse_level(p, level)

	//t.Error(prettyPrint(node))

	vis := NewVisualizer("", "   ", verb, true)

	//representation, _ := vis.VisP(node)
	//fmt.Print(prettyPrint(node))

	vis.visualizeNode(node)
	representation := vis.out.String()
	t.Error("\n", representation)
}

func parse_level(p *parser.Parser, level string) ast.Node {

	switch level {
	case "e", "expression":
		return p.ParseExpression()
	case "s", "statement":
		return p.ParseStatement()
	case "p", "program":
		return p.ParseProgram()
	default:
		return nil
	}
}
