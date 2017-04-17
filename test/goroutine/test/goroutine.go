package main

import (
	"fmt"
	//"time"
	"runtime"
	"runtime/debug"
)
func goroutine(){
	name := "eric"
	go func() {
		fmt.Printf("hello, %s\n", name)
	}()
	//runtime.Gosched()
	name = "Harry"
	runtime.Gosched()
}

func sayhi(name string) {
	fmt.Printf("hi, %s\n",name)
}

func sayHi(){
	names := []string{"eric","harry", "robert","jim", "mark"}
	for _, name := range  names {
		go sayhi(name)
	}
	runtime.Gosched()
}

func goExit(){
	//runtime.Goexit()
	defer fmt.Println("hahahah,after Goexit")
	runtime.Goexit()
	fmt.Println("this is not defer")
}

func init(){
	//设置过小,会报栈溢出
	//debug.SetMaxStack(1)
	fmt.Println("thread nunber:",debug.SetMaxThreads(200))
}

func main() {

	//go fmt.Println("hello world")
	//go fmt.Println("hello world and how are you")
	//go func(i int){
	//	println("Go routine ",i)
	//}(10)
	//runtime.Gosched()
	//time.Sleep(1 * time.Second)


	//time.Sleep(1 * time.Second)
	//goroutine()
	//sayHi()
	//time.Sleep(1 * time.Second)

	/*
	names := []string{"eric", "harry", "robert", "jim", "mark"}
	for _, name := range names{
		go func(name string){
			fmt.Printf("Hello, %s\n", name)
		}(name)
	}
	runtime.Gosched()
	fmt.Println("-----------------------------")

	for _, name :=range names{
		go func(){
			fmt.Printf("Hello, %s\n", name)

		}()
		runtime.Gosched()
	}
	fmt.Println(runtime.Version())
	fmt.Println(runtime.GOMAXPROCS(0))
	go goExit()
	//runtime.Gosched()
	fmt.Println("goroutine number:",runtime.NumGoroutine())
	*/


	defer func(){
		if r := recover(); r != nil {
			fmt.Println("we catch panic error")
		}
	}()

	var ch = make(chan int)
	var done = make(chan bool)
	go func() {
		v := <-ch
		fmt.Println("i am done, too:",v)
		done <- true
		ch <- 11

	}()
	go func(){
		fmt.Println("i am done")
		ch <- 10
	}()
	<-done
	v,ok := <- ch
	if !ok{
		fmt.Println("channel has been close")
	}
	fmt.Println("i am done after both of you:",v)

	var t = make(chan string, 3)

	go func(){
		v := <- t
		fmt.Println("got value from t:",v)
		v1 := <- t
		fmt.Println("got value from t:",v1)
	}()
	t <- "a"
	t <- "b"
	t <- "c"
	t <- "d"

}
