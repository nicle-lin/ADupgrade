package main

import (
	//"net"
	"fmt"
	"os"
	"log"
	"strings"
	"sync"
	"time"
	"regexp"
	"math/rand"
	"path/filepath"
)

func main() {
	/*
	conn, err := net.Dial("tcp","192.168.1.100:51111")
	if err != nil{
		fmt.Println("Dial error:",err)
	}
	var readdata []byte
	conn.Read(readdata)
	fmt.Println("readdata:",readdata)
	conn.Write([]byte("from go"))
	*/
	ip := "192.168.1.100"
	port := "5000"
	fmt.Println("ip + port:",ip+":"+port)

	file, err := os.OpenFile(
		"test.txt",
		os.O_WRONLY|os.O_TRUNC|os.O_CREATE,
		0666,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	writebyte := []byte{0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x3a, 0x30, 0xa, 0x6d, 0x73, 0x67, 0x3a, 0x33, 0x35, 0xa, 0x0}
	fmt.Println("string:",string(writebyte))
	ret := strings.Contains(string(writebyte),"result:1")
	fmt.Println("ret:",ret)
	// 写字节到文件中
	//byteSlice := []byte("Bytes!\n")
	_ , err = file.Write(writebyte)
	if err != nil {
		log.Fatal(err)
	}

	name := Name{i:10}
	var once sync.Once
	for i:=0; i <5;i++{
		go func(i int){
			once.Do(name.OnceDo)
			fmt.Println("iii:",i)
		}(i)
		//fmt.Println("okokiiii:",i)
	}
	time.Sleep(1*time.Second)
	fmt.Println("-----------------------------------------")
	text := "\nverversion:450\n"
	reg := regexp.MustCompile(`version:[\d]+`)
	fmt.Printf("%q\n", reg.FindAllString(text, -1))
	str := reg.FindAllString(text,-1)[0]
	version := strings.Split(str,":")[1]
	fmt.Println("str:",version)
	var all []byte
	one := make([]byte,10)
	one[9] = 2
	two := make([]byte, 20)
	two[19] = 1
	all = one
	all = append(all,two...)
	fmt.Println("all:",all)
	fmt.Println("######################################################")
	text1 := "\nSANGFOR-M5400-AD-5.4R1\nBUILD20160830"
	reg1 := regexp.MustCompile(`[\w]+-[\w]+\.[\w]+`)
	//fmt.Printf("%q\n",reg1.FindAllString(text1,-1))
	str1 := reg1.FindAllString(text1,-1)[0]
	version1 := strings.Split(str1,"-")[1]
	fmt.Println("version AD:",version1)
	fmt.Println("#######################################################")
	fmt.Println("random string:",GetRandomString(32))

}

//生成随机字符串
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

func getCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}


type Name struct {
	i int
}

func (n *Name)OnceDo(){
	fmt.Println("i:",n.i)
}