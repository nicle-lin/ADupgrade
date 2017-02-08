package update

import (
	"net"
)

func DoCmd(S *Session, cmdType, params string) {
	cmdstr, err := MakeCmdPacket(cmdType,params)
	S.err = S.WritePacket(cmdstr)
	S.ReadPacket()
}

func Login(S *Session)bool{
	S.Conn, S.err = net.Dial("tcp4", S.IP+":"+S.Port)
	if S.err != nil {
		return false
	}
	return true
}