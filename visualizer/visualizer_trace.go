package visualizer

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"monkey/evaluator"
	"monkey/object"
	"reflect"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func VisTraceInteractive(t *evaluator.Trace, out io.Writer, scanner *bufio.Scanner, verbosity int, goObjType bool) { // before: TraceEvalConsole

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
			fmt.Fprintf(out, "%v %v", consColorNode(call.Node, verbosity), call.Node)
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
			fmt.Fprintf(out, "%v %v", consColorNode(exit.Node, verbosity), exit.Node)
			val := "nil"
			if exit.Val != nil {
				val = strings.ReplaceAll(exit.Val.Inspect(), "\n", " ")
			}
			fmt.Fprintf(out, " -> %v %v ", VisObjectType(exit.Val, verbosity, goObjType), val)

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
			env_rep := visualizeEnvTable(cur_env_snap, "   ", verbosity, goObjType)
			cur_step++
			fmt.Fprint(out, env_rep)
			continue
		default:
			fmt.Fprint(out, "\tUnknown option (h for help)\n")

		}
	}
}

func VisTraceTable(t *evaluator.Trace, out io.Writer, verbosity int, goObjType bool) { // before: RepresentEvalTraceConsole
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
				consColorNode(call.Node, verbosity),
				fmt.Sprintf("%v", call.Node)})
		} else if exit, ok := exits[i]; ok {
			val := "nil"
			if exit.Val != nil {
				val = strings.ReplaceAll(exit.Val.Inspect(), "\n", " ")
			}
			tab.AppendRow([]interface{}{
				consColorize(fmt.Sprintf("exit %v", exit.Depth), Green),
				consColorNode(exit.Node, verbosity),
				fmt.Sprintf("%v", exit.Node),
				VisObjectType(exit.Val, verbosity, goObjType),
				val,
			})
		} else {
			fmt.Fprint(out, "We have a problem")
		}

	}
	tab.Render()
}

func visualizeEnvTable(env *object.Environment, indent string, verbosity int, goObjType bool) string {

	var temp_out bytes.Buffer

	table := GetStoreTable(env, verbosity, goObjType)
	lines := strings.Split(table, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintln(&temp_out, indent, line)
		}
	}
	if env.Outer == nil {
		fmt.Fprintln(&temp_out, indent, "--> outer: nil")
		return temp_out.String()
	}

	fmt.Fprintln(&temp_out, indent, "--> outer: ")
	table = visualizeEnvTable(env.Outer, indent, verbosity, goObjType)
	lines = strings.Split(table, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintln(&temp_out, indent, line)
		}
	}

	return temp_out.String()
}

func GetStoreTable(env *object.Environment, verbosity int, goObjType bool) string {

	var temp_out bytes.Buffer
	store := env.Store

	//sort alphabetically
	keys := make([]string, 0, len(store))
	for k := range store {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	t := table.NewWriter()
	t.SetOutputMirror(&temp_out)
	t.AppendHeader(table.Row{"Identifier", "Type", "Value"})
	t.AppendSeparator()

	for _, key := range keys {
		object := store[key]
		objecttype := VisObjectType(object, verbosity, goObjType)
		var value string
		if object == nil {
			value = "<nil>"
		} else {
			value = object.Inspect() //strings.ReplaceAll(object.Inspect(), "\n", "\n\t  ")
		}
		t.AppendRow([]interface{}{key, objecttype, value})
	}

	t.Render()
	return temp_out.String()
}
