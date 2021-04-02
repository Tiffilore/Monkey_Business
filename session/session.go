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
	"sort"
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

type Session struct {
	scanner     *bufio.Scanner
	out         io.Writer
	environment *object.Environment
	prompt      string
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

	//historyExpr		[]ast.Expression
	//historyStmsts		[]ast.Statement
	//historyPrograms	[]ast.Programs
	// --> maybe not needed, maybe we should put the Stmts programs consist of into historyStmts
}

const ( //default settings
	prompt_default       = ">> "
	treefile_default     = "tree.pdf"
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

	s := &Session{
		scanner:     bufio.NewScanner(in),
		out:         out,
		prompt:      prompt_default,
		environment: object.NewEnvironment(),
		level:       inputLevel_default,
		process:     inputProcess_default,
		logtype:     logtype_default,
		logtrace:    logtrace_default,
		logparse:    logparse_default,
		paste:       paste_default,
		incltoken:   incltoken_default,
		treefile:    treefile_default,
	}

	s.init()
	return s
}

type command struct {
	name     string
	single   func()
	with_arg func(string) // initialized here --> end msg about potential cycle
	usage    []struct {
		args string
		msg  string
	}
}

var commands = make(map[string]command)

func (s *Session) init() { // to avoid cycle

	c_quit := &command{
		name:   "q[uit]",
		single: s.exec_quit,
		usage: []struct {
			args string
			msg  string
		}{
			{"~", "quit the session"},
		},
	}
	commands["quit"] = *c_quit
	commands["q"] = commands["quit"]

	c_settings := &command{
		name:   "settings",
		single: s.exec_settings,
		usage: []struct {
			args string
			msg  string
		}{
			{"~", "list all settings with their current values and defaults"},
		},
	}
	commands["settings"] = *c_settings

	c_clear := &command{
		name:   "clear",
		single: s.exec_clear,
		usage: []struct {
			args string
			msg  string
		}{
			{"~", "clear the environment"},
		},
	}
	commands["clear"] = *c_clear

	c_clearscreen := &command{
		name:   "cl[earscreen]",
		single: s.exec_clearscreen,
		usage: []struct {
			args string
			msg  string
		}{
			{"~", "clear the terminal screen"},
		},
	}
	commands["clearscreen"] = *c_clearscreen
	commands["cl"] = commands["clearscreen"]

	c_list := &command{
		name:   "l[ist]",
		single: s.exec_list,
		usage: []struct {
			args string
			msg  string
		}{
			{"~", "list all identifiers in the environment alphabetically \n\t with types and values"},
		},
	}
	commands["list"] = *c_list
	commands["l"] = commands["list"]

	c_paste := &command{
		name:     "paste",
		with_arg: s.exec_paste,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "evaluate multiline <input> (terminated by blank line)"},
		},
	}
	commands["paste"] = *c_paste

	c_set := &command{
		name:     "set",
		with_arg: s.exec_set,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ process <p>", "<p> must be: [e]val, [p]arse, [t]ype"},
			{"~ level <l>", "<l> must be: [p]rogram, [s]tatement, [e]xpression"},
			{"~ logparse", "additionally output ast-string"},
			{"~ logtype", "additionally output objecttype"},
			{"~ logtrace", "additionally output evaluation trace"},
			{"~ incltoken", "include tokens in representations of asts"},
			{"~ paste", "enable multiline support"},
			{"~ prompt <prompt>", "set prompt string to <prompt>"},
			{"~ treefile <f>", "set file that outputs pdfs to <f>"},
		},
	}
	commands["set"] = *c_set

	c_reset := &command{
		name:     "reset",
		with_arg: s.exec_reset,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <setting>", "set <setting> to default\n\t for an overview consult :settings and/or :h set"},
		},
	}
	commands["reset"] = *c_reset

	c_unset := &command{
		name:     "unset",
		with_arg: s.exec_unset,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <setting>", "set boolean <setting> to false\n\t for an overview consult :settings and/or :h set"},
			//	{"~ logparse", "don't additionally output ast-string"},
			//	{"~ logtype", "don't additionally output objecttype"},
			//	{"~ paste", "disable multiline support"},
			//incltoken logtrace
		},
	}
	commands["unset"] = *c_unset

	c_help := &command{
		name:     "h[elp]",
		single:   s.exec_help_all,
		with_arg: s.exec_help,
		usage: []struct {
			args string
			msg  string
		}{
			{"~", "list all commands with usage"},
			{"~ <cmd>", "print usage command <cmd>"},
		},
	}
	commands["help"] = *c_help
	commands["h"] = commands["help"]

	c_type := &command{
		name:     "t[ype]",
		with_arg: s.exec_type,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "show objecttype <input> evaluates to"},
		},
	}
	commands["type"] = *c_type
	commands["t"] = commands["type"]

	c_trace := &command{
		name:     "trace",
		with_arg: s.exec_trace,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "show evaluation trace"},
		},
	}
	commands["trace"] = *c_trace

	c_parse := &command{
		name:     "p[arse]",
		with_arg: s.exec_parse,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "parse <input>"},
		},
	}
	commands["parse"] = *c_parse
	commands["p"] = commands["parse"]

	c_eval := &command{
		name:     "e[val]",
		with_arg: s.exec_eval,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "print out value of object <input> evaluates to"},
		},
	}
	commands["eval"] = *c_eval
	commands["e"] = commands["eval"]

	c_expr := &command{
		name:     "expr[ession]",
		with_arg: s.exec_expression,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "expect <input> to be an expression"},
		},
	}
	commands["expression"] = *c_expr
	commands["expr"] = commands["expression"]

	c_stmt := &command{
		name:     "stmt|statement",
		with_arg: s.exec_statement,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "expect <input> to be a statement"},
		},
	}
	commands["statement"] = *c_stmt
	commands["stmt"] = commands["statement"]

	c_prog := &command{
		name:     "prog[ram]",
		with_arg: s.exec_program,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "expect <input> to be a program"},
		},
	}
	commands["program"] = *c_prog
	commands["prog"] = commands["program"]

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
	if c_entry, ok := commands[cmd]; ok {
		if len(slice) == 1 {
			if c_entry.single != nil {
				c_entry.single()
				return
			}
		} else {
			arg := slice[1]
			if c_entry.with_arg != nil {
				c_entry.with_arg(arg)
				return
			}
		}
	}

	s.exec_help(cmd)
}

// quit
func (s *Session) exec_quit() {
	os.Exit(0)
}

// clear the screen
func (s *Session) exec_clearscreen() {

	_, err := exec.LookPath("clear")
	if err != nil {
		fmt.Fprintln(s.out, "command clearscreen is not available to you")
	}

	cmd := exec.Command("clear")
	cmd.Stdout = s.out
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}

// environment
func (s *Session) exec_clear() {
	s.environment = object.NewEnvironment()
}

func (s *Session) exec_list() {
	store := s.environment.Store()

	//sort alphabetically
	keys := make([]string, 0, len(store))
	for k := range store {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	t := table.NewWriter()
	t.SetOutputMirror(s.out)
	t.AppendHeader(table.Row{"Identifier", "Type", "Value"})
	t.AppendSeparator()

	for _, key := range keys {
		object := store[key]
		nodetype := reflect.TypeOf(object)
		var value string
		if object == nil {
			value = "nil"
		} else {

			value = object.Inspect() //strings.ReplaceAll(object.Inspect(), "\n", "\n\t  ")
		}
		t.AppendRow([]interface{}{key, nodetype, value})
	}
	//t.AppendFooter(table.Row{"", "", "Total", 10000})
	//t.SetStyle(table.StyleColoredBright)
	t.Render()
}

// commands
func (s *Session) exec_help(cmd string) {

	if command, ok := commands[cmd]; ok {

		//print
		t := table.NewWriter()
		t.SetOutputMirror(s.out)

		usage := command.usage
		if len(usage) == 0 {
			t.AppendRow([]interface{}{command.name, "no usage message provided"})
		} else {
			for i, msg := range usage {
				if i == 0 {
					t.AppendRow([]interface{}{command.name, msg.args, msg.msg})
				} else {
					t.AppendRow([]interface{}{"", msg.args, msg.msg})
				}
			}

		}
		t.Render()

		return
	}

	fmt.Fprintln(s.out, "unknown command: ", cmd)
}

func (s *Session) exec_help_all() {

	// print out all commands, but only once (--> quit, q)
	// always in the same order

	//remove duplicates
	var set = make(map[string]command)
	for _, command := range commands {
		set[command.name] = command
	}

	//sort
	names := make([]string, len(set))
	i := 0
	for name := range set {
		names[i] = name
		i++
	}
	sort.Strings(names)

	//print
	t := table.NewWriter()
	t.SetOutputMirror(s.out)
	t.AppendHeader(table.Row{"Name", "", "Usage"})
	t.AppendSeparator()

	for _, name := range names {
		command := set[name]
		usage := command.usage
		if len(usage) == 0 {
			t.AppendRow([]interface{}{name, "no usage message provided"})
		} else {
			for i, msg := range usage {
				if i == 0 {
					t.AppendRow([]interface{}{name, msg.args, msg.msg})
				} else {
					t.AppendRow([]interface{}{"", msg.args, msg.msg})
				}
			}

		}
	}
	t.Render()
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
		}
	}
	s.exec_help("set")
}

// input processing
func (s *Session) exec_process(line string) {
	s.process_input_dim(s.paste, s.level, s.process, s.logtrace, line)
}

func (s *Session) exec_paste(line string) {
	s.process_input_dim(true, s.level, s.process, s.logtrace, line)
}

func (s *Session) exec_expression(line string) {
	s.process_input_dim(s.paste, ExpressionL, s.process, s.logtrace, line)
}

func (s *Session) exec_statement(line string) {
	s.process_input_dim(s.paste, StatementL, s.process, s.logtrace, line)
}

func (s *Session) exec_program(line string) {
	s.process_input_dim(s.paste, ProgramL, s.process, s.logtrace, line)
}

func (s *Session) exec_eval(line string) {
	s.process_input_dim(s.paste, s.level, EvalP, s.logtrace, line)
}

func (s *Session) exec_type(line string) {
	s.process_input_dim(s.paste, s.level, TypeP, s.logtrace, line)
}

func (s *Session) exec_trace(line string) {
	s.process_input_dim(s.paste, s.level, EvalP, true, line)
}

func (s *Session) exec_parse(line string) {
	s.process_input_dim(s.paste, s.level, ParseP, s.logtrace, line)
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
		//fmt.Fprintln(s.out, visualizer.QTree(node, !s.incltoken))
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

	if trace {
		evaluator.StartTracer()
	}

	evaluated := evaluator.Eval(node, s.environment)

	if trace {
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

	if trace {
		visualizer.RepresentEvalConsole(evaluator.T, s.out)
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
