package main

import (
	"sync"
	"os"
	"syscall"
	"fmt"
	"os/signal"
	"time"
)

func main() {
	sigRecv := make(chan os.Signal, 1)
	sigs := []os.Signal{syscall.SIGINT,syscall.SIGQUIT}
	fmt.Printf("Set notification for %s... [sigRecv1]\n",sigs)
	signal.Notify(sigRecv,sigs...)


	sigRecv1 := make(chan os.Signal, 1)
	sigs1 := []os.Signal{syscall.SIGQUIT}
	fmt.Printf("Set notification for %s... [sigRecv1]\n",sigs1)
	signal.Notify(sigRecv1,sigs1...)

	var wg sync.WaitGroup
	wg.Add(2)
	go func(){
		for sig := range sigRecv1 {
			fmt.Printf("received a signal from sigrecv1:%s\n",sig)
		}
		fmt.Printf("end. [sigRecv1]\n")
		wg.Done()
	}()
	go func(){
		for sig := range sigRecv {
			fmt.Printf("received a signal from sigrecv1:%s\n",sig)
		}
		fmt.Printf("end. [sigRecv]\n")
		wg.Done()
	}()

	fmt.Println("wait for 2 seconds...")
	time.Sleep(2 * time.Second)
	signal.Stop(sigRecv)
	fmt.Printf("done.[sigRecv]\n")
	wg.Wait()
}
