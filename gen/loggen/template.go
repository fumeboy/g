package loggen

import ast "github.com/dave/dst"

/*
	loggen.CtxInfo("fn `fn_name` begin")
*/
func logAtStart(fn_name string, fn_type *ast.FuncType){

}
/*
	loggen.CtxInfo("fn `fn_name` end with error = `%v`, arg0 =`%v`")
*/
func logAtErrorReturn(){

}
/*
	loggen.CtxInfo("fn `fn_name` end")
*/
func logAtNormalReturn(){

}
