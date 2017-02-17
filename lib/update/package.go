package update

import (
	"fmt"
	"path/filepath"
	"sync"
	"time"
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

func InitEnviroment(U *Update) error {
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

func PrepareUpgrade(S *Session, U *Update) error {
	fmt.Println("init to upgrade or restore  the package:%s", U.SSUPackage)
	if U.UpdatingFlag && (time.Now().Sub(U.UpdateTime) < UPD_TIMEOUT * time.Second ) {
		fmt.Errorf("now update the package:%s,begin at %v\n ....",U.UpdateTime)
	}

	return nil
}
