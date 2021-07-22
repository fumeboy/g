package loggen

import (
	ast "github.com/dave/dst"
)

type funcBodyStmtVisitor struct {
	fn_name string
	fn_type *ast.FuncType
	*ctx
	returnStmtVisited map[*ast.ReturnStmt]struct{}
}

func (v *funcBodyStmtVisitor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.FuncDecl, *ast.FuncLit:
		return nil
	case *ast.IfStmt:
		return (*ifErrBlockVisitor)(v)
	case *ast.ReturnStmt:
		if _,ok := v.returnStmtVisited[node.(*ast.ReturnStmt)]; !ok{

		}
	}
	return v
}

type ifErrBlockVisitor funcBodyStmtVisitor

func (v *ifErrBlockVisitor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.FuncDecl, *ast.FuncLit:
		return nil
	case *ast.ReturnStmt:
		v.returnStmtVisited[node.(*ast.ReturnStmt)] = struct{}{}

	}
	return v
}
