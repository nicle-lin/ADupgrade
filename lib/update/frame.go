package update

import "fmt"

/* base frame
* a frame begin with "0xDB0xF3".....
 */

const (
	FRAME_HEADER_LEN = 4    //a frame header is 4 bytes
	SECDATA_HEADER_LEN = 5  //a secData header len
	FRAMEFLAG       = 0xf3db //a frame is started with "0xf3db"
	MAX_DATA_LEN     = 1024 // the max frame data length
	MAX_FRAME_LEN    = 1080 + FRAME_HEADER_LEN
	CMDFRAME         = 0    //command frame flag
	DATAFRAME        = 1    //data frame flag

	// used by cmd data.
	LOGIN   = iota
	EXEC     //exec cmd,be care!!!!!
	GET
	GETOVER
	PUT
	PUTOVER
	VERSION
	MaxCmdLen
)

var CMD = [MaxCmdLen]string{
	LOGIN:   "login",
	EXEC:    "excute",
	GET:     "get",
	GETOVER: "getover",
	PUT:     "putover",
	PUTOVER:  "putover",
	VERSION:  "version",
}

type Frame struct {
	flag uint16 
	length uint16
	data []byte

}

type SecData struct {
	flag uint16
	length uint16
	typ uint8
	data []byte
}


//命令格式组合
func JoinCmd(cmd string,params [][2]string)[]byte{
	var b []byte
	b = append(b,[]byte(cmd)...)
	b = append(b,[]byte("\n")...)
	for _, v := range params{
		b = append(b,[]byte(v[0])...)
		b = append(b, []byte(":")...)
		b = append(b, []byte(v[1])...)
		b = append(b,[]byte("\n")...)
	}
	//in go lang,it must octal express Null character
	//b = append(b, []byte("\000")...)
	length := len(b)
	return b[:length]
}

func MakeCmdStr(cmdType,command string)[]byte{
	switch cmdType {
	case CMD[LOGIN]:
		return JoinCmd(CMD[LOGIN],[][2]string{{"passwd", command},{"flage","HandleVersion"}})
	case CMD[EXEC]:
		return JoinCmd(CMD[EXEC],[][2]string{{"cmd",command}})
	case CMD[GET]:
		return JoinCmd(CMD[GET],[][2]string{{"file",command}})
	case CMD[PUT]:
		return JoinCmd(CMD[PUT],[][2]string{{"file",command}})
	case CMD[PUTOVER]:
		return JoinCmd(CMD[PUTOVER],[][2]string{{}})
	case CMD[VERSION]:
		return JoinCmd(CMD[VERSION],[][2]string{{"value","1280"}})
	default:
		return nil
	}

}

func MakeCmdPacket(cmdType string, params string) ([]byte,int){
	cmdByte := MakeCmdStr(cmdType,params)
	fmt.Printf("cmdByte:%#v\n",cmdByte)

	fmt.Println("-------------------------------------")
	return BuildFrame(CMDFRAME,cmdByte)
}

func MakeDataPacket(content []byte)([]byte,int){
	return BuildFrame(DATAFRAME,content)
}


/*        1 byte      1 byte   　2 byte    1byte       1 byte      2byte  1 byte  less than 1024 byte
 * 格式:[FRAMEFLAG0][FRAMEFLAG0][length]([FRAMEFLAG0][FRAMEFLAG0][length][flag][data........])
 *      前两2个是协议开头标志　　后面数据字节数　　　括号里是加密的数据　flag表示是数据还是命令　data为真实数据　　　　　　　　
 */
func BuildFrame(flag byte, content[]byte)([]byte, int){
	secBuff := make([]byte,MAX_DATA_LEN)
	sec := NewLEStream(secBuff)
	sec.WriteUint16(FRAMEFLAG)
	contentLength := len(content)
	sec.WriteUint16(uint16(contentLength))
	sec.WriteByte(flag)
	sec.WriteBuff(content)

	fmt.Printf("secData:%#v\n",sec.buff[:sec.pos])

	//fmt.Println("secData:",secData)

	frameBuff := make([]byte,MAX_FRAME_LEN)
	frame := NewLEStream(frameBuff)
	frame.WriteUint16(FRAMEFLAG)
	secLength := EncLen(sec.pos)
	fmt.Printf("before enc secData len:%x\n",sec.pos)
	fmt.Printf("after enc secData len:%x\n",secLength)
	frame.WriteUint16(uint16(secLength))

	tempBuff := make([]byte,MAX_FRAME_LEN)
	//function Encrypt will combine secData and FrameData

	buf,_ := Encrypt(sec.buff[:sec.pos],tempBuff)

	frame.WriteBuff(buf)
	fmt.Printf("whole frame length:%x\n",frame.pos)
	return frame.buff[:frame.pos],frame.pos
}

func ReadPacket(data []byte){

}




func NewFrame(flag int, content []byte) *Frame {



	return &Frame{

	}
}

