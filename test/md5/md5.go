package main

import (
	"os"
	"io"
	"crypto/md5"
	"fmt"
	"bufio"
	"encoding/hex"
	"path/filepath"
	"hash"
)


func md5sum3(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	r := bufio.NewReader(f)

	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))

}

func md5sumstring() string {
	md5ctx := md5.New()
	md5ctx.Write([]byte("gubl"))
	fmt.Println("string md5:",string(md5ctx.Sum(nil)))
	fmt.Println(hex.EncodeToString(md5ctx.Sum(nil)))
	return "gubl"
}

/*
func ComposePackageMd5(ssuPath string)error{
	file, err := os.Open(ssuPath)
	if err != nil{
		return err
	}
	buf := make([]byte,40)
	n, err := io.ReadFull(file,buf)
	if err != nil && n != 40{
		return err
	}
	//SSUMd5 := string(buf[8:])
	//md5.New()
	fmt.Println("ture md5:",string(buf[8:]))
	return nil

}
*/

func Md5SumFile(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	_,_ = f.Seek(48,1)
	r := bufio.NewReader(f)

	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

func Digest(msg []byte) hash.Hash{
	var h hash.Hash = md5.New()
	h.Write(msg)
	return h
}

func Md5SumString(data []byte) string {
	h := md5.New()
	h.Write(data)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func Md5Sum(arg interface{}) string {
	switch v := arg.(type) {
	case []byte:
		fmt.Println("[]byte")
		return Md5SumString(v)
	case string:
		fmt.Println("string")
		return Md5SumFile(string(v))
	default:
		fmt.Println("default")
		return fmt.Sprintf("%x is wrong type",arg)
	}
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
	fmt.Println("real md5:",string(ssuMd5))
	fmt.Println("data md5:",Md5Sum(ssuMd5))
	if string(ssuMd5) == Md5Sum(ssuMd5) {
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

func SinglePackageMd5(ssuPath string) error {
	ssuMd5, err := ReadMd5FromPackage(ssuPath,1,32)
	if err != nil {
		return err
	}
	if string(ssuMd5) == Md5Sum(ssuMd5) {
		return nil
	} else {
		return fmt.Errorf("single package md5 don't match\n")
	}
}


//prameter offset only need one prameter
func md5SumFileOffset(file string,offset ...int64) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()

	if len(offset) <= 1{
		f.Seek(offset[0],1)
	}else {
		fmt.Println("prameter offset only need one prameter")
	}
	r := bufio.NewReader(f)

	h := md5.New()

	_, err = io.Copy(h, r)
	if err != nil {
		return ""
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}


func main() {
	/*
	fmt.Println("md5:",md5sum3(os.Args[1]))
	if err := ComposePackageMd5(os.Args[1]); err !=nil{
		fmt.Println("wrong")
	}else{
		fmt.Println("right")
	}
	*/
	/*
	h := md5.New()
	//h.Write([]byte("sharejs.com")) // 需要加密的字符串为 sharejs.com
	//fmt.Printf("%s\n", hex.EncodeToString(h.Sum(nil))) // 输出加密结果
	io.WriteString(h, "sharejs.com")
	fmt.Println(fmt.Sprintf("%x", h.Sum(nil)))


	//md5sumstrin()

	fmt.Println(md5sum3(os.Args[1]))
	*/
	//fmt.Println(Md5Sum(1))

	data,_:= ReadMd5FromPackage(os.Args[1],8,40)
	data2, _ := ReadMd5FromPackage(os.Args[1],0,32)
	fmt.Println("second:",data)
	fmt.Println("second:",md5SumFileOffset(os.Args[1],48))
	fmt.Println("third:",data2)
	fmt.Println("third:",md5SumFileOffset(os.Args[1],33))

}
