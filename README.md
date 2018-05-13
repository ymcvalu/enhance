# Enhance
a micro lib to provide aop for func

# Example

```go
    func main() {
        en := enhance.New()
        en.Register("log1", logger1)
        en.Register("log2", logger2)
        en.Enhance2(fn, "log1","log2").(func())()
    }

    func fn() {
        fmt.Printf("hello,enhance\n")
    }

    func logger1(c enhance.Context) {
        fmt.Println("before1")
	    c.Call()
        fmt.Println("after1")
    }
    func logger2(c enhance.Context) {
        fmt.Println("before2")
	    c.Call()
        fmt.Println("after2")
    }

output:
    before1
    before2
    hello,enhance
    after2
    after1
```




```go 
    func main() {
        en := enhance.New()
        en.Register("intercept", Intercept)
        en.Enhance2(fn, "intercept").(func())()
    }

    func fn() {
        fmt.Printf("hello,enhance\n")
    }

    func Intercept(c enhance.Context) {
        fmt.Println("abort")
        c.Abort()
    }
    
output:
    abort
```


```go
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

output:
    hello,everyone
```

```go
    func main() {
        en := enhance.New()
        en.Register("hook", hook)
        ret := en.Enhance2(fn, "hook").(func() string)()
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

output:
    hello,Jim
```
    