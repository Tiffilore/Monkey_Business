package visualizer

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

// preconditions:
//   non-circular
//   nodes are structs
//   Tokens have fields Type and Literal

const (
	prefixTeX = ""   //prefix for teX-trees
	indentTeX = "  " //indentation for teX-trees
	DEBUG     = false
)

func ConsParseTree(
	node ast.Node,
	verbosity int,
	inclToken bool,
	prefix string,
	indent string,
) string {

	v := NewVisRun(
		prefix,
		indent,
		getVerbosity(verbosity),
		CONSOLE,
		PARSE,
		inclToken,
		false,
		false,
	)

	tree := v.tree(node, nil)

	return tree
}

func TeXParseTree(
	input string,
	node ast.Node,
	verbosity int,
	inclToken bool,
	file string,
	path string,
) error {

	v := NewVisRun(
		prefixTeX,
		indentTeX,
		getVerbosity(verbosity),
		TEX,
		PARSE,
		inclToken,
		false,
		false,
	)

	tree := v.tree(node, nil)
	if DEBUG {
		fmt.Println(tree)
	}
	document := makeStandalone(texInput(input) + "\n" + tree)

	err := tex2pdf(document, file, path)

	strT := ""
	if inclToken {
		strT = ", with Tokens"
	}
	fmt.Printf("representation of parsetree of %v with verbosity %v%v in file %v\n", node, verbosity, strT, file)
	return err
}

func ConsEvalTree(
	trace *evaluator.Trace,
	verbosity int,
	inclToken bool,
	goObjType bool,
	inclEnv bool,
	prefix string,
	indent string,
) string {

	v := NewVisRun(
		prefix,
		indent,
		getVerbosity(verbosity),
		CONSOLE,
		EVAL,
		inclToken,
		inclEnv,
		goObjType,
	)
	tree := v.tree(trace.GetRoot(), trace)

	return tree
}

func TeXEvalTree(
	input string,
	trace *evaluator.Trace,
	verbosity int,
	inclToken bool,
	goObjType bool,
	inclEnv bool,
	file string,
	path string,
) error {

	v := NewVisRun(
		prefixTeX,
		indentTeX,
		getVerbosity(verbosity),
		TEX,
		EVAL,
		inclToken,
		inclEnv,
		goObjType,
	)
	tree := v.tree(trace.GetRoot(), trace)

	if DEBUG {
		fmt.Println(tree)
	}

	content := texInput(input) + "\n" + tree
	if inclEnv {
		// new buffer
		var out bytes.Buffer
		v.out = &out
		envs := v.envs(trace)
		if DEBUG || true {
			fmt.Println(envs)
		}
		content = content + "\n" + envs
	}

	document := makeStandalone(content)

	err := tex2pdf(document, file, path)

	return err
}

func (v *visRun) tree(node ast.Node, trace *evaluator.Trace) string { //replaces VisualizeQTree and VisualizeConsTree

	v.visitedNodes = make(map[ast.Node]bool)
	if v.process == EVAL {
		v.visitedObjects = make(map[object.Object]bool)
	}
	v.visualizeNode(node, trace, COLLECT)

	v.visitedNodes = make(map[ast.Node]bool)
	if v.process == EVAL {
		v.visitedObjects = make(map[object.Object]bool)
	}
	v.visualizeNode(node, trace, WRITE) // was: visualizeFieldValue(node)

	switch v.display {
	case TEX:
		return makeTikz("\\Tree " + v.out.String())
	case CONSOLE:
		return v.prefix + v.out.String()
	default:
		return "unknown display"
	}
}

func (v *visRun) visualizeNode(node ast.Node, trace *evaluator.Trace, mode mode) {

	_, visited := v.visitedNodes[node]

	// case nil
	if node == nil {
		if mode == WRITE {
			v.visualizeNil()
		}
		return
	}

	if visited && mode == COLLECT { // we do not need to ask whether it is a pointer
		v.createNodeName(node)
	}

	// label node
	v.beginNode(node, trace, visited, mode) // case CONS: with objects!

	if reflect.ValueOf(node).IsNil() {
		if mode == WRITE {
			v.visualizeNilValue()
		}
		v.endNode(visited, mode)
		return
	}

	// if node, ok := node.(*ast.Identifier); ok && v.verbosity < VVV && v.display == TEX { // also if it has already been displayed!
	// 	v.visualizeRoofed(node.String(), mode)
	// }

	// children
	if !visited { // we do not need to ask whether it is a pointer
		v.visitedNodes[node] = true

		// if _, ok := node.(*ast.Identifier); ok && v.verbosity < VVV && v.display == TEX { // also if it has already been displayed!
		// 	v.endNode(visited, mode)

		// 	return
		// }
		nodeContVal := reflect.ValueOf(node).Elem()
		//if nodeContVal.Kind() != reflect.Struct {
		//	v.printW(" NO STRUCT VALUE") // TODO: might be an err ? für Erweiterungen
		//	return
		//}

		nodeContType := nodeContVal.Type()

		for i := 0; i < nodeContVal.NumField(); i++ {
			f := nodeContVal.Field(i)
			// label: fieldname
			fieldname := nodeContType.Field(i).Name
			if !v.inclToken && fieldname == "Token" {
				continue
			}

			v.beginField(fieldname, mode)

			v.visualizeFieldValue(f.Interface(), trace, mode)
			v.endField(mode)

		}

		if v.process == EVAL && v.display == TEX {
			//add objects
			_, exits := v.getCallsAndExits(node, trace)
			for _, exit := range exits {
				no := exit.No
				v.beginVal(no, mode)                     // label
				v.visualizeObject(exit.Val, trace, mode) // Name oder Darstellung
				v.endVal(mode)                           // Klammer zu

			}
		}
	}

	v.endNode(visited, mode)

	//TODO? error: any node should be either a Statement an Expression or a Program
}

func (v *visRun) endVal(mode mode) {
	if mode == COLLECT {
		return
	}

	switch v.display {
	case TEX: //nix
	case CONSOLE: //nix
	}

}
func (v *visRun) beginVal(no int, mode mode) {
	if mode == COLLECT {
		return
	}

	switch v.display {
	case TEX:
		//v.incrIndent()
		v.printInd("\\edge node[auto=left]{\\tiny ", no, "};  ")
	case CONSOLE:
		v.printInd("val ", no, ": ") //TODO

	}

}

//TODO: Methode in den Tracer verschieben!
func (v *visRun) getCallsAndExits(node ast.Node, t *evaluator.Trace) ([]evaluator.Call, []evaluator.Exit) {
	if v.process != EVAL { //should never happen
		return nil, nil
	}

	calls := make([]evaluator.Call, 0)
	exits := make([]evaluator.Exit, 0)

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
func (v *visRun) visualizeFieldValue(i interface{}, trace *evaluator.Trace, mode mode) { //visualize field

	//case nil
	if i == nil {
		if mode == WRITE {
			v.visualizeNil() // fieldvalue
		}
		return
	}

	// case slice
	if reflect.TypeOf(i).Kind() == reflect.Slice {

		values := reflect.Indirect(reflect.ValueOf(i))

		v.beginList(values.Len(), mode)

		for i := 0; i < values.Len(); i++ {
			if v.display == CONSOLE {
				v.printInd()
			}
			v.visualizeFieldValue(values.Index(i).Interface(), trace, mode)
		}
		v.endList(mode)
		return
	}

	switch i := i.(type) {

	case ast.Node:
		v.visualizeNode(i, trace, mode)
		return
	case token.Token:
		v.visualizeToken(i, trace, mode)
		return
	case *object.Environment:
		v.visualizeEnv(i, mode)
		return
	case object.Object:
		v.visualizeObject(i, trace, mode)
		return
	default:
		v.visualizeLeaf(i, false, mode)
		return

	}
}
func (v *visRun) visualizeEnv(env *object.Environment, mode mode) {
	if mode == COLLECT {

		return
	}
	name := v.getEnvName(env)

	switch v.display {
	case TEX:
		v.printW("[.", name, " ]")
	case CONSOLE:
		v.printW(name)
	}
}
func (v *visRun) visualizeObject(obj object.Object, trace *evaluator.Trace, mode mode) {

	// case nil
	if obj == nil {
		if mode == WRITE {
			v.visualizeNilObject()
		}
		return
	}

	if v.verbosity < VVV &&
		(obj == evaluator.FALSE ||
			obj == evaluator.TRUE ||
			obj == evaluator.NULL) {
		v.visualizeSimpleObj(obj, mode)
		return

	}

	if reflect.ValueOf(obj).IsNil() { // so far never happens
		if mode == WRITE {
			v.visualizeNilValue()
		}
		v.endObject(mode)
		return
	}

	if _, ok := v.visitedObjects[obj]; ok && mode == COLLECT { // we do not need to ask whether it is a pointer
		v.createObjectName(obj)
	} // TODO: evtl nach begin object verschieben?

	// label node
	v.beginObject(obj, mode)

	if obj, ok := obj.(*object.Integer); ok && v.verbosity < VVV { // also if it has already been displayed!
		v.visualizeRoofed(obj.Inspect(), mode)

	}
	if obj, ok := obj.(*object.Error); ok && v.verbosity < VVV && v.display == CONSOLE { // also if it has already been displayed!
		v.visualizeRoofed(obj.Inspect(), mode)
	}

	// children --> Nilvalue
	if _, ok := v.visitedObjects[obj]; !ok { // we do not need to ask whether it is a pointer
		v.visitedObjects[obj] = true
		if _, ok := obj.(*object.Integer); ok && v.verbosity < VVV {
			v.endObject(mode)
			return
		}
		if _, ok := obj.(*object.Error); ok && v.verbosity < VVV && v.display == CONSOLE {
			v.endObject(mode)
			return
		}

		if obj, ok := obj.(*object.Error); ok && v.verbosity < VVV && v.display == TEX {
			v.visualizeErrorMsgShort(obj, mode)
			v.endObject(mode)
			return
		}

		if v.display == CONSOLE && mode == WRITE {
			v.printW(consColorize("{", Green))
		}

		objContVal := reflect.ValueOf(obj).Elem()
		//if nodeContVal.Kind() != reflect.Struct {
		//	v.printW(" NO STRUCT VALUE") // TODO: might be an err ? für Erweiterungen
		//	return
		//}

		nodeContType := objContVal.Type()

		for i := 0; i < objContVal.NumField(); i++ {
			f := objContVal.Field(i)
			// label: fieldname
			fieldname := nodeContType.Field(i).Name

			if fieldname == "Ennv" {
				continue
			}

			v.beginField(fieldname, mode)
			//v.printW("%")
			v.visualizeFieldValue(f.Interface(), trace, mode)
			_ = f
			v.endField(mode)

		}

		if v.display == CONSOLE && mode == WRITE {
			v.printInd(consColorize("}", Green))
		}

	}

	v.endObject(mode)
}

func (v *visRun) visualizeErrorMsgShort(obj *object.Error, mode mode) {
	if mode == COLLECT {
		return
	}
	switch v.display {
	case TEX:
		message := obj.Message
		tex_msg, _ := teXify(message)
		message = strings.ReplaceAll(tex_msg, ":", ":\\\\\\small\\it")
		if v.verbosity < VV {
			message = strings.ReplaceAll(message, "INTEGER", "INT")
			message = strings.ReplaceAll(message, "BOOLEAN", "BOOL")
		}
		texStr := fmt.Sprintf("\\small\\it %v", message)
		v.printInd(roofify(texStr))

	case CONSOLE: //TODO
		v.printW(v.colorObj(strings.ToUpper(obj.Inspect()), mode))

	}

}

func (v *visRun) beginObject(obj object.Object, mode mode) {
	if mode == COLLECT {
		return
	}

	switch v.display {
	case TEX:

		v.printW("[.", v.representObjectType(obj, mode))
	case CONSOLE:
		v.printW(v.representObjectType(obj, mode))
	}
	v.incrIndent()

}

func (v *visRun) colorObj(str string, mode mode) string {
	if mode == COLLECT {
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

func (v *visRun) ObjLabel(obj object.Object) string {

	if name, ok := v.getObjectName(obj); ok {
		return name
	}
	return v.objectType(obj)
}

func (v *visRun) objectType(obj object.Object) string {
	vOT := visObjectType(obj, v.verbosity, v.goObjType)
	switch v.display {
	case TEX:
		tex_label, _ := teXify(vOT)
		return tex_label
	case CONSOLE:
		return vOT
	default:
		return "unimplemented display"
	}
}

func (v *visRun) getObjectName(obj object.Object) (string, bool) {

	typestr := v.objectType(obj)
	name, ok := v.namesObjects[typestr][obj]
	return name, ok
}

func (v *visRun) visualizeSimpleObj(obj object.Object, mode mode) {
	if mode == COLLECT {
		return
	}
	v.printW(v.colorObj(strings.ToUpper(obj.Inspect()), mode))

}

func (v *visRun) endObject(mode mode) { // == endNode
	if mode == COLLECT {
		return
	}
	v.decrIndent()
	switch v.display {
	case TEX:
		v.printInd("]")
	case CONSOLE: //nix
	}
}

func (v *visRun) visualizeToken(t token.Token, trace *evaluator.Trace, mode mode) { //TODO: überarbeiten!
	if mode == COLLECT {
		return
	}
	// super-verbose:
	if v.verbosity == VVV && v.display == TEX {
		// label
		v.printW("[.{\\tt ", reflect.TypeOf(t), "}")
		v.incrIndent()

		// fields
		v.printInd("[.")
		v.representFieldName("Type")

		v.incrIndent()
		v.visualizeFieldValue(t.Type, trace, mode) // take method for strings
		v.decrIndent()
		v.printInd("]")

		v.printInd("[.")
		v.representFieldName("Literal")
		v.incrIndent()
		v.visualizeFieldValue(t.Literal, trace, mode) // take method for strings
		v.decrIndent()
		v.printInd("]")

		//
		v.decrIndent()
		v.printInd("]")
	} else {
		switch v.verbosity {

		case VV, VVV:
			v.visualizeLeaf(t, true, mode)
		case V:
			v.visualizeLeaf(t.Literal, false, mode)
		}
	}

}

func (v *visRun) visualizeLeaf(i interface{}, roof bool, mode mode) {
	if mode == COLLECT {
		return
	}
	// string - dependent on verbosity,
	leafValue := fmt.Sprintf("%+v", i)
	if v.display == TEX {
		leafValue, _ = teXify(leafValue)
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

func (v *visRun) beginList(len int, mode mode) {
	if mode == COLLECT {
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

func (v *visRun) endList(mode mode) {
	if mode == COLLECT {
		return
	}
	switch v.display {
	case TEX: //nix
	case CONSOLE:
		v.decrIndent()
		v.printInd("]")
	}
}

func (v *visRun) beginNode(node ast.Node, trace *evaluator.Trace, visited bool, mode mode) {

	switch v.display {
	case TEX:
		if mode == WRITE {
			v.beginNodeTEX(node, trace, visited)
		}
	case CONSOLE:
		v.beginNodeCONSOLE(node, trace, visited, mode)
	}
}

func (v *visRun) getEnvName(env *object.Environment) string {
	if env == nil {
		return "nil"
	}
	for i, e := range v.envsOrdered {
		if e == env {
			switch v.display {
			case CONSOLE:
				return fmt.Sprintf("e%v", i)
			case TEX:
				return fmt.Sprintf("e$_{%v}$", i)
			default:
				return "unknown display"
			}
		}
	}
	v.envsOrdered = append(v.envsOrdered, env)
	return v.getEnvName(env)
}

// only to be called if v.display == TEX and mode == WRITE
func (v *visRun) beginNodeTEX(node ast.Node, trace *evaluator.Trace, visited bool) {

	//display eval-calls and exits only for first occurence
	if v.process == EVAL && !visited {

		left, right := "", ""
		calls, exits := v.getCallsAndExits(node, trace)
		for _, call := range calls {
			eName := v.getEnvName(call.Env)
			left = left + fmt.Sprintf("%v,%v $\\downarrow$ ", call.No, eName)
		}
		for _, exit := range exits {
			eName := v.getEnvName(exit.Env)
			right = right + fmt.Sprintf(" $\\uparrow$%v,%v ", exit.No, eName)
		}
		v.printW("[.{{\\small ", left, "}", v.representNodeType(node), " {\\small ", right, "}}")

	} else {
		v.printW("[.{", v.representNodeType(node), "}")
	}
	v.incrIndent()
}

// only to be called if v.display == CONSOLE
func (v *visRun) beginNodeCONSOLE(node ast.Node, trace *evaluator.Trace, visited bool, mode mode) {

	if mode == WRITE {
		v.printW(v.representNodeType(node))
		v.incrIndent()
	}

	if visited {
		return
	}
	if mode == WRITE {
		v.printW(" {")
	}
	//represent evaluation - steps + objects!
	if v.process == EVAL {
		calls, exits := v.getCallsAndExits(node, trace)
		for _, call := range calls {
			for _, exit := range exits {
				if exit.Id == call.Id {
					if mode == WRITE {
						eName := v.getEnvName(call.Env)
						v.printInd(fmt.Sprintf("[\u2193%v\u2191%v],%v: ", call.No, exit.No, eName))
					}
					v.visualizeObject(exit.Val, trace, mode)
				}

			}
		}

	}

}

func (v *visRun) beginField(fieldname string, mode mode) {
	if mode == COLLECT {
		return
	}
	str := v.representFieldName(fieldname)
	switch v.display {
	case TEX:
		v.printInd("[.", str, " ")
		v.incrIndent()
		v.printInd()
	case CONSOLE:
		v.printInd(str, ": ")

	}
}

func (v *visRun) endField(mode mode) {
	if mode == COLLECT {
		return
	}
	switch v.display {
	case TEX:
		v.decrIndent()
		v.printInd("]")
	case CONSOLE: //nix
	}
}

func (v *visRun) representFieldName(str_fieldname string) string {

	if v.verbosity < VV {
		str_fieldname = abbreviateFieldName(str_fieldname)
	}
	switch v.display {
	case TEX:
		tex_fieldname, _ := teXify(str_fieldname)
		return "{\\small " + tex_fieldname + "}"
	default:
		return str_fieldname
	}
}

func (v *visRun) endNode(visited bool, mode mode) {
	if mode == COLLECT {
		return
	}
	v.decrIndent()
	switch v.display {
	case TEX:
		v.printInd("]")
	case CONSOLE:
		if !visited {
			v.printInd("}")
		}
	}
}

func (v *visRun) visualizeRoofed(str string, mode mode) { // if CONSOLE: use it only for objects, please!
	if mode == COLLECT {
		return
	}
	switch v.display {
	case TEX:
		t_str, _ := teXify(str)
		texStr := fmt.Sprintf("\\small\\it %v", t_str)
		v.printInd(roofify(texStr))

	case CONSOLE: //TODO
		v.printW(consColorize(" { "+str+" }", Green))

	}

}

// only to be called if mode == WRITE
func (v *visRun) representNodeType(node ast.Node) string {
	switch v.display {
	case TEX:
		return texColorNodeStr(v.NodeLabel(node), node)
	case CONSOLE:
		return consColorNodeStr(v.NodeLabel(node), node)
	default:
		return "unknown"
	}
}

func (v *visRun) representObjectType(obj object.Object, mode mode) string {
	if mode == COLLECT { // should not be called
		return ""
	}
	switch v.display {
	case TEX:
		return v.colorObj(v.ObjLabel(obj), mode)
	case CONSOLE:
		return v.colorObj(v.ObjLabel(obj), mode)
	default:
		return "unknown"
	}
}

func (v *visRun) NodeLabel(node ast.Node) string {

	if name, ok := v.getNodeName(node); ok {
		return name
	}

	switch v.display {
	case TEX:
		tex_label, _ := teXify(visNodeType(node, v.verbosity))
		return tex_label
	case CONSOLE:
		return visNodeType(node, v.verbosity)
	default:
		return "unknown"
	}
}

func (v *visRun) createNodeName(node ast.Node) {
	if node == nil { //should never happen
		return
	}
	if _, ok := v.getNodeName(node); ok {
		return
	}

	typestr := visNodeType(node, v.verbosity)
	if _, ok := v.namesNodes[typestr]; !ok {
		v.namesNodes[typestr] = make(map[ast.Node]string)
	}

	switch v.display {
	case CONSOLE:
		v.namesNodes[typestr][node] = typestr + fmt.Sprint(len(v.namesNodes[typestr]))

	case TEX:
		v.namesNodes[typestr][node] = typestr + fmt.Sprintf("$_{%v}$", len(v.namesNodes[typestr]))

	}
}

func (v *visRun) createObjectName(obj object.Object) {
	if obj == nil {
		return
	} // should not happen

	if _, ok := v.getObjectName(obj); ok {
		return
	}

	typestr := v.objectType(obj)
	if _, ok := v.namesObjects[typestr]; !ok {
		v.namesObjects[typestr] = make(map[object.Object]string)
	}

	switch v.display {
	case CONSOLE:
		v.namesObjects[typestr][obj] = typestr + fmt.Sprint(len(v.namesObjects[typestr]))

	case TEX:
		v.namesObjects[typestr][obj] = typestr + fmt.Sprintf("$_{%v}$", len(v.namesObjects[typestr]))

	}

}
func (v *visRun) getNodeName(node ast.Node) (string, bool) {
	typestr := visNodeType(node, v.verbosity)
	name, ok := v.namesNodes[typestr][node]
	return name, ok
}

func (v *visRun) visualizeNil() {
	switch v.display {
	case TEX:
		v.printW("[.", texColorize("nil", "red", "black"), " ]")
	case CONSOLE:
		v.printW(consColorize("nil", Red))
	}
}

func (v *visRun) visualizeNilObject() {
	switch v.display {
	case TEX:
		v.printW("[.", texColorize("nil", "black", "red"), " ]")
	case CONSOLE:
		v.printW(consColorize("nil", Red))
	}
}

func (v *visRun) visualizeNilValue() {

	switch v.display {
	case TEX:
		v.printInd(texColorize("$\\emptyset$", "red", "black"))
	case CONSOLE:
		v.printInd(consColorize("is nil", Red))
	}
}

type visRun struct {
	prefix         string
	indent         string
	curIndent      string
	depth          int
	out            *bytes.Buffer
	verbosity      verbosity
	display        display
	process        process
	inclToken      bool
	inclEnv        bool // vielleicht kein setting?
	goObjType      bool
	visitedNodes   map[ast.Node]bool
	namesNodes     map[string]map[ast.Node]string
	visitedObjects map[object.Object]bool
	namesObjects   map[string]map[object.Object]string
	envsOrdered    []*object.Environment

	// visited --> to avoid printing out cycles
	//-> only for those things that are not ends = don*t call the visualize-Method again
	//nodeFilters []func() bool
}

func NewVisRun(
	prefix string,
	indent string,
	verbosity verbosity,
	display display,
	process process,
	inclToken bool,
	inclEnv bool,
	goObjType bool,
) *visRun {

	var out bytes.Buffer
	namesNodes := make(map[string]map[ast.Node]string)
	namesObjects := make(map[string]map[object.Object]string)
	envsOrdered := make([]*object.Environment, 0)

	return &visRun{
		prefix:       prefix,
		indent:       indent,
		depth:        0,
		curIndent:    prefix,
		out:          &out,
		verbosity:    verbosity,
		display:      display,
		process:      process,
		inclToken:    inclToken,
		inclEnv:      inclEnv,
		goObjType:    goObjType,
		namesNodes:   namesNodes,
		namesObjects: namesObjects,
		envsOrdered:  envsOrdered,
	}
}

type mode int

const (
	COLLECT mode = iota
	WRITE
)

type verbosity int

const (
	V verbosity = iota
	VV
	VVV
)

func getVerbosity(v int) verbosity {

	switch v {
	case 0:
		return V
	case 1:
		return VV
	default:
		return VVV
	}
}

type display int

const (
	TEX display = iota
	CONSOLE
)

type process int

const (
	PARSE process = iota
	EVAL
)

func (v *visRun) incrIndent() {
	v.depth++
	v.curIndent = v.curIndent + v.indent
}

func (v *visRun) decrIndent() {
	v.depth--
	v.curIndent = v.prefix + strings.Repeat(v.indent, v.depth)
}

func (v *visRun) printInd(a ...interface{}) {
	fmt.Fprint(v.out, "\n", v.curIndent)
	fmt.Fprint(v.out, a...)
}

func (v *visRun) printW(a ...interface{}) {
	fmt.Fprint(v.out, a...)
}
