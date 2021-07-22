package recovergen

import (
	"github.com/fumeboy/g/define"
	ast "github.com/dave/dst"
)

// 本文件的程序搜索函数签名里带有 (err safeerr) 的函数声明、函数字面量

func errTypeMatch(s string) bool{
	return s == define.SAFEERR
}

type funcVisitor struct {
	*ctx
	funcLitVisited map[*ast.FuncLit]struct{}
}

func (v *funcVisitor) Visit(node ast.Node) ast.Visitor {
	var fn_name string
	var fn_type *ast.FuncType
	var target *ast.BlockStmt
	switch node.(type) {
	case *ast.FuncDecl:
		fn := node.(*ast.FuncDecl)
		if fn.Type.Results != nil && fn.Type.Results.List != nil{
			for _, n := range fn.Type.Results.List {
				if id,ok := n.Type.(*ast.Ident); ok && errTypeMatch(id.Name) {
					for _, nn := range n.Names {
						if nn.Name == "err" {
							fn_name = fn.Name.Name
							fn_type = fn.Type
							target = fn.Body
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
						if v3.Type.Results != nil && v3.Type.Results.List != nil{
							for _, n := range v3.Type.Results.List {
								if id,ok := n.Type.(*ast.Ident); ok && errTypeMatch(id.Name){
									for _, nn := range n.Names {
										if nn.Name == "err" {
											fn_name = s2.Names[i].Name
											fn_type = v3.Type
											target = v3.Body
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
			if fn.Type.Results != nil && fn.Type.Results.List != nil{
				for _, n := range fn.Type.Results.List {
					if id,ok := n.Type.(*ast.Ident); ok && errTypeMatch(id.Name) {
						for _, nn := range n.Names {
							if nn.Name == "err" {
								fn_name = define.AnonymousFn
								fn_type = fn.Type
								target = fn.Body
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
	v2 := &funcBodyStmtVisitor{ctx: v.ctx}
	ast.Walk(v2, target)
	if !v2.recoverExsited{
		// insert at func first line
		temp, err := tempAST(tempString(fn_name, fn_type), v.dec)
		if err != nil{
			v.err = err
			return nil
		}
		v.AddCallback(func() error{
			var l = []ast.Stmt{temp}
			target.List = append(l, target.List...)
			return nil
		})
	}
	return v
}
