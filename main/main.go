package main

import "fmt"
import enhance "aop"

func main() {
	fn := func(msg string) int {
		fmt.Printf("hello,%s\n", msg)
		return 0
	}
	logger1 := func(c enhance.Context) {
		fmt.Println("log1 before")
		c.SetInParam(0, "enhancer")
		c.Call()
		fmt.Println("log1 after")
	}
	logger2 := func(c enhance.Context) {
		fmt.Println("log2 before")
		c.Call()
		c.SetOutParam(0, 1)
		fmt.Println("log2 after")
	}
	en := enhance.New()
	en.Register("log1", logger1)
	en.Register("log2", logger2)
	ret := en.Enhance(fn, "log1", "log2").(func(string) int)("haha")
	fmt.Println(ret)
}
