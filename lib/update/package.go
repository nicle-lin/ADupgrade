package update

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"
	"os"
	"io"
	"crypto/md5"
	"runtime"
	"os/exec"
	"io/ioutil"
)

var Flag uint16
var m *sync.RWMutex = new(sync.RWMutex)

//return false,the caller have to unpack the SSU,and inc Flag
func GetFlag() bool {
	m.RLock()
	defer m.RUnlock()
	if Flag == 0 {
		return false
	} else {
		return true
	}
}

//when unpack SSU done, it should call this function
func IncFlag() {
	m.Lock()
	defer m.Unlock()
	Flag++
}

//when upgrade success, it should call this function
func DecFlag() {
	m.Lock()
	defer m.Unlock()
	if Flag > 0 {
		Flag--
	}
}

//相同的版本的SSU只能解压一次,在没有解压完成之前其它goroute只能等待解压完成，需要channel来通信
var once sync.Once

func (S *Session) unpackSSU(ssu string) {

}

func UnpackSSU() {
	if !GetFlag() {
		IncFlag()
		//don't have to unpack SSU,because it has been unpacked
		return
	}
	//var name string
	var S Session
	once.Do(S.unpackSSU)

	IncFlag()
}
func unpackPackage(U *Update)error {
	// function InitEnvironment has been init the path U.SingleUnpkg
	fmt.Println("begin to unpack the package")
	var UnpackTool string
	if runtime.GOOS	 == "windows"{
		UnpackTool = filepath.Join(U.CurrentWorkFolder,"tool","7z.exe")
	}else{
		UnpackTool = "7za"
	}
	newArgs := []string{
		0: "x",
		1: "-y",
		2: "-p"+SSU_DEC_PASSWD,
		3: U.SSUPackage,
		4: "-o"+ filepath.Join(U.CurrentWorkFolder,U.SingleUnpkg),
	}

	new := exec.Command(UnpackTool,newArgs...)
	stdout, _ := new.StdoutPipe()
	if err := new.Start(); err != nil {return err}
	data, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("unpack log has been lost")
	}
	ioutil.WriteFile(filepath.Join(U.CurrentWorkFolder,"7z.log"),data,0664)
	if err := new.Wait(); err != nil {
		return err
	}else{
		fmt.Println("unpack success")
		return nil
	}


	oldArgs := []string{
		0: "x",
		1: "-y",
		2: "-p"+SSU_DEC_PASSWD_OLD,
		3: U.SSUPackage,
		4: "-o"+ filepath.Join(U.CurrentWorkFolder,U.SingleUnpkg),
	}
	old := exec.Command(UnpackTool,oldArgs...)
	stdout, _ = old.StdoutPipe()
	if err := old.Start(); err != nil {return err}
	data, err = ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("unpack log has been lost:",err)
	}
	if err := ioutil.WriteFile(filepath.Join(U.CurrentWorkFolder,"7z.log"),data,0664); err != nil {
		fmt.Println("unpack log can't write it to 7z.log:",err)
	}
	if err := old.Wait(); err != nil {
		return err
	}else{
		fmt.Println("unpack success")
		return nil
	}


}

func UnpackCfg()  {

}

func FreeUpdateDir(){

}


func FreeCfgDir(){
	
}

func UnpackPackage(U *Update)error{
	if U.SSUType == PACKAGE_TYPE || U.SSUType == RESTORE_TYPE {
		return unpackPackage(U)
	}
	return nil
}

func ObtainApps(AppPath string)[]string  {

	return
}

func DesApps(AppPath string) []string{

	return
}

func LoadAppData (AppPath string) {

}

func UpdateApps(S *Session,apps, path string) {
	return
}

func UpdateSinglePacket(){
	CheckUpdateCondition()
}


func CheckUpdateCondition()error{
	Put()
	Exec()
}

func InitClient(appVersion []byte) *Update {
	U := new(Update)
	U.FolderPrefix = GetRandomString(32)
	U.CurrentWorkFolder = GetCurrentDirectory()
	if IsArmChip(appVersion) {
		U.TempExecFile, U.TempRstFile = ARM_LINUX_BASIC[0], ARM_LINUX_BASIC[1]
		U.CustomErrFile, U.TempRetFile = ARM_LINUX_BASIC[2], ARM_LINUX_BASIC[3]
		U.LoginPwdFile, U.Compose = ARM_LINUX_BASIC[4], ARM_LINUX_BASIC[5]

		U.ServerAppRe, U.ServerAppSh = ARM_LINUX_UPDATE[0], ARM_LINUX_UPDATE[1]
		U.ServerCfgPre, U.ServerCfgSh = ARM_LINUX_UPDATE[2], ARM_LINUX_UPDATE[3]

		U.LocalBackSh = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/arm_bin/bakcfgsh")
		U.LocalPreCfgSh = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/arm_bin/prercovcfgsh")
		U.LocalCfgSh = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/arm_bin/rcovcfgsh")
		U.LocalUpdHistory = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/arm_bin/updhistory.sh")
		U.LocalUpdCheck = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/arm_bin/updatercheck.sh")

		fmt.Println("The device is a arm platform,init arm info.")
	} else {
		U.TempExecFile, U.TempRstFile = X86_LINUX_BASIC[0], X86_LINUX_BASIC[1]
		U.CustomErrFile, U.TempRetFile = X86_LINUX_BASIC[2], X86_LINUX_BASIC[3]
		U.LoginPwdFile, U.Compose = X86_LINUX_BASIC[4], X86_LINUX_BASIC[5]

		U.ServerAppRe, U.ServerAppSh = X86_LINUX_UPDATE[0], X86_LINUX_UPDATE[1]
		U.ServerCfgPre, U.ServerCfgSh = X86_LINUX_UPDATE[2], X86_LINUX_UPDATE[3]

		U.LocalBackSh = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/bin/bakcfgsh")
		U.LocalPreCfgSh = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/bin/prercovcfgsh")
		U.LocalCfgSh = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/bin/rcovcfgsh")
		U.LocalUpdHistory = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/bin/updhistory.sh")
		U.LocalUpdCheck = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/bin/updatercheck.sh")

		fmt.Println("The device is a x86 platform,init x86 info.")
	}
	return U
}

func InitEnvironment(U *Update) error {
	fmt.Println("now init enviroment for update or restore")
	U.SingleUnpkg = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/unpkg/")
	U.ComposeUnpkg = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/compose_unpkg/")
	U.PkgTemp = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/pkg_tmp/")
	U.Download = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/download/")
	U.AutoBak = filepath.Join(U.CurrentWorkFolder, U.FolderPrefix, "/autobak/")
	if err := InitDirectory(U.SingleUnpkg); err != nil {return err}
	if err := InitDirectory(U.ComposeUnpkg); err != nil {return err}
	if err := InitDirectory(U.PkgTemp); err != nil {return err}
	if err := InitDirectory(U.Download); err != nil {return err}
	if err := InitDirectory(U.AutoBak); err != nil {return err}
	return nil
}
//read file  from start to end
func ReadMd5FromPackage(ssuPath string, start,end int64) (string,error){
	if start < 0 || end < 0 || start > end {
		fmt.Println("params start or end is wrong")
		return "",fmt.Errorf("params start or end is wrong\n")
	}
	file, err := os.Open(ssuPath)
	if err != nil{
		return "",err
	}
	length := end-start
	buf := make([]byte,length)
	_,err = file.Seek(start,1)
	n, err := io.ReadFull(file,buf)
	if err != nil && int64(n) != length{
		return "",err
	}
	return string(buf),nil
}

func ComposePackageMd5(ssuPath string)error{
	ssuMd5, err := ReadMd5FromPackage(ssuPath,8,40)
	if err != nil {
		return err
	}
	if ssuMd5 == Md5Sum(ssuMd5,48) {
		return nil
	} else {
		return fmt.Errorf("compose package md5 don't match\n")
	}
}



func ComposePackage(ssuPath string) bool{
	if ComposePackageMd5(ssuPath) == nil{
		if filepath.Ext(ssuPath) == ".cssu" {
			return true
		}else {
			fmt.Println("The package is a cssu file,but not have a .cssu extname.")
			return false
		}
	}else {
		return false
	}
}

func InitComposePackageArr(ssuPath string) []string {
	return
}

func SinglePackageMd5(ssuPath string) error {
	ssuMd5, err := ReadMd5FromPackage(ssuPath,0,32)
	if err != nil {
		return err
	}
	if ssuMd5 == Md5Sum(ssuMd5,33) {
		return nil
	} else {
		return fmt.Errorf("single package md5 don't match\n")
	}
}


func PrepareUpgrade(S *Session, U *Update) error {
	fmt.Println("init to upgrade or restore  the package:%s", U.SSUPackage)
	if U.UpdatingFlag && (time.Now().Sub(U.UpdateTime) < UPD_TIMEOUT * time.Second ) {
		fmt.Errorf("now update the package:%s,begin at %v\n ....",U.UpdateTime)
	}
	if err := InitEnvironment(U); err != nil {return err}
	if err := FtpDownloadSSUPackage(U.SSUPackage); err != nil {return err}
	if !IsPathExist(U.SSUPackage){
		return fmt.Errorf("can't find the SSU package,please check it\n");
	}

	if ComposePackage(U.SSUPackage){
		InitComposePackageArr(U.SSUPackage) //TODO: not done yet
	}else if SinglePackageMd5(U.SSUPackage){
		//TODO:
		/*
		@package_arr = Array.new
		packhash = {"packet" => now_package, "type" => "1"}
		@package_arr<<packhash
		*/
		U.SSUType = PACKAGE_TYPE
	}else {
		return fmt.Errorf("The package is not a valid package,please check first. if your use a ftp path,please download it to local and try again.")
	}

	return nil
}



