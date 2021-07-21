// +build dev

package example

type safeerr = error

func use_safeerr() (err safeerr){
	defer func(){
		if a := recover(); a != nil{}
	}()
	return
}

func use_safeerr2(a,b string) (err safeerr){
	return
}
