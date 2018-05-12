package enhance

import (
	"fmt"
	"reflect"
	"sync"
)

type Handler = func(c Context)

func New() *Enhancer {
	return &Enhancer{
		fns: make(map[string]Handler),
	}
}

type Enhancer struct {
	fns map[string]Handler
	mu  sync.RWMutex
}

func (e *Enhancer) Register(name string, handler Handler) {
	if handler == nil {
		panic("the handler is nil.")
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	if _, exists := e.fns[name]; !exists {
		e.fns[name] = handler
		return
	}
	panic(fmt.Errorf("the name %s has exists.", name))
}

func (e *Enhancer) Enhance(fn interface{}, handlers ...string) interface{} {
	fnTyp := reflect.TypeOf(fn)
	if fnTyp.Kind() != reflect.Func {
		panic("the param fn must be a function")
	}

	fnVal := reflect.ValueOf(fn)
	e.mu.RLock()
	defer e.mu.RUnlock()
	for i := len(handlers) - 1; i >= 0; i-- {
		handler, exists := e.fns[handlers[i]]
		if !exists {
			continue
		}
		fnVal = enhance(fnVal, fnTyp, handler)
	}
	return fnVal.Interface()
}
