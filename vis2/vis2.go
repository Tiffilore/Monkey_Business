package vis2

import (
	"bytes"
	"fmt"
	"monkey/ast"
	"monkey/token"
	"reflect"
	"strings"
)

type Display int

const (
	TEX Display = iota
	CONSOLE
	PLAIN
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
	verbosity Verbosity
	exclToken bool
	out       *bytes.Buffer
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
	v.depth++
	v.curIndent = v.curIndent + v.indent
}

func (v *Visualizer) decrIndent() {
	v.depth--
	v.curIndent = v.prefix + strings.Repeat(v.indent, v.depth)
}

func (v *Visualizer) printInd(a ...interface{}) {
	fmt.Fprint(v.out, "\n", v.curIndent)
	fmt.Fprint(v.out, a...)
}

func (v *Visualizer) printW(a ...interface{}) {
	fmt.Fprint(v.out, a...)
}

func (v *Visualizer) beginNode(node ast.Node) {
	//	str := v.representNodeType(node)
	switch v.display {
	case TEX:
		v.printW("[.", v.representNodeType(node))
		//fmt.Fprint(v.out, "\n", v.curIndent)
		//fmt.Fprint(v.out, "[.", str)
	case CONSOLE:
		v.printW(v.representNodeType(node), " {")
		//		fmt.Fprint(v.out, str, " {")
		// default:
		// 	fmt.Fprint(v.out, "\n", v.curIndent)
		// 	fmt.Fprint(v.out, "BEGIN ", str)
	}
	v.incrIndent()
}

func (v *Visualizer) beginField(fieldname string) {
	str := v.representFieldName(fieldname)

	//fmt.Fprint(v.out, "\n", v.curIndent)
	switch v.display {
	case TEX:
		//fmt.Fprint(v.out, "[.", str)
		v.printInd("[.", str, " ")
		v.incrIndent()
	case CONSOLE:
		v.printInd(str, ": ")
		//fmt.Fprint(v.out, str, ": ") //keine neue Zeile danach!
		// kein indent!!
		// default:
		// 	fmt.Fprint(v.out, "BEGIN ", str)
		// 	v.incrIndent()
	}
}

func (v *Visualizer) beginList() {
	//	str := v.representNodeType(node)
	switch v.display {
	case TEX: //nix
	case CONSOLE:
		v.printW("[")
		v.incrIndent()
	}
}

func (v *Visualizer) endList() {
	//	str := v.representNodeType(node)
	switch v.display {
	case TEX: //nix
	case CONSOLE:
		v.decrIndent()
		v.printInd("]")
	}
}

func (v *Visualizer) endField() {
	switch v.display {
	case TEX:
		v.decrIndent()
		v.printInd("]")
	case CONSOLE: //nix
	}
}

func (v *Visualizer) endNode() {
	v.decrIndent()
	switch v.display {
	case TEX:
		v.printInd("]")
	case CONSOLE:
		v.printInd("}")
	}
}

func (v *Visualizer) representNodeType(node ast.Node) string {
	str_nodetype := reflect.TypeOf(node).String()

	if v.verbosity < VVV {
		str_nodetype = strings.TrimLeft(str_nodetype, "*ast.")
	}
	if v.verbosity < VV {
		str_nodetype = abbreviateNodeType(str_nodetype)
	}

	return v.colorNode(str_nodetype, node)
}

func (v *Visualizer) representFieldName(str_fieldname string) string {

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

func (v *Visualizer) visualizeNil() {

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

func (v *Visualizer) visualizeNilValue() {
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

func (s *Visualizer) colorNode(str string, node ast.Node) string {
	switch s.display {
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
