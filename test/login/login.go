package main

import (
	//"net"
	"fmt"
	"os"
	"log"
	"strings"
	"sync"
	"time"
)

func main() {
	/*
	conn, err := net.Dial("tcp","192.168.1.100:51111")
	if err != nil{
		fmt.Println("Dial error:",err)
	}
	var readdata []byte
	conn.Read(readdata)
	fmt.Println("readdata:",readdata)
	conn.Write([]byte("from go"))
	*/
	ip := "192.168.1.100"
	port := "5000"
	fmt.Println("ip + port:",ip+":"+port)

	file, err := os.OpenFile(
		"test.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writebyte := []byte{0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x3a, 0x30, 0xa, 0x6d, 0x73, 0x67, 0x3a, 0x33, 0x35, 0xa, 0x0}
	fmt.Println("string:",string(writebyte))
	ret := strings.Contains(string(writebyte),"result:1")
	fmt.Println("ret:",ret)
	// 写字节到文件中
	//byteSlice := []byte("Bytes!\n")
	_ , err = file.Write(writebyte)
	if err != nil {
		log.Fatal(err)
	}

	name := Name{i:10}
	var once sync.Once
	for i:=0; i <5;i++{
		go func(i int){
			once.Do(name.OnceDo)
			fmt.Println("iii:",i)
		}(i)
		//fmt.Println("okokiiii:",i)
	}
	time.Sleep(2*time.Second)
}

type Name struct {
	i int
}

func (n *Name)OnceDo(){
	fmt.Println("i:",n.i)
}