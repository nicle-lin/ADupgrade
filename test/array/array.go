package main

import "fmt"

func test(b []byte){
	b = append(b, []byte("hello")...)
}

func main() {

	var b = []byte("abcdefghijklmnopqrstuvwabcdefghi")
	fmt.Println("before b:",b)
	test(b)
	fmt.Println("after test:",b)

	a := 8
	c := 5
	if a&c > 0 {
		fmt.Println("it is true")
	}else{
		fmt.Println("it is false")
	}
	fmt.Println("a&c",a&c)
}
