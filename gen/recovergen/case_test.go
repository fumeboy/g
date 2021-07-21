package recovergen

import (
	"github.com/dave/dst/decorator"
	"go/token"
	"testing"
)

func TestCase1(t *testing.T){
	var code = `
package main

func a() (err safeerr){
	defer func(){
		if a := recover(); a != nil{}
	}()
}
`
	var filename = "a.go"
	dec := decorator.NewDecorator(token.NewFileSet())
	file, _ := dec.Parse(code)
	Edit(filename, dec, file)
	decorator.Print(file)
}

func TestCase2(t *testing.T){
	var code = `
package main

func a() (err safeerr){

}
`
	var filename = "a.go"
	dec := decorator.NewDecorator(token.NewFileSet())
	file, _ := dec.Parse(code)
	Edit(filename, dec, file)
	decorator.Print(file)
}