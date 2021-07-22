package loggen

import (
	"github.com/fumeboy/g/define"
	ast "github.com/dave/dst"
)

func ctxTypeMatch(s string) bool{
	return s == define.CTXLOG
}

type funcVisitor struct {
	*ctx
	funcLitVisited map[*ast.FuncLit]struct{}
}

func (v *funcVisitor) Visit(node ast.Node) ast.Visitor {
	var fn_name string
	var fn_type *ast.FuncType
	var next ast.Node
	switch node.(type) {
	case *ast.FuncDecl:
		fn := node.(*ast.FuncDecl)
		if fn.Type.Params != nil && fn.Type.Params.List != nil{
			for _, n := range fn.Type.Params.List {
				if ctxTypeMatch(n.Type.(*ast.Ident).Name) {
					for _, nn := range n.Names {
						if nn.Name == "ctx" {
							fn_name = fn.Name.Name
							next = fn.Body
							goto HIT
						}
					}
				}
			}
		}
	case *ast.GenDecl:
		v1 := node.(*ast.GenDecl)
		for _, s := range v1.Specs {
			if s2, ok := s.(*ast.ValueSpec); ok {
				for i, v2 := range s2.Values {
					if v3, ok := v2.(*ast.FuncLit); ok {
						v.funcLitVisited[v3] = struct{}{}
						if v3.Type.Params != nil && v3.Type.Params.List != nil{
							for _, n := range v3.Type.Params.List {
								if ctxTypeMatch(n.Type.(*ast.Ident).Name) {
									for _, nn := range n.Names {
										if nn.Name == "ctx" {
											fn_name = s2.Names[i].Name
											next = v3.Body
											goto HIT
										}
									}
								}
							}
						}
					}
				}
			}
		}
	case *ast.FuncLit:
		fn := node.(*ast.FuncLit)
		if _,ok := v.funcLitVisited[fn]; !ok {
			v.funcLitVisited[fn] = struct{}{}
			if fn.Type.Params != nil && fn.Type.Params.List != nil{
				for _, n := range fn.Type.Params.List {
					if ctxTypeMatch(n.Type.(*ast.Ident).Name) {
						for _, nn := range n.Names {
							if nn.Name == "ctx" {
								fn_name = define.AnonymousFn
								next = fn.Body
								goto HIT
							}
						}
					}
				}
			}
		}
	}
	return v
HIT:
	// insert at func first line
	ast.Walk(&funcBodyStmtVisitor{fn_name: fn_name, fn_type: fn_type, ctx: v.ctx, returnStmtVisited: map[*ast.ReturnStmt]struct{}{}}, next)
	return v
}
