package main

import (
	"github.com/nicle-lin/ADCM/lib/update"
	//"os"
	"fmt"
)
func main() {
	if err := update.GetFile("192.168.1.101","admin","51111","ibdata1","/aclog/database/data/ibdata1");err != nil{
		fmt.Println(err)
	}

	if err := update.GetFile("192.168.1.101","admin","51111","","/aclog/database/data/mysql-bin.index");err != nil{
		fmt.Println(err)
	}
}
