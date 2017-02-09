package update

import (
	"net"
	"fmt"
)

func DoCmd(S *Session, cmdType, params string) {
	cmdstr, err := MakeCmdPacket(cmdType,params)
	if err != nil{
		fmt.Errorf("MakeCmdPacket error:",err)
		return
	}
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

func Logout(S *Session) error{
	return S.Conn.Close()
}