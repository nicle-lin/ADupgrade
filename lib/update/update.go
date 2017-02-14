package update

import (
	"net"
	"fmt"
	"strings"
	"regexp"
	"io/ioutil"
	"os"
	"bufio"
	"io"
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

// Get the Server Version(updateme program version)
func VersionResult(S *Session){
	reg := regexp.MustCompile(`version:[\d]+`)
	str := reg.FindAllString(string(S.data),-1)[0]
	S.SerVersion = int(strings.Split(str,":")[1])
}

//Get AD Version
func GetAppVersion(S *Session,appVersion []byte){
	reg := regexp.MustCompile(`[\w]+-[\w]+\.[\w]+`)
	str := reg.FindAllString(string(appVersion),-1)[0]
	S.AppVersion = int(strings.Split(str,"-")[1])
	fmt.Println("The first line of appversion of the current device is:",S.AppVersion)
}

func IsArmChip(appVersion []byte)bool{
	str := strings.ToLower(string(appVersion))
	if strings.Contains(str,"-ac-") || strings.Contains(str,"sinfor-m") || strings.Contains(str,"-ad-"){
		return true
	}
	if strings.Contains(str,"-bm-") || strings.Contains(str,"-bc-") || strings.Contains(str,"-iam"){
		return true
	}

	if strings.Contains(str,"-nag") || strings.Contains(str,"sinfor--") || strings.Contains(str,"sangfor--"){
		return true
	}
	if strings.Contains(str,"ar") || strings.Contains(str,"xp") || strings.Contains(str,"plus"){
		return false
	}
	return false
}

//Get file from Server, and download,write it to the LocalFile
func Get(S *Session,RemoteFile ,LocalFile string)([]byte,error){
	if !DoCmd(S,CMD[GET],RemoteFile){
		return nil,fmt.Errorf("the server can't send the file:%s.check the file exists.\n",RemoteFile)
	}
	var allData []byte
	S.ReadPacket()
	allData = S.data
	if S.typ == DATAFRAME{
		S.ReadPacket()
		allData = append(allData,S.data...)
	}
	if !IsGetOver(S){
		return nil, fmt.Errorf("Not found getover flag while get the file:%s\n",RemoteFile)
	}
	if LocalFile == ""{
		return allData,nil
	}

	err := ioutil.WriteFile(LocalFile, allData, 0666)
	return nil,err
}

func InitClient(S *Session,appVersion []byte){
	S.FolderPrefix = GetRandomString(32)
	S.CurrentWorkFolder = GetCurrentDirectory()
	if IsArmChip(appVersion){
		S.TempExecFile,S.TempRstFile = ARM_LINUX_BASIC[0],ARM_LINUX_BASIC[1]
		S.CustomErrFile,S.TempRetFile = ARM_LINUX_BASIC[2],ARM_LINUX_BASIC[3]
		S.LoginPwdFile,S.Compose = ARM_LINUX_BASIC[4],ARM_LINUX_BASIC[5]

		S.ServerAppRe,S.ServerAppSh = ARM_LINUX_UPDATE[0],ARM_LINUX_UPDATE[1]
		S.ServerCfgPre,S.ServerCfgSh = ARM_LINUX_UPDATE[2],ARM_LINUX_UPDATE[3]

		S.LocalBackSh = S.CurrentWorkFolder + "/" + S.FolderPrefix + "/arm_bin/bakcfgsh"
		S.LocalPreCfgSh = S.CurrentWorkFolder + "/" +  S.FolderPrefix + "/arm_bin/prercovcfgsh"
		S.LocalCfgSh = S.CurrentWorkFolder + "/" + S.FolderPrefix + "/arm_bin/rcovcfgsh"
		S.LocalUpdHistory = S.CurrentWorkFolder + "/" + S.FolderPrefix + "/arm_bin/updhistory.sh"
		S.LocalUpdCheck = S.CurrentWorkFolder + "/" + S.FolderPrefix + "/arm_bin/updatercheck.sh"


		fmt.Println("The device is a arm platform,init arm info.")
	}else{
		S.TempExecFile,S.TempRstFile = X86_LINUX_BASIC[0],X86_LINUX_BASIC[1]
		S.CustomErrFile,S.TempRetFile = X86_LINUX_BASIC[2],X86_LINUX_BASIC[3]
		S.LoginPwdFile,S.Compose = X86_LINUX_BASIC[4],X86_LINUX_BASIC[5]

		S.ServerAppRe,S.ServerAppSh = X86_LINUX_UPDATE[0],X86_LINUX_UPDATE[1]
		S.ServerCfgPre,S.ServerCfgSh = X86_LINUX_UPDATE[2],X86_LINUX_UPDATE[3]


		S.LocalBackSh = S.CurrentWorkFolder + "/" + S.FolderPrefix + "/bin/bakcfgsh"
		S.LocalPreCfgSh = S.CurrentWorkFolder + "/" +  S.FolderPrefix + "/bin/prercovcfgsh"
		S.LocalCfgSh = S.CurrentWorkFolder + "/" + S.FolderPrefix + "/bin/rcovcfgsh"
		S.LocalUpdHistory = S.CurrentWorkFolder + "/" + S.FolderPrefix + "/bin/updhistory.sh"
		S.LocalUpdCheck = S.CurrentWorkFolder + "/" + S.FolderPrefix + "/bin/updatercheck.sh"

		fmt.Println("The device is a x86 platform,init x86 info.")
	}
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

func Exec(S *Session,Command string)(string,error){
	doRet := DoCmd(S,CMD[EXEC],Command)
	getReturn,err := Get(S,S.TempRetFile,"")
	if err != nil {
		return nil,err
	}
	getResult,err1 := Get(S,S.TempRstFile,"")
	if err1 != nil {
		return nil,err1
	}
	if strings.TrimSpace(string(getReturn)) != "0" || !doRet{
		return nil,fmt.Errorf("DoCmd error or return result is 0\n")
	}
	return string(getResult),nil
}

func Put(S *Session, RemoteFile,LocalFile string)error{
	if !DoCmd(S,CMD[PUT],RemoteFile){
		return fmt.Errorf("DoCmd fail, put %s fail",RemoteFile)
	}
	file ,err := os.Open(LocalFile)
	if err != nil{
		return  err
	}
	defer file.Close()

	buf := make([]byte,MAX_DATA_LEN)
	bufRead := bufio.NewReader(file)

	for{
		n,err1 := bufRead.Read(buf)
		data,_ := MakeDataPacket(buf[:n])
		S.WritePacket(data)
		if err1 != nil{
			if err1 == io.EOF{
				break
			}
			return err1
		}
	}
	//TODO: send putover
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
	var appVersion []byte
	appVersion,err = Get(S,APPVERSION_FILE,"")
	if err != nil{
		return err
	}
	GetAppVersion(S,appVersion)
	InitClient(S,appVersion)
	fmt.Println("login success")
	return nil
}

func Logout(S *Session) error{
	return S.Conn.Close()
}

func UpgradeCheck(S *Session){
	result, err := Exec(S,"ls " + UPDATE_CHECK_SCRIPT)
	if err != nil{
		Put(S,S.LocalUpdCheck,UPDATE_CHECK_SCRIPT)
	}
}

func Upgrade(ip,port,password,ssu string){

}

func ThreadUpgrade(){

}