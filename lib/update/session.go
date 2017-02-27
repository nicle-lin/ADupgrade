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
	SerVersion int
	AppVersion string
}

type SSU struct {
	Flag       bool   //Same Version SSU packet has been unpack or not
	Version    string //AD version
	SSUPackage string //SSU packet name
	SSUType int8    /*PACKAGE_TYPE = 1 RESTORE_TYPE = 2 EXECUTE_TYPE  = 3 AUTOBAK_NUMS  = 10 */
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
	frameHeaderBuf := make([]byte, FRAME_HEADER_LEN)
	var n int
	var err error
	n, err = S.Conn.Read(frameHeaderBuf)
	if n != FRAME_HEADER_LEN || err != nil {
		fmt.Println("read frame len > Max frame len")
		//return nil, fmt.Errorf("frame len is wrong:#%#v\n",n)
		return fmt.Errorf("frame len is wrong:#%#v\n", n)
	}
	frameHeader := NewLEStream(frameHeaderBuf)
	frameFlag, _ := frameHeader.ReadUint16()
	secDataLen, _ := frameHeader.ReadUint16()
	if frameFlag != FRAMEFLAG {
		fmt.Printf("frameflage is wrong:0x%x", frameFlag)
		//return nil, fmt.Errorf("frameflage is wrong:#%#v\n",frameFlag)
		return fmt.Errorf("frameflage is wrong:#%#v\n", frameFlag)
	}

	if secDataLen > MAX_DATA_LEN {
		//return nil, fmt.Errorf("sec data len is wrong:#%#v\n",secDataLen)
		return fmt.Errorf("sec data len is wrong:#%#v\n", secDataLen)
	}

	encSecData := make([]byte, secDataLen)
	n, err = S.Conn.Read(encSecData)

	fmt.Printf("read frame enc data:%#v\n", encSecData)

	var decSecData []byte
	outSecData := make([]byte, MAX_DATA_LEN)
	decSecData, err = Decrypt(encSecData, outSecData)
	if err != nil {
		fmt.Println("dec sec data error:", err)
		//return nil,fmt.Errorf("dec sec data error:\n",err)
		return fmt.Errorf("dec sec data error %s:\n", err)
	}

	secDataHeader := NewLEStream(decSecData)
	secDataFlag, _ := secDataHeader.ReadUint16()
	if secDataFlag != FRAMEFLAG {
		fmt.Printf("sec Data flag is wrong:0x%x\n", secDataFlag)
		//return nil,fmt.Errorf("sec Data flag is wrong:0x%x\n",secDataFlag)
		return fmt.Errorf("sec Data flag is wrong:0x%x\n", secDataFlag)
	}
	dataLen, _ := secDataHeader.ReadUint16()
	secDataType, _ := secDataHeader.ReadByte()
	realDataLen := uint16(len(decSecData[secDataHeader.pos:]))
	if dataLen != realDataLen {
		//fmt.Printf("sec Data len is wrong:0x%x\n", dataLen, "receive data len:ox%x\n", realDataLen)
		//return nil, fmt.Errorf("sec Data len is wrong:0x%x\n",dataLen,"receive data len:ox%x\n",realDataLen)
		return fmt.Errorf("sec Data len is wrong:0x%x\n,receive data len:ox%x\n", dataLen, realDataLen)
	}
	if secDataType != CMDFRAME && secDataType != DATAFRAME {
		fmt.Printf("sec data type is wrong:0x%x\n", secDataType)
		//return nil, fmt.Errorf("sec data type is wrong:0x%x\n",secDataType)
		return fmt.Errorf("sec data type is wrong:0x%x\n", secDataType)
	}
	S.typ = secDataType
	S.length = secDataLen
	S.data = decSecData[secDataHeader.pos:]
	return nil
}

//just send data to peer
func (S *Session) WritePacket(data []byte) error {
	_, err := S.Conn.Write(data)
	return err
}
