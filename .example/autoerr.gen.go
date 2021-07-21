// Code generated. DO NOT EDIT.
// +build !dev


package example

import "strconv"

type autoerr = error

// comment
func use_autoerr1() (i int, err autoerr) { // comment
	_, err = strconv.Atoi("abc")
	if err != nil {
		err = errors.Wrap(err, "fn `use_autoerr1` failed at autoerr:11 when invoking `Atoi`"+fmt.Sprintf(" with arg0 = %v; ", "abc"))
		return
	}

	return
}

var use_autoerr2 = func() (err autoerr) {
	_, err = strconv.Atoi("a")
	if err != nil {
		err = errors.Wrap(err, "fn `use_autoerr2` failed at autoerr:16 when invoking `Atoi`"+fmt.Sprintf(" with arg0 = %v; ", "a"))
		return
	}

	return
}

var use_autoerr3 = func() (i int, err autoerr) {
	_, err = strconv.Atoi("a")
	if err != nil {
		err = errors.Wrap(err, "fn `use_autoerr3` failed at autoerr:21 when invoking `Atoi`"+fmt.Sprintf(" with arg0 = %v; ", "a"))
		return
	}

	var inner = func() (err autoerr) {
		_, err = strconv.Atoi("b")
		if err != nil {
			err = errors.Wrap(err, "fn `inner` failed at autoerr:23 when invoking `Atoi`"+fmt.Sprintf(" with arg0 = %v; ", "b"))
			return
		}

		var inner_er = func() (err autoerr) {
			_, err = strconv.Atoi("c")
			if err != nil {
				err = errors.Wrap(err, "fn `inner_er` failed at autoerr:25 when invoking `Atoi`"+fmt.Sprintf(" with arg0 = %v; ", "c"))
				return
			}

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
	if err != nil {
		err = errors.Wrap(err, "fn `use_autoerr4` failed at autoerr:40 when invoking `Atoi`"+fmt.Sprintf(" with arg0 = %v; ", "a"))
		return
	}

	return
}

func init() {
	err := func() (err autoerr) {
		_, err = strconv.Atoi("a")
		if err != nil {
			err = errors.Wrap(err, "fn `[anonymous_fn]` failed at autoerr:46 when invoking `Atoi`"+fmt.Sprintf(" with arg0 = %v; ", "a"))
			return
		}

		return
	}()
	_ = err
}
