package update

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"regexp"
	"strings"
	"github.com/go-ini/ini"
)

//if S.data contains string "result:1",it means command executed fail by AD
func IsResultOK(S *Session) bool {
	return !strings.Contains(string(S.data), "result:1")
}

func IsGetOver(S *Session) bool {
	return strings.Contains(string(S.data), CMD[GETOVER])
}

func QueryVersion(S *Session) bool {
	return strings.Contains(string(S.data), "result:7629414")
}

// Get the Server Version(updateme program version)
func VersionResult(S *Session) {
	reg := regexp.MustCompile(`version:[\d]+`)
	str := reg.FindAllString(string(S.data), -1)[0]
	S.SerVersion = int(strings.Split(str, ":")[1])
}

//Get AD Version
func GetAppVersion(S *Session, appVersion []byte) {
	reg := regexp.MustCompile(`[\w]+-[\w]+\.[\w]+`)
	str := reg.FindAllString(string(appVersion), -1)[0]
	S.AppVersion = strings.Split(str, "-")[1]
	fmt.Println("The first line of appversion of the current device is:", S.AppVersion)
}

func IsArmChip(appVersion []byte) bool {
	str := strings.ToLower(string(appVersion))
	if strings.Contains(str, "-ac-") || strings.Contains(str, "sinfor-m") || strings.Contains(str, "-ad-") {
		return true
	}
	if strings.Contains(str, "-bm-") || strings.Contains(str, "-bc-") || strings.Contains(str, "-iam") {
		return true
	}

	if strings.Contains(str, "-nag") || strings.Contains(str, "sinfor--") || strings.Contains(str, "sangfor--") {
		return true
	}
	if strings.Contains(str, "ar") || strings.Contains(str, "xp") || strings.Contains(str, "plus") {
		return false
	}
	return false
}

//Get file from Server, and download,write it to the LocalFile
func Get(S *Session, RemoteFile, LocalFile string) ([]byte, error) {
	if !DoCmd(S, CMD[GET], RemoteFile) {
		return nil, fmt.Errorf("the server can't send the file:%s.check the file exists.\n", RemoteFile)
	}
	var allData []byte
	S.ReadPacket()
	allData = S.data
	if S.typ == DATAFRAME {
		S.ReadPacket()
		allData = append(allData, S.data...)
	}
	if !IsGetOver(S) {
		return nil, fmt.Errorf("Not found getover flag while get the file:%s\n", RemoteFile)
	}
	if LocalFile == "" {
		return allData, nil
	}

	err := ioutil.WriteFile(LocalFile, allData, 0666)
	return nil, err
}

//return true,it mean command execute success by peer
//return false, it mean command execute fail by peer
func DoCmd(S *Session, cmdType, params string) error {
	cmdStr, err := MakeCmdPacket(cmdType, params)
	if err != nil {
		return fmt.Errorf("MakeCmdPacket error:%v", err)

	}
	err = S.WritePacket(cmdStr)
	if err != nil {
		return err
	}
	err = S.ReadPacket()
	if err != nil {
		return err
	}
	if IsResultOK(S) {
		return nil
	} else {
		return fmt.Errorf("result is not ok")
	}

}

func Exec(S *Session, U *Update, Command string) (string, error) {
	doRet := DoCmd(S, CMD[EXEC], Command)
	getReturn, err := Get(S, U.TempRetFile, "")
	if err != nil {
		return nil, err
	}
	getResult, err1 := Get(S, U.TempRstFile, "")
	if err1 != nil {
		return nil, err1
	}
	if strings.TrimSpace(string(getReturn)) != "0" || doRet != nil {
		return nil, fmt.Errorf("DoCmd error or return result is 0\n")
	}
	return string(getResult), nil
}

func Put(S *Session, LocalFile, RemoteFile string) error {
	if DoCmd(S, CMD[PUT], RemoteFile) != nil {
		return fmt.Errorf("DoCmd fail, put %s fail\n", RemoteFile)
	}
	file, err := os.Open(LocalFile)
	if err != nil {
		return err
	}
	defer file.Close()

	buf := make([]byte, MAX_DATA_LEN)
	bufRead := bufio.NewReader(file)

	for {
		n, err1 := bufRead.Read(buf)
		if err1 != nil && err1 != io.EOF {
			return err1
		}
		if 0 == n {
			break
		}
		data, _ := MakeDataPacket(buf[:n])
		S.WritePacket(data)

	}
	if DoCmd(S, CMD[PUTOVER], "") != nil {
		return fmt.Errorf("DoCmd fail, PUTOVER fail\n")
	}
	return nil
}

func PutFile(ip, port, passwd, LocalFile, RemoteFile string) error {
	if !IsPathExist(LocalFile) {
		return fmt.Errorf("%s don't exist", LocalFile)
	}
	S, loginErr := Login(ip, port, passwd)
	if loginErr != nil {
		return loginErr
	}
	defer Logout(S)
	return Put(S, LocalFile, RemoteFile)
}

func GetFile(ip, passwd, port, LocalFile, RemoteFile string) error {
	S, loginErr := Login(ip, port, passwd)
	if loginErr != nil {
		return loginErr
	}
	defer Logout(S)
	_, err := Get(S, RemoteFile, LocalFile)
	return err
}

func Login(ip, port, passwd string) (*Session, error) {
	conn, err := net.Dial("tcp4", ip+":"+port)
	if err != nil {
		return nil, err
	}
	S := new(Session)
	S.Conn = conn
	if !DoCmd(S, CMD[LOGIN], passwd) {
		return fmt.Errorf("Login fail,please check the passwd\n")
	}
	if QueryVersion(S) {
		if !DoCmd(S, CMD[VERSION], "") {
			return nil, fmt.Errorf("DoCmd %s fail\n", CMD[VERSION])
		}
		VersionResult(S)
	} else {
		S.SerVersion = 300
		fmt.Println("server version lower than v300. nothing to do.")
	}
	fmt.Println("login success")
	return S, nil
}

func Logout(S *Session) error {
	return S.Conn.Close()
}

func UpgradeCheck(S *Session, U *Update) error {
	_, err := Exec(S, U, "ls "+UPDATE_CHECK_SCRIPT)
	if err != nil {
		Put(S, U.LocalUpdCheck, UPDATE_CHECK_SCRIPT)
	}
	//execute /usr/sbin/updatercheck.sh, check it pass or fail
	msgVersion, resultVersion := Exec(S, U, UPDATE_CHECK_SCRIPT)
	if resultVersion != nil {
		return fmt.Errorf("Upgrade failed!!!,error msg:%s", msgVersion)
	}

	//check upgrade sn valid or invalid
	msgSn, resultSn := Exec(S, U, CHECK_UPGRADE_SN)
	if resultSn != nil {
		return fmt.Errorf("Upgrade failed!!!,error msg:%s", msgSn)
	}
	return nil
}


//TODO only support to update single package right now
func ThreadUpdateAllPackages(S *Session,U *Update)error  {
	switch U.SSUType {
	case PACKAGE_TYPE:
		if err := UpdateSinglePacket(S,U);err != nil {return err}
	case RESTORE_TYPE:
		if err := RestoreDefaultPriv(); err != nil {return err}
	case EXECUTE_TYPE:
		if err := Put(S,U.SSUPackage,U.Compose); err != nil {return err}
		if _, err:= Exec(S,U,U.Compose); err != nil {return err}
	default:
		fmt.Println("unknown type packet:",U.SSUType)
		return fmt.Errorf("unknown type packet %s:",U.SSUType)
	}
	return nil
}

func UpdateUpgradeHistory(S *Session,U *Update)error  {
	_, err := Exec(S,U,"ls "+UPDHISTORY_SCRIPT)
	if err != nil {
		if err := Put(S,U,U.LocalUpdHistory);err != nil {return err}
		Exec(S,U,"sync") //TODO: not done yet
	}
	if _, err := Exec(S,U,UPDHISTORY_SCRIPT + " " + U.SSUPackage);err != nil{
		return err
	}
	return nil
}

//TODO: ini format file
//TODO: now
func ConfirmRebootDevice(S *Session,U *Update)error{

}

func Upgrade(ip, port, password, ssu string) error {

	S, err := Login(ip, port, password)
	if err != nil {
		return err
	}
	var appVersion []byte
	appVersion, err = Get(S, APPVERSION_FILE, "")
	if err != nil {
		return err
	}
	GetAppVersion(S, appVersion)

	U := InitClient(appVersion)
	U.SSUPackage = ssu
	err = UpgradeCheck(S, U)
	if err != nil {
		return err
	}
	if PrepareUpgrade(S, U) != nil {
		return err
	}

	apps := GetApps(U.SingleUnpkg)
	for _, v := range apps {
		if err := EncFile(v, v+"_des"); err != nil {return err}
	}
	
	if err := ThreadUpdateAllPackages(S,U); err != nil {return err}
	if err := UpdateUpgradeHistory(S,U);err != nil {return err}
	if err := ConfirmRebootDevice(S, U); err != nil {return err}

	defer FreeUpdateDir()
	defer FreeCfgDir()
	defer Logout(S)


	return nil
}

func ThreadUpgrade(ip []string, port []string, passwd []string, ssu []string) {

}
