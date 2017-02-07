package update

import "net"

type session struct {
	Conn net.Conn
	Ip string
}