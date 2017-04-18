package main

import (
	"net"

	"fmt"
	"github.com/nicle-lin/ADupgrade/test/mblb/proto"
	"io"
)

var (
	network = "tcp"
	address = "127.0.0.1:7777"
	factory = func() (net.Conn, error) { return net.Dial(network, address) }
)

func main() {
	client()
	server()
}
func handleClient() error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	defer conn.Close()

	for i := 0; i < 10; i++ {
		_, err1 := proto.WriteFrame([]byte("hi,this is from client"), conn)
		if err1 != nil {
			return err1
		}
		_, err2 := proto.WriteFrame([]byte("how are you"), conn)
		if err2 != nil {
			return err2
		}
		_, err3 := proto.WriteFrame([]byte("what is your name"), conn)
		if err3 != nil {
			return err3
		}
	}
	for {
		_, err := proto.ReadFrame(conn)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
	}

	return nil
}

func client() error {

	for i := 0; i < 10; i++ {
		go handleClient()
	}
	return nil
}

func handleServer(l net.Listener) error {
	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
	}

	go func() {
		for {
			_, err := proto.ReadFrame(conn)
			if err == io.EOF {
				conn.Close()
				break
			} else if err != nil {
				fmt.Println(err)
			}
			_, err1 := proto.WriteFrame([]byte("hi,this is from server"), conn)
			if err1 != nil {
				fmt.Println(err1)
			}
		}
	}()
	return nil
}

func server() {
	l, err := net.Listen(network, address)
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()

	for {
		go handleServer(l)
	}
}
