package visualizer

import (
	"bytes"
	"fmt"
	"monkey/evaluator"
	"monkey/object"
	"sort"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
)

// precondition: env != nil
func VisEnvStoreCons(env *object.Environment, verbosity int, goObjType bool) string {
	return consStoreTable(env.Store, getVerbosity(verbosity), goObjType)
}

// prints envs indented!
func consEnvTables(env *object.Environment, indent string, verbosity verbosity, goObjType bool) string {

	var temp_out bytes.Buffer

	if env == nil {
		return indent + "nil"
	}

	table := consStoreTable(env.Store, verbosity, goObjType)
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

	table = consEnvTables(env.Outer, indent, verbosity, goObjType)
	lines = strings.Split(table, "\n")
	for _, line := range lines {
		if line != "" {
			fmt.Fprintln(&temp_out, indent, line)
		}
	}

	return temp_out.String()
}

func consStoreTable(store map[string]object.Object, verbosity verbosity, goObjType bool) string {

	var temp_out bytes.Buffer

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

func texEnvTables(env *object.Environment, verbosity verbosity, goObjType bool) string {

	content := texEnvTableNodes(env, 0, verbosity, goObjType)
	return makeTikz(content)
}

func texEnvTableNodes(env *object.Environment, depth int, verbosity verbosity, goObjType bool) string {

	var temp_out bytes.Buffer

	if env == nil {
		return makeTikzNode("nil", depth)
	}

	table := texStoreTable(env.Store, verbosity, goObjType)
	node := makeTikzNode(table, depth)

	fmt.Fprintln(&temp_out, node)

	outer := texEnvTableNodes(env.Outer, depth+1, verbosity, goObjType)
	fmt.Fprintln(&temp_out, outer)

	return temp_out.String()
}

func texStoreTable(store map[string]object.Object, verbosity verbosity, goObjType bool) string {

	var temp_out bytes.Buffer

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

// precondition: walk through tree; visit envs & name them
// + parameter: env-Liste
func (v *visRun) usedEnvs(trace *evaluator.Trace) string {
	// new buffer
	var out bytes.Buffer
	v.out = &out

	// durch Liste iterieren!
	for _, env := range v.envsOrdered {
		name := v.getEnvName(env)
		k_tex, _ := teXify(fmt.Sprintf("%v", env))
		v.printInd(name, ":\t", k_tex, "\n")

		//Name über get-env-name

		// stelle Abhängigkeiten dar e0 --> e1 --> nil
		v.visEnvDep(env)

		// iteriere durch Trace und stelle Intervalle dar!

	}
	return v.out.String()
}

func (v *visRun) visEnvDep(env *object.Environment) {
	if env == nil {
		v.printW("nil")
	}

}

//
