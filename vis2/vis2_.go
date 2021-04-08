package vis2

import (
	"fmt"
	"monkey/ast"
	"monkey/token"
	"reflect"
)

// preconditions:
//   non-circular
//   nodes are structs
//   Tokens have fields Type and Literal
func (v *Visualizer) VisualizeQTree(node ast.Node) string {
	v.display = TEX
	//v.visualizeNode(node)
	v.visualizeFieldValue(node)

	return "\\Tree " + v.out.String()
}

func (v *Visualizer) VisualizeConsTree(node ast.Node) string {
	v.display = CONSOLE
	//v.visualizeNode(node)
	v.visualizeFieldValue(node)
	return v.prefix + v.out.String()
}

func (v *Visualizer) visualizeFieldValue(i interface{}) { //visualize field

	//case nil
	if i == nil {
		v.visualizeNil() // fieldvalue

		//***tex: NIL	done
		//***con: NIL	done

		return
	}
	if reflect.TypeOf(i).Kind() == reflect.Ptr {
		fmt.Printf("!!!!!!!!!! %T %v\n", i, i)
	}
	// case slice
	if reflect.TypeOf(i).Kind() == reflect.Slice {

		//***con: [		done
		//***con: i++	done

		v.beginList()
		values := reflect.Indirect(reflect.ValueOf(i))

		for i := 0; i < values.Len(); i++ {
			//***tex: nli + VALUE	done
			//***con: nli + VALUE	done
			if i > 0 || v.display == CONSOLE {
				v.printInd()
			}
			v.visualizeFieldValue(values.Index(i).Interface())
		}
		v.endList()
		//***con: i--	done
		//***con: nli + ]		done

		return
	}

	switch i := i.(type) {

	case ast.Node:
		v.visualizeNode(i)
		return
	case token.Token:
		v.visualizeToken(i)
		return
	default:
		v.visualizeLeaf(i, false)
		return

	}
}

func (v *Visualizer) visualizeNode(node ast.Node) {

	// case nil
	if node == nil { // unnötig, wenn wir mit visualizeFieldValue starten!
		v.visualizeNil()
		return
	}

	// not nil

	// label()
	v.beginNode(node)

	// children
	if reflect.ValueOf(node).IsNil() {
		v.visualizeNilValue()

	} else { // visualize fields

		nodeContVal := reflect.ValueOf(node).Elem()
		if nodeContVal.Kind() != reflect.Struct {
			v.printW(" NO STRUCT VALUE") // TODO: might be an err ? für Erweiterungen
			return
		}
		nodeContType := nodeContVal.Type()

		for i := 0; i < nodeContVal.NumField(); i++ {
			f := nodeContVal.Field(i)
			// label: fieldname
			fieldname := nodeContType.Field(i).Name
			if v.exclToken && fieldname == "Token" {
				continue
			}

			v.beginField(fieldname)

			// field value
			//fmt.Printf("%d: %s %s = %v\n", i,
			//	nodeContType.Field(i).Name, f.Type(), f.Interface())

			v.visualizeFieldValue(f.Interface())
			v.endField()

		}
	}

	v.endNode()
	//TODO error: any node should be either a Statement an Expression or a Program
}
