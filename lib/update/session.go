package update

import (
	"fmt"
	"net"
	"time"
)

type Address struct {
	IP   string
	Port string
}

type PeerInfo struct {
	SerVersion string
	AppVersion string
}

type SSUSlice struct {
	SSUPacket string
	SSUType int
}

type SSU struct {
	Flag       bool   //Same Version SSU packet has been unpack or not
	Version    string //AD version
	SSUPackage string //SSU packet name
	SSUType int    /*PACKAGE_TYPE = 1 RESTORE_TYPE = 2 EXECUTE_TYPE  = 3 AUTOBAK_NUMS  = 10 */
	SSUInfo []SSUSlice
}

type Unpack struct {
	FolderPrefix      string //random string
	CurrentWorkFolder string
	LocalBackSh       string
	LocalPreCfgSh     string
	LocalCfgSh        string
	LocalUpdHistory   string
	LocalUpdCheck     string
	ServerAppRe       string
	ServerAppSh       string
	ServerCfgPre      string
	ServerCfgSh       string
	TempExecFile      string
	TempRstFile       string
	TempRetFile       string
	CustomErrFile     string
	LoginPwdFile      string
	Compose           string

	SingleUnpkg  string
	ComposeUnpkg string
	PkgTemp      string
	Download     string
	AutoBak      string

	UpdatePath  string


}

type Package struct {
	UpdatingFlag bool      //updating or not
	UpdateTime   time.Time //when to update
	RestoringFlag bool
}

type Cfg struct {
	CfgPath string
	CfgPathTmp string
}

type Session struct {
	Conn net.Conn
	*PeerInfo
	*SecData
}

type Update struct {
	*SSU
	*Package
	*Unpack
	*Cfg
}

//read data from peer and decrypt data, and return data
func (S *Session) ReadPacket() error {
	//step 1: 分配frame长度的大小的空间
	frameHeaderBuf := make([]byte, FRAME_HEADER_LEN)
	var n int
	var err error
	//step 2:　读取frame长度大小的数据
	n, err = S.Conn.Read(frameHeaderBuf)
	if n != FRAME_HEADER_LEN || err != nil {
		fmt.Println("read frame len is FRAME_HEADER_LEN,it is ",n)
		//return nil, fmt.Errorf("frame len is wrong:#%#v\n",n)
		return fmt.Errorf("frame len is wrong:%d,err msg is:%s\n", n,err)
	}
	frameHeader := NewLEStream(frameHeaderBuf)
	frameFlag, errFlag := frameHeader.ReadUint16()
	if errFlag != nil {
		return errFlag
	}
	secDataLen, errDataLen := frameHeader.ReadUint16()
	if errDataLen != nil {
		return errDataLen
	}
	if frameFlag != FRAMEFLAG {
		fmt.Printf("frameflage is wrong:0x%x", frameFlag)
		//return nil, fmt.Errorf("frameflage is wrong:#%#v\n",frameFlag)
		return fmt.Errorf("frameflage is wrong:#%#v\n", frameFlag)
	}

	if secDataLen > MAX_DATA_LEN {
		//return nil, fmt.Errorf("sec data len is wrong:#%#v\n",secDataLen)
		return fmt.Errorf("sec data len is wrong:#%#v\n", secDataLen)
	}
	//step 3: 分配加了密的sec Data的长度的空间
	encSecData := make([]byte, secDataLen)
	n, err = S.Conn.Read(encSecData)

	fmt.Printf("read frame enc data:%#v\n", encSecData)

	var decSecData []byte
	//step 4: 由于暂时没法知道解密之后的数据是多大，所以直接先分配最大的
	//TODO:   当然是可以通过EncLen这个函数反过来推知，暂时不做　　
	outSecData := make([]byte, MAX_DATA_LEN)
	decSecData, err = Decrypt(encSecData, outSecData)
	if err != nil {
		fmt.Println("dec sec data error:", err)
		//return nil,fmt.Errorf("dec sec data error:\n",err)
		return fmt.Errorf("dec sec data error %s:\n", err)
	}
	//fmt.Println("dec sec data:",string(decSecData))


	secDataHeader := NewLEStream(decSecData)
	secDataFlag, errSecDataFlag := secDataHeader.ReadUint16()
	if errSecDataFlag != nil {
		return errSecDataFlag
	}
	if secDataFlag != FRAMEFLAG {
		fmt.Printf("sec Data flag is wrong:0x%x\n", secDataFlag)
		//return nil,fmt.Errorf("sec Data flag is wrong:0x%x\n",secDataFlag)
		return fmt.Errorf("sec Data flag is wrong:0x%x\n", secDataFlag)
	}
	dataLen, errSecDataLen := secDataHeader.ReadUint16()
	if errSecDataLen != nil {
		return errSecDataLen
	}
	secDataType, errSecDataType := secDataHeader.ReadByte()
	if errSecDataType != nil {
		return errSecDataType
	}
	realDataLen := uint16(len(decSecData[secDataHeader.pos:]))
	//fmt.Println("##############################################")

	if dataLen != realDataLen {
		//fmt.Printf("sec Data len is wrong:0x%x\n", dataLen, "receive data len:ox%x\n", realDataLen)
		//return nil, fmt.Errorf("sec Data len is wrong:0x%x\n",dataLen,"receive data len:ox%x\n",realDataLen)
		return fmt.Errorf("sec Data len is wrong:0x%x\n,receive data len:ox%x\n", dataLen, realDataLen)
	}
	//fmt.Println("----------------befor pos------------------------")
	if secDataType != CMDFRAME && secDataType != DATAFRAME {
		fmt.Printf("sec data type is wrong:0x%x\n", secDataType)
		//return nil, fmt.Errorf("sec data type is wrong:0x%x\n",secDataType)
		return fmt.Errorf("sec data type is wrong:0x%x\n", secDataType)
	}
	//fmt.Println("################almost to pos ################################")
	//fmt.Printf("##############################secDataType:0x%x\n",secDataType)
	S.typ = secDataType
	S.length = secDataLen

	//fmt.Println("###############pos#######################")
	//fmt.Println("decSecData[secDataHeader.pos:]:",decSecData[secDataHeader.pos:])
	S.data = secDataHeader.buff[secDataHeader.pos:]
	//fmt.Println("#################read data seen like is ok#############")
	fmt.Println("###################################")
	fmt.Println(string(S.data))
	fmt.Println("###################################\n\n")
	return nil
}

//just send data to peer
func (S *Session) WritePacket(data []byte) error {
	_, err := S.Conn.Write(data)
	return err
}
