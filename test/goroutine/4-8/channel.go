package main

import (
	"fmt"
	"time"
)

type Counter struct {
	count int
}

var mapChan = make(chan map[string]Counter,1)
var mapCh = make(chan map[string]*Counter,1)
func main() {
	syncChan := make(chan struct{},2)
	go func(){
		for{
			if elem, ok := <-mapChan; ok{
				counter := elem["count"]
				counter.count++

				co := elem["bad"]
				co.count++
			}else{
				break
			}


		}

		for {
			if elem, ok := <-mapCh; ok{
				counter := elem["count"]
				counter.count++

				co := elem["bad"]
				co.count++
			}else{
				break
			}
		}
		fmt.Println("Stopped.[receiver]")
		syncChan <- struct {}{}
	}()

	go func() {
		countMap := map[string]Counter{
			"count": Counter{},
			"bad": Counter{1},
		}
		for i := 0; i < 5; i++{
			mapChan <- countMap
			time.Sleep(time.Millisecond)
			fmt.Printf("the count map: %v. [sender]\n",countMap)
		}

		countM := map[string]*Counter{
			"count": &Counter{},
			"bad": &Counter{1},
		}
		for i := 0; i < 5; i++{
			mapCh <- countM
			time.Sleep(time.Millisecond)
			fmt.Printf("the count map: %v. [sender]\n",countM)
		}

		close(mapChan)
		close(mapCh)
		syncChan <- struct {}{}
	}()

	<-syncChan
	<-syncChan
}
