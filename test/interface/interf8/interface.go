package main

import "fmt"

func main() {
	var x interface{} = func(x int) string{
		return fmt.Sprintf("data:%d",x)
	}

	switch v := x.(type) {
	case nil:
		println("nil")
	case *int:
		println(*v)
	case func(int)string:
		println(v(100))
	case fmt.Stringer:
		fmt.Println(v)
	default:
		println("unknown")

	}
}
