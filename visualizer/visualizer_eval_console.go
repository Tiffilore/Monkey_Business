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

func RepresentEvalTraceConsole(t *evaluator.Trace, out io.Writer) {
	tab := table.NewWriter()
	tab.SetOutputMirror(out)
	tab.AppendHeader(table.Row{"", "Nodetype", "Node", "Valuetype", "Value"})
	tab.AppendSeparator()

	calls := t.Calls
	exits := t.Exits
	for i := 0; i < t.Steps(); i++ {

		if call, ok := calls[i]; ok {
			tab.AppendRow([]interface{}{
				Red + fmt.Sprintf("call %v", call.Depth) + Reset,
				representNodeType(call.Node, 1),
				fmt.Sprintf("%v", call.Node)})

		} else if exit, ok := exits[i]; ok {
			val := "nil"
			if exit.Val != nil {
				val = exit.Val.Inspect()
			}
			tab.AppendRow([]interface{}{
				Green + fmt.Sprintf("exit %v", exit.Depth) + Reset,
				representNodeType(exit.Node, 1),
				fmt.Sprintf("%v", exit.Node),
				representObjectType(exit.Val, 0),
				val,
			})

		} else {
			fmt.Fprint(out, "We have a problem")
		}

	}
	tab.Render()
}
func TraceEvalConsole(t *evaluator.Trace, out io.Writer, scanner *bufio.Scanner) {

	calls := t.Calls
	exits := t.Exits
	envs := make(map[*object.Environment]int)

	cur_step := 0
	var cur_env_snap *object.Environment
	var cur_env *object.Environment
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
			fmt.Fprintf(out, Red+fmt.Sprintf("call %v", call.Depth)+Reset)
			fmt.Fprint(out, ",")
			if envChanged {
				fmt.Fprintf(out, Red+fmt.Sprintf("e%v: ", envNo)+Reset)
			} else {
				fmt.Fprintf(out, "e%v: ", envNo)
			}
			fmt.Fprintf(out, "%v %v", representNodeType(call.Node, 2), call.Node)
		} else if exit, ok := exits[cur_step]; ok {
			if cur_env != exit.Env || !reflect.DeepEqual(exit.EnvSnap, cur_env_snap) {
				envChanged = true
			}
			cur_env = exit.Env
			cur_env_snap = exit.EnvSnap
			envNo = envs[cur_env]
			fmt.Fprintf(out, Green+fmt.Sprintf("exit %v", exit.Depth)+Reset)
			fmt.Fprint(out, ",")
			if envChanged {
				fmt.Fprintf(out, Red+fmt.Sprintf("e%v: ", envNo)+Reset)
			} else {
				fmt.Fprintf(out, "e%v: ", envNo)
			}
			fmt.Fprintf(out, "%v %v", representNodeType(exit.Node, 2), exit.Node)
			val := "nil"
			if exit.Val != nil {
				val = strings.ReplaceAll(exit.Val.Inspect(), "\n", " ")
			}
			fmt.Fprintf(out, " -> %T %v ", exit.Val, val)
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
			fmt.Fprintf(out, "\tOptions: a: abort, c: continue, e: display environment\n")
		case "c", "":
			cur_step++
			continue
		case "e":
			env_rep := visualizeEnvTable(cur_env_snap, "   ")
			cur_step++
			fmt.Fprint(out, env_rep)
			continue
		default:
			fmt.Fprint(out, "\tUnknown option (h for help)\n")

		}
	}
}

func visualizeEnvTable(env *object.Environment, indent string) string {

	var temp_out bytes.Buffer

	table := GetStoreTable(env)
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
	table = visualizeEnvTable(env.Outer, indent)
	lines = strings.Split(table, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintln(&temp_out, indent, line)
		}
	}

	return temp_out.String()
}
func GetStoreTable(env *object.Environment) string {

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
		nodetype := fmt.Sprintf("%T", object)
		var value string
		if object == nil {
			value = "nil"
		} else {
			value = object.Inspect() //strings.ReplaceAll(object.Inspect(), "\n", "\n\t  ")
		}
		t.AppendRow([]interface{}{key, nodetype, value})
	}

	t.Render()
	return temp_out.String()
}
