package main

import (
	"github.com/go-ini/ini"
	//"strings"
	"fmt"
	"os"
	"strings"
)

func ConfirmRebootDevice(iniConf string)error{
	cfg, err := ini.Load(iniConf)
	if err != nil {
		return err
	}
	value := cfg.Section("restart").Key("needrestart").String()


	if strings.ToLower(value) == "yes" {
		fmt.Println("yes")
	}

	fmt.Println(value)
	return nil
}


func main() {
	fmt.Println(ConfirmRebootDevice(os.Args[1]))
}
