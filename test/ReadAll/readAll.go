package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	fd, err := os.Open("test.txt")
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := ioutil.ReadAll(fd) //需要注意ReadAll()函数读取最大是512字节,但能把所有都读出来
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(len(data))
}
