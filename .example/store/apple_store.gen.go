// Code generated. DO NOT EDIT.
// +build !dev


package store

import (
	"fmt"

	"github.com/pkg/errors"
)

type autoerr = error

type Apple struct {
	Banana *banana
}

type AppleStore struct {
}

func NewAppleStore() AppleStore {
	return AppleStore{}
}

func (s *AppleStore) Buy(money int) (apple Apple, err autoerr) {
	err = checkMoney(money)
	if err != nil {
		err = errors.Wrap(err, "fn `Buy` failed at apple_store:24 when invoking `checkMoney`"+fmt.Sprintf(" with arg0 = %v; ", money))
		return
	}

	b, err := bananaBuy(money - 6)
	if err != nil {
		err = errors.Wrap(err, "fn `Buy` failed at apple_store:25 when invoking `bananaBuy`"+fmt.Sprintf(" with arg0 = %v; ", money-6))
		return
	}

	apple.Banana = b
	return
}

func checkMoney(money int) (err error) {
	switch money {
	case 1, 3, 5:
		return errors.New("dont like this num")
	case 2, 4, 6:
		return errors.New("still dont like this num")
	}
	return nil
}

type banana struct {
}

func bananaBuy(money int) (b *banana, err autoerr) {
	if money < 2 {
		return nil, errors.New("not enough to buy")
	}
	return &banana{}, nil
}
