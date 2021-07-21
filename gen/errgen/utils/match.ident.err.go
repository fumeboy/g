package utils

import ast "github.com/dave/dst"

type errMatch struct {
	matched bool
}

func (v *errMatch) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.Ident:
		if node.(*ast.Ident).Name == "err" {
			v.matched = true
			return nil
		}
	}
	return v
}

func ErrMatch(n ast.Node) bool {
	if n == nil{
		return false
	}
	v := &errMatch{}
	ast.Walk(v, n)
	return v.matched
}
