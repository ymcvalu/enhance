package main

import (
	enhance "aop"
	"fmt"
)

func main() {
	en := enhance.New()
	en.Register("hook", hook)
	msg := en.Enhance(fn, "hook").(func(string) string)("Jim")
	fmt.Println(msg)
}

func fn(name string) string {
	return fmt.Sprintf("hello,%s", name)
}

func hook(c enhance.Context) {

	c.SetInParam(0, "everyone")
}
