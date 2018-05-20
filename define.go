package enhance

var enhancer *Enhancer

func init() {
	enhancer = New()
}

func Register(key string, hook Hook) {
	enhancer.Register(key, hook)
}

func Enhance(fn interface{}, hooks ...string) interface{} {
	return enhancer.Enhance(fn, hooks...)
}

func Enhance2(fn interface{}, hooks ...string) interface{} {
	return enhancer.Enhance2(fn, hooks...)
}
