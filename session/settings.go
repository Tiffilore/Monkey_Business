package session

import (
<<<<<<< HEAD
	"bytes"
=======
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
	"fmt"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type settings struct {
	prompt    string
	paste     bool
<<<<<<< HEAD
	level     inputLevel
	process   inputProcess
	logs      logs
	displays  visDisplays
=======
	level     InputLevel
	process   InputProcess
	logtree   bool
	logtype   bool
	logtrace  bool
	displays  VisDisplays
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
	verbosity int
	inclToken bool
	inclEnv   bool
	file      string
}

<<<<<<< HEAD
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
=======
func NewSettings() *settings {

	s := settings{
		prompt:  ">> ",
		paste:   false,
		level:   ProgramL,
		process: EvalP,

		logtree:  false,
		logtype:  false,
		logtrace: false,

		displays:  ConsD,
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
		verbosity: 0,
		inclToken: false,
		inclEnv:   false,
		file:      "tree.pdf",
	}

	return &s
}

var defaultSettings *settings
<<<<<<< HEAD
var currentSettings *settings

func init() {
	defaultSettings = newSettings()
	currentSettings = newSettings()
}

func menuSettings() string {

	var out bytes.Buffer

	t := table.NewWriter()
	t.SetOutputMirror(&out)
=======

var currentSettings *settings

func init() {
	defaultSettings = NewSettings()
	currentSettings = NewSettings()
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
	case "logtree":
		currentSettings.logtree = defaultSettings.logtree
	case "logtype":
		currentSettings.logtype = defaultSettings.logtype
	case "logtrace":
		currentSettings.logtrace = defaultSettings.logtrace
	case "displays":
		currentSettings.displays = defaultSettings.displays
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
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
	t.AppendHeader(table.Row{"setting", "current value", "default value"})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"prompt", currentSettings.prompt, defaultSettings.prompt})
	t.AppendRow([]interface{}{"paste", currentSettings.paste, defaultSettings.paste})
	t.AppendRow([]interface{}{"level", currentSettings.level, defaultSettings.level})
	t.AppendRow([]interface{}{"process", currentSettings.process, defaultSettings.process})
<<<<<<< HEAD
	t.AppendRow([]interface{}{"logs", currentSettings.logs, defaultSettings.logs})
=======
	t.AppendRow([]interface{}{"logtree", currentSettings.logtree, defaultSettings.logtree})
	t.AppendRow([]interface{}{"logtype", currentSettings.logtype, defaultSettings.logtype})
	t.AppendRow([]interface{}{"logtrace", currentSettings.logtrace, defaultSettings.logtrace})
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
	t.AppendRow([]interface{}{"displays", currentSettings.displays, defaultSettings.displays})
	t.AppendRow([]interface{}{"verbosity", currentSettings.verbosity, defaultSettings.verbosity})
	t.AppendRow([]interface{}{"inclToken", currentSettings.inclToken, defaultSettings.inclToken})
	t.AppendRow([]interface{}{"inclEnv", currentSettings.inclEnv, defaultSettings.inclEnv})
	t.AppendRow([]interface{}{"file", currentSettings.file, defaultSettings.file})
	//t.SetStyle(table.StyleColoredBright)
	t.Render()
<<<<<<< HEAD
	return out.String()
}

func unset(input string) bool {
=======
}

func (s *Session) exec_unset(input string) {
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b

	switch strings.Trim(input, " ") {
	case "paste":
		currentSettings.paste = false
<<<<<<< HEAD
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
=======
	case "logtree":
		currentSettings.logtree = false
	case "logtype":
		currentSettings.logtype = false
	case "logtrace":
		currentSettings.logtrace = false
	case "inclToken":
		currentSettings.inclToken = false
	case "inclEnv":
		currentSettings.inclEnv = false
	default:
		s.exec_help("reset")
	}
}

func (s *Session) exec_set(input string) {
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b

	splits := strings.SplitN(strings.Trim(input, " "), " ", 2)
	setting := splits[0]

<<<<<<< HEAD
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

=======
	if len(splits) == 0 {

		switch setting {
		case "paste":
			currentSettings.paste = true
			return
		case "logtree":
			currentSettings.logtree = true
			return
		case "logtype":
			currentSettings.logtype = true
			return
		case "logtrace":
			currentSettings.logtrace = true
			return
		case "inclToken":
			currentSettings.inclToken = true
			return
		case "inclEnv":
			currentSettings.inclEnv = true
			return
		}
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
	} else {
		arg := splits[1]
		switch setting {
		case "prompt":
			currentSettings.prompt = arg + " "
<<<<<<< HEAD
			return true
=======
			return
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
		case "level":
			level, ok := getInputLevel(arg)
			if ok {
				currentSettings.level = level
<<<<<<< HEAD
				return true
=======
				return
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
			}
		case "process":
			process, ok := getInputProcess(arg)
			if ok {
				currentSettings.process = process
<<<<<<< HEAD
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
=======
				return
			}
		case "displays": //TODO displays  VisDisplays
			displays, ok := getDisplays(currentSettings.displays, arg)
			if ok {
				currentSettings.displays = displays
				return
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
			}
		case "verbosity": //TODO
			i, err := strconv.Atoi(arg)
			if err == nil && 0 <= i && i <= 2 {
				currentSettings.verbosity = i
<<<<<<< HEAD
				return true
			}
=======
				return
			}

>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
		case "file":
			if !strings.HasSuffix(arg, ".pdf") {
				arg = arg + ".pdf"
			}
			currentSettings.file = arg
<<<<<<< HEAD
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
=======
			return
		}
	}

	s.exec_help("set")
}

type InputLevel int

const (
	ProgramL InputLevel = iota
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
	StatementL
	ExpressionL
)

<<<<<<< HEAD
func (i inputLevel) String() string {
=======
func (i InputLevel) String() string {
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
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

<<<<<<< HEAD
func getInputLevel(s string) (inputLevel, bool) {
=======
func getInputLevel(s string) (InputLevel, bool) {
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b

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

<<<<<<< HEAD
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
=======
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
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
	default:
		return fmt.Sprintf("%d", int(i))
	}
}

<<<<<<< HEAD
func getInputProcess(s string) (inputProcess, bool) {
=======
func getInputProcess(s string) (InputProcess, bool) {
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b

	switch s {
	case "p", "parse":
		return ParseP, true
<<<<<<< HEAD
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
=======
	case "e", "eval":
		return EvalP, true
	case "t", "type":
		return TypeP, true
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
	default:
		return EvalP, false
	}
}

<<<<<<< HEAD
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
=======
type VisDisplays int

const (
	ConsD VisDisplays = iota
	PdfD
	BothD
)

func (v VisDisplays) String() string {
	switch v {
	case ConsD:
		return "[console]"
	case PdfD:
		return "[pdf]"
	case BothD:
		return "[console, pdf]"
	default:
		return fmt.Sprintf("%d", int(v))
	}
}

func getDisplays(current VisDisplays, arg string) (VisDisplays, bool) {
>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b

	args := strings.Split(strings.Trim(arg, " "), " ")

	for _, arg := range args {
<<<<<<< HEAD
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
		case "c", "cons", "console":
			current[ConsD] = val
		case "p", "pdf":
			current[PdfD] = val
		default:
			return false
		}

	}
	return true

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
=======
		switch arg {
		case "+c", "+cons", "+console":
			if current == PdfD {
				current = BothD
			}
		case "-c", "-cons", "-console":
			if current == BothD {
				current = PdfD
			}
		case "+p", "+pdf":
			if current == ConsD {
				current = BothD
			}
		case "-p", "-pdf":
			if current == BothD {
				current = ConsD
			}
		default:
			return BothD, false
		}

	}
	return current, true

>>>>>>> b5cef29d44e884792ea17ab26e469f3298a0cd4b
}
