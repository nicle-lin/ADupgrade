package main

import (
	"os"
	"os/signal"
	"fmt"
	"syscall"
)

func main() {

	/*
	flag := make(chan bool)
	c := make(chan os.Signal)
	go func() {
		signal.Notify(c)
		s := <-c
		fmt.Println("get signal:",s)
		flag <- true
	}()
	<- flag
	fmt.Println("time to exit")
	var h HaHa
	h.hello()
	h.hi()
	*/


	sigRecv := make(chan os.Signal, 1)
	sigs := []os.Signal{syscall.SIGINT, syscall.SIGQUIT}
	signal.Notify(sigRecv, sigs...)
	for sig := range sigRecv {
		fmt.Println("received a signal:",sig)
		//signal.Stop(sigRecv)
		close(sigRecv)
	}
}




type HeHe interface {
	hi()
	hello()
}

type HaHa int

func (h HaHa) hi(){
	fmt.Println("hi")
}

func (h HaHa) hello(){
	fmt.Println("hello")
}
