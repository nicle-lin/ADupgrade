package main

import (
	"fmt"
	"github.com/nicle-lin/ADupgrade/lib/update"
)

func main() {
	encData := []byte{0x0f, 0x00, 0x21, 0x1d,
		0x14, 0x0c, 0x05, 0x60, 0xdd, 0xf0, 0x95, 0x3c,
		0xb3, 0x81, 0x16, 0xd4, 0x87, 0xa4}
	//update.Decrypt()
	outSecData := make([]byte, 1024)
	decSecData, err := update.Decrypt(encData, outSecData)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(update.EncLen(1024))
	fmt.Println(string(decSecData[5:]))
}
