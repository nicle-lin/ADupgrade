package main

import (
	"net"

	"fmt"
)

var (
	network    = "tcp"
	address    = "127.0.0.1:7777"
	factory    = func() (net.Conn, error) { return net.Dial(network, address) }
)





func main() {

	conn := net.DialTCP()
}


func client(){

}

func server(){
	l, err := net.Listen(network, address)
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}

		go func() {
			buffer := make([]byte, 256)
			conn.Read(buffer)
		}()
	}
}