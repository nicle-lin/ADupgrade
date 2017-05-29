package main

import "fmt"

var intChan1 chan int// = make(chan int,1)
var intChan2 chan int //= make(chan int,1)
//如果intChan1 intChan2没有初始化，就没有make(chan int,1)
//向未初始化的通道发送会造成永久阻塞，所以first与second没有被选中，而default被选中
var channels = []chan int{intChan1,intChan2}
var number = []int{1,2,3,4,5}
func main() {
	select {
	case getChan(0) <- getNumber(0):
		fmt.Println("first case is selected")
	case getChan(1) <- getNumber(1):
		fmt.Println("second case is selected")
	default:
		fmt.Println("default case is selected")

	}
}
func getNumber(i int)int {
	fmt.Printf("number[%d]\n",i)
	return number[i]
}
func getChan(i int)chan int {
	fmt.Printf("channels[%d]\n",i)
	return channels[i]
}
