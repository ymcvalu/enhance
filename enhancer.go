package enhance

import (
	"fmt"
	"reflect"
	"sync"
)

type Hook = func(c Context)

func New() *Enhancer {
	return &Enhancer{
		hooks: make(map[string]Hook),
	}
}

type Enhancer struct {
	hooks map[string]Hook
	mu    sync.RWMutex
}

func (e *Enhancer) Register(name string, hook Hook) {
	if hook == nil {
		panic("the hook is nil.")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.hooks[name]; !exists {
		e.hooks[name] = hook
		return
	}
	panic(fmt.Errorf("the name %s has exists.", name))
}

func (e *Enhancer) Enhance(fn interface{}, hooks ...string) interface{} {

	fnTyp := reflect.TypeOf(fn)
	if fnTyp.Kind() != reflect.Func {
		panic("the param fn must be a function")
	}

	fnVal := reflect.ValueOf(fn)
	e.mu.RLock()
	defer e.mu.RUnlock()
	for i := len(hooks) - 1; i >= 0; i-- {
		hook, exists := e.hooks[hooks[i]]
		if !exists {
			continue
		}
		fnVal = e.enhance(fnVal, fnTyp, hook)
	}
	return fnVal.Interface()
}

func (e *Enhancer) enhance(fnVal reflect.Value, fnTyp reflect.Type, hook Hook) reflect.Value {
	proxy := func(params []reflect.Value) []reflect.Value {
		ctx := &context{
			fn:  fnVal,
			ins: params,
		}
		hook(ctx)
		if !ctx.hasFlag(abort) && !ctx.hasFlag(executed) {
			return fnVal.Call(params)
		}
		return ctx.outs
	}

	newFn := reflect.MakeFunc(fnTyp, proxy)
	return newFn
}
