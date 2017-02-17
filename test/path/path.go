package main

import (
	"os"
	"fmt"
)

func IsPathExist(path string)bool{
	_,err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}
func main() {
	if IsPathExist("haha.txt"){
		fmt.Println("exist")
	}else{
		fmt.Println("don't exist")
	}
}
