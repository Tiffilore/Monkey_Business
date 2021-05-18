package visualizer

import (
	"bytes"
	"fmt"
	"monkey/object"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

func VisEnvStoreCons(env *object.Environment, verbosity int, goObjType bool) string {
	return consStoreTable(env, getVerbosity(verbosity), goObjType)
}

func consStoreTable(env *object.Environment, verbosity verbosity, goObjType bool) string {

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
		objecttype := visObjectType(object, verbosity, goObjType)
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

func consEnvTable(env *object.Environment, indent string, verbosity verbosity, goObjType bool) string {

	var temp_out bytes.Buffer

	table := consStoreTable(env, verbosity, goObjType)
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
	table = consEnvTable(env.Outer, indent, verbosity, goObjType)
	lines = strings.Split(table, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintln(&temp_out, indent, line)
		}
	}

	return temp_out.String()
}

func texEnvTables(env *object.Environment, verbosity verbosity, goObjType bool) string {

	content := texEnvNodes(env, 0, verbosity, goObjType)
	return makeTikz(content)
}

func texEnvNodes(env *object.Environment, depth int, verbosity verbosity, goObjType bool) string {

	var temp_out bytes.Buffer

	table := texStoreTable(env, verbosity, goObjType)
	node := makeTikzNode(table, depth)

	fmt.Fprintln(&temp_out, node)

	if env.Outer == nil {
		outer := makeTikzNode("nil", depth+1)
		fmt.Fprintln(&temp_out, outer)
	} else {
		outer := texEnvNodes(env.Outer, depth+1, verbosity, goObjType)
		fmt.Fprintln(&temp_out, outer)
	}

	return temp_out.String()
}

func texStoreTable(env *object.Environment, verbosity verbosity, goObjType bool) string {

	var temp_out bytes.Buffer
	store := env.Store

	//sort alphabetically
	keys := make([]string, 0, len(store))
	for k := range store {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	table_start := `
	\begin{tabular}{|l|l|l|}\hline
	\textbf{Identifier} & \textbf{Type} & \textbf{Value}\\\hline
	`
	table_end := `	
	\end{tabular}
	`
	fmt.Fprint(&temp_out, table_start)

	for _, key := range keys {
		object := store[key]
		objecttype := visObjectType(object, verbosity, goObjType)
		var value string
		if object == nil {
			value = "<nil>"
		} else {
			value = object.Inspect() //strings.ReplaceAll(object.Inspect(), "\n", "\n\t  ")
		}
		tex_key, _ := teXify(key)
		tex_objecttype, _ := teXify(objecttype)
		tex_value, _ := teXify(value)

		row := fmt.Sprintf("%v & %v & %v \\\\\\hline\n", tex_key, tex_objecttype, tex_value)
		fmt.Fprint(&temp_out, row)
	}

	fmt.Fprint(&temp_out, table_end)

	return temp_out.String()
}

func makeTikzNode(content string, nodeNumber int) string {

	prefix := `\node [draw] `

	if nodeNumber > 0 {
		prefix = fmt.Sprintf("\\node [draw, right=of %v]", nodeNumber-1)
	}

	suffix := "\n};"

	node := prefix + fmt.Sprintf("(%v) {\n", nodeNumber) + content + suffix

	if nodeNumber == 0 {
		return node
	}
	arrow := fmt.Sprintf("\n\\draw [->] (%v) -- (%v);", nodeNumber-1, nodeNumber)

	return node + "\n" + arrow

}
