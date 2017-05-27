package main

type Ner interface {
	a()
	b(int)
	c(string) string
}

type N int
func (N) a() {
	println("aaaaaaaaa")
}
func(N) b(int) {
	println("bbbbbbbbb")
}
func(N) c(string) string{
	println("cccccccccccc")
	return "cccccccccc"
}

func main() {
	var n N
	var t  Ner = &n
	t.a()
}
