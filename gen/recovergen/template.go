package recovergen

import (
	"fmt"
	ast "github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/fumeboy/g/util"
	"github.com/pkg/errors"
	"go/token"
)

/*
	defer func(){
		if panic_reason := recovergen(); panic_reason != nil {
			err = errors.New(fmt.Sprintf("fn `fn_name` recovergen the panic(%v) with arg0 = %v", panic_reason, arg0))
		}
	}()
*/

func tempString(fn_name string, typ *ast.FuncType) string {
	var arg1, arg2 = "", ""
	var num_a = 0
	if typ.Params != nil && len(typ.Params.List) > 0 {
		arg1 = "with "
		for _, item := range typ.Params.List {
			for _, n := range item.Names {
				arg1 += fmt.Sprintf("arg%d = %%v; ", num_a)
				arg2 += n.Name + ","
				num_a++
			}
		}
	}
	return fmt.Sprintf(`defer func(){
		if panic_reason := recover(); panic_reason != nil {
			err = errors.New(fmt.Sprintf("fn %s recover the panic(%%v) %s", panic_reason, %s))
		}
}()`, "`"+fn_name+"`", arg1, arg2)
}

func tempAST(s string, dec *decorator.Decorator) (ast.Stmt, error) {
	var code = `
package main

import (
	"github.com/pkg/errors"
	"fmt"
)

func a(){
%s
}
`
	var f = fmt.Sprintf(code, s)

	dec2 := decorator.NewDecoratorWithImports(token.NewFileSet(), "main", goast.New())
	file, err := dec2.Parse(f)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("tempAST() 失败, code = `%s`", s))
	}

	util.DecCopy(dec2, dec)
	return file.Decls[1].(*ast.FuncDecl).Body.List[0], nil
}