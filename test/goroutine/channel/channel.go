package main

import "fmt"


func Count(ch chan int) {
	fmt.Println("Counting...")
	ch <- 10 * 10
}

func main() {

	chs := make([]chan int,13)
	for i := 0; i < 13; i++ {
		chs[i] = make(chan int)
		go Count(chs[i])
	}
	var sum int = 0
	for _, ch := range chs {
		sum = <-ch + sum
	}
	fmt.Println("sum:",sum)

	add()
}

func add()  {
	var sum int64 = 0
	for i := 0; i < 1000000; i++ {
		sum = sum + 123456789 * 987654321
	}
	fmt.Println("sum:",sum)
}