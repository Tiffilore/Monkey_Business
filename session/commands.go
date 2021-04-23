package session

import (
	"bytes"
	"fmt"

	"github.com/jedib0t/go-pretty/v6/table"
)

type command struct {
	name     string
	single   func()
	with_arg func(string) // initialized here --> end msg about potential cycle
	usage    []struct {
		args string
		msg  string
	}
}

type commandSet struct {
	m map[string]*command
	l []*command
}

func (c *commandSet) register(cmd_string string, cmd *command) {
	_, ok := c.m[cmd_string]
	if ok {
		fmt.Printf("warning: command %v has already been defined!\n", cmd_string)
		return
	}

	c.m[cmd_string] = cmd

	for _, command := range c.l {
		if command == cmd {
			return
		}
	}

	c.l = append(c.l, cmd)

}

// func (c commandSet) has_command(cmd_string string) bool {
// 	_, ok := c.m[cmd_string]
// 	return ok
// }

func (c commandSet) get_exec_single(cmd_string string) (func(), bool) {
	cmd, ok := c.m[cmd_string]
	if !ok {
		return nil, false
	}
	if cmd.single == nil {
		return nil, false
	}
	return cmd.single, true
}

func (c commandSet) get_exec_with_arg(cmd_string string) (func(string), bool) {
	cmd, ok := c.m[cmd_string]
	if !ok {
		return nil, false
	}
	if cmd.with_arg == nil {
		return nil, false
	}
	return cmd.with_arg, true
}

func (c commandSet) usage(cmd_string string) (string, bool) {
	command, ok := c.m[cmd_string]
	if !ok {
		return "", false
	}

	var out bytes.Buffer

	t := table.NewWriter()
	t.SetOutputMirror(&out)

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
	return out.String(), true
}

func (c commandSet) menu() string {

	var out bytes.Buffer

	t := table.NewWriter()
	t.SetOutputMirror(&out)

	t.AppendHeader(table.Row{"Name", "", "Usage"})
	t.AppendSeparator()

	for _, command := range c.l {
		name := command.name
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
	return out.String()
}

func newCommandSet() *commandSet {
	m := make(map[string]*command)
	l := make([]*command, 0)

	return &commandSet{m: m, l: l}
}

var commands *commandSet

func (s *Session) init() { // to avoid cycle

	commands = newCommandSet()

	// help
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

	commands.register("help", c_help)
	commands.register("h", c_help)

	//quit
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

	commands.register("quit", c_quit)
	commands.register("q", c_quit)

	// clearscreen
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

	commands.register("clearscreen", c_clearscreen)
	commands.register("cl", c_clearscreen)

	// environment: list
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

	commands.register("list", c_list)
	commands.register("l", c_list)

	// environment: clear
	c_clear := &command{
		name:   "c[lear]",
		single: s.exec_clear,
		usage: []struct {
			args string
			msg  string
		}{
			{"~", "clear the environment"},
		},
	}

	commands.register("clear", c_clear)
	commands.register("c", c_clear)

	// paste
	c_paste := &command{
		name:     "paste",
		with_arg: s.exec_paste,
		single:   s.exec_paste_empty_arg,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "evaluate multiline <input> (terminated by blank line)"},
		},
	}
	commands.register("paste", c_paste)

	// level: expression
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

	commands.register("expression", c_expr)
	commands.register("expr", c_expr)

	// level: statement
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

	commands.register("statement", c_stmt)
	commands.register("stmt", c_stmt)

	// level: program
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

	commands.register("program", c_prog)
	commands.register("prog", c_prog)

	// process: parse
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

	commands.register("parse", c_parse)
	commands.register("p", c_parse)

	// process: parsetree: TODO
	c_parsetree := &command{
		name:     "p[arse]tree",
		with_arg: s.exec_parse, //TODO
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "parse <input>"}, //TODO
		},
	}

	commands.register("parsetree", c_parsetree)
	commands.register("ptree", c_parsetree)

	// process: eval
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

	commands.register("eval", c_eval)
	commands.register("e", c_eval)

	// process: type
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

	commands.register("type", c_type)
	commands.register("t", c_type)

	// process: trace
	c_trace := &command{
		name:     "trace",
		with_arg: s.exec_trace,
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "show evaluation trace step by step"},
		},
	}
	commands.register("trace", c_trace)

	// process: evaltree TODO
	c_evaltree := &command{
		name:     "e[val]tree",
		with_arg: s.exec_eval, //TODO
		usage: []struct {
			args string
			msg  string
		}{
			{"~ <input>", "print out value of object <input> evaluates to"}, //TODO
		},
	}

	commands.register("evaltree", c_evaltree)
	commands.register("etree", c_evaltree)

	// settings: settings
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
	commands.register("settings", c_settings)

	// settings: set
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
			{"~ treefile <f>", "set file that outputs ast-pdfs to <f>"},
			{"~ evalfile <f>", "set file that outputs eval-pdfs to <f>"},
		},
	}
	commands.register("set", c_set)

	// settings: reset
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
	commands.register("reset", c_reset)

	// settings: unset
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
	commands.register("unset", c_unset)
}
