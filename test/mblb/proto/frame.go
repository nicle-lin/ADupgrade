package proto

import (
	"fmt"
	"io"
	"net"
)

/* base frame
* a frame begin with "0xDB0xF3".....
 */

const (
	FRAME_HEADER_LEN = 4      //a frame header is 4 bytes
	FRAMEFLAG        = 0xf3db //a frame is started with "0xf3db"
	MAX_FRAME_LEN    = 1020 + FRAME_HEADER_LEN
)

type Frame interface {
	ReadFrame(b []byte, conn net.Conn) (n int, err error)

	WriteFrame(b []byte, conn net.Conn) (n int, err error)
}

type MBLB struct {
	net.Conn
}

/*        1 byte      1 byte   　2 byte     less than 1024 byte
 * 格式:[FRAMEFLAG0][FRAMEFLAG0][length](data........])
 *      前两2个是协议开头标志　　后面数据字节数　　　data为真实数据
 */
func WriteFrame(b []byte, conn net.Conn) (n int, err error) {
	length := len(b)
	if length > MAX_FRAME_LEN {
		return 0, fmt.Errorf("message too long\n")
	}
	frameHeader := make([]byte, FRAME_HEADER_LEN+length)
	f := NewLEStream(frameHeader)
	f.WriteUint16(FRAMEFLAG)
	f.WriteUint16(uint16(len(b)))
	err = f.WriteBuff(b)
	if err != nil {
		return 0, err
	}
	return conn.Write(f.buff[:f.pos])
}

func ReadFrame(conn net.Conn) (n int, err error) {
	//step 1: 分配frame长度的大小的空间
	frameHeader := make([]byte, FRAME_HEADER_LEN)
	var n int
	var err error
	var realNeed int = 0
	//step 2:　读取frame长度大小的数据
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

	f := NewLEStream(frameHeader)
	frameFlag, errFlag := f.ReadUint16()
	if errFlag != nil {
		return FRAME_HEADER_LEN, fmt.Errorf("read frame flag is wrong:%s\n", frameFlag)
	}
	frameLength, errLength := f.ReadUint16()
	if errLength != nil {
		return FRAME_HEADER_LEN, fmt.Errorf("read frame length is wrong: %s\n", errLength)
	}
	if frameFlag != FRAMEFLAG {
		return FRAME_HEADER_LEN, fmt.Errorf("frame flag is wrong:0x%x", frameFlag)
	}

	if frameLength > MAX_FRAME_LEN {
		return FRAME_HEADER_LEN, fmt.Errorf("frameLength wrong:0x%x", frameLength)
	}
	//step 3: 读取frame中的数据
	frameBody := make([]byte, frameLength)

	for {
		n, err = conn.Read(frameBody[realNeed:])
		if err != nil && err != io.EOF {
			return fmt.Errorf("[Readpacket] read Sec Data error:", err)
		} else if err == io.EOF {
			return 0, err
		}
		realNeed = realNeed + n
		if realNeed == int(frameLength) || n == 0 {
			realNeed = 0
			break
		}

	}

	return frameLength, nil
}
