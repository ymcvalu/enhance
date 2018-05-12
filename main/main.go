package main

import (
	enhance "aop"
	"fmt"
)

func main() {
	en := enhance.New()
	en.Register("hook", hook)
	ret := en.Enhance(fn, "hook").(func() string)()
	fmt.Println(ret)
}

func fn() string {
	return "Jim"
}

func hook(c enhance.Context) {
	c.Call()
	str := c.OutParam(0).(string)
	c.SetOutParam(0, "hello,"+str)
}
