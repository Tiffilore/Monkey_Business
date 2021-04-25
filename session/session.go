package session

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"monkey/ast"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"monkey/visualizer"
	"os"
	"os/exec"
	"reflect"
	"strings"
)

func Start(in io.Reader, out io.Writer) error {

	s, err := NewSession(in, out)
	if err != nil {
		return err
	}

	for {
		fmt.Fprint(out, currentSettings.prompt) // Fprint instead of Fprintf due to SA1006

		scanned := s.scanner.Scan()
		if !scanned {
			return errors.New("problem with scanning")
		}

		line := s.scanner.Text()
		s.exec_cmd(line)
	}
}

type Session struct {
	scanner       *bufio.Scanner
	out           io.Writer
	environment   *object.Environment
	path_pdflatex string
}

// NewSession creates a new Session.
func NewSession(in io.Reader, out io.Writer) (*Session, error) {

	path, err := exec.LookPath("pdflatex")
	if err != nil {
		path = ""
	}

	s := &Session{
		scanner:       bufio.NewScanner(in),
		out:           out,
		environment:   object.NewEnvironment(),
		path_pdflatex: path,
	}

	if err := s.init_commands(); err != nil {
		return nil, err
	}
	return s, nil
}

// decide which function
func (s *Session) exec_cmd(line string) {
	if !strings.HasPrefix(line, ":") {
		s.exec_process(line)
		return
	}
	line = strings.TrimPrefix(line, ":")
	slice := strings.SplitN(line, " ", 2)

	cmd := slice[0]
	if len(slice) == 1 {
		if exec, ok := commands.get_exec_single(cmd); ok {
			exec()
			return
		}
	} else {
		arg := slice[1]
		if exec, ok := commands.get_exec_with_arg(cmd); ok {
			exec(arg)
			return
		}
	}

	s.exec_help(cmd)
}

// quit
func (s *Session) exec_quit() {
	os.Exit(0)
}

// clear the screen

func (s *Session) exec_clearscreen() func() {
	if _, err := exec.LookPath("clear"); err != nil {
		return func() {
			fmt.Fprintln(s.out, "command clearscreen is not available for you")
		}
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = s.out
		return func() {
			err := cmd.Run()
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

// environment
func (s *Session) exec_clear() {
	s.environment = object.NewEnvironment()
}

func (s *Session) exec_list() {

	table := visualizer.GetStoreTable(s.environment)
	lines := strings.Split(table, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintln(s.out, line)
		}
	}
}

// commands
func (s *Session) exec_help(cmd string) {

	if usage, ok := commands.usage(cmd); ok {
		fmt.Fprint(s.out, usage)
		return
	}

	fmt.Fprintln(s.out, "unknown command: ", cmd)
}

func (s *Session) exec_help_all() {
	fmt.Fprint(s.out, commands.menu())
}

// input processing

func (s *Session) exec_parsetree(line string) {
	fmt.Fprint(s.out, "not yet implemented!\n")
}

func (s *Session) exec_evaltree(line string) {
	fmt.Fprint(s.out, "not yet implemented!\n")
}

func (s *Session) exec_process(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, currentSettings.process, false, line)
}

func (s *Session) exec_paste_empty_arg() {
	s.process_input_dim(true, currentSettings.level, currentSettings.process, false, "")
}

func (s *Session) exec_paste(line string) {
	s.process_input_dim(true, currentSettings.level, currentSettings.process, false, line)
}

func (s *Session) exec_expression(line string) {
	s.process_input_dim(currentSettings.paste, ExpressionL, currentSettings.process, false, line)
}

func (s *Session) exec_statement(line string) {
	s.process_input_dim(currentSettings.paste, StatementL, currentSettings.process, false, line)
}

func (s *Session) exec_program(line string) {
	s.process_input_dim(currentSettings.paste, ProgramL, currentSettings.process, false, line)
}

func (s *Session) exec_eval(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, EvalP, false, line)
}

func (s *Session) exec_type(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, TypeP, false, line)
}

func (s *Session) exec_trace(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, EvalP, true, line)
}

func (s *Session) exec_parse(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, ParseP, false, line)
}

func (s *Session) process_input_dim(paste bool, level InputLevel, process InputProcess, trace bool, input string) {

	if paste {
		input = s.multiline_input(input)
	}

	l := lexer.New(input)
	p := parser.New(l)

	node := parse_level(p, level)

	if currentSettings.logtree {
		fmt.Fprint(s.out, "log ast:\t")
	}
	if currentSettings.logtree || process == ParseP {
		fmt.Fprintln(s.out, node)
		fmt.Fprintln(s.out, visualizer.RepresentNodeConsoleTree(node, "|   ", !currentSettings.inclToken))
		//	fmt.Fprintln(s.out, visualizer.QTree(node, !s.incltoken))
		path, err := exec.LookPath("pdflatex")
		if err != nil {
			fmt.Fprintln(s.out, "Displaying trees as pdfs is not available to you, since you have not installed pdflatex.")
		} else {
			visualizer.Ast2pdf(node, !currentSettings.inclToken, currentSettings.file, path)
		}
	}

	if len(p.Errors()) != 0 {
		s.printParserErrors(p.Errors(), level)
		return
	}

	if process == ParseP {
		return
	}

	if trace || currentSettings.logtrace {
		evaluator.StartTracer()
	}

	evaluated := evaluator.Eval(node, s.environment)

	if trace || currentSettings.logtrace {
		evaluator.StopTracer()
	}

	if currentSettings.logtype {
		fmt.Fprint(s.out, "log type:\t")
	}
	if currentSettings.logtype || process == TypeP {
		fmt.Fprintln(s.out, reflect.TypeOf(evaluated))
	}
	if process == TypeP {
		return
	}

	if currentSettings.logtrace {
		visualizer.RepresentEvalConsole(evaluator.T, s.out)
		//fmt.Fprint(s.out, visualizer.QTreeEval(evaluator.T))
		path, err := exec.LookPath("pdflatex")
		if err != nil {
			fmt.Fprintln(s.out, "Displaying evaluation trees as pdfs is not available to you, since you have not installed pdflatex.")
		} else {
			visualizer.EvalTree2pdf(evaluator.T, currentSettings.file, path)
		}

	}

	if trace {
		visualizer.TraceEvalConsole(evaluator.T, s.out, s.scanner)
		return
	}

	if evaluated != nil { // TODO: Umgang mit nil würdig?
		fmt.Fprintln(s.out, evaluated.Inspect())
	}

	//	} else {
	//		fmt.Fprintln(s.out, nil)
	//	}
	/*
		io.WriteString vs fmt.Fprint ?????
			The difference is that fmt.Fprint is formatting the arguments provided first in a buffer before calling w.Write.
			And io.WriteString is checking if w provides the StringWriter interface and calls that instead.
	*/
}

func parse_level(p *parser.Parser, level InputLevel) ast.Node {

	switch level {
	case ExpressionL:
		return p.ParseExpression()
	case StatementL:
		return p.ParseStatement()
	case ProgramL:
		return p.ParseProgram()
	default:
		return nil
	}
}

func (s *Session) multiline_input(input string) string {
	for {
		scanned := s.scanner.Scan()
		if !scanned {
			return input //TODO!!
		}
		line := s.scanner.Text()
		if line == "" {
			return input
		}
		input += " " + line
	}
}
func (s *Session) printParserErrors(errors []string, level InputLevel) {

	fmt.Fprintf(s.out, "... cannot be parsed as %v\n", level)
	//io.WriteString(s.out, " parser errors:\n")
	for _, msg := range errors {
		fmt.Fprintf(s.out, "\t%v\n", msg)
	}
}
