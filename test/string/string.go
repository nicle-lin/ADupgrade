package main

import (
	"fmt"
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
}
