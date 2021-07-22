package example

import (
	"fmt"
	"testing"
)

// go test -v -run=Test
func Test1(t *testing.T){
	bob,err := NewBobWithWire("a12")
	if err != nil{
		fmt.Println(err)
	}
	_ = bob
}

func Test2(t *testing.T){
	bob,err := NewBobWithWire("12")
	if err != nil{
		fmt.Println(err)
	}
	apple, err := bob.Buy(3)
	if err != nil{
		fmt.Println(err)
	}
	_ = apple
}

func Test3(t *testing.T){
	bob,err := NewBobWithWire("12")
	if err != nil{
		fmt.Println(err)
	}
	apple, err := bob.Buy(7)
	if err != nil{
		fmt.Println(err)
	}
	_ = apple
}

func Test4(t *testing.T){
	bob,err := NewBobWithWire("12")
	if err != nil{
		fmt.Println(err)
	}
	apple, err := bob.Buy(8)
	if err != nil{
		fmt.Println(err)
	}
	err = bob.EatApple(&apple)
	if err != nil{
		fmt.Println(err)
	}
}
