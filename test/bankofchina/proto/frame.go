package proto

import (
	"fmt"
	"io"
	"net"
	"strings"
	"time"
	"math/rand"
)

/* base frame
* a frame begin with "0xDB0xF3".....
 */

const (
	FRAME_HEADER_LEN = 10                 //a frame header is 10 bytes
	FRAMEFLAG        = 0x31323334f1f2b1b2 //a frame is started with "0xf3db"
	FRAMEFLAG1        = 0x35363738d3d4a3a4 //a frame is started with "0xf3db"

	MAX_FRAME_LEN    = 65535
)

type Frame interface {
	ReadFrame(b []byte, conn net.Conn) (n int, err error)

	WriteFrame(b []byte, conn net.Conn) (n int, err error)
}

type MBLB struct {
	net.Conn
}

func GetRandomNumber(number int) int {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Intn(number)
}

func GetRandomFloat() float32{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Float32()
}

func GetRandom64Number() uint64{
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return r.Uint64()
}

/*
 * 格式:[2byte][8byte][data]
 *      前两2个是协议长度(10+data的长度)　　后面8个是标志　　　data为真实数据
        [2byte][8byte][8byte][data]为了返回时能够校验8byte的校验码，特地把校验码加到校验码的后面，
        当然服务端回复时必须把检验码带回来
 */
func BuildFrame(data []byte, randomNum int, isServer uint64) (MultiFrame []byte,err error) {
	for i := 0; i < randomNum; i++ {
		lenScale := GetRandomNumber(500)
		length := len(data) * lenScale
		if length > MAX_FRAME_LEN {
			return nil, fmt.Errorf("message too long\n")
		}
		frameHeader := make([]byte, FRAME_HEADER_LEN + length)
		f := NewBEStream(frameHeader)
		f.WriteUint16(uint16(length) + 10 + 8)
		frameFlag := GetRandom64Number()
		f.WriteUint64(frameFlag)
		if isServer == 0 {   //如果是server端，需要把读到frameFlag传进来
			f.WriteUint64(isServer) //写两遍
		}else {
			//如果client端，则不需要传进来，传个0就行了
			f.WriteUint64(frameFlag) //写两遍
		}

		for i := 0; i < lenScale; i++{
			err = f.WriteBuff(data)
			if err != nil {
				return nil, err
			}
		}
		MultiFrame = append(MultiFrame,f.buff[:f.pos]...)
	}

	return MultiFrame, nil
}


func WriteFrame(data[]byte,randomNum int,isServer uint64, conn net.Conn) (int, error) {
	frame , err := BuildFrame(data,randomNum,isServer)
	if err != nil{
		return 0, err
	}
	return conn.Write(frame)
}

func ReadFrame(conn net.Conn,randomNum int,flag bool) (frameFlag2 uint64,err error) {
	for i := 0; i < randomNum; i++ {
		//step 1: 分配frame头部长度的大小的空间
		frameHeader := make([]byte, FRAME_HEADER_LEN)
		var n int
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

			frameFlag2,_ = f.ReadUint64() //服务端方向，把读到的frameFlag返回，然后再发送回客户端
			ip := []string{
				fmt.Sprintf("%d",uint8(a)),
				fmt.Sprintf("%d",uint8(b)),
				fmt.Sprintf("%d",uint8(c)),
				fmt.Sprintf("%d",uint8(d)),
			}
			sip := strings.Join(ip,".")


			fmt.Println("Got ip:",sip)
			fmt.Println("Got port:",sport)
			fmt.Printf("other:0x%x\n",other)
			fmt.Printf("the whole frame flag:0x%x%x%x%x%x%x\n",a,b,c,d,sport,other)


		}else{
			frameFlag, errFlag := f.ReadUint64()   //第一个frame会被处理，
			if errFlag != nil {
				return FRAME_HEADER_LEN, fmt.Errorf("read frame flag is wrong:%s\n", frameFlag)
			}
			realFrameFlag, _ := f.ReadUint64()   //第二个frame不会被处理，是用于比较第一个的，防止第一个在传输过程中被处理错了
			if frameFlag != realFrameFlag {
				fmt.Printf("expect flag:0x%x Got flag: 0x%x\n",realFrameFlag,frameFlag)
				return FRAME_HEADER_LEN, fmt.Errorf("expect flag:0x%x Got flag: 0x%x\n",realFrameFlag,frameFlag)
			}
			fmt.Printf("expect flag:0x%x == Got flag: 0x%x\n",realFrameFlag,frameFlag)
		}

		if frameLength > MAX_FRAME_LEN {
			return FRAME_HEADER_LEN, fmt.Errorf("frameLength wrong:0x%x", frameLength)
		}
		//step 3: 读取frame中的数据,data的长度是frameLength长度减去协议头部长度
		frameBody := make([]byte, frameLength - 10 - 8)

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
	}

	return frameFlag2, nil
}
