package update

import (
	"net"
	"fmt"
	"strings"
	"regexp"
	"io/ioutil"
)

//if S.data contains string "result:1",it means command executed fail by AD
func IsResultOK(S *Session) bool{
	return !strings.Contains(string(S.data),"result:1")
}

func IsGetOver(S *Session) bool{
	return strings.Contains(string(S.data),CMD[GETOVER])
}


func QueryVersion(S *Session)bool{
	return strings.Contains(string(S.data),"result:7629414")
}

// Get the Server Version
func VersionResult(S *Session){
	reg := regexp.MustCompile(`version:[\d]+`)
	str := reg.FindAllString(string(S.data),-1)[0]
	S.SerVersion = int(strings.Split(str,":")[1])
}

//Get file from Server, and download,write it to the LocalFile
func Get(S *Session,RemoteFile ,LocalFile string)([]byte,error){
	if !DoCmd(S,CMD[GET],RemoteFile){
		return nil,fmt.Errorf("the server can't send the file:%s.check the file exists.\n",RemoteFile)
	}
	var alldata []byte
	S.ReadPacket()
	alldata = S.data
	if S.typ == DATAFRAME{
		S.ReadPacket()
		alldata = append(alldata,S.data...)
	}
	if !IsGetOver(S){
		return nil, fmt.Errorf("Not found getover flag while get the file:%s\n",RemoteFile)
	}
	if LocalFile == ""{
		return alldata,nil
	}

	err := ioutil.WriteFile(LocalFile, alldata, 0666)
	return nil,err
}

func InitClient(S *Session,appversion []byte){

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

func Login(S *Session,passwd string)(err error){
	S.Conn, err = net.Dial("tcp4", S.IP+":"+S.Port)
	if err != nil {
		return err
	}
	if !DoCmd(S,CMD[LOGIN],passwd){
		return fmt.Errorf("Login fail,please check the passwd\n")
	}
	if QueryVersion(S){
		if !DoCmd(S,CMD[VERSION],"") {
			return fmt.Errorf("DoCmd %s fail\n",CMD[VERSION])
		}
		VersionResult(S)
	}else{
		S.SerVersion = 300
		fmt.Println("server version lower than v300. nothing to do.")
	}
	var appversion []byte
	appversion,err = Get(S,APPVERSION_FILE,"")
	if err != nil{
		return err
	}
	InitClient(S,appversion)
	fmt.Println("login success")
	return nil
}

func Logout(S *Session) error{
	return S.Conn.Close()
}