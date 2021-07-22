package errgen

import (
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/goast"
	"github.com/fumeboy/g/util"
	"github.com/pkg/errors"
	"go/token"
	"strconv"
	"strings"
)

func tempString(filename string, fn_name string, line_number int, err_fn_name string, err_fn_args []string) string {
	/*
		所在函数名
		行号
		调用出错的函数名
		调用出错的函数参数变量名
	*/

	var text = ""
	text += "if err != nil {\n"
	text += "	err = errors.Wrap(err, \"" + fmt.Sprintf("fn `%s` failed at %s:%d", fn_name, filename, line_number)
	if err_fn_name != "" {
		text += fmt.Sprintf(" when invoking `%s`\"", err_fn_name)
		// 打印 args 的部分有风险，考虑优化或关闭
		if len(err_fn_args) > 0 {
			var text_arg = `+ fmt.Sprintf(" with `
			var text_arg2 = []string{}
			for i, arg := range err_fn_args {
				text_arg += "arg" + strconv.Itoa(i) + ` = %v; `
				text_arg2 = append(text_arg2, arg)
			}
			text_arg += "\"," + strings.Join(text_arg2, ",") + ")"
			text += text_arg
		}
	} else {
		text += "\""
	}
	text += ")\n	return\n"
	text += "}\n"
	return text
}

func tempAST(s string, dec *decorator.Decorator) (dst.Stmt, error) {
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
	return file.Decls[1].(*dst.FuncDecl).Body.List[0], nil
}
