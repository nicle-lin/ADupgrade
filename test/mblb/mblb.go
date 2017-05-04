package main

import (
	"net"

	"fmt"
	"github.com/nicle-lin/ADupgrade/test/mblb/proto"
	"io"
	"os"
	"time"
	"flag"
)

var (
	network = "tcp"
	usage = `usage: mblb client|server ip:port [options]
	it is designed to test AD mblb.
options:
	-c: Number of requests to run concurrently per second (client),default is 50
	-t: how many second to latest to run (client),default is 60s
	-s: what message to send (less than 1020) (client), default is hi,this is from client
	-r: what message to response (less than 1020) (server), default is hi,this is server
	`

	c = flag.Int("-c",50, "number of requests to run")
	t = flag.Int64("-t",60, "time")
	s = flag.String("-s","hi,this is from client","send message")
	r = flag.String("-r", "hi,this is server","response message")
)

func usageAndExit(msg string) {
	if msg != "" {
		fmt.Fprintf(os.Stderr, msg)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	flag.Usage = func(){
		fmt.Fprint(os.Stderr,usage)
	}
	flag.Parse()

	if flag.NArg() < 2 {
		usageAndExit("")
	}
	typ := flag.Args()[0]
	address := flag.Args()[1]

	if typ== "client"{
		fmt.Println("start client....")
		client(address)
	}else if typ == "server"{
		fmt.Println("start server....")
		server(address)
	}



}
func handleClient(ch chan <- bool, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	defer func(){
		if r := recover(); r != nil {
			fmt.Println("receive message time out")
		}
	}()
	defer conn.Close()
	defer func(){
		ch <- true
	}()
	conn.SetDeadline(time.Now().Add(3*time.Second))
	for i := 0; i < 10; i++ {
		_, err1 := proto.WriteFrame([]byte(*s), conn)
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

func client(address string) error {
	ch := make(chan bool,*c)
	for i := 0; i < *c; i++ {
		go handleClient(ch,address)
	}
	for i := 0; i < *c; i++{
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
		_, err = proto.WriteFrame([]byte(*r), conn)
		if err != nil {
			fmt.Println(err)
			return err
		}
	}
	fmt.Println("has close connection:",err)
	return nil
}

func server(address string) {
	l, err := net.Listen(network, address)
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	fmt.Println("Listening on %s",address)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleServer(conn)

	}
}
