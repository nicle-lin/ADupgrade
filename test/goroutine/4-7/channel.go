package main

import (
	"fmt"
	"time"
)

var mapChan = make(chan map[string]int, 1)
func main() {
	var mapStr = make(map[string]int,1)
	mapStr["good"]++
	mapStr["bad"]++
	mapStr["bad"]++
	bad, ok := mapStr["bad"]
	if ok{
		fmt.Printf("got bad:%d\n",bad)
	}
	fmt.Printf("mapStr:%v\n",mapStr)
	syncChan := make(chan struct{}, 2)
	go func(){
		for{
			if elem, ok := <- mapChan; ok{
				fmt.Printf("mapChan:%v",elem)
				elem["count"]++
				elem["test"]++
			}else{
				break
			}
		}
		fmt.Println("Stopped. [receiver]")
		syncChan <- struct {}{}
	}()

	go func(){
		countMap := make(map[string]int)
		for i := 0; i < 5; i++{
			mapChan <- countMap
			time.Sleep(time.Millisecond)
			fmt.Printf("the count map: %v.[sender]\n", countMap)
		}
		close(mapChan)
		syncChan <- struct {}{}
	}()

	<- syncChan
	<- syncChan
}
