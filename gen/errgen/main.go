package errgen

import (
	"fmt"
	"github.com/fumeboy/g/util"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/pkg/errors"
)

type ctx struct {
	dec      *decorator.Decorator
	file     *dst.File
	fileName string
	err      error
	util.Callbacker
}

func Edit(filename string, dec *decorator.Decorator, file *dst.File) error {
	v := &funcVisitor{ctx: &ctx{dec: dec, fileName: filename, file: file}, funcLitVisited: map[*dst.FuncLit]struct{}{}}
	dst.Walk(v, file)
	if v.err != nil{
		return errors.Wrap(v.err, fmt.Sprintf("Edit(), filename = %s", filename))
	}
	if err := v.DO(); err != nil{
		return errors.Wrap(err, fmt.Sprintf("Edit(), filename = %s", filename))
	}
	return nil
}

func Gen(files []*dst.File, filenames []string, pkg *decorator.Package) (err error) {
	for i,f := range files{
		if err := Edit(filenames[i], pkg.Decorator, f); err != nil{
			return errors.Wrap(err, fmt.Sprintf("errgen failed, pkg = %s", pkg.PkgPath))
		}
	}
	return nil
}