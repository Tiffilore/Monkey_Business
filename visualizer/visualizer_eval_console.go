package visualizer

import (
	"fmt"
	"io"
	"monkey/evaluator"
	"reflect"

	"github.com/jedib0t/go-pretty/v6/table"
)

// we want to output the Tracer

//Schritt 1: textuell

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
				RepresentType(call.Node),
				call.Node})

		} else if exit, ok := exits[i]; ok {
			val := "nil"
			if exit.Val != nil {
				val = exit.Val.Inspect()
			}
			tab.AppendRow([]interface{}{
				Green + fmt.Sprintf("exit %v", exit.Depth) + Reset,
				RepresentType(exit.Node),
				exit.Node,
				reflect.TypeOf(exit.Val),
				val,
			})

		} else {
			fmt.Fprint(out, "We have a problem")
		}

	}
	tab.Render()

}
