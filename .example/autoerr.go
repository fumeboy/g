// +build dev

package example

import "strconv"

type autoerr = error

// comment
func use_autoerr1() (i int, err autoerr) { // comment
	_, err = strconv.Atoi("abc")
	return
}

var use_autoerr2 = func() (err autoerr) {
	_, err = strconv.Atoi("a")
	return
}

var use_autoerr3 = func() (i int, err autoerr) {
	_, err = strconv.Atoi("a")
	var inner = func() (err autoerr) {
		_, err = strconv.Atoi("b")
		var inner_er = func() (err autoerr) {
			_, err = strconv.Atoi("c")
			return
		}
		_ = inner_er
		return
	}
	_ = inner()
	return
}

type r struct {

}

func (receiver r) use_autoerr4() (err autoerr) {
	_, err = strconv.Atoi("a")
	return
}

func init(){
	err := func()(err autoerr){
		_, err = strconv.Atoi("a")
		return
	}()
	_ = err
}