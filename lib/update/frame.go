package update

import (

)



/* base frame
* a frame begin with "0xDB0xF3".....
 */

const (
	FRAME_HEADER_LEN = 4    //a frame header is 4 bytes
	FRAMEFLAG0       = 0xDB //a frame is started with "\xdb\xf3"
	FRAMEFLAG1       = 0xF3 //
	MAX_DATA_LEN     = 1024 // the max frame data length
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
	}
	b = append(b, []byte("\0")...)
	return b
}

func MakeCmdStr(cmdtype,command string)[]byte{
	switch cmdtype {
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

func MakeCmdPacket(cmd string, params...string){

}

func MakeDataPacket(content []byte){
	BuildPacket(DATAFRAME,content)
}

//写错了
//TODO:协议格式写错了
func BuildPacket(flag byte, content[]byte)[]byte{
	sec := NewLEStream()
	sec.WriteByte(FRAMEFLAG0)
	sec.WriteByte(FRAMEFLAG1)
	contentLength := len(content)
	//sec.WriteByte(contentLength%256)
	//sec.WriteByte(contentLength/256)
	sec.WriteUint16(contentLength)
	sec.WriteByte(flag)
	sec.WriteBuff(content)
	SecData := des_enc()


	frame := NewLEStream()
	frame.WriteByte(FRAMEFLAG0)
	frame.WriteByte(FRAMEFLAG1)
	secLength := len(sec.buff)
	//frame.WriteByte(secLength%256)
	//frame.WriteByte(secLength/256)
	frame.WriteUint16(secLength)
	frame.WriteBuff(SecData)

	totalLength := len(frame.buff)
	return frame.buff[:totalLength]
}

func ReadPacket(data []byte){

}




func NewFrame(flag int, content []byte) *Frame {



	return &Frame{

	}
}

func HeaderLen() int{
	return FRAME_HEADER_LEN
}

func MaxDataLen() int{
	return MAX_DATA_LEN
}

//TODO:
func SplitByLength(str string,len int)[]string{

	return []string("")
}

