package main

import (
	"github.com/nicle-lin/ADCM/lib/update"
	"os"
	"fmt"
)

func main() {
	S ,err := update.Login(os.Args[1],os.Args[2],os.Args[3])
	if err != nil {
		fmt.Println("Login fail:",err)
	}
	U := update.InitClient("SANGFOR-M12000-AD-6.6")
	msg,err1 := update.Exec(S,U,os.Args[4])
	if err1 != nil {
		fmt.Println("exec command:",os.Args[4],"fail")
	}else{
		fmt.Println("exec result:",msg)
	}
}
