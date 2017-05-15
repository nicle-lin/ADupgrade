package main

import (
	"fmt"
	"time"
)

var strChan = make(chan string, 3)
func main() {
	syncChan1 := make(chan struct{},1)
	syncChan2 := make(chan struct{}, 2)
	go receive(strChan, syncChan1, syncChan2)
	go send(strChan, syncChan1,syncChan2)
	<- syncChan2
	<- syncChan2
}

func receive(strChan <- chan  string,syncChan1 <- chan struct{}, syncChan2 chan <- struct{}){
	<- syncChan1
	fmt.Println("receive a sync signal and wait a second...[receiver]")
	time.Sleep(time.Second)
	for elem := range strChan{
		fmt.Println("receive:", elem, "[receive]")
	}
	/* 这个与上面的range 一样
	for {

		if elem, ok := <-strChan; ok{
			fmt.Println(elem)
		}else{
			break
		}
	}
	*/

	fmt.Println("stopped. [receive]")
	syncChan2 <- struct {}{}
}

func send(strChan chan <- string, syncChan1 chan <- struct{}, syncChan2 chan <- struct{}){
	for _, elem := range []string {"a", "b","c","d"}{
		strChan <- elem
		fmt.Println("sent:", elem,"[sender]")
		if elem == "c"{
			syncChan1 <- struct {}{}
			fmt.Println("sent a sync signal. [sender]")
		}
	}
	fmt.Println("wait 2 seconds....[sender]")
	time.Sleep(2 * time.Second)
	close(strChan)
	syncChan2 <- struct {}{}
}