package main

import (
	"github.com/nicle-lin/ADCM/lib/update"
	"os"
	"fmt"
)
func main() {
	if err := update.GetFile(os.Args[1],"admin","51111","ibdata1","/aclog/database/data/ibdata1");err != nil{
		fmt.Println(err)
	}
}
