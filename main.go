package g

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/gopackages"
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

func G(wd string, env []string, tags string, patterns []string) *Generator {
	cfg := &packages.Config{
		Context:    context.TODO(),
		Mode:       packages.LoadAllSyntax,
		Dir:        wd,
		Env:        env,
		BuildFlags: []string{"-tags=dev"},
		// TODO(light): Use ParseFile to skip function bodies and comments in indirect packages.
	}
	if len(tags) > 0 {
		cfg.BuildFlags[0] += " " + tags
	}
	escaped := make([]string, len(patterns))
	for i := range patterns {
		escaped[i] = "pattern=" + patterns[i]
	}
	pkgs, err := decorator.Load(cfg, escaped...)
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

func (g *Generator) Gen2(fns ...func(files []*dst.File, filenames []string, pkg *decorator.Package) error) ([]GenerateResult, error) {
	generated := []GenerateResult{}
	for _, pkg := range g.pkgs {
		outDir, err := detectOutputDir(pkg.GoFiles)
		if err != nil {
			return nil, err
		}
		files := []*dst.File{}
		files_name := []string{}
		for i, f := range pkg.Syntax {
		L:
			for _, c := range pkg.Decorator.Ast.Nodes[f].(*ast.File).Comments {
				m := re2.Match([]byte(c.Text()))
				if m {
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
		for _, fn := range fns {
			if err := fn(files, files_name, pkg); err != nil {
				return nil, err
			}
		}
		for i, f := range files {
			g := GenerateResult{}
			g.PkgPath = pkg.PkgPath
			g.OutputPath = filepath.Join(outDir, files_name[i]+".gen.go")
			buf := bytes.NewBuffer([]byte{})
			r := decorator.NewRestorerWithImports(pkg.PkgPath, gopackages.New(pkg.Dir))
			if err := r.Fprint(buf, f); err != nil {
				panic(err)
			}
			byt := buf.Bytes()
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
