package errgen

import (
	"bytes"
	"fmt"
	"github.com/andreyvit/diff"
	"github.com/dave/dst/decorator"
	"github.com/dave/dst/decorator/resolver/guess"
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
			r := decorator.NewRestorerWithImports("main", guess.New())
			if err := r.Fprint(buf, file); err != nil {
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

func TestCase1(t *testing.T){
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
	r := decorator.NewRestorerWithImports("main", guess.New())
	r.Print(file)
}

func TestCaseWarn(t *testing.T){
	var code = `
package main

func fn() (i int, err autoerr) {
	a,err := strconv.Atoi("abc")
	{
		if err != nil{}
		c,err := strconv.Atoi("abc")
	}
	if err != nil{}
	b,err := strconv.Atoi("abc")
}
`
	var filename = "a.go"
	dec := decorator.NewDecorator(token.NewFileSet())
	file, _ := dec.Parse(code)
	err := Edit(filename, dec, file)
	fmt.Println(err)
}

func TestCaseWarn2(t *testing.T){
	var code = `
package main

func fn() (i int, err autoerr) {
	a,err := call(call1(), other_args)
}
`
	var filename = "a.go"
	dec := decorator.NewDecorator(token.NewFileSet())
	file, _ := dec.Parse(code)
	Edit(filename, dec, file)
	r := decorator.NewRestorerWithImports("main", guess.New())
	r.Print(file)

	/*
	上面的 src 将会生成得到下面这句

	err = errors.Wrap(err, "fn `fn` failed at a.go:5 when invoking `call`"+fmt.Sprintf(" with arg0 = %v; arg1 = %v; ", call1(), other_args))

	可以发现 call1() 被原封不动出现在了 err 里 。。这意味着 call1 会被调用两次，这是不应该的

	后续考虑优化掉或者去掉打印 args
	*/
}
