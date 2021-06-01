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
	"strings"
)

const (
	prefixCons = ""    //prefix for trees in console
	indentCons = "   " //indentation for trees in console
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
		s.exec_process(line) //default
		return
	}
	line = strings.TrimPrefix(line, ":")
	slice := strings.SplitN(line, " ", 2)

	cmd := slice[0]
	var arg string
	if len(slice) == 1 {
		if exec, ok := commands.get_exec_single(cmd); ok {
			exec()
			return
		} else {
			arg = ""
		}
	} else {
		arg = slice[1]
	}
	if exec, ok := commands.get_exec_with_arg(cmd); ok {
		exec(arg)
		return
	}
	s.exec_help(cmd)
}

// quit
func (s *Session) exec_quit() {
	os.Exit(0)
}

// clear the screen

func (s *Session) f_exec_clearscreen() func() {
	if _, err := exec.LookPath("clear"); err != nil {
		return func() {
			fmt.Fprintln(s.out, "command clearscreen is not available for you")
		}
	} else {
		return func() {
			cmd := exec.Command("clear")
			cmd.Stdout = s.out
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

	table := visualizer.VisEnvStoreCons(s.environment, currentSettings.verbosity, currentSettings.goObjType)
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

// settings

func (s *Session) exec_settings() {
	fmt.Fprint(s.out, menuSettings())
}

func (s *Session) exec_reset_all() {
	currentSettings = newSettings()
}

func (s *Session) exec_reset(input string) {
	if ok := reset(input); !ok {
		s.exec_help("reset")
	}
}

func (s *Session) exec_set(input string) {
	if ok := set(input); !ok {
		s.exec_help("set")
	}
}

func (s *Session) exec_unset(input string) {
	if ok := unset(input); !ok {
		s.exec_help("unset")
	}
}

// input processing

func (s *Session) exec_process(line string) { // if no command is used
	s.process_input_dim(currentSettings.paste, currentSettings.level, currentSettings.process, line)
}

// PASTE
func (s *Session) exec_paste(line string) {
	s.process_input_dim(true, currentSettings.level, currentSettings.process, line)
}

// LEVEL
func (s *Session) exec_expression(line string) {
	s.process_input_dim(currentSettings.paste, ExpressionL, currentSettings.process, line)
}

func (s *Session) exec_statement(line string) {
	s.process_input_dim(currentSettings.paste, StatementL, currentSettings.process, line)
}

func (s *Session) exec_program(line string) {
	s.process_input_dim(currentSettings.paste, ProgramL, currentSettings.process, line)
}

// PROCESS

func (s *Session) exec_parse(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, ParseP, line)
}

func (s *Session) exec_parsetree(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, ParseTreeP, line)
}

func (s *Session) exec_eval(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, EvalP, line)
}

func (s *Session) exec_type(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, TypeP, line)
}

func (s *Session) exec_trace(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, TraceP, line)
}

func (s *Session) exec_evaltree(line string) {
	s.process_input_dim(currentSettings.paste, currentSettings.level, EvalTreeP, line)
}

// input processing
func (s *Session) process_input_dim(paste bool, level inputLevel, process inputProcess, input string) {

	// get input dependent on PASTE
	if paste {
		input = s.multiline_input(input)
	}

	// determine what needs to be logged
	logPtree := false
	logType := false
	logTrace := false
	logEtree := false

	if process == ParseP {
		logPtree = currentSettings.logs[ParseTreeP]
	}
	if process == EvalP {
		logType = currentSettings.logs[TypeP]
		logTrace = currentSettings.logs[TraceP]
		logPtree = currentSettings.logs[ParseTreeP]
		logEtree = currentSettings.logs[EvalTreeP]
	}

	// parse input dependent on LEVEL
	l := lexer.New(input)
	p := parser.New(l)

	node := parse_level(p, level)

	// PROCESS / [ LOGs ]

	if process == ParseTreeP || logPtree {
		if currentSettings.displays[ConsD] {
			if process != ParseTreeP {
				fmt.Fprint(s.out, "log parsetree:\n")
			}
			consPtree := visualizer.ConsParseTree(
				node,
				currentSettings.verbosity,
				currentSettings.inclToken,
				prefixCons,
				indentCons,
			)
			//fmt.Fprintln(s.out, "display ptree in console: ")
			fmt.Fprintln(s.out, consPtree)

		}
		if currentSettings.displays[PdfD] {
			if !s.supportsPdflatex() {
				fmt.Fprintln(s.out, "Displaying trees as pdfs is not available to you, since you have not installed pdflatex.")
			} else {
				err := visualizer.TeXParseTree(input, node, currentSettings.verbosity, currentSettings.inclToken, currentSettings.pfile, s.path_pdflatex)
				if err != nil {
					fmt.Fprintln(s.out, err)
				} else {
					fmt.Fprintf(s.out, "parsetree is printed to %v\n", currentSettings.pfile)

				}
			}
		}
	}

	if process == ParseP {
		fmt.Fprintln(s.out, node) // Stringer method
	}

	if len(p.Errors()) != 0 {
		s.printParserErrors(p.Errors(), level)
		return
	}

	if process == ParseP || process == ParseTreeP { // in these cases, we do not care about logging related to evaluation
		return
	}

	// evaluate ast - trace dependent on process + DISPLAYED logs

	trace_required := false
	if process == TraceP || process == EvalTreeP || logTrace || logEtree {
		trace_required = true
	}

	obj, trace := s.eval_process(node, trace_required)

	if process == TraceP {
		visualizer.TraceInteractive(trace, s.out, s.scanner, currentSettings.verbosity, currentSettings.goObjType)
		return // no additional evaluation logging !
	}

	if logTrace {
		visualizer.TraceTable(trace, s.out, currentSettings.verbosity, currentSettings.goObjType)
	}

	if process == TypeP || logType {

		if process != TypeP {
			fmt.Fprint(s.out, "log type:\t")
		}
		fmt.Fprintln(s.out, visualizer.VisObjectType(obj, currentSettings.verbosity, currentSettings.goObjType))
	}

	if process == EvalTreeP || logEtree {

		if currentSettings.displays[ConsD] {
			if process != EvalTreeP {
				fmt.Fprint(s.out, "log evaltree:\n")
			}
			consEtree := visualizer.ConsEvalTree(
				trace,
				currentSettings.verbosity,
				currentSettings.inclToken,
				currentSettings.goObjType,
				currentSettings.inclEnv,
				prefixCons,
				indentCons,
			)
			//fmt.Fprintln(s.out, "display etree in console")
			fmt.Fprintln(s.out, consEtree)

		}
		if currentSettings.displays[PdfD] {
			if !s.supportsPdflatex() {
				fmt.Fprintln(s.out, "Displaying trees as pdfs is not available to you, since you have not installed pdflatex.")
			} else {
				err := visualizer.TeXEvalTree(
					input,
					trace,
					currentSettings.verbosity,
					currentSettings.inclToken,
					currentSettings.goObjType,
					currentSettings.inclEnv,
					currentSettings.efile,
					s.path_pdflatex)
				if err != nil {
					fmt.Fprintln(s.out, err)
				} else {
					fmt.Fprintf(s.out, "evaltree is printed to %v\n", currentSettings.efile)
				}

			}
		}

		if process == EvalTreeP {
			return
		}
	}

	if process == EvalP {
		if obj != nil { // TODO: Umgang mit nil würdig?
			fmt.Fprintln(s.out, obj.Inspect())
		}
		// } else {
		// 	fmt.Fprintln(s.out, nil)
		// }
	}

}

func parse_level(p *parser.Parser, level inputLevel) ast.Node {
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

func (s *Session) eval_process(node ast.Node, trace_required bool) (object.Object, *evaluator.Trace) {

	return evaluator.EvalT(node, s.environment, trace_required)

}

func (s *Session) supportsPdflatex() bool {
	return s.path_pdflatex != ""
}

func (s *Session) multiline_input(input string) string {
	for {
		fmt.Fprint(s.out, "... ") // secondary prompt

		scanned := s.scanner.Scan()
		if !scanned {
			return input //TODO!! when can that happen anyway?
		}
		line := s.scanner.Text()
		if line == "" {
			return input
		}
		input += " " + line
	}
}

func (s *Session) printParserErrors(errors []string, level inputLevel) {

	fmt.Fprintf(s.out, "... cannot be parsed as %v\n", level)
	//io.WriteString(s.out, " parser errors:\n")
	for _, msg := range errors {
		fmt.Fprintf(s.out, "\t%v\n", msg)
	}
}
