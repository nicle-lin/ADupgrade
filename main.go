package main

import (
	"github.com/nicle-lin/ADupgrade/lib/update"
	"fmt"
	"net"
	"time"
)

func main() {
	msg  := []byte{'h','e','l','l','0'}
	//msgstring := "login\npasswd:dlanrecover\nflage:HandleVersion\n\x00"
	//msgstring := "hello"
	//msg := []byte(msgstring)
	//fmt.Println("msg:",msg)
	fmt.Printf("msg:%#v\n",msg)
	fmt.Println("len:",update.EncLen(len(msg)))

	var outmsg []byte
	var err error
	outmsg, err = update.Encrypt(msg,outmsg)
	if err == nil{
		//fmt.Println("enc msg:",outmsg)
		fmt.Printf("enc msg:%#v\n",outmsg)
	}else{
		fmt.Println("err:",err)
	}

	decmsg, _ := update.Decrypt(outmsg,outmsg)
	fmt.Println("decmsg:",decmsg)

	cmdstr ,outlen:= update.MakeCmdPacket("login","admin")
	//fmt.Println("cmdstr:",cmdstr)
	fmt.Printf("cmdstr len:%f\n",outlen)
	fmt.Printf("cmdstr:%#v\n",cmdstr)

	conn, err := net.Dial("tcp4", "192.168.1.100:51111")
	if err != nil {
		fmt.Println("dial error:",err)
		return
	}

	/**********************************************************
	capture data from real connection
	//login

	************************************************************
	real_login := []byte{0xdb, 0xf3, 0x32, 0x00, 0x2d, 0x00, 0x32, 0x0d,
		0x0f, 0x95, 0xee, 0x94, 0x81, 0x00, 0x93, 0xf1,
		0x49, 0xca, 0xf9, 0x1c, 0x3e, 0xb9, 0xa2, 0x25,
		0x38, 0x83, 0x22, 0xc9, 0xdb, 0xe2, 0x31, 0x89,
		0x72, 0xe2, 0x10, 0x7d, 0xc2, 0xa0, 0x7d, 0xe9,
		0x85, 0x42, 0xfd, 0x0b, 0xd0, 0x2e, 0x31, 0x34,
		0x88, 0xb5, 0x72, 0x6b, 0x43, 0xe6 }


	fate_login := []byte{0xdb, 0xf3, 0x32, 0x00, 0x2a, 0x00, 0x61, 0x29,
		0xe7, 0x23, 0xf5, 0x03, 0x6a, 0x8f, 0x93, 0xf1,
		0x49, 0xca, 0xf9, 0x1c, 0x3e, 0xb9, 0x8c, 0xe6,
		0x50, 0x63, 0x2e, 0xdf, 0xc4, 0xec, 0x1b, 0xba,
		0x6b, 0x6f, 0xa7, 0x36, 0xc9, 0xe9, 0x6c, 0xee,
		0x2b, 0x84, 0xc9, 0x75, 0xf0, 0xbd, 0xb9, 0x9f,
		0x64, 0x53, 0x0a, 0x58, 0xe6, 0xea }
	*/


	_, err = conn.Write(cmdstr)
	if err != nil{
		fmt.Println("write buff:",err)
	}
	readBuf3, _  := update.ReadPacket(conn)
	fmt.Println("read dec data string:",string(readBuf3[4:]))

	//cmdstr2 ,_:= update.MakeCmdPacket("version","")
	//_, err = conn.Write(cmdstr2)

	cmdstr1 ,_:= update.MakeCmdPacket("get","/app/appversion")
	_, err = conn.Write(cmdstr1)

	//readbuf := make([]byte,update.MAX_DATA_LEN)
	//var n int
	err = conn.SetReadDeadline(time.Now().Add(time.Second * 10))
	if err != nil{
		fmt.Printf("set readdeadline error:",err)
	}

	//n , err = conn.Read(readbuf)
	readBuf, err2  := update.ReadPacket(conn)
	if err2 != nil{
		fmt.Println("read error:",err2)
	}
	//fmt.Println("read ",n ,"byte")
	fmt.Printf("read sec dec data:%#v\n",readBuf)
	fmt.Println("read dec data string:",string(readBuf[4:]))

	readBuf, err2  := update.ReadPacket(conn)
	fmt.Println("appversion data")

	time.Sleep(10*time.Second)
	_ = conn.Close()
}
