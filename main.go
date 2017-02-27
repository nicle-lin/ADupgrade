package main

import (
	"github.com/nicle-lin/ADupgrade/lib/update"
	"fmt"

)

func main() {
	if err := update.Upgrade("192.168.1.41","5111","admin","AD6.5.ssu");err != nil{
		fmt.Println("err:",err)
	}else {
		fmt.Println("success")
	}
}
