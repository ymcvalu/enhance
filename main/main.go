package main

import (
	enhance "aop"
	"fmt"
)

func main() {
	en := enhance.New()
	en.Register("hook", hook)
	en.Enhance2(fn, "hook").(func(string))("Jim")
}

func fn(name string) {
	fmt.Printf("hello,%s\n", name)
}

func hook(c enhance.Context) {
	c.SetInParam(0, "everyone")
}
