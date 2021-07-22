package errgen

import (
	"bytes"
	"fmt"
	"github.com/andreyvit/diff"
	"github.com/dave/dst/decorator"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const ext1 = ".src.txt"
const ext2 = ".target.txt"
const dir = "./case/"

func TestCase(t *testing.T) {
	var files = map[string]int{}

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		n := info.Name()
		if len(n) > len(ext2) {
			if n[len(n)-len(ext1):] == ext1 {
				files[n[:len(n)-len(ext1)]]++
			} else if n[len(n)-len(ext2):] == ext2 {
				files[n[:len(n)-len(ext2)]]++
			}
		}
		return nil
	})

	for k, v := range files {
		if v == 2 {
			bytes1, err := ioutil.ReadFile(dir + k + ext1)
			if err != nil {
				panic(err)
			}
			bytes2, err := ioutil.ReadFile(dir + k + ext2)
			if err != nil {
				panic(err)
			}

			dec := decorator.NewDecorator(token.NewFileSet())
			file, _ := dec.Parse(string(bytes1))
			Edit("a.go", dec, file)

			var buf = bytes.NewBuffer([]byte{})
			if err := decorator.Fprint(buf, file); err != nil {
				panic(err)
			}
			if strings.TrimSpace(string(bytes2)) != strings.TrimSpace(buf.String()) {
				fmt.Println(k)
				fmt.Println(len(buf.String()), len(string(bytes2)))
				fmt.Println(diff.LineDiff(string(bytes2), buf.String()))
				panic("src -> target failed")
			}
		}
	}
}

func TestCase1(t *testing.T){ // should panic
	var code = `
package main

func fn() (i int, err autoerr) {
	a, err := strconv.Atoi("abc")

	return
}
`
	var filename = "a.go"
	dec := decorator.NewDecorator(token.NewFileSet())
	file, _ := dec.Parse(code)
	Edit(filename, dec, file)
	decorator.Print(file)
}

func TestCaseWarn(t *testing.T){ // should panic
	var code = `
package main

func fn() (i int, err autoerr) {
	a,err := fn1()
	{
		if err != nil{}
		c,err := fn2()
	}
	if err != nil{}
	b,err := fn3()
}
`
	var filename = "a.go"
	dec := decorator.NewDecorator(token.NewFileSet())
	file, _ := dec.Parse(code)
	Edit(filename, dec, file)
	decorator.Print(file)
}
