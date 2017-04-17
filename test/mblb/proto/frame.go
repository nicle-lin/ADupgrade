package update

import "fmt"

/* base frame
* a frame begin with "0xDB0xF3".....
 */

const (
	FRAME_HEADER_LEN   = 4      //a frame header is 4 bytes
	SECDATA_HEADER_LEN = 5      //a secData header len
	FRAMEFLAG          = 0xf3db //a frame is started with "0xf3db"
	MAX_DATA_LEN       = 1024   // the max frame data length
	MAX_FRAME_LEN      = 1080 + FRAME_HEADER_LEN
	CMDFRAME           = 0 //command frame flag
	DATAFRAME          = 1 //data frame flag


)

type Frame interface {
	Write()
	Read()
}



/*        1 byte      1 byte   　2 byte    1byte       1 byte      2byte  1 byte  less than 1024 byte
 * 格式:[FRAMEFLAG0][FRAMEFLAG0][length]([FRAMEFLAG0][FRAMEFLAG0][length][flag][data........])
 *      前两2个是协议开头标志　　后面数据字节数　　　括号里是加密的数据　flag表示是数据还是命令　data为真实数据
 */
func Write(flag byte, content []byte) (buf []byte, err error) {
	contentLength := len(content)
	secBuff := make([]byte, contentLength+SECDATA_HEADER_LEN)
	sec := NewLEStream(secBuff)
	sec.WriteUint16(FRAMEFLAG)
	sec.WriteUint16(uint16(contentLength))
	sec.WriteByte(flag)
	err = sec.WriteBuff(content)
	if err != nil {
		return nil, err
	}

	return
}


/read data from peer and decrypt data, and return data
func Read() error {
	//step 1: 分配frame长度的大小的空间
	frameHeaderBuf := make([]byte, FRAME_HEADER_LEN)
	var n int
	var err error
	var realNeed int = 0
	//step 2:　读取frame长度大小的数据
	for {
		n, err = S.Conn.Read(frameHeaderBuf[realNeed:])
		if err != nil && err != io.EOF {
			log.Error("[ReadPacket]Read Frame error:%s", err)
			return fmt.Errorf("[ReadPacket]Read Frame error:%s", err)
		}
		realNeed = realNeed + n
		if realNeed == FRAME_HEADER_LEN || 0 == n {
			realNeed = 0
			break
		}
	}

	frameHeader := NewLEStream(frameHeaderBuf)
	frameFlag, errFlag := frameHeader.ReadUint16()
	if errFlag != nil {

		log.Error("[ReadPacket]read frame flag fail:%s", errFlag)
		return fmt.Errorf("[ReadPacket]frame flag is wrong:0x%x", frameFlag)
	}
	secDataLen, errDataLen := frameHeader.ReadUint16()
	if errDataLen != nil {
		log.Error("[ReadPacket]read frame secDataLen fail:%s", errDataLen)
		return errDataLen
	}
	if frameFlag != FRAMEFLAG {
		log.Error("[ReadPacket]frame flag is wrong:0x%x", frameFlag)
		return fmt.Errorf("[ReadPacket]frame flag is wrong:0x%x", frameFlag)
	}

	if secDataLen > MAX_FRAME_LEN {
		log.Error("[ReadPacket]SecDataLen wrong:0x%x", secDataLen)
		return fmt.Errorf("[ReadPacket]SecDataLen wrong:0x%x", secDataLen)
	}
	//step 3: 分配加了密的sec Data的长度的空间
	encSecData := make([]byte, secDataLen)

	for {
		n, err = S.Conn.Read(encSecData[realNeed:])
		if err != nil && err != io.EOF {
			return fmt.Errorf("[Readpacket] read Sec Data error:", err)
		}
		realNeed = realNeed + n
		if realNeed == int(secDataLen) || n == 0 {
			realNeed = 0
			break
		}

	}

	var decSecData []byte
	//step 4: 由于暂时没法知道解密之后的数据是多大，所以直接先分配最大的
	//TODO:   当然是可以通过EncLen这个函数反过来推知，暂时不做
	outSecData := make([]byte, MAX_DATA_LEN)
	decSecData, err = Decrypt(encSecData, outSecData)
	if err != nil {
		log.Error("[ReadPacket]dec sec data error:%s", err)
		return fmt.Errorf("[ReadPacket]dec sec data error:%s", err)
	}

	secDataHeader := NewLEStream(decSecData)
	secDataFlag, errSecDataFlag := secDataHeader.ReadUint16()
	if errSecDataFlag != nil {
		log.Error("[ReadPacket]Read Sec Data Flag error:%s", errSecDataFlag)
		return fmt.Errorf("[ReadPacket]Read Sec Data Flag error:%s", errSecDataFlag)
	}
	if secDataFlag != FRAMEFLAG {
		log.Error("[ReadPacket]Sec Data Flag wrong:0x%x", secDataFlag)
		return fmt.Errorf("[ReadPacket]Sec Data Flag wrong:0x%x", secDataFlag)
	}
	dataLen, errSecDataLen := secDataHeader.ReadUint16()
	if errSecDataLen != nil {
		log.Error("[ReadPacket]Read Sec Data Len error:%s", errSecDataLen)
		return fmt.Errorf("[ReadPacket]Read Sec Data Len error:%s", errSecDataLen)
	}
	secDataType, errSecDataType := secDataHeader.ReadByte()
	if errSecDataType != nil {
		log.Error("[ReadPacket]Read Sec Data Type error:%s", errSecDataType)
		return fmt.Errorf("[ReadPacket]Read Sec Data Type error:%s", errSecDataType)
	}

	if secDataType != CMDFRAME && secDataType != DATAFRAME {
		log.Error("[ReadPacket]Sec Data Type wrong:%d", secDataType)
		return fmt.Errorf("[ReadPacket]Sec Data Type wrong:%d", secDataType)
	}

	realDataLen := uint16(len(decSecData[secDataHeader.pos:]))
	if dataLen != realDataLen {
		log.Error("[ReadPacket]Read Sec Data len %d is not equal need Read Sec Data len %d", realDataLen, dataLen)
		return fmt.Errorf("[ReadPacket]Read Sec Data len %d is not equal need Read Sec Data len %d", realDataLen, dataLen)
	}

	S.typ = secDataType
	S.length = secDataLen
	S.data = secDataHeader.buff[secDataHeader.pos:]
	return nil
}
