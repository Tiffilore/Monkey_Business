package session

import (
	"bufio"
	"fmt"
	"io"
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

type Session struct {
	scanner     *bufio.Scanner
	out         io.Writer
	prompt      string
	environment *object.Environment
	logtype     bool
	//historyExpr		[]ast.Expression
	//historyStmsts		[]ast.Statement
	//historyPrograms	[]ast.Programs
	// --> maybe not needed, maybe we should put the Stmts programs consist of into historyStmts
}

const prompt_default = ">> "
const logtype_default = false

// NewSession creates a new Session.
func NewSession(in io.Reader, out io.Writer) *Session {

	s := &Session{
		scanner:     bufio.NewScanner(in),
		out:         out,
		prompt:      prompt_default,
		environment: object.NewEnvironment(),
		logtype:     logtype_default,
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

	c_set := &command{
		name:     "set",
		with_arg: s.exec_set,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ logtype", "when eval, additionally output objecttype"},
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
			{"~ logtype", "set logtype to default"},
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
		s.exec_eval(line) //default action
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

func (s *Session) exec_reset(input string) {
	// todo: datastructure for settings
	// "prompt" und die, die folgen, sind Steuer[ungs]zeichen

	switch input {
	case "prompt":
		s.prompt = prompt_default
	case "logtype":
		s.logtype = logtype_default
	default:
		s.exec_help("reset")
	}
}

func (s *Session) exec_unset(input string) {

	switch input {
	case "logtype":
		s.logtype = false
	default:
		s.exec_help("reset")
	}
}

func (s *Session) exec_set(input string) {
	// todo: datastructure for settings
	slice := strings.SplitN(input, " ", 2)
	//	{"~ logtype", "when eval, additionally output objecttype "},
	setting := slice[0]
	if len(slice) == 1 {
		if setting == "logtype" {
			s.logtype = true
			return
		}
	}
	if len(slice) == 2 {
		if setting == "prompt" {
			s.prompt = slice[1] + " "
			return
		}
	}
	s.exec_help("set")
}

func (s *Session) exec_type(line string) {
	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()

	//visualizer.Ast2pdf(program, "show")
	if len(p.Errors()) != 0 {
		printParserErrors(s.out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, s.environment)
	fmt.Fprintln(s.out, reflect.TypeOf(evaluated))

}
func (s *Session) exec_eval(line string) {
	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()

	//visualizer.Ast2pdf(program, "show")
	if len(p.Errors()) != 0 {
		printParserErrors(s.out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, s.environment)
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

const MONKEY_FACE = `            __,__
   .--.  .-"     "-.  .--.
  / .. \/  .-. .-.  \/ .. \
 | |  '|  /   Y   \  |'  | |
 | \   \  \ 0 | 0 /  /   / |
  \ '- ,\.-"""""""-./, -' /
   ''-' /_   ^ ^   _\ '-''
       |  \._   _./  |
       \   \ '~' /   /
        '._ '-=-' _.'
           '-----'
`

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, MONKEY_FACE)
	io.WriteString(out, "Woops! We ran into some monkey business here!\n")
	io.WriteString(out, " parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
