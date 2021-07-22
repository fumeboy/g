// Code generated. DO NOT EDIT.
// +build !dev


package example

import (
	"example/store"
	"fmt"
	"strconv"

	"github.com/pkg/errors"
)

type autoerr = error
type safeerr = error

type Human struct {
	age int
}

func NewHuman(age string) (h Human, err autoerr) {
	a, err := strconv.Atoi(age)
	if err != nil {
		err = errors.Wrap(err, "fn `NewHuman` failed at human:19 when invoking `strconv.Atoi`"+fmt.Sprintf(" with arg0 = %v; ", age))
		return
	}

	h.age = a
	return
}

func (c Human) EatApple(a *store.Apple) (err safeerr) {
	defer func() {
		if panic_reason := recover(); panic_reason != nil {
			err = errors.New(fmt.Sprintf("fn `EatApple` recover the panic(%v) with arg0 = %v; ", panic_reason, a))
		}
	}()
	if a.Banana != nil {
		panic("not apple, is banana")
	}
	return
}

type BobWantEatApple struct {
	Human
	store.AppleStore
}

func NewBob(h Human, a store.AppleStore) BobWantEatApple {
	return BobWantEatApple{h, a}
}

func NewBobWithWire(age string) (BobWantEatApple, error) {
	human, err := NewHuman(age)
	if err != nil {
		return BobWantEatApple{}, err
	}
	appleStore := store.NewAppleStore()
	bobWantEatApple := NewBob(human, appleStore)
	return bobWantEatApple, nil
}
