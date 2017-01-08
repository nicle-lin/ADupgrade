package update

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
	LOGIN   = "login"
	EXEC    = "excute" //exec cmd,be care!!!!!
	GET     = "get"
	GETOVER = "getover"
	PUT     = "put"
	PUTOVER = "putover"
	VERSION = "version"
)


type PacketHeader struct {
	ID         uint32
	PacketType uint32
	Len        uint32
	Version    uint32
	Ack        uint32
	Token      uint32
}

type Frame struct {

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

