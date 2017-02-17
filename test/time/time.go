package main

import (
	"time"
	"fmt"
)

func main() {
	before := time.Now()
	time.Sleep(2*time.Second)
	dur := time.Now().Sub(before)
	if dur > 1 *time.Second {
		fmt.Println("dur:",dur)
	}else {
		fmt.Println("haha:",dur)
	}

}
