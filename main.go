package main

import (
	"github.com/nicle-lin/ADupgrade/lib/update"
	"fmt"

	"os"
)

func main() {
	if err := update.Upgrade(os.Args[1],"51111","admin",os.Args[2]);err != nil{
		fmt.Println("err:",err)
	}else {
		fmt.Println("success")
	}


}
