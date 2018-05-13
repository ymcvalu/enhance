package enhance

import (
	"reflect"
)

func (e *Enhancer) Enhance2(fn interface{}, hooks ...string) interface{} {
	if fn == nil || len(hooks) == 0 {
		return fn
	}

	fnTyp := reflect.TypeOf(fn)
	if fnTyp.Kind() != reflect.Func {
		panic("the param fn must be a function")
	}

	_hooks := make([]Hook, 0, len(hooks))
	e.mu.RLock()
	for i := range hooks {
		hook, exists := e.hooks[hooks[i]]
		if exists {
			_hooks = append(_hooks, hook)
		}
	}
	e.mu.RUnlock()
	fnVal := reflect.ValueOf(fn)
	proxy := func(params []reflect.Value) []reflect.Value {
		ctx := &context2{
			fn:    fnVal,
			hooks: _hooks,
			ins:   params,
		}
		for ctx.idx <= len(ctx.hooks) {
			ctx.Call()
		}
		return ctx.outs
	}
	newVal := reflect.MakeFunc(fnTyp, proxy)
	return newVal.Interface()
}
