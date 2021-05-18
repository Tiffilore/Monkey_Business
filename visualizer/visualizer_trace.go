package visualizer

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/object"
	"reflect"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func TraceInteractive(t *evaluator.Trace, out io.Writer, scanner *bufio.Scanner, verbosity int, goObjType bool) {
	traceInteractive(t, out, scanner, getVerbosity(verbosity), goObjType)
}

func traceInteractive(t *evaluator.Trace, out io.Writer, scanner *bufio.Scanner, verbosity verbosity, goObjType bool) { // before: TraceEvalConsole

	calls := t.Calls
	exits := t.Exits
	envs := make(map[*object.Environment]int)
	var cur_env_snap *object.Environment
	var cur_env *object.Environment

	cur_step := 0

	for cur_step < t.Steps() {
		var envNo int
		envChanged := false
		if call, ok := calls[cur_step]; ok {
			if cur_step > 0 && (cur_env != call.Env || !reflect.DeepEqual(call.EnvSnap, cur_env_snap)) {
				envChanged = true
			}
			cur_env = call.Env
			cur_env_snap = call.EnvSnap
			if no, ok := envs[cur_env]; ok {
				envNo = no
			} else {
				envNo = len(envs)
				envs[cur_env] = envNo
			}
			fmt.Fprint(out, consColorize(fmt.Sprintf("call %v", call.Depth), Red))
			fmt.Fprint(out, ",")
			if envChanged {
				fmt.Fprint(out, consColorize(fmt.Sprintf(" e%v: ", envNo), Red))
			} else {
				fmt.Fprintf(out, " e%v: ", envNo)
			}
			fmt.Fprintf(out, "%v %v", consNode(call.Node, verbosity), call.Node)
		} else if exit, ok := exits[cur_step]; ok {
			if cur_env != exit.Env || !reflect.DeepEqual(exit.EnvSnap, cur_env_snap) {
				envChanged = true
			}
			cur_env = exit.Env
			cur_env_snap = exit.EnvSnap
			envNo = envs[cur_env]
			fmt.Fprint(out, consColorize(fmt.Sprintf("exit %v", exit.Depth), Green))

			fmt.Fprint(out, ",")
			if envChanged {
				fmt.Fprint(out, consColorize(fmt.Sprintf("e%v: ", envNo), Red))
			} else {
				fmt.Fprintf(out, " e%v: ", envNo)
			}
			fmt.Fprintf(out, "%v %v", consNode(exit.Node, verbosity), exit.Node)
			val := "nil"
			if exit.Val != nil {
				val = strings.ReplaceAll(exit.Val.Inspect(), "\n", " ")
			}
			fmt.Fprintf(out, " -> %v %v ", visObjectType(exit.Val, verbosity, goObjType), val)

		} else {
			fmt.Fprint(out, "We have a problem")
		}
		fmt.Fprint(out, " ? ")

		scanned := scanner.Scan()
		if !scanned {
			return
		}
		reply := scanner.Text()
		switch reply {
		case "a":
			return
		case "h":
			fmt.Fprintf(out, "\tOptions: a: abort, [c]: continue, e: display environment\n")
		case "c", "":
			cur_step++
			continue
		case "e":
			env_rep := consEnvTable(cur_env_snap, "   ", verbosity, goObjType)
			cur_step++
			fmt.Fprint(out, env_rep)
			continue
		default:
			fmt.Fprint(out, "\tUnknown option (h for help)\n")

		}
	}
}

func TraceTable(t *evaluator.Trace, out io.Writer, verbosity int, goObjType bool) { // before: RepresentEvalTraceConsole
	traceTable(t, out, getVerbosity(verbosity), goObjType)
}

func traceTable(t *evaluator.Trace, out io.Writer, verbosity verbosity, goObjType bool) { // before: RepresentEvalTraceConsole
	tab := table.NewWriter()
	tab.SetOutputMirror(out)
	tab.AppendHeader(table.Row{"", "Nodetype", "Node", "Objecttype", "Value"})
	tab.AppendSeparator()

	calls := t.Calls
	exits := t.Exits

	for i := 0; i < t.Steps(); i++ {

		if call, ok := calls[i]; ok {
			tab.AppendRow([]interface{}{
				consColorize(fmt.Sprintf("call %v", call.Depth), Red),

				consNode(call.Node, verbosity),
				fmt.Sprintf("%v", call.Node)})
		} else if exit, ok := exits[i]; ok {
			val := "nil"
			if exit.Val != nil {
				val = strings.ReplaceAll(exit.Val.Inspect(), "\n", " ")
			}
			tab.AppendRow([]interface{}{
				consColorize(fmt.Sprintf("exit %v", exit.Depth), Green),
				consNode(exit.Node, verbosity),
				fmt.Sprintf("%v", exit.Node),
				visObjectType(exit.Val, verbosity, goObjType),
				val,
			})
		} else {
			fmt.Fprint(out, "We have a problem")
		}

	}
	tab.Render()
}
