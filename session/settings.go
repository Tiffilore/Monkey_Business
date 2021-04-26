package session

import (
	"bytes"
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

func menuSettings() string {

	var out bytes.Buffer

	t := table.NewWriter()
	t.SetOutputMirror(&out)
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
	return out.String()
}

func unset(input string) bool {

	switch strings.Trim(input, " ") {
	case "paste":
		currentSettings.paste = false
		return true
	case "inclToken":
		currentSettings.inclToken = false
		return true
	case "inclEnv":
		currentSettings.inclEnv = false
		return true
	default:
		return false
	}
}

func set(input string) bool {

	splits := strings.SplitN(strings.Trim(input, " "), " ", 2)
	setting := splits[0]

	if len(splits) == 1 {
		switch setting {
		case "paste":
			currentSettings.paste = true
			return true
		case "inclToken":
			currentSettings.inclToken = true
			return true
		case "inclEnv":
			currentSettings.inclEnv = true
			return true
		}

	} else {
		arg := splits[1]
		switch setting {
		case "prompt":
			currentSettings.prompt = arg + " "
			return true
		case "level":
			level, ok := getInputLevel(arg)
			if ok {
				currentSettings.level = level
				return true
			}
		case "process":
			process, ok := getInputProcess(arg)
			if ok {
				currentSettings.process = process
				return true
			}
		case "displays":
			ok := setDisplays(currentSettings.displays, arg)
			if ok {
				return true
			}
		case "logs":
			ok := setLogs(currentSettings.logs, arg)
			if ok {
				return true
			}
		case "verbosity": //TODO
			i, err := strconv.Atoi(arg)
			if err == nil && 0 <= i && i <= 2 {
				currentSettings.verbosity = i
				return true
			}
		case "file":
			if !strings.HasSuffix(arg, ".pdf") {
				arg = arg + ".pdf"
			}
			currentSettings.file = arg
			return true
		}
	}
	return false
}

func reset(input string) bool {
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
		return false
	}
	return true
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

func setDisplays(current visDisplays, arg string) bool {

	args := strings.Split(strings.Trim(arg, " "), " ")

	for _, arg := range args {
		op := arg[0]
		var val bool
		switch op {
		case '+':
			val = true
		case '-':
			val = false
		default:
			return false
		}
		arg = arg[1:]
		switch arg {
<<<<<<< HEAD
		case "c", "cons", "console":
			current[ConsD] = val
		case "p", "pdf":
			current[PdfD] = val
=======
		case "+c", "+cons", "+console":
			current[ConsD] = true
			return true
		case "-c", "-cons", "-console":
			current[ConsD] = false
			return true
		case "+p", "+pdf":
			current[PdfD] = true
			return true
		case "-p", "-pdf":
			current[PdfD] = false
			return true
>>>>>>> 7f9bb427fb102452adbca5c6692355e30a36484f
		default:
			return false
		}

	}
<<<<<<< HEAD
	return true

=======
>>>>>>> 7f9bb427fb102452adbca5c6692355e30a36484f
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

func setLogs(current logs, arg string) bool {

	args := strings.Split(strings.Trim(arg, " "), " ")

	for _, arg := range args {

		op := arg[0]
		var val bool
		switch op {
		case '+':
			val = true
		case '-':
			val = false
		default:
			return false
		}
		arg = arg[1:]
		switch arg {
		case "ptree", "parsetree":
			current[ParseTreeP] = val
		case "etree", "evaltree":
			current[EvalTreeP] = val
		case "t", "type":
			current[TypeP] = val
		case "tr", "trace":
			current[TraceP] = val
		default:
			return false
		}
	}
	return true
}
