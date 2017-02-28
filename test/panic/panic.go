package main

import "net"

type SecData struct {
	flag   uint16
	length uint16
	typ    byte
	data   []byte
}
type Session struct {
	Conn net.Conn
	//*PeerInfo
	*SecData
}
//嵌套的struct如果没有初始化就直接赋值会panic
func main() {
	S := new(Session)
	var conn net.Conn
	S.Conn = conn
	S.flag = 16
}
