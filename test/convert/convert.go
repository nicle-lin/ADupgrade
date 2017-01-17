package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"strings"
)

func convert() {
	stringSlice := []string{"通知中心", "perfect!"}

	stringByte := "\x00" + strings.Join(stringSlice, "\x20\x00") // x20 = space and x00 = null

	fmt.Println([]byte(stringByte))

	fmt.Println(string([]byte(stringByte)))
}

func convert2() {
	stringSlice := []string{"通知中心", "perfect!"}

	buffer := &bytes.Buffer{}

	gob.NewEncoder(buffer).Encode(stringSlice)
	byteSlice := buffer.Bytes()
	fmt.Printf("%q\n", byteSlice)

	fmt.Println("---------------------------")

	backToStringSlice := []string{}
	gob.NewDecoder(buffer).Decode(&backToStringSlice)
	fmt.Printf("%v\n", backToStringSlice)
}

func login(){
	login := "login"
	fmt.Println("login:",login)
	var b []byte
	b = append(b,[]byte(login)...)
	fmt.Println("byte login:",b)
}


func array(str [][2]string){
	fmt.Println("length:",len(str))
	for k,v := range str{
		fmt.Println("key:",k, "value0:",v[0], "value1:",v[1])
	}

}

func JoinCmd(cmd string,params [][2]string)[]byte{
	var b []byte
	b = append(b,[]byte(cmd)...)
	b = append(b,[]byte("\n")...)
	for _, v := range params{
		b = append(b,[]byte(v[0])...)
		b = append(b, []byte(":")...)
		b = append(b, []byte(v[1])...)
	}
	return b
}

func main() {
	convert()
	fmt.Println("####################################")
	convert2()
	login()
	str := [][2]string{{"hi", "hello"},{"nicle","lillian"}}
	array(str)
	b := JoinCmd("login",[][2]string{{"passwd", "haha"},{"flage","HandleVersion"}})
	fmt.Println("b:",b)
	fmt.Println("b str:",string(b))
}
