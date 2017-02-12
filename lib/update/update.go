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

func QueryVersion(S *Session)bool{
	return strings.Contains(string(S.data),"result:7629414")
}

// Get the Server Version
func VersionResult(S *Session){
	if strings.Contains(string(S.data),"version")&&
	S.SerVersion =
}

//Get file from Server, and download
func Get(S *Session,RemoteFile ,LocalFile string){

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

func Login(S *Session,passwd string)error{
	S.Conn, S.err = net.Dial("tcp4", S.IP+":"+S.Port)
	if S.err != nil {
		return S.err
	}
	if !DoCmd(S,CMD[LOGIN],passwd){
		return fmt.Errorf("Login fail,please check the passwd\n")
	}
	if QueryVersion(S){
		if !DoCmd(S,CMD[VERSION],"") {
			return fmt.Errorf("DoCmd %s fail\n",CMD[VERSION])
		}
		VersionResult(S)
	}
}

func Logout(S *Session) error{
	return S.Conn.Close()
}