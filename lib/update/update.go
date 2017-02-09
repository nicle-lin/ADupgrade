package update

import (
	"net"
	"fmt"
	"strings"
)

//if S.data contains string "result:1",it means command executed fail by AD
func IsResultOK(S *Session) bool{
	return !strings.Contains(string(S.data),"result:1")
}

//return true,it mean command execute success by peer
//return false, it mean command execute fail by peer
func DoCmd(S *Session, cmdType, params string) bool{
	cmdStr, err := MakeCmdPacket(cmdType,params)
	if err != nil{
		fmt.Errorf("MakeCmdPacket error:",err)
		return
	}
	S.err = S.WritePacket(cmdStr)
	S.ReadPacket()
	return IsResultOK(S)

}

func Login(S *Session,passwd string)bool{
	S.Conn, S.err = net.Dial("tcp4", S.IP+":"+S.Port)
	if S.err != nil {
		return false
	}
	return DoCmd(S,CMD[LOGIN],passwd)
}

func Logout(S *Session) error{
	return S.Conn.Close()
}