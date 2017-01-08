package main

import (
	"net"
	"fmt"
)

func main() {
	conn, err := net.Dial("tcp","192.168.1.100:51111")
	if err != nil{
		fmt.Println("Dial error:",err)
	}
	var readdata []byte
	conn.Read(readdata)
	fmt.Println("readdata:",readdata)
	conn.Write([]byte("from go"))
}
