package main

import (
	"net"

	"flag"
	"fmt"
	"github.com/nicle-lin/ADupgrade/test/mblb/proto"
	"io"
	"os"
	"sync"
	"time"
)

var (
	network = "tcp"
	usage   = `usage: mblb [options] client|server ip:port
	it is designed to test AD mblb.
options:
	-c: Number of requests to run concurrently per second (client),default is 50
	-q: how many request of every connection.... (client), default is 3
	-t: how many second to latest to run (client),default is 60s
	-p: connections/per second (client), default is 10
	-s: what message to send (less than 1020) (client), default is hi,this is from client
	-r: what message to response (less than 1020) (server), default is hi,this is server
	-d: response to client after x second  (server), default is without delay
	`

	c = flag.Int("c", 50, "number of requests to run")
	q = flag.Int("q", 1, "how many request of every connection....")
	t = flag.Int64("t", 60, "time")
	p = flag.Int("p", 10, "connections/per second")
	s = flag.String("s", "hi,this is from client", "send message")
	r = flag.String("r", "hi,this is server", "response message")
	d = flag.Int64("d", 0, "delay")
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
	flag.Usage = func() {
		fmt.Fprint(os.Stderr, usage)
	}
	flag.Parse()

	if flag.NArg() < 2 {
		usageAndExit("")
	}
	typ := flag.Args()[0]
	address := flag.Args()[1]

	if typ == "client" {
		fmt.Println("start client....")
		client(address)
	} else if typ == "server" {
		fmt.Println("start server....")
		server(address)
	}

}
func handleClient(ch chan<- bool, address string) error {
	conn, err := net.Dial(network, address)
	if err != nil {
		return err
	}
	defer conn.Close()
	defer func() {
		ch <- true
	}()

	for i := 0; i < *q; i++ {
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

	for i := 0; i < *q; i++{
		_, err := proto.ReadFrame(conn)
		if err == io.EOF {
			fmt.Println("connection has been close....")
			break
		} else if err != nil {
			fmt.Println("read frome server error:", err)
			return err
		}

	}

	return nil
}

func client(address string) error {
	ch := make(chan bool, *c)
	var lr LimitRate
	lr.SetRate(*p)

	for i := 0; i < *c; i++ {
		if lr.Limit() {
			go handleClient(ch, address)
		}
	}
	for i := 0; i < *c; i++ {
		<-ch
		fmt.Println("a connection has been close.....")
	}
	return nil
}

func handleServer(conn net.Conn) (err error) {
	fmt.Printf("Established a connection with a client(remote address:%s)\n", conn.RemoteAddr())
	_, err = proto.ReadFrame(conn)
	if err == io.EOF {
		fmt.Println("connection has been closed, can't read")
	} else if err != nil {
		fmt.Println(err)
		return err
	}
	time.Sleep(time.Duration(*d))
	_, err = proto.WriteFrame([]byte(*r), conn)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer conn.Close()
	fmt.Println("closing the connection.....")
	return nil
}

func server(address string) {
	l, err := net.Listen(network, address)
	if err != nil {
		fmt.Println(err)
	}
	defer l.Close()
	fmt.Println("Listening on %s", address)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
		}
		go handleServer(conn)

	}
}

type LimitRate struct {
	rate       int
	interval   time.Duration
	lastAction time.Time
	lock       sync.Mutex
}

func (l *LimitRate) Limit() bool {
	result := false
	for {
		l.lock.Lock()
		//判断最后一次执行的时间与当前的时间间隔是否大于限速速率
		if time.Now().Sub(l.lastAction) > l.interval {
			l.lastAction = time.Now()
			result = true
		}
		l.lock.Unlock()
		if result {
			return result
		}
		time.Sleep(l.interval)
	}
}

//SetRate 设置Rate
func (l *LimitRate) SetRate(r int) {
	l.rate = r
	l.interval = time.Microsecond * time.Duration(1000*1000/l.rate)
}

//GetRate 获取Rate
func (l *LimitRate) GetRate() int {
	return l.rate
}
