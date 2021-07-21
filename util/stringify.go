package util

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

func StringAST(n ast.Node) string {
	if n == nil {
		return ""
	}
	dst := bytes.NewBuffer([]byte{})
	err := format.Node(dst, token.NewFileSet(), n)
	if err != nil {
		panic(err)
	}
	return (dst.String())
}
