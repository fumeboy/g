// +build dev

package example

import (
	"example/store"
	"github.com/google/wire"
	"strconv"
)

type autoerr = error
type safeerr = error

type Human struct {
	age int
}

func NewHuman(age string) (h Human, err autoerr) {
	a, err := strconv.Atoi(age)
	h.age = a
	return
}

func (c Human) EatApple(a *store.Apple) (err safeerr) {
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
	wire.Build(NewHuman, store.NewAppleStore, NewBob)
	return BobWantEatApple{}, nil
}
