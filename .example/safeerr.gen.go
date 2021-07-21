// Code generated. DO NOT EDIT.
// +build !dev


package example

type safeerr = error

func use_safeerr() (err safeerr) {
	defer func() {
		if a := recover(); a != nil {
		}
	}()
	return
}

func use_safeerr2(a, b string) (err safeerr) {
	defer func() {
		if panic_reason := recover(); panic_reason != nil {
			err = errors.New(fmt.Sprintf("fn `use_safeerr2` recover the panic(%v) with arg0 = %v; arg1 = %v; ", panic_reason, a, b))
		}
	}()
	return
}
