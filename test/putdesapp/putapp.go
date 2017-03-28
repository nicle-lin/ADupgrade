package main

import (
	"github.com/nicle-lin/ADCM/lib/update"
	"os"
	"fmt"
)

func main() {
	if err := update.PutFile(os.Args[1],"51111","admin",os.Args[2],os.Args[3]);err !=nil {
		fmt.Println("fail:",err)
	}else{
		fmt.Println("success")
	}
}
