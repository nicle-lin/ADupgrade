package main

import (
	"time"
	"fmt"
)

func main() {
	timeout := make(chan bool,1)
	ch := make(chan bool,1)
	go func() {
		time.Sleep(1e9)
		timeout <- true
	}()

	select {
	case <-ch:
	case <-timeout:
	}
	fmt.Println("select done")
}
