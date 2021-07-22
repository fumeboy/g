// +build dev

package store

import (
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

func (s *AppleStore) Buy(money int) (apple Apple,err autoerr) {
	err = checkMoney(money)
	b, err := bananaBuy(money - 6)
	apple.Banana = b
	return
}

func checkMoney(money int) (err error) {
	switch money {
	case 1,3,5:
		return errors.New("dont like this num")
	case 2,4,6:
		return errors.New("still dont like this num")
	}
	return nil
}

type banana struct {

}

func bananaBuy(money int) (b *banana, err autoerr) {
	if money < 2{
		return nil, errors.New("not enough to buy")
	}
	return &banana{}, nil
}