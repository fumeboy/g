package define

/*
use as
	type autoerr = error
*/
const (
	AUTOERR = "autoerr" // 自动生成 if err != nil
	SAFEERR = "safeerr" // 自动生成 defer recover，并有 autoerr 的功能
)

const AnonymousFn = "[anonymous_fn]"

/*
use as
	type ctxlog = context.Context
*/
const (
	CTXLOG = "ctxlog"
	CTXLOGERR = "ctxlogerr"
)
