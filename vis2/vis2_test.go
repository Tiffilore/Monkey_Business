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

var inputs = []string{ // for parsing
	"",    //empty program
	"@",   //nil
	"let", //nil value
	"1",
	"1 2",
	"let a = 5",
	"return true",
	"{}", //empty BlkStmt (level s)
	"if(false){}",
	"if(1){} if(2){}", // 2 nil-Alternatives
	"if(true){a}",
	"!a",
	"if(1){2+3}",
	"fn(){}",
	"fn(x,y){}",
	"id()",
	"fn(+,-,!,/,*,==,=,!=,<,>){}", // test operator display as long as it is not fixed
}

var inputs_objects = []string{
	"1", //int
	// "true",            //TRUE
	// "false",           //FALSE
	// "fn(){}",          //function
	"if(1){}", //nil
	// "if(!1){}",        //NULL
	// "return 1",        //return int
	// "return true",     //return TRUE
	// "return fn(){}",   //return function
	// "return if(1){}",  //return nil
	// "return if(!1){}", //return NULL

	// "1+true",       //error
	// "true + false", //error
	// "id(2)",        //error
	// "1(1)",         //error
	//...
}

var inputs_closures = []string{
	"let cl_m = fn(x){fn(y){x+y}}",
	"let cl_m = fn(x){fn(y){x+y}}; let cl = cl_m(2)",
	"let cl_m = fn(x){fn(y){x+y}}; let cl = cl_m(2); cl(3)",
	"let cl_m = fn(x){fn(y){x+y}}; let cl = cl_m(2); let cl_ = cl_m(2); cl(3)+cl_(4)",
}

func Test_eval_tex(t *testing.T) {
	exclToken := true
	// exclEnv

	for index, input := range inputs_closures { //inputs_objects { //
		// if index != 0 {
		// 	continue
		// }
		for _, level := range []string{"p", "s", "e"} {
			if level != "p" {
				continue
			}
			l := lexer.New(input)
			p := parser.New(l)

			node := parse_level(p, level)
			for _, err := range p.Errors() {
				t.Error(err)
			}

			for _, verb := range []Verbosity{V, VV, VVV} {
				if verb != V {
					continue
				}

				path := "/usr/bin/pdflatex"
				file := "test/test" + fmt.Sprint(index) + level + fmt.Sprint(verb) + ".pdf"
				err := Eval2pdf(node, file, path, verb, exclToken)
				if err != nil {
					t.Errorf("Problem for %v level %v", input, level)
				} else {
					t.Errorf("a")
				}

			}
		}
	}
}

func _Test_ast_cons(t *testing.T) {

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

func _Test_ast_tex(t *testing.T) {

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
