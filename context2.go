package enhance

import (
	"fmt"
	"reflect"
)

//assert
var _ Context = &context2{}

type context2 struct {
	ins   []reflect.Value
	outs  []reflect.Value
	fn    reflect.Value
	hooks []Hook
	idx   int
}

func (c *context2) InParam(idx int) interface{} {
	if idx < 0 || idx >= len(c.ins) {
		panic("InParam: out of range")
	}

	return c.ins[idx].Interface()

}

func (c *context2) InParams() []interface{} {
	params := make([]interface{}, 0, len(c.ins))
	for i := range c.ins {
		params = append(params, c.ins[i].Interface())
	}
	return params
}

func (c *context2) SetInParam(idx int, in interface{}) {
	if idx < 0 || idx >= len(c.ins) {
		panic("SetInParam: out of range")
	}
	inVal := c.ins[idx]
	inTyp := inVal.Type()
	//the in is nil
	if in == nil {
		c.ins[idx] = reflect.Zero(inTyp)
		return
	}

	typ := reflect.TypeOf(in)
	if typ.AssignableTo(inTyp) {
		c.ins[idx] = reflect.ValueOf(in)
		return
	}
	panic(fmt.Errorf("the type of param is :%s,and need :%s", typ.String(), inTyp.String()))
}

func (c *context2) SetInParams(ins []interface{}) {
	if len(ins) != len(c.ins) {
		panic(fmt.Errorf("SetInParams:len of ins is %d,and need %d", len(ins), len(c.ins)))
	}
	for i := range ins {
		c.SetInParam(i, ins[i])
	}
}

func (c *context2) InParamsLen() int {
	return len(c.ins)
}
func (c *context2) OutParam(idx int) interface{} {
	if c.outs == nil {
		panic("the out params can't be access now")
	}
	if idx < 0 || idx >= len(c.outs) {
		panic("OutParam: out of range")
	}

	return c.outs[idx].Interface()
}

func (c *context2) OutParams() []interface{} {
	if c.outs == nil {
		panic("the out params can't be access now")
	}
	outs := make([]interface{}, 0, len(c.outs))
	for i := range c.outs {
		outs = append(outs, c.outs[i].Interface())
	}
	return outs
}
func (c *context2) OutParamsLen() int {
	if c.outs == nil {
		panic("the out params can't be access now")
	}
	return len(c.outs)

}
func (c *context2) SetOutParam(idx int, out interface{}) {
	if c.outs == nil {
		panic("the out params can't be access now")
	}
	if idx < 0 || idx >= len(c.outs) {
		panic("SetInParam: out of range")
	}
	outVal := c.outs[idx]
	outTyp := outVal.Type()
	//the in is nil
	if out == nil {
		c.outs[idx] = reflect.Zero(outTyp)
		return
	}

	typ := reflect.TypeOf(out)
	if typ.AssignableTo(outTyp) {
		c.outs[idx] = reflect.ValueOf(out)
		return
	}
	panic(fmt.Errorf("the type of param is :%s,and need :%s", typ.String(), outTyp.String()))
}
func (c *context2) SetOutParams(outs []interface{}) {
	if c.outs == nil {
		panic("the out params can't be access now")
	}
	if len(outs) != len(c.outs) {
		panic(fmt.Errorf("SetInParams:len of ins is %d,and need %d", len(outs), len(c.outs)))
	}
	for i := range outs {
		c.SetOutParam(i, outs[i])
	}
}
func (c *context2) Call() {
	idx := c.idx
	if idx < len(c.hooks) {
		c.idx++
		c.hooks[idx](c)
	} else if idx == len(c.hooks) {
		c.idx++
		c.outs = c.fn.Call(c.ins)
	}
}

func (c *context2) Abort() {
	c.idx = len(c.hooks) + 1
	fnTyp := c.fn.Type()
	c.outs = make([]reflect.Value, 0, fnTyp.NumOut())
	for i := 0; i < fnTyp.NumOut(); i++ {
		c.outs = append(c.outs, reflect.Zero(fnTyp.Out(i)))
	}
}
