package session

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

type settings struct {
	prompt    string
	paste     bool
	level     InputLevel
	process   InputProcess
	logtree   bool
	logtype   bool
	logtrace  bool
	displays  VisDisplays
	verbosity int
	inclToken bool
	inclEnv   bool
	file      string
}

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
	t.AppendHeader(table.Row{"setting", "current value", "default value"})
	t.AppendSeparator()
	t.AppendRow([]interface{}{"prompt", currentSettings.prompt, defaultSettings.prompt})
	t.AppendRow([]interface{}{"paste", currentSettings.paste, defaultSettings.paste})
	t.AppendRow([]interface{}{"level", currentSettings.level, defaultSettings.level})
	t.AppendRow([]interface{}{"process", currentSettings.process, defaultSettings.process})
	t.AppendRow([]interface{}{"logtree", currentSettings.logtree, defaultSettings.logtree})
	t.AppendRow([]interface{}{"logtype", currentSettings.logtype, defaultSettings.logtype})
	t.AppendRow([]interface{}{"logtrace", currentSettings.logtrace, defaultSettings.logtrace})
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

	splits := strings.SplitN(strings.Trim(input, " "), " ", 2)
	setting := splits[0]

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
		case "displays": //TODO displays  VisDisplays
			displays, ok := getDisplays(currentSettings.displays, arg)
			if ok {
				currentSettings.displays = displays
				return
			}
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

func getInputLevel(s string) (InputLevel, bool) {

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

func getInputProcess(s string) (InputProcess, bool) {

	switch s {
	case "p", "parse":
		return ParseP, true
	case "e", "eval":
		return EvalP, true
	case "t", "type":
		return TypeP, true
	default:
		return EvalP, false
	}
}

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

	args := strings.Split(strings.Trim(arg, " "), " ")

	for _, arg := range args {
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

}
