package session

import (
	"bufio"
	"fmt"
	"io"
	"monkey/ast"
	"monkey/evaluator"
	"monkey/lexer"
	"monkey/object"
	"monkey/parser"
	"os"
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
	prompt      string
	environment *object.Environment
	logtype     bool
	paste       bool
	level       InputLevel
	process     InputProcess

	//historyExpr		[]ast.Expression
	//historyStmsts		[]ast.Statement
	//historyPrograms	[]ast.Programs
	// --> maybe not needed, maybe we should put the Stmts programs consist of into historyStmts
}

const prompt_default = ">> "
const logtype_default = false
const paste_default = false
const inputLevel_default = ProgramL
const inputProcess_default = EvalP

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
		paste:       paste_default,
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
			{"~ logtype", "when eval, additionally output objecttype"},
			{"~ paste", "enable multiline support"},
			{"~ prompt <prompt>", "set prompt string to <prompt>"},
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
			{"~ process", "set process to default"},
			{"~ level", "set level to default"},
			{"~ logtype", "set logtype to default"},
			{"~ paste", "set multiline support to default"},
			{"~ prompt", "set prompt to default"},
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
			{"~ logtype", "when eval, don't additionally output objecttype"},
			{"~ paste", "disable multiline support"},
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
}

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

func (s *Session) exec_quit() {
	os.Exit(0)
}

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

/*
	settings:
		(set|reset|[unset])
		set bool vs set value

*/

func (s *Session) exec_settings() {

	t := table.NewWriter()
	t.SetOutputMirror(s.out)
	t.AppendHeader(table.Row{"setting", "current value", "default value"})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"process", s.process, inputProcess_default})
	t.AppendRow([]interface{}{"level", s.level, inputLevel_default})
	t.AppendRow([]interface{}{"logtype", s.logtype, logtype_default})
	t.AppendRow([]interface{}{"paste", s.paste, paste_default})
	t.AppendRow([]interface{}{"prompt", s.prompt, prompt_default})

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
	//	{"~ logtype", "when eval, additionally output objecttype "},
	setting := slice[0]
	if len(slice) == 1 {
		switch setting {
		case "logtype":
			s.logtype = true
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
		}
	}
	s.exec_help("set")
}

func (s *Session) exec_paste(input string) {
	s.helper_paste(input, s.eval)
}

func (s *Session) helper_paste(input string, f func(string)) {
	for {
		scanned := s.scanner.Scan()
		if !scanned {
			return
		}
		line := s.scanner.Text()
		if line == "" {
			f(input)
			return
		}
		input += " " + line
	}
}

func (s *Session) exec_process(line string) {
	switch s.process {
	case EvalP:
		s.exec_eval(line)
	case ParseP:
		s.exec_parse(line)
	case TypeP:
		s.exec_type(line)
	default:
		fmt.Println(s.out, "uppsi")

	}
}

func (s *Session) exec_eval(line string) {
	if s.paste {
		s.helper_paste(line, s.eval)
		return
	}
	s.eval(line)
}

func (s *Session) exec_type(line string) {
	if s.paste {
		s.helper_paste(line, s.det_type)
		return
	}
	s.det_type(line)
}

func (s *Session) det_type(line string) {

	l := lexer.New(line)
	p := parser.New(l)

	ast := s.parseL(p, line)

	//visualizer.Ast2pdf(program, "show")
	if len(p.Errors()) != 0 {
		s.printParserErrors(p.Errors())
		return
	}

	evaluated := evaluator.Eval(ast, s.environment)
	fmt.Fprintln(s.out, reflect.TypeOf(evaluated))
}

func (s *Session) exec_parse(line string) {
	if s.paste {
		s.helper_paste(line, s.parse)
		return
	}
	s.parse(line)
}

func (s *Session) parse(line string) {

	l := lexer.New(line)
	p := parser.New(l)

	ast := s.parseL(p, line)

	//visualizer.Ast2pdf(program, "show")
	if len(p.Errors()) != 0 {
		s.printParserErrors(p.Errors())
		return
	}

	fmt.Fprintln(s.out, ast)

}

func (s *Session) parseL(p *parser.Parser, line string) ast.Node {

	switch s.level {
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

func (s *Session) eval(line string) {

	l := lexer.New(line)
	p := parser.New(l)

	ast := s.parseL(p, line)

	//visualizer.Ast2pdf(program, "show")
	if len(p.Errors()) != 0 {
		s.printParserErrors(p.Errors())
		return
	}

	evaluated := evaluator.Eval(ast, s.environment)
	if s.logtype {
		fmt.Fprintln(s.out, reflect.TypeOf(evaluated))
	}
	if evaluated != nil {
		fmt.Fprintln(s.out, evaluated.Inspect())
	}
	/*
		io.WriteString vs fmt.Fprint ?????
			The difference is that fmt.Fprint is formatting the arguments provided first in a buffer before calling w.Write.
			And io.WriteString is checking if w provides the StringWriter interface and calls that instead.
	*/
}

func (s *Session) printParserErrors(errors []string) {

	fmt.Fprintf(s.out, "... cannot be parsed as %v\n", s.level)
	//io.WriteString(s.out, " parser errors:\n")
	for _, msg := range errors {
		fmt.Fprintf(s.out, "\t%v\n", msg)
	}
}
