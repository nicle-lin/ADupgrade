package main

import (
	"context"
	"fmt"
)

func main() {
	gen := func(ctx context.Context) <- chan int{
		dst := make(chan int)
		n := 1
		go func(){
			for{
				select {
				case <- ctx.Done():
					return
					case dst <- n:
					n++
				}
			}
		}()
		//close(dst) can't close the channel
		return dst
	}
	//ctx, cancel := context.WithCancel(context.Background())
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	for n := range gen(ctx){
		fmt.Println(n)
		if n == 5{
			break
		}
	}
}
