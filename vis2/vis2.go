package vis2

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"monkey/evaluator"
	"monkey/object"
	"monkey/token"
	"reflect"
	"strings"
)

type Display int

const (
	TEX Display = iota
	CONSOLE
)

type Mode int

const (
	COLLECT Mode = iota
	WRITE
)

type Process int

const (
	PARSE Process = iota
	EVAL
)

type Verbosity int

const (
	V Verbosity = iota
	VV
	VVV
)

type Visualizer struct {
	indent    string
	prefix    string
	curIndent string
	depth     int
	display   Display
	mode      Mode
	process   Process
	verbosity Verbosity
	exclToken bool
	out       *bytes.Buffer
	tracer    *evaluator.Tracer
	// visited --> to avoid printing out cycles
	//-> only for those things that are not ends = don*t call the visualize-Method again
	//nodeFilters []func() bool
}

func NewVisualizer(
	prefix string,
	indent string,
	//display Display,
	verb Verbosity,
	exclToken bool) *Visualizer {

	var out bytes.Buffer

	return &Visualizer{
		indent:    indent,
		prefix:    prefix,
		depth:     0,
		out:       &out,
		curIndent: prefix,
		//display:   display,
		exclToken: exclToken,
		verbosity: verb,
	}
}

func (v *Visualizer) incrIndent() {
	if v.mode != COLLECT {
		v.depth++
		v.curIndent = v.curIndent + v.indent
	}
}

func (v *Visualizer) decrIndent() {
	if v.mode != COLLECT {
		v.depth--
		v.curIndent = v.prefix + strings.Repeat(v.indent, v.depth)
	}
}

func (v *Visualizer) printInd(a ...interface{}) {
	if v.mode != COLLECT {
		fmt.Fprint(v.out, "\n", v.curIndent)
		fmt.Fprint(v.out, a...)
	}
}

func (v *Visualizer) printW(a ...interface{}) {
	if v.mode != COLLECT {
		fmt.Fprint(v.out, a...)
	}
}

//TODO: Methode in den Tracer verschieben!
func (v *Visualizer) getCallsAndExits(node ast.Node) ([]evaluator.Call, []evaluator.Exit) {
	if v.process != EVAL { //should never happen
		return nil, nil
	}

	calls := make([]evaluator.Call, 0)
	exits := make([]evaluator.Exit, 0)
	t := v.tracer

	for i := 0; i < t.Steps(); i++ {
		if call, ok := t.Calls[i]; ok {
			if call.Node == node {
				calls = append(calls, call)
			}
		}
		if exit, ok := t.Exits[i]; ok {
			if exit.Node == node {
				exits = append(exits, exit)

			}
		}
	}

	return calls, exits

}
func (v *Visualizer) beginNode(node ast.Node) {
	left, right := "", ""
	if _, ok := visitedNodes[node]; v.process == EVAL && !ok {
		calls, exits := v.getCallsAndExits(node)
		switch v.mode {
		case COLLECT:
			// visit environments
			// only calls relevant
			for _, call := range calls {
				env := call.Env
				// if not visited env, give it a name [if not nil]
				if _, ok := visitedEnvs[env]; !ok {
					v.createEnvName(env)
					visitedEnvs[env] = true
				}
			}
			return
		case WRITE:
			for _, call := range calls { //TODO: different approach for console
				eName, _ := v.getEnvName(call.Env)
				left = left + fmt.Sprintf("%v,%v$\\downarrow$ ", call.No, eName)
			}
			for _, exit := range exits { //TODO: different approach for console
				eName, _ := v.getEnvName(exit.Env)
				right = right + fmt.Sprintf(" $\\uparrow$%v,%v ", exit.No, eName)
			}
		}
	}

	if v.mode == COLLECT {
		return
	}

	switch v.display {
	case TEX:
		v.printW("[.{{\\small ", left, "}", v.representNodeType(node), " {\\small ", right, "}}")
		//	NodeLabel
		v.incrIndent()
	case CONSOLE:
		v.printW(left, v.representNodeType(node), right, " {") //TODO, exp. uparrow!!
		v.incrIndent()
	}
}

func (v *Visualizer) beginObject(obj object.Object) {
	if v.mode == COLLECT {
		return
	}

	switch v.display {
	case TEX:
		v.printW("[.", v.representObjectType(obj))
		v.incrIndent()
	case CONSOLE:
		v.printW(v.representObjectType(obj), " {")
		v.incrIndent()
	}
}

func (v *Visualizer) beginField(fieldname string) {
	if v.mode == COLLECT {
		return
	}
	str := v.representFieldName(fieldname)
	switch v.display {
	case TEX:
		//fmt.Fprint(v.out, "[.", str)
		v.printInd("[.", str, " ")
		v.incrIndent()
		v.printInd()
	case CONSOLE:
		v.printInd(str, ": ")
		//fmt.Fprint(v.out, str, ": ") //keine neue Zeile danach!
		// kein indent!!
		// default:
		// 	fmt.Fprint(v.out, "BEGIN ", str)
		// 	v.incrIndent()
	}
}

func (v *Visualizer) beginVal(no int) {
	if v.mode == COLLECT {
		return
	}

	switch v.display {
	case TEX:
		v.incrIndent()
		v.printInd("\\edge node[auto=left]{\\tiny ", no, "};  ")
	case CONSOLE:
		v.printInd("val ", no, ": ") //TODO

	}

}

func (v *Visualizer) beginList(len int) {
	if v.mode == COLLECT {
		return
	}
	switch v.display {
	case TEX:
		if len == 0 {
			v.printW("%") // empty lines are not allowed
		}
	case CONSOLE:
		v.printW("[")
		v.incrIndent()
	}
}

func (v *Visualizer) endList() {
	if v.mode == COLLECT {
		return
	}
	switch v.display {
	case TEX: //nix
	case CONSOLE:
		v.decrIndent()
		v.printInd("]")
	}
}

func (v *Visualizer) endVal() {
	if v.mode == COLLECT {
		return
	}

	switch v.display {
	case TEX:
		v.decrIndent()
	//	v.printInd("]")
	case CONSOLE: //nix
	}

}

func (v *Visualizer) endField() {
	if v.mode == COLLECT {
		return
	}
	switch v.display {
	case TEX:
		v.decrIndent()
		v.printInd("]")
	case CONSOLE: //nix
	}
}

func (v *Visualizer) endNode() {
	if v.mode == COLLECT {
		return
	}
	v.decrIndent()
	switch v.display {
	case TEX:
		v.printInd("]")
	case CONSOLE:
		v.printInd("}")
	}
}

func (v *Visualizer) endObject() { // == endNode
	if v.mode == COLLECT {
		return
	}
	v.decrIndent()
	switch v.display {
	case TEX:
		v.printInd("]")
	case CONSOLE:
		v.printInd("}")
	}
}

func (v *Visualizer) getNodeName(node ast.Node) (string, bool) {
	typestr := v.stringNodeType(node)
	name, ok := namesNodes[typestr][node]
	return name, ok
}

func (v *Visualizer) getObjectName(obj object.Object) (string, bool) {

	typestr := v.stringObjectType(obj)
	name, ok := namesObjects[typestr][obj]
	return name, ok
}

func (v *Visualizer) getEnvName(env *object.Environment) (string, bool) {
	if env == nil {
		return "nil", true
	}
	name, ok := namesEnvs[env] // every environment must have a name!
	return name, ok
}

func (v *Visualizer) createEnvName(env *object.Environment) {
	if env == nil {
		return
	}

	if _, ok := v.getEnvName(env); ok {
		return
	}

	switch v.display {
	case CONSOLE:
		namesEnvs[env] = fmt.Sprintf("e%v", len(namesEnvs))
	case TEX:
		namesEnvs[env] = fmt.Sprintf("e$_{%v}$", len(namesEnvs))
	}
}

func (v *Visualizer) createObjectName(obj object.Object) {
	if obj == nil ||
		obj == evaluator.NULL ||
		obj == evaluator.FALSE ||
		obj == evaluator.TRUE {
		return
	}

	if _, ok := v.getObjectName(obj); ok {
		return
	}

	typestr := v.stringObjectType(obj)
	if _, ok := namesObjects[typestr]; !ok {
		namesObjects[typestr] = make(map[object.Object]string)
	}

	switch v.display {
	case CONSOLE:
		namesObjects[typestr][obj] = typestr + fmt.Sprint(len(namesObjects[typestr]))

	case TEX:
		namesObjects[typestr][obj] = typestr + fmt.Sprintf("$_{%v}$", len(namesObjects[typestr]))

	}

}

func (v *Visualizer) createNodeName(node ast.Node) {
	if node == nil { //should never happen
		return
	}
	if _, ok := v.getNodeName(node); ok {
		return
	}

	typestr := v.stringNodeType(node)
	if _, ok := namesNodes[typestr]; !ok {
		namesNodes[typestr] = make(map[ast.Node]string)
	}

	switch v.display {
	case CONSOLE:
		//names[node] = v.stringNodeType(node) + fmt.Sprint(len(names))
		namesNodes[typestr][node] = typestr + fmt.Sprint(len(namesNodes[typestr]))

	case TEX:
		namesNodes[typestr][node] = typestr + fmt.Sprintf("$_{%v}$", len(namesNodes[typestr]))

	}
}

func (v *Visualizer) stringNodeType(node ast.Node) string {

	var str_nodetype string

	str_nodetype = reflect.TypeOf(node).String()

	if v.verbosity < VVV {
		str_nodetype = strings.TrimLeft(str_nodetype, "*ast.")
	}
	if v.verbosity < VV {
		str_nodetype = abbreviateNodeType(str_nodetype)
	}

	return str_nodetype
}

func (v *Visualizer) stringObjectType(obj object.Object) string {
	// nil sollte nicht passieren
	var str_objtype string
	str_objtype = reflect.TypeOf(obj).String()

	if v.verbosity < VVV {
		str_objtype = strings.TrimLeft(str_objtype, "*object.")
	}
	if v.verbosity < VV {
		str_objtype = abbreviateObjectType(str_objtype)
	}

	return str_objtype
}

func (v *Visualizer) representObjectType(obj object.Object) string {
	if v.mode == COLLECT { // should not be called
		return ""
	}

	return v.colorObj(v.ObjLabel(obj))
}

func (v *Visualizer) ObjLabel(obj object.Object) string {

	if name, ok := v.getObjectName(obj); ok {
		return name
	}
	return v.stringObjectType(obj)
}

func (v *Visualizer) representNodeType(node ast.Node) string {
	if v.mode == COLLECT { // should not be called
		return ""
	}

	return v.colorNode(v.NodeLabel(node), node)
}

func (v *Visualizer) NodeLabel(node ast.Node) string {

	if name, ok := v.getNodeName(node); ok {
		return name
	}
	return v.stringNodeType(node)
}

func (v *Visualizer) representFieldName(str_fieldname string) string {
	if v.mode == COLLECT {
		return ""
	}
	if v.verbosity < VV {
		str_fieldname = abbreviateFieldName(str_fieldname)
	}
	switch v.display {
	case TEX:
		return "{\\small " + str_fieldname + "}"
	default:
		return str_fieldname
	}
}

func (v *Visualizer) visualizeEnv(env *object.Environment) {
	if v.mode == COLLECT {
		return
	}
	name, _ := v.getEnvName(env)
	switch v.display {
	case TEX:
		v.printW("[.", name, " ]")
	case CONSOLE:
		v.printW(name)
	}
}

func (v *Visualizer) visualizeNil() {
	if v.mode == COLLECT {
		return
	}
	// display ?
	switch v.display {
	case TEX:
		v.printW("[.", texColorize("nil", "red", "black"), " ]")
		//fmt.Fprint(v.out, "\n", v.curIndent)
		//fmt.Fprint(v.out, texColorize("black", "red", "nil"))
	case CONSOLE:
		v.printW(consColorize("nil", Red))

		//fmt.Fprint(v.out, consColorize("nil", Red))
		// default:
		// 	fmt.Fprint(v.out, "\n", v.curIndent)

		// 	v.print("nil")

	}
}

func (v *Visualizer) visualizeSimpleObj(obj object.Object) {
	if v.mode == COLLECT {
		return
	}

	if obj == nil {
		switch v.display {
		case TEX:
			v.printW("[.", texColorize("nil", "black", "red"), " ]")
		case CONSOLE:
			v.printW(consColorize("nil", Red)) //TODO
		}
		return
	}
	v.printW(v.colorObj(strings.ToUpper(obj.Inspect())))

}

func (v *Visualizer) visualizeNilValue() {
	if v.mode == COLLECT {
		return
	}
	switch v.display {
	case TEX:
		v.printInd(texColorize("$\\emptyset$", "red", "black"))
	default:
		v.printInd(consColorize("is nil", Red))
	}
}

/*
	type Token struct {
		Type    TokenType
		Literal string
	}
*/
func (v *Visualizer) visualizeToken(t token.Token) { //TODO: überarbeiten!
	if v.mode == COLLECT {
		return
	}
	// super-verbose:
	if v.verbosity == VVV && v.display == TEX {
		// label
		v.printInd("[.{\\tt ", reflect.TypeOf(t), "}")
		v.incrIndent()
		// fields
		v.printInd("[.Type")
		v.incrIndent()
		v.visualizeFieldValue(t.Type) // take method for strings
		v.decrIndent()
		v.printInd("]")
		v.printInd("[.Literal")
		v.incrIndent()
		v.visualizeFieldValue(t.Literal) // take method for strings
		v.decrIndent()
		v.printInd("]")
		//
		v.decrIndent()
		v.printInd("]")
	} else {
		switch v.verbosity {

		case VV, VVV:
			v.visualizeLeaf(t, true)
		case V:
			v.visualizeLeaf(t.Literal, false)
		}
	}

}

func (v *Visualizer) visualizeLeaf(i interface{}, roof bool) {
	if v.mode == COLLECT {
		return
	}
	// string - dependent on verbosity
	leafValue := fmt.Sprintf("%+v", i)
	if v.display == TEX {
		leafValue = teXify(leafValue)
	}
	leafType := fmt.Sprintf("%T", i)
	var leafStr string
	if v.verbosity < VVV {
		leafStr = leafValue
	} else {
		leafStr = leafType + " " + leafValue
	}

	// display
	switch v.display {
	case TEX:
		texStr := fmt.Sprintf("\\underline{\\it %v}", leafStr)
		if roof {
			texStr = roofify(texStr)
		}
		v.printW(texStr)
	case CONSOLE:
		v.printW(leafStr)
	}
}

func (v *Visualizer) colorNode(str string, node ast.Node) string {
	if v.mode == COLLECT {
		return ""
	}
	switch v.display {
	case TEX:
		if _, ok := node.(ast.Expression); ok {
			return texColorize(str, "bluish", "black")
		} else if _, ok := node.(ast.Statement); ok {
			return texColorize(str, "yellish", "black")
		} else if _, ok := node.(*ast.Program); ok {
			return texColorize(str, "dbluish", "white")
		} else { //new nodes that fall under neither of these cases
			return texColorize(str, "red", "black")
		}
	case CONSOLE:
		if _, ok := node.(ast.Expression); ok {
			return consColorize(str, Cyan)
		} else if _, ok := node.(ast.Statement); ok {
			return consColorize(str, Yellow)
		} else if _, ok := node.(*ast.Program); ok {
			return consColorize(str, Blue)
		} else { //new nodes that fall under neither of these cases
			return consColorize(str, Red)
		}
	default:
		return str
	}
}

func (v *Visualizer) colorObj(str string) string {
	if v.mode == COLLECT {
		return ""
	}
	switch v.display {
	case TEX:
		return texColorize(str, "black", "white")
	case CONSOLE:
		return consColorize(str, Green) //TODO
	default:
		return str
	}
}
