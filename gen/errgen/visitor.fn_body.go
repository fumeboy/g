package errgen

import (
	"fmt"
	"github.com/fumeboy/g/define"
	"github.com/fumeboy/g/gen/errgen/utils"
	"github.com/fumeboy/g/util"
	"github.com/dave/dst"
	"github.com/pkg/errors"
	"go/token"
)

type funcBodyStmtVisitor struct {
	fn_name string
	*ctx
}

func (v *funcBodyStmtVisitor) Visit(node dst.Node) dst.Visitor {
	switch node.(type) {
	case *dst.FuncDecl, *dst.FuncLit:
		return nil
	case *dst.BlockStmt:
		b := node.(*dst.BlockStmt)
		var inse []struct {
			index            int
			err_call_fn_name string
			err_call_fn_args []string
			token_pos        token.Pos
		}
	L:
		for i := 0; i < len(b.List); i++ {
			sent := b.List[i]
			var err_call dst.Node
			var err_call_fn_name string
			var err_call_args []string
			var end_pos token.Pos
			switch sent.(type) {
			case *dst.AssignStmt:
				s := sent.(*dst.AssignStmt)
				for j, n := range s.Lhs {
					if n.(*dst.Ident).Name == "err" {
						if len(s.Rhs) > 1 {
							err_call = s.Rhs[j]
						} else {
							err_call = s.Rhs[0]
						}
						end_pos = (v.dec.Ast.Nodes[s].End())
						goto HIT
					}
				}
			case *dst.DeclStmt:
				s := sent.(*dst.DeclStmt).Decl
				if s2, ok := s.(*dst.GenDecl); ok {
					for _, n := range s2.Specs {
						vs := n.(*dst.ValueSpec)
						for j, n2 := range vs.Names {
							if n2.Name == "err" {
								if len(vs.Values) > 1 {
									err_call = vs.Values[j]
								} else {
									err_call = vs.Values[0]
								}
								end_pos = (v.dec.Ast.Nodes[s].End())
								goto HIT
							}
						}
					}
				}
			case *dst.IfStmt:
				if s, ok := sent.(*dst.IfStmt); ok && utils.ErrMatch(s.Cond) && !utils.ErrMatch(s.Init) {
					if len(inse) >= 1 {
						inse = inse[:len(inse)-1]
					} else {
						v.err = errors.New(fmt.Sprintf("遇到 assign err 压栈， 遇到 if err 出栈；栈空时出栈. fn_name = %v", v.fn_name) )
						return nil
					}
				}
			}
			continue L
		HIT:
			if rv, ok := err_call.(*dst.CallExpr); ok {
				switch rv.Fun.(type) {
				case *dst.SelectorExpr:
					err_call_fn_name = util.StringAST(v.dec.Ast.Nodes[rv.Fun])
				case *dst.Ident:
					err_call_fn_name = rv.Fun.(*dst.Ident).Name
				case *dst.FuncLit:
					err_call_fn_name = define.AnonymousFn
				}
				for _, a := range rv.Args {
					err_call_args = append(err_call_args, util.StringAST(v.dec.Ast.Nodes[a]))
				}
			}
			inse = append(inse, struct {
				index            int
				err_call_fn_name string
				err_call_fn_args []string
				token_pos        token.Pos
			}{index: i, err_call_fn_name: err_call_fn_name, err_call_fn_args: err_call_args, token_pos: end_pos})

		}
		if len(inse) > 0 {
			v.AddCallback(func() error{
				var newList = make([]dst.Stmt, 0, len(b.List)+len(inse))
				for i, ins := range inse {
					if i == 0 {
						newList = append(newList, b.List[:inse[0].index+1]...)
					} else {
						newList = append(newList, b.List[inse[i-1].index+1:inse[i].index+1]...)
					}
					new1, err := tempAST(tempString(v.fileName, v.fn_name, v.dec.Fset.Position(ins.token_pos).Line, ins.err_call_fn_name, ins.err_call_fn_args), v.dec)
					if err != nil{
						return errors.Wrap(err,fmt.Sprintf("callback failed, fileName = %s, fn_name = %s, line_no = %d, err_call_fn_name = %s, err_call_fn_args = %v", v.fileName, v.fn_name, v.dec.Fset.Position(ins.token_pos).Line, ins.err_call_fn_name, ins.err_call_fn_args))
					}
					newList = append(newList, new1.(dst.Stmt))
				}
				b.List = append(newList, b.List[inse[len(inse)-1].index+1:]...)
				return nil
			})
		}
	}
	return v
}
