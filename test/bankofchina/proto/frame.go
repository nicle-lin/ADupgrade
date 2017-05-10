package proto

import (
	"fmt"
	"io"
	"net"
	"strings"
)

/* base frame
* a frame begin with "0xDB0xF3".....
 */

const (
	FRAME_HEADER_LEN = 10                 //a frame header is 10 bytes
	FRAMEFLAG        = 0x01020304f2b1b2 //a frame is started with "0xf3db"
	MAX_FRAME_LEN    = 1020 + FRAME_HEADER_LEN
)

type Frame interface {
	ReadFrame(b []byte, conn net.Conn) (n int, err error)

	WriteFrame(b []byte, conn net.Conn) (n int, err error)
}

type MBLB struct {
	net.Conn
}

/*
 * 格式:[2byte][8byte][data]
 *      前两2个是协议长度(10+data的长度)　　后面8个是标志　　　data为真实数据
 */
func WriteFrame(b []byte, conn net.Conn) (n int, err error) {
	length := len(b)
	if length > MAX_FRAME_LEN {
		return 0, fmt.Errorf("message too long\n")
	}
	frameHeader := make([]byte, FRAME_HEADER_LEN + length)
	f := NewBEStream(frameHeader)
	f.WriteUint16(uint16(len(b)) + 10)
	f.WriteUint64(FRAMEFLAG)
	err = f.WriteBuff(b)
	if err != nil {
		return 0, err
	}
	return conn.Write(f.buff[:f.pos])
}

func ReadFrame(conn net.Conn,flag bool) (int, error) {
	//step 1: 分配frame头部长度的大小的空间
	frameHeader := make([]byte, FRAME_HEADER_LEN)
	var n int
	var err error
	var realNeed int = 0
	//step 2:　读取frame头部长度大小的数据
	for {
		n, err = conn.Read(frameHeader[realNeed:])
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("Read Frame error:%s", err)
		} else if err == io.EOF {
			return 0, err
		}
		realNeed = realNeed + n
		if realNeed == FRAME_HEADER_LEN || 0 == n {
			realNeed = 0
			break
		}
	}

	f := NewBEStream(frameHeader)

	frameLength, errLength := f.ReadUint16()
	if errLength != nil {
		return FRAME_HEADER_LEN, fmt.Errorf("read frame length is wrong: %s\n", errLength)
	}
	//服务端收数据不需要检验标志,把标志打出来
	if flag {
		a, _ := f.ReadByte()
		b, _ := f.ReadByte()
		c, _ := f.ReadByte()
		d, _ := f.ReadByte()
		sport, _ := f.ReadUint16()
		other, _ := f.ReadUint16()
		ip := []string{
			fmt.Sprintf("%s",string(a)),
			fmt.Sprintf("%s",string(b)),
			fmt.Sprintf("%s",string(c)),
			fmt.Sprintf("%s",string(d)),
		}
		sip := strings.Join(ip,".")


		fmt.Println("Got ip:",sip)
		fmt.Println("Got port:",sport)
		fmt.Println("other:",other)


	}else{
		frameFlag, errFlag := f.ReadUint64()
		if errFlag != nil {
			return FRAME_HEADER_LEN, fmt.Errorf("read frame flag is wrong:%s\n", frameFlag)
		}
		if frameFlag != FRAMEFLAG {
			return FRAME_HEADER_LEN, fmt.Errorf("frame flag is wrong:0x%x", frameFlag)
		}
	}

	if frameLength > MAX_FRAME_LEN {
		return FRAME_HEADER_LEN, fmt.Errorf("frameLength wrong:0x%x", frameLength)
	}
	//step 3: 读取frame中的数据,data的长度是frameLength长度减去协议头部长度
	frameBody := make([]byte, frameLength - 10)

	for {
		n, err = conn.Read(frameBody[realNeed:])
		if err != nil && err != io.EOF {
			return 0, fmt.Errorf("read frameBody error:", err)
		} else if err == io.EOF {
			return 0, err
		}
		realNeed = realNeed + n
		if realNeed == int(frameLength) || n == 0 {
			realNeed = 0
			break
		}

	}
	fmt.Printf("Got message:%s\n", string(frameBody))
	return int(frameLength), nil
}
