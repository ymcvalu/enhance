package enhance

import "reflect"

func enhance(fnVal reflect.Value, fnTyp reflect.Type, handler Handler) reflect.Value {
	proxy := func(params []reflect.Value) []reflect.Value {
		ctx := new(context)
		ctx.fn = fnVal
		ctx.ins = params
		handler(ctx)
		if !ctx.hasFlag(abort) && !ctx.hasFlag(executed) {
			return fnVal.Call(params)
		}
		return ctx.outs
	}

	newFn := reflect.MakeFunc(fnTyp, proxy)
	return newFn
}
