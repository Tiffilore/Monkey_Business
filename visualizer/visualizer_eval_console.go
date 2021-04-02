package visualizer

import (
	"bufio"
	"fmt"
	"io"
	"monkey/evaluator"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// TODO: environment, side-effects ?
func RepresentEvalConsole(t *evaluator.Tracer, out io.Writer) {
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
func TraceEvalConsole(t *evaluator.Tracer, out io.Writer, scanner *bufio.Scanner) {

	calls := t.Calls
	exits := t.Exits

	cur_step := 0
	for cur_step < t.Steps() {
		if call, ok := calls[cur_step]; ok {
			fmt.Fprintf(out, Red+fmt.Sprintf("call %v: ", call.Depth)+Reset)
			fmt.Fprintf(out, "%v %v", representNodeType(call.Node, 2), call.Node)

		} else if exit, ok := exits[cur_step]; ok {
			fmt.Fprintf(out, Green+fmt.Sprintf("exit %v: ", exit.Depth)+Reset)
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
			fmt.Fprintf(out, "\tOptions: a: abort, c:continue\n")
		case "c", "":
			cur_step++
			continue
		default:
			fmt.Fprint(out, "\tUnknown option (h for help)\n")

		}
	}
}
