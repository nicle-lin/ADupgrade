package main

import (
	//"github.com/nicle-lin/ADupgrade/lib/update"
	"os"
	"fmt"
	//"syscall"
	"path/filepath"
	//"strings"
	"time"
	"math/rand"
	"syscall"
)

func IsPathExist(path string)bool{
	_,err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		return false
	}
	return true
}

func InitDirectory(path string)error {
	if IsPathExist(path){
		if err := os.RemoveAll(path); err != nil {return err}

	}else {
		fmt.Println("error:not exist")
	}
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	if err := os.MkdirAll(path,0775); err != nil {return err }
	return nil
}

func GetCurrentDirectory() string {
	/* it only work in linux correctly

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		fmt.Println("Get Current Directory wrong:",err)
		return "/tmp"
	}
	return strings.Replace(dir, "\\", "/", -1)
	*/

	pwd, _ :=  os.Getwd()
	return pwd
}

func GetRandomString(length int) string{
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

func main() {
	fmt.Println("initdir:",InitDirectory(os.Args[1]))
	mask := syscall.Umask(0)
	defer syscall.Umask(mask)
	os.Mkdir(os.Args[2],0777)
	fmt.Println(os.Getwd())

	fmt.Println("----------------------------")
	fmt.Println("current:",GetCurrentDirectory())

	pwd := filepath.Join(GetCurrentDirectory(),"test","/lillin","/gubl/","")
	fmt.Println("pwd:",pwd)
}
