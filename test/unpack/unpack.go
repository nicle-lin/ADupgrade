package main

import (
	"fmt"
	"runtime"
	"path/filepath"
	"os/exec"
	"os"
)

var SSU_DEC_PASSWD        = "sangforupd~!@#$%"
var SSU_DEC_PASSWD_OLD    = "greatsinfor"

func GetCurrentDirectory() string {
	/* it only work in linux correctly

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Get Current Directory wrong:",err)
		return "/tmp"
	}
	return strings.Replace(dir, "\\", "/", -1)
	*/

	pwd, _ := os.Getwd()
	return pwd
}

type Update struct {
	CurrentWorkFolder string
	SSUPackage string
	SingleUnpkg string
}

func unpackPackage(U *Update)error {
	// function InitEnvironment has been init the path U.SingleUnpkg
	fmt.Println("begin to unpack the package")
	var UnpackTool string
	if runtime.GOOS	 == "windows"{
		UnpackTool = filepath.Join(U.CurrentWorkFolder,"tool","7z.exe")
	}else{
		UnpackTool = "/usr/bin/7za"
	}
	//7za x -y -psangforupd\~\!\@\#\$\%  /home/gubl/Desktop/sangfor/ADUpgrade/AD6.5\(20160809\).ssu  -o/home/gubl/web/shipyard/src/github.com/nicle-lin/ADupgrade/test/unpack/singleunpkg
	//oldPasswdCommand := UnpackTool + " x -y -p" + SSU_DEC_PASSWD_OLD + " " + U.SSUPackage + " -o" + U.SingleUnpkg + " > 7z.log"
	//fmt.Println("old:",oldPasswdCommand)
	newPasswdCommand := UnpackTool + " x -y -p" + SSU_DEC_PASSWD + " " + "\\\"" + U.SSUPackage + "\\\"" + " -o" + U.SingleUnpkg + " > 7z.log"
	fmt.Println("new:",newPasswdCommand)
	//old := exec.Command(oldPasswdCommand)
	new := exec.Command(newPasswdCommand)
	errnew := new.Run()

	if errnew == nil {
		return nil
	}else {
		/*
		errold := old.Run()
		if errold != nil {
			return errold
		}else{
			return nil
		}
		*/
		return errnew
	}
}

func main() {

	U := new(Update)
	U.CurrentWorkFolder = GetCurrentDirectory()
	U.SSUPackage = os.Args[1]
	//U.SingleUnpkg = filepath.Join(U.CurrentWorkFolder,"singleunpkg")
	U.SingleUnpkg = "singleunpkg"
	if err := unpackPackage(U); err != nil {
		fmt.Println("error:",err)
	}
}
