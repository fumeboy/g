package g

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/guess"
	"go/ast"
	"golang.org/x/tools/go/packages"
	"io/ioutil"
	"path/filepath"
	"regexp"
)

var re, _ = regexp.Compile(`//.*\+build.*[^!]dev(\s|\S)`)
var re2, _ = regexp.Compile(`.*\+build.*[^!]dev(\s|\S)`)

type GenerateResult struct {
	PkgPath string
	OutputPath string
	Content []byte
}

func (gen GenerateResult) Commit() error {
	if len(gen.Content) == 0 {
		return nil
	}
	return ioutil.WriteFile(gen.OutputPath, gen.Content, 0666)
}

type Generator struct {
	pkgs []*decorator.Package
}

// G 是 Generator 的构造函数
// 传入参数构造 编译条件 cfg，用于获取该条件下能编译的文件
func G(wd string, env []string, tags string, patterns []string) *Generator {
	cfg := &packages.Config{
		Context:    context.TODO(),
		Mode:       packages.LoadAllSyntax,
		Dir:        wd,
		Env:        env,
		BuildFlags: []string{"-tags=dev"},
	}
	if len(tags) > 0 {
		cfg.BuildFlags[0] += " " + tags
	}
	escaped := make([]string, len(patterns))
	for i := range patterns {
		escaped[i] = "pattern=" + patterns[i]
	}
	pkgs, err := decorator.Load(cfg, escaped...) // 根据 cfg 获取满足编译条件的 pkgs
	fmt.Println(pkgs)
	if err != nil {
		panic(err)
	}
	var errs []error
	for _, p := range pkgs {
		for _, e := range p.Errors {
			errs = append(errs, e)
		}
	}
	if len(errs) > 0 {
		panic(errs)
	}
	return &Generator{pkgs: pkgs}
}

// Gen 是 Gen2 的包装
func (g *Generator) Gen(fns ...func(files []*dst.File, filenames []string, pkg *decorator.Package) error) {
	gs, err := g.Gen2(fns...)
	if err != nil {
		panic(err)
	}
	for _, gg := range gs {
		if err := gg.Commit(); err != nil {
			panic(err)
		}
	}
}

// Gen2 接受一组 特定函数签名的函数 gardeners(修剪 AST 的人)，这些函数会获取 ast.File 并修改 AST
// 顺序执行这些 gardener
// 需要注意，只有第一个 gardener 可以准确获取 AST Node 对应的的文件位置信息，因为经过第一个 editor 修改后，AST 的位置就会发生偏移了
// 所以需要获取位置信息的 errgen 应该放在第一位执行
func (g *Generator) Gen2(gardeners ...func(files []*dst.File, filenames []string, pkg *decorator.Package) error) ([]GenerateResult, error) {
	generated := []GenerateResult{}
	for _, pkg := range g.pkgs {
		outDir, err := detectOutputDir(pkg.GoFiles)
		if err != nil {
			return nil, err
		}
		files := []*dst.File{} // files 和 files_name 是平行的两个slice，一一对应
		files_name := []string{}
		for i, f := range pkg.Syntax {
		L:
			for _, c := range pkg.Decorator.Ast.Nodes[f].(*ast.File).Comments {
				m := re2.Match([]byte(c.Text()))
				if m { // 如果文件带有 // +build dev 就加入到 files 这个 slice
					files = append(files, f)
					_, fileName := filepath.Split(pkg.CompiledGoFiles[i])
					files_name = append(files_name, fileName[:len(fileName)-3])
					break L
				}
			}
		}
		if len(files) == 0 {
			continue
		}
		for _, fn := range gardeners { // 将 files 传给 gardeners
			if err := fn(files, files_name, pkg); err != nil {
				return nil, err
			}
		}
		for i, f := range files { // 将修改后的 files 输出到磁盘文件，后缀为 .gen.go
			g := GenerateResult{}
			g.PkgPath = pkg.PkgPath
			g.OutputPath = filepath.Join(outDir, files_name[i]+".gen.go")

			// 下面三行的逻辑是将 ast.File 对应的 code 写入到 buf
			buf := bytes.NewBuffer([]byte{})
			r := decorator.NewRestorerWithImports(pkg.PkgPath, guess.New())
			if err := r.Fprint(buf, f); err != nil {
				panic(err)
			}

			byt := buf.Bytes() // 下面四行关于 byt 的逻辑是，找到 +build dev 删除，添加 +build !dev
			pos := re.FindIndex(byt)
			byt1, byt2 := byt[:pos[0]], byt[pos[1]:]
			g.Content = append([]byte("// Code generated. DO NOT EDIT.\n// +build !dev\n\n"), append(byt1, byt2...)...)

			generated = append(generated, g)
		}
	}
	return generated, nil
}

func detectOutputDir(paths []string) (string, error) {
	if len(paths) == 0 {
		return "", errors.New("no files to derive output directory from")
	}
	dir := filepath.Dir(paths[0])
	for _, p := range paths[1:] {
		if dir2 := filepath.Dir(p); dir2 != dir {
			return "", fmt.Errorf("found conflicting directories %q and %q", dir, dir2)
		}
	}
	return dir, nil
}
