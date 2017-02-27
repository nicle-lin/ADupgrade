package main

import (
	"github.com/nicle-lin/ADupgrade/lib/update"
	"fmt"

	"os"
)

func main() {
	if err := update.Upgrade("192.168.1.100","51111","admin",os.Args[1]);err != nil{
		fmt.Println("err:",err)
	}else {
		fmt.Println("success")
	}


}
