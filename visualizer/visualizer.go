package visualizer

import "reflect"

func abbreviateObjectType(objecttype string) string {
	return abbreviateGoObjectType(objecttype)
}

func isLiterallyNil(i interface{}) bool {
	return i == nil
}

func hasNilValue(i interface{}) bool {
	if isLiterallyNil(i) {
		return false
	}
	return reflect.ValueOf(i).IsNil()
}

/* fooling around:

func hasNilValue2(i ast.Node) bool {
	if isLiterallyNil(i) {
		return false
	}
	t := i.(ast.Node)
	return reflect.ValueOf(t).IsNil()
}

func WhatAmI(i interface{}) string {

	if i == nil {
		return "no type"
	}

	nodetype := reflect.TypeOf(i)
	return nodetype.String()
}

func WhatNodeInterfaceAmI(n ast.Node) string {

	if n == nil {
		return "no interface"
	}

	nodetype := reflect.TypeOf(n)

	expr_interface := reflect.TypeOf((*ast.Expression)(nil)).Elem()
	stmt_interface := reflect.TypeOf((*ast.Statement)(nil)).Elem()
	node_interface := reflect.TypeOf((*ast.Node)(nil)).Elem()

	if nodetype.Implements(expr_interface) {
		return expr_interface.String()
	}
	if nodetype.Implements(stmt_interface) {
		return stmt_interface.String()
	}
	if nodetype.Implements(node_interface) {
		return node_interface.String()
	}

	return nodetype.String()
}

func IsNode(i interface{}) bool {
	nodetype := reflect.TypeOf(i)
	node_interface := reflect.TypeOf((*ast.Node)(nil)).Elem()
	return nodetype.Implements(node_interface)
}

func RepresentAsJson(i interface{}, indent string) string {
	json, err := json.MarshalIndent(i, "", indent)
	if err == nil {
		return string(json)
	}
	return ""
}

*/
