package main

import (
	"fmt"
	"time"
)

func main() {
	/*
		ch4 := make(chan int)//正常的channel
		ch5 := <- chan int(ch4) //接收channel, 只读的
		ch6 := chan<- int(ch4)  //发送channel,　只写的
		//close(ch5)
		close(ch6)
		<- ch5
		fmt.Println("channel:")
	*/
	ch := make(chan int, 5)
	sign := make(chan byte, 2)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(1 * time.Second)
		}
		close(ch)
		fmt.Println("the channel is closed.")
		sign <- 0
	}()
	go func() {
		for {
			e, ok := <-ch
			fmt.Printf("%d(%v)\n", e, ok)
			if !ok {
				break
			}
			time.Sleep(2 * time.Second)

		}
		fmt.Println("done")
		sign <- 1
		close(sign)
	}()
	c := make(chan int)
	go justReceive(c)
	go justSend(c)

	/*
	v1 := <- sign
	v2 := <- sign
	fmt.Println("v1 :",v1)
	fmt.Println("v2 :", v2)
	*/
	if sign != nil{
		for s := range sign{
			fmt.Println("got sign value:",s)
		}
	}

	/*
	for {
		select {
		case v := <-sign:
			fmt.Println("the value is ", v)
		default:
			//fmt.Println("break the select")
			break
		}
	}
	*/

}

//one way channel
func justSend(ch chan <- int)  {
	ch <- 10
	fmt.Println("send message to channel 10")
}
//one way channel
func justReceive(ch <- chan int){
	v := <- ch
	fmt.Println("receive value ",v)
}