package main

import (
	"fmt"
	"sync/atomic"
	"time"
)

func main() {
	var status uint32 = 1
	fmt.Println(status)
	for i := 0; i < 10; i++ {
		go atomic.CompareAndSwapUint32(&status, 1, uint32(i))
	}
	time.Sleep(time.Second)
	fmt.Println(status)
}
