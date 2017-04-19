package main

import (
	"net"

	"fmt"
	"github.com/nicle-lin/ADupgrade/test/mblb/proto"
	"io"
	"os"
	"time"
)

var (
	network = "tcp"
	address = os.Args[2]
)

func main() {
	if os.Args[1] == "client"{
		fmt.Println("start client....")
		client()
	}else if os.Args[1] == "server"{
		fmt.Println("start server....")
		server()
	}



}
func handleClient(ch chan <- bool) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer func(){
		ch <- true
	}()
	conn.SetDeadline(time.Now().Add(3*time.Second))
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
			fmt.Println("connection has been close....")
			break
		} else if err != nil {
			fmt.Println("read frome server error:",err)
			return err
		}

	}

	fmt.Println("close the connection.....")
	return nil
}

func client() error {
	ch := make(chan bool,10)
	for i := 0; i < 10; i++ {
		go handleClient(ch)
	}
	for i := 0; i < 10; i++{
		<- ch
		fmt.Println("a connection has been close.....")
	}
	return nil
}

func handleServer(conn net.Conn)( err error) {
	fmt.Printf("Established a connection with a client(remote address:%s)\n",conn.RemoteAddr())
	defer conn.Close() //we close conn after peer close conn
	for {
		_, err = proto.ReadFrame(conn)
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println(err)
			return err
		}
		_, err = proto.WriteFrame([]byte("hi,this is from server"), conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	fmt.Println("has close connection:",err)
	return nil
}

func server() {
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
		go handleServer(conn)

	}
}
