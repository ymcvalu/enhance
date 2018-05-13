package enhance

import (
	"fmt"
	"reflect"
)

const (
	executed uint8 = 1 << iota
	abort
)

type Context interface {
	InParam(int) interface{}
	InParams() []interface{}
	SetInParam(int, interface{})
	SetInParams([]interface{})
	InParamsLen() int
	OutParam(int) interface{}
	OutParams() []interface{}
	OutParamsLen() int
	SetOutParam(int, interface{})
	SetOutParams([]interface{})
	Call()
	Abort()
}

//assert
var _ Context = &context{}

type context struct {
	ins  []reflect.Value
	outs []reflect.Value
	fn   reflect.Value
	flag uint8
}

func (c *context) setFlag(flag uint8) {
	c.flag |= flag
}

func (c *context) hasFlag(flag uint8) bool {
	return c.flag&flag != 0
}

func (c *context) InParam(idx int) interface{} {
	if idx < 0 || idx >= len(c.ins) {
		panic("InParam: out of range")
	}

	return c.ins[idx].Interface()

}

func (c *context) InParams() []interface{} {
	params := make([]interface{}, 0, len(c.ins))
	for i := range c.ins {
		params = append(params, c.ins[i].Interface())
	}
	return params
}

func (c *context) SetInParam(idx int, in interface{}) {
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

func (c *context) SetInParams(ins []interface{}) {
	if len(ins) != len(c.ins) {
		panic(fmt.Errorf("SetInParams:len of ins is %d,and need %d", len(ins), len(c.ins)))
	}
	for i := range ins {
		c.SetInParam(i, ins[i])
	}
}

func (c *context) InParamsLen() int {
	return len(c.ins)
}
func (c *context) OutParam(idx int) interface{} {
	if c.outs == nil {
		panic("the out params can't be access now")
	}
	if idx < 0 || idx >= len(c.outs) {
		panic("OutParam: out of range")
	}

	return c.outs[idx].Interface()
}

func (c *context) OutParams() []interface{} {
	if c.outs == nil {
		panic("the out params can't be access now")
	}
	outs := make([]interface{}, 0, len(c.outs))
	for i := range c.outs {
		outs = append(outs, c.outs[i].Interface())
	}
	return outs
}
func (c *context) OutParamsLen() int {
	if c.outs == nil {
		panic("the out params can't be access now")
	}
	return len(c.outs)

}
func (c *context) SetOutParam(idx int, out interface{}) {
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
func (c *context) SetOutParams(outs []interface{}) {
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
func (c *context) Call() {
	if c.hasFlag(abort) {
		return
	}
	outs := c.fn.Call(c.ins)
	c.outs = outs
	c.setFlag(executed)
}
func (c *context) Abort() {
	c.setFlag(abort)
	fnTyp := c.fn.Type()
	c.outs = make([]reflect.Value, 0, fnTyp.NumOut())
	for i := 0; i < fnTyp.NumOut(); i++ {
		c.outs = append(c.outs, reflect.Zero(fnTyp.Out(i)))
	}
}
