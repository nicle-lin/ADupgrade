package main

import (
	"github.com/nicle-lin/ADupgrade/lib/update"
	"fmt"
	"net"
)

func main() {
	msg  := []byte{'h','e','l','l','0'}
	//msgstring := "login\npasswd:dlanrecover\nflage:HandleVersion\n\x00"
	//msgstring := "hello"
	//msg := []byte(msgstring)
	fmt.Println("msg:",msg)
	fmt.Println("len:",update.EncLen(len(msg)))

	var outmsg []byte
	var err error
	outmsg, err = update.Encrypt(msg,outmsg)
	if err == nil{
		fmt.Println("enc msg:",outmsg)
	}else{
		fmt.Println("err:",err)
	}

	decmsg, _ := update.Decrypt(outmsg,outmsg)
	fmt.Println("decmsg:",decmsg)

	cmdstr ,_:= update.MakeCmdPacket("login","admin")
	fmt.Println("cmdstr:",cmdstr)

	conn, err := net.Dial("tcp4", "192.168.1.100:51111")
	if err != nil {
		fmt.Println("dial error:",err)
		return
	}
	_, err = conn.Write(cmdstr)
	if err != nil{
		fmt.Println("write buff:",err)
	}
/*
	cmdstr1 ,_:= update.MakeCmdPacket("get","/app/appversion")
	_, err = conn.Write(cmdstr1)
*/
	var readbuf []byte
	var n int
	n , err = conn.Read(readbuf)
	if err != nil{
		fmt.Println("read error:",err)
	}
	fmt.Println("read ",n ,"byte")
}
