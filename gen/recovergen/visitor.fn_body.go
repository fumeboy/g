package recovergen

import (
	ast "github.com/dave/dst"
)

// 本文件的程序检查函数体内是否有调用过 recover，如果有就跳过（不是 code generation 的目标）

type funcBodyStmtVisitor struct {
	*ctx
	recoverExsited bool
}

func (v *funcBodyStmtVisitor) Visit(node ast.Node) ast.Visitor {
	switch node.(type) {
	case *ast.FuncDecl, *ast.FuncLit:
		return nil
	case *ast.DeferStmt:
		return (*recoverCallExprVisitor)(v)
	}
	return v
}

type recoverCallExprVisitor struct {
	*ctx
	recoverExsited bool
}

func (v *recoverCallExprVisitor) Visit(node ast.Node) ast.Visitor{
	switch node.(type) {
	case *ast.CallExpr:
		ce := node.(*ast.CallExpr)
		if i,ok := ce.Fun.(*ast.Ident); ok{
			if i.Name == "recover" {
				v.recoverExsited = true
				return nil
			}
		}
	}
	return v
}
