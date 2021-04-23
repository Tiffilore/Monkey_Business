package session

import (
	"bufio"
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

	"github.com/jedib0t/go-pretty/v6/table"
)

func Start(in io.Reader, out io.Writer) {

	s := NewSession(in, out)

	for {

		fmt.Fprint(out, s.prompt) // Fprint instead of Fprintf due to SA1006

		scanned := s.scanner.Scan()
		if !scanned {
			return
		}

		line := s.scanner.Text()
		s.exec_cmd(line)
	}
}

// TODO: type Settings

type Session struct {
	scanner       *bufio.Scanner
	out           io.Writer
	environment   *object.Environment
	path_pdflatex string
	//
	prompt string
	//
	process InputProcess
	level   InputLevel
	paste   bool
	// levels of verbosity / amount of logging:
	logparse  bool
	logtype   bool
	logtrace  bool
	incltoken bool
	treefile  string
	evalfile  string
}
type InputLevel int

const (
	ProgramL InputLevel = iota
	StatementL
	ExpressionL
)

func (i InputLevel) String() string {
	switch i {
	case ProgramL:
		return "program"
	case StatementL:
		return "statement"
	case ExpressionL:
		return "expression"
	default:
		return fmt.Sprintf("%d", int(i))
	}
}

type InputProcess int

const (
	EvalP InputProcess = iota
	ParseP
	TypeP
)

func (i InputProcess) String() string {
	switch i {
	case EvalP:
		return "eval"
	case ParseP:
		return "parse"
	case TypeP:
		return "type"
	default:
		return fmt.Sprintf("%d", int(i))
	}
}

const ( //default settings
	prompt_default       = ">> "
	treefile_default     = "tree.pdf"
	evalfile_default     = "eval.pdf"
	inputProcess_default = EvalP
	inputLevel_default   = ProgramL
	paste_default        = false

	logparse_default = false
	logtype_default  = false
	logtrace_default = false

	incltoken_default = false
)

// NewSession creates a new Session.
func NewSession(in io.Reader, out io.Writer) *Session {

	path, err := exec.LookPath("pdflatex")
	if err != nil {
		path = ""
	}

	s := &Session{
		scanner:       bufio.NewScanner(in),
		out:           out,
		prompt:        prompt_default,
		environment:   object.NewEnvironment(),
		path_pdflatex: path,
		level:         inputLevel_default,
		process:       inputProcess_default,
		logtype:       logtype_default,
		logtrace:      logtrace_default,
		logparse:      logparse_default,
		paste:         paste_default,
		incltoken:     incltoken_default,
		treefile:      treefile_default,
		evalfile:      evalfile_default,
	}

	s.init()
	return s
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

// settings
func (s *Session) exec_settings() {

	t := table.NewWriter()
	t.SetOutputMirror(s.out)
	t.AppendHeader(table.Row{"setting", "current value", "default value"})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"paste", s.paste, paste_default})
	t.AppendRow([]interface{}{"process", s.process, inputProcess_default})
	t.AppendRow([]interface{}{"level", s.level, inputLevel_default})
	t.AppendRow([]interface{}{"logparse", s.logparse, logparse_default})
	t.AppendRow([]interface{}{"logtype", s.logtype, logtype_default})
	t.AppendRow([]interface{}{"logtrace", s.logtrace, logtrace_default})
	t.AppendRow([]interface{}{"prompt", s.prompt, prompt_default})
	t.AppendRow([]interface{}{"incltoken", s.incltoken, incltoken_default})
	t.AppendRow([]interface{}{"treefile", s.treefile, treefile_default})
	t.AppendRow([]interface{}{"evalfile", s.evalfile, evalfile_default})

	//t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func (s *Session) exec_reset(input string) {
	// todo: datastructure for settings
	// "prompt" und die, die folgen, sind Steuer[ungs]zeichen

	switch input {
	case "prompt":
		s.prompt = prompt_default
	case "logtype":
		s.logtype = logtype_default
	case "logparse":
		s.logparse = logparse_default
	case "logtrace":
		s.logtrace = logtrace_default
	case "incltoken":
		s.incltoken = incltoken_default
	case "treefile":
		s.treefile = treefile_default
	case "evalfile":
		s.evalfile = evalfile_default
	case "paste":
		s.paste = paste_default
	case "level":
		s.level = inputLevel_default
	case "process":
		s.process = inputProcess_default
	default:
		s.exec_help("reset")
	}
}

func (s *Session) exec_unset(setting string) {

	switch setting {
	case "logtype":
		s.logtype = false
	case "logparse":
		s.logparse = false
	case "logtrace":
		s.logtrace = false
	case "incltoken":
		s.incltoken = false
	case "paste":
		s.paste = false
	default:
		s.exec_help("reset")
	}
}

func (s *Session) exec_set(input string) {

	levelmap := make(map[string]InputLevel)
	levelmap["program"] = ProgramL
	levelmap["p"] = ProgramL
	levelmap["statement"] = StatementL
	levelmap["s"] = StatementL
	levelmap["expression"] = ExpressionL
	levelmap["e"] = ExpressionL

	processmap := make(map[string]InputProcess)
	processmap["parse"] = ParseP
	processmap["p"] = ParseP
	processmap["eval"] = EvalP
	processmap["e"] = EvalP
	processmap["type"] = TypeP
	processmap["t"] = TypeP

	// todo: datastructure for settings
	slice := strings.SplitN(input, " ", 2)
	setting := slice[0]
	if len(slice) == 1 {
		switch setting {
		case "logtype":
			s.logtype = true
			return
		case "logparse":
			s.logparse = true
			return
		case "logtrace":
			s.logtrace = true
			return
		case "incltoken":
			s.incltoken = true
			return
		case "paste":
			s.paste = true
			return
		}
	}
	if len(slice) == 2 {
		arg := slice[1]
		switch setting {
		case "prompt":
			s.prompt = arg + " "
			return
		case "level":
			level, ok := levelmap[arg]
			if ok {
				s.level = level
				return
			}
		case "process":
			process, ok := processmap[arg]
			if ok {
				s.process = process
				return
			}
		case "treefile":
			//TODO: maybe check whether that's a valid filename nö
			if !strings.HasSuffix(arg, ".pdf") {
				arg = arg + ".pdf"
			}
			s.treefile = arg
			return

		case "evalfile":
			//TODO: maybe check whether that's a valid filename nö
			if !strings.HasSuffix(arg, ".pdf") {
				arg = arg + ".pdf"
			}
			s.evalfile = arg
			return

		}

	}
	s.exec_help("set")
}

// input processing

func (s *Session) exec_parsetree(line string) {
	fmt.Fprint(s.out, "not yet implemented!\n")
}

func (s *Session) exec_evaltree(line string) {
	fmt.Fprint(s.out, "not yet implemented!\n")
}

func (s *Session) exec_process(line string) {
	s.process_input_dim(s.paste, s.level, s.process, false, line)
}

func (s *Session) exec_paste_empty_arg() {
	s.process_input_dim(true, s.level, s.process, false, "")
}

func (s *Session) exec_paste(line string) {
	s.process_input_dim(true, s.level, s.process, false, line)
}

func (s *Session) exec_expression(line string) {
	s.process_input_dim(s.paste, ExpressionL, s.process, false, line)
}

func (s *Session) exec_statement(line string) {
	s.process_input_dim(s.paste, StatementL, s.process, false, line)
}

func (s *Session) exec_program(line string) {
	s.process_input_dim(s.paste, ProgramL, s.process, false, line)
}

func (s *Session) exec_eval(line string) {
	s.process_input_dim(s.paste, s.level, EvalP, false, line)
}

func (s *Session) exec_type(line string) {
	s.process_input_dim(s.paste, s.level, TypeP, false, line)
}

func (s *Session) exec_trace(line string) {
	s.process_input_dim(s.paste, s.level, EvalP, true, line)
}

func (s *Session) exec_parse(line string) {
	s.process_input_dim(s.paste, s.level, ParseP, false, line)
}

func (s *Session) process_input_dim(paste bool, level InputLevel, process InputProcess, trace bool, input string) {

	if paste {
		input = s.multiline_input(input)
	}

	l := lexer.New(input)
	p := parser.New(l)

	node := parse_level(p, level)

	if s.logparse {
		fmt.Fprint(s.out, "log ast:\t")
	}
	if s.logparse || process == ParseP {
		fmt.Fprintln(s.out, node)
		fmt.Fprintln(s.out, visualizer.RepresentNodeConsoleTree(node, "|   ", !s.incltoken))
		//	fmt.Fprintln(s.out, visualizer.QTree(node, !s.incltoken))
		path, err := exec.LookPath("pdflatex")
		if err != nil {
			fmt.Fprintln(s.out, "Displaying trees as pdfs is not available to you, since you have not installed pdflatex.")
		} else {
			visualizer.Ast2pdf(node, !s.incltoken, s.treefile, path)
		}
	}

	if len(p.Errors()) != 0 {
		s.printParserErrors(p.Errors(), level)
		return
	}

	if process == ParseP {
		return
	}

	if trace || s.logtrace {
		evaluator.StartTracer()
	}

	evaluated := evaluator.Eval(node, s.environment)

	if trace || s.logtrace {
		evaluator.StopTracer()
	}

	if s.logtype {
		fmt.Fprint(s.out, "log type:\t")
	}
	if s.logtype || process == TypeP {
		fmt.Fprintln(s.out, reflect.TypeOf(evaluated))
	}
	if process == TypeP {
		return
	}

	if s.logtrace {
		visualizer.RepresentEvalConsole(evaluator.T, s.out)
		//fmt.Fprint(s.out, visualizer.QTreeEval(evaluator.T))
		path, err := exec.LookPath("pdflatex")
		if err != nil {
			fmt.Fprintln(s.out, "Displaying evaluation trees as pdfs is not available to you, since you have not installed pdflatex.")
		} else {
			visualizer.EvalTree2pdf(evaluator.T, s.evalfile, path)
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
