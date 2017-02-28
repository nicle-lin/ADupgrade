package main

import (
	"net"
	"fmt"
)

type SecData struct {
	flag   uint16
	length uint16

}
type Session struct {
	Conn net.Conn
	//*PeerInfo
	*SecData
}

func main() {
	S := new(Session)
	var conn net.Conn
	S.Conn = conn
	//S.flag = 16  //嵌套的struct如果没有初始化就直接赋值会panic

	SS := &Session{Conn:conn,SecData:&SecData{}}
	SS.flag = 16
	fmt.Println("flag",SS.flag)
}
