package main

import (
	"fmt"
	"reflect"
	"unsafe"
)

const SSU_DEC_PASSWD  = "sangforupd~!@#$%"

//const X86_LINUX_UPDATE = ["/etc/dlancmd/apppre","/etc/dlancmd/appsh","/etc/dlancmd/cfgpre","/etc/dlancmd/cfgsh"]
/*
const X86 = map[int] string{
	0: "appre",
	1: "appsh",
	2: "cfgpre",
	3: "cfgsh",

}
*/
var KK = [2]string{0: "appre", 1: "appsh"}

func main() {
//	fmt.Println("x86:",X86_LINUX_UPDATE)
	fmt.Println("ssu_dec_passwd:",SSU_DEC_PASSWD)
//	fmt.Println("x86:",X86)
	fmt.Println("KK:",KK)
	//var e,sh string
	Re,Sh := KK[0],KK[1]
	fmt.Println("re,sh:",Re,Sh)
	println("---------------------------------")
	str()
	str2()
	str3()

}
func str(){
	s := "ab" +
		"cd"
	println(s == "abcd")
	println(s > "abc")
}

func str2(){
	s := "dabc"

	println(s[1])
	//s[1] = "f"
	//println(s[1])
	//println(&s[1])
	println(&s)
}

func pp(format string, ptr interface{}){
	p := reflect.ValueOf(ptr).Pointer()
	h := (*uintptr)(unsafe.Pointer(p))
	fmt.Printf(format,*h)
}

func str3(){
	s := "hello world!"
	pp("s:%x\n",&s )

	bs := []byte(s)
	bs[1] = 98
	bs[2] = 'g'
	s2 := string(bs)
	fmt.Println("s2:", s2)

	pp("string to []byte, bs: %x\n",&bs)
	pp("[]byte to string, s2: %x\n", &s2)

	rs := []rune(s)
	s3 := string(rs)

	pp("string to []rune, rs :%x\n", &rs)
	pp("[]rune to string, s3:%x\n", &s3)
}