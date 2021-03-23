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
	"sort"
	"strings"
)

func Start(in io.Reader, out io.Writer) {

	s := NewSession(in, out)

	for {

		fmt.Fprint(out, s.PROMPT) // Fprint instead of Fprintf due to SA1006

		scanned := s.scanner.Scan()
		if !scanned {
			return
		}

		line := s.scanner.Text()

		s.exec_cmd(line)
	}
}

type Session struct {
	PROMPT      string
	environment *object.Environment
	scanner     *bufio.Scanner
	out         io.Writer
	//historyExpr		[]ast.Expression
	//historyStmsts		[]ast.Statement
	//historyPrograms	[]ast.Programs
	// --> maybe not needed, maybe we should put the Stmts programs consist of into historyStmts
}

// NewSession creates a new Session.
func NewSession(in io.Reader, out io.Writer) *Session {

	s := &Session{
		PROMPT:      ">> ",
		environment: object.NewEnvironment(),
		scanner:     bufio.NewScanner(in),
		out:         out,
	}

	s.init()

	return s
}

type command struct {
	name     string
	usage    string
	single   func()
	with_arg func(string) // initialized here --> end msg about potential cycle
}

var commands = make(map[string]command)

func (s *Session) init() { // to avoid cycle

	c_quit := &command{
		name:   "q[uit]",
		single: s.exec_quit,
		usage:  "quit the session",
	}
	commands["quit"] = *c_quit
	commands["q"] = commands["quit"]

	c_clear := &command{
		name:   "clear",
		single: s.exec_clear,
		usage:  "clear the environment",
	}
	commands["clear"] = *c_clear

	c_set := &command{
		name:     "set",
		with_arg: s.exec_set,
		usage:    "~ prompt <prompt> \t set prompt string to <prompt>",
	}
	commands["set"] = *c_set

	c_help := &command{
		name:     "h[elp]",
		single:   s.exec_help_all,
		with_arg: s.exec_help,
		usage:    "list all commands with usage \n\t ~ <cmd> \t usage command <cmd>",
	}

	commands["help"] = *c_help
	commands["h"] = commands["help"]

}

func (s *Session) exec_cmd(line string) {
	if !strings.HasPrefix(line, ":") {
		s.exec_default(line)
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

func (s *Session) exec_help(cmd string) {

	if command, ok := commands[cmd]; ok {
		fmt.Fprintln(s.out, "usage", command.name+":\t", command.usage)
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
	for _, name := range names {
		command := set[name]
		fmt.Fprintln(s.out, name, "\t", command.usage)
	}
}

func (s *Session) exec_set(input string) {
	// todo: datastructure for settings
	slice := strings.SplitN(input, " ", 2)
	if len(slice) == 2 {
		if slice[0] == "prompt" {
			s.PROMPT = slice[1] + " "
			return
		}
	}
	s.exec_help("set")
}

func (s *Session) exec_default(line string) {
	l := lexer.New(line)
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(s.out, p.Errors())
		return
	}

	evaluated := evaluator.Eval(program, s.environment)
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
