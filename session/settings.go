package session

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type settings struct {
	prompt    string
	paste     bool
	level     inputLevel
	process   inputProcess
	logs      logs
	displays  visDisplays
	verbosity int
	inclToken bool
	inclEnv   bool
	file      string
}

func newSettings() *settings {

	displays := visDisplays{
		ConsD: true,
		PdfD:  false,
	}
	logs := logs{
		ParseTreeP: false,
		EvalTreeP:  false,
		TypeP:      false,
		TraceP:     false,
	}

	s := settings{
		prompt:    ">> ",
		paste:     false,
		level:     ProgramL,
		process:   EvalP,
		logs:      logs,
		displays:  displays,
		verbosity: 0,
		inclToken: false,
		inclEnv:   false,
		file:      "tree.pdf",
	}

	return &s
}

var defaultSettings *settings
var currentSettings *settings

func init() {
	defaultSettings = newSettings()
	currentSettings = newSettings()
}

func (s *Session) exec_reset_all() {
	currentSettings = newSettings()
}
func (s *Session) exec_reset(input string) {
	switch strings.Trim(input, " ") {
	case "prompt":
		currentSettings.prompt = defaultSettings.prompt
	case "paste":
		currentSettings.paste = defaultSettings.paste
	case "level":
		currentSettings.level = defaultSettings.level
	case "process":
		currentSettings.process = defaultSettings.process
	case "logs":
		for key := range currentSettings.logs {
			currentSettings.logs[key] = defaultSettings.logs[key]
		}
	case "displays":
		for key := range currentSettings.displays {
			currentSettings.displays[key] = defaultSettings.displays[key]
		}
	case "verbosity":
		currentSettings.verbosity = defaultSettings.verbosity
	case "inclToken":
		currentSettings.inclToken = defaultSettings.inclToken
	case "inclEnv":
		currentSettings.inclEnv = defaultSettings.inclEnv
	case "file":
		currentSettings.file = defaultSettings.file
	default:
		s.exec_help("reset")
	}
}

func (s *Session) exec_settings() {

	t := table.NewWriter()
	t.SetOutputMirror(s.out)
	t.AppendHeader(table.Row{"setting", "current value", "default value"})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"prompt", currentSettings.prompt, defaultSettings.prompt})
	t.AppendRow([]interface{}{"paste", currentSettings.paste, defaultSettings.paste})
	t.AppendRow([]interface{}{"level", currentSettings.level, defaultSettings.level})
	t.AppendRow([]interface{}{"process", currentSettings.process, defaultSettings.process})
	t.AppendRow([]interface{}{"logs", currentSettings.logs, defaultSettings.logs})
	t.AppendRow([]interface{}{"displays", currentSettings.displays, defaultSettings.displays})
	t.AppendRow([]interface{}{"verbosity", currentSettings.verbosity, defaultSettings.verbosity})
	t.AppendRow([]interface{}{"inclToken", currentSettings.inclToken, defaultSettings.inclToken})
	t.AppendRow([]interface{}{"inclEnv", currentSettings.inclEnv, defaultSettings.inclEnv})
	t.AppendRow([]interface{}{"file", currentSettings.file, defaultSettings.file})
	//t.SetStyle(table.StyleColoredBright)
	t.Render()
}

func (s *Session) exec_unset(input string) {

	switch strings.Trim(input, " ") {
	case "paste":
		currentSettings.paste = false
	case "inclToken":
		currentSettings.inclToken = false
	case "inclEnv":
		currentSettings.inclEnv = false
	default:
		s.exec_help("reset")
	}
}

func (s *Session) exec_set(input string) {

	splits := strings.SplitN(strings.Trim(input, " "), " ", 2)
	setting := splits[0]

	if len(splits) == 0 {

		switch setting {
		case "paste":
			currentSettings.paste = true
			return
		case "inclToken":
			currentSettings.inclToken = true
			return
		case "inclEnv":
			currentSettings.inclEnv = true
			return
		}
	} else {
		arg := splits[1]
		switch setting {
		case "prompt":
			currentSettings.prompt = arg + " "
			return
		case "level":
			level, ok := getInputLevel(arg)
			if ok {
				currentSettings.level = level
				return
			}
		case "process":
			process, ok := getInputProcess(arg)
			if ok {
				currentSettings.process = process
				return
			}
		case "displays":
			err := setDisplays(currentSettings.displays, arg)
			if err != nil {
				fmt.Fprintln(s.out, err)
			}
			return
		case "logs":
			err := setLogs(currentSettings.logs, arg)
			if err != nil {
				fmt.Fprintln(s.out, err)
			}
			return
		case "verbosity": //TODO
			i, err := strconv.Atoi(arg)
			if err == nil && 0 <= i && i <= 2 {
				currentSettings.verbosity = i
				return
			}

		case "file":
			if !strings.HasSuffix(arg, ".pdf") {
				arg = arg + ".pdf"
			}
			currentSettings.file = arg
			return
		}
	}

	s.exec_help("set")
}

type inputLevel int

const (
	ProgramL inputLevel = iota
	StatementL
	ExpressionL
)

func (i inputLevel) String() string {
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

func getInputLevel(s string) (inputLevel, bool) {

	switch s {
	case "p", "program":
		return ProgramL, true
	case "s", "statement":
		return StatementL, true
	case "e", "expression":
		return ExpressionL, true
	default:
		return ProgramL, false
	}
}

type inputProcess int

const (
	ParseP inputProcess = iota
	ParseTreeP
	EvalP
	EvalTreeP
	TypeP
	TraceP
)

func (i inputProcess) String() string {
	switch i {
	case ParseP:
		return "parse"
	case ParseTreeP:
		return "parsetree"
	case EvalP:
		return "eval"
	case EvalTreeP:
		return "evaltree"
	case TypeP:
		return "type"
	case TraceP:
		return "trace"
	default:
		return fmt.Sprintf("%d", int(i))
	}
}

func getInputProcess(s string) (inputProcess, bool) {

	switch s {
	case "p", "parse":
		return ParseP, true
	case "ptree", "parsetree":
		return ParseTreeP, true
	case "e", "eval":
		return EvalP, true
	case "etree", "evaltree":
		return EvalTreeP, true
	case "t", "type":
		return TypeP, true
	case "tr", "trace":
		return TraceP, true
	default:
		return EvalP, false
	}
}

type Display int

const (
	ConsD Display = iota
	PdfD
)

func (d Display) String() string {
	switch d {
	case ConsD:
		return "console"
	case PdfD:
		return "pdf"
	default:
		return fmt.Sprintf("%d", int(d))
	}
}

type visDisplays map[Display]bool

func (v visDisplays) String() string {

	displays := make([]string, 0)

	for display, in := range v {
		if in {
			displays = append(displays, display.String())
		}
	}
	return "[" + strings.Join(displays, ", ") + "]"
}

func setDisplays(current visDisplays, arg string) error {

	args := strings.Split(strings.Trim(arg, " "), " ")

	for _, arg := range args {
		switch arg {
		case "+c", "+cons", "+console":
			current[ConsD] = true
		case "-c", "-cons", "-console":
			current[ConsD] = false
		case "+p", "+pdf":
			current[PdfD] = true
		case "-p", "-pdf":
			current[PdfD] = false
		default:
			return errors.New("unknown display " + arg)
		}

	}
	return nil
}

type logs map[inputProcess]bool

func (l logs) String() string {

	ps := make([]string, 0)

	for p, in := range l {
		if in {
			ps = append(ps, p.String())
		}
	}
	return "[" + strings.Join(ps, ", ") + "]"
}

func setLogs(current logs, arg string) error {

	args := strings.Split(strings.Trim(arg, " "), " ")

	for _, arg := range args {
		switch arg {
		case "+type":
			current[TypeP] = true
		case "-type":
			current[TypeP] = false
		case "+evaltree":
			current[EvalTreeP] = true
		case "-evaltree":
			current[EvalTreeP] = false
		case "+parsetree":
			current[ParseTreeP] = true
		case "-parsetree":
			current[ParseTreeP] = false
		case "+trace":
			current[TraceP] = true
		case "-trace":
			current[TraceP] = false
		default:
			return errors.New("unknown display " + arg)
		}

	}
	return nil
}
