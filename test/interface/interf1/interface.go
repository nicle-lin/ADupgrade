package main

import "fmt"

type Co interface {
	Print()
}

type co struct {
	name string
}

func (c *co)Print(){
	fmt.Println("hahah")
}

func main() {
	var h Co
	h.Print()
}
