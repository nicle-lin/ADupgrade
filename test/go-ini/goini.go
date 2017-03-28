package main

import (
	"github.com/go-ini/ini"
	"fmt"
	"github.com/nicle-lin/ADCM/lib/update"
	"sync"
	"time"
)

func WriteIni() error {
	cfg, err := ini.Load("ssu.conf")
	if err != nil {
		return err
	}
	if _,err := cfg.NewSection("ssu"); err !=nil {
		return err
	}
	if err := cfg.SaveTo("ssu.conf"); err != nil {
		return err
	}
	return nil
}

func GetKeyValue()  {
	cfg, err := ini.Load("ssu.conf")
	if err != nil {
		fmt.Println("Load error:",err)
	}
	value := cfg.Section("restart").Key("reboot").String()
	fmt.Println("value:",value)

	if _, err := cfg.NewSection("ssuu"); err != nil {
		fmt.Println("newsection error:",err)
	}
	names := cfg.SectionStrings()
	fmt.Println("section:",names)
}



func main() {
	fmt.Println(WriteIni())
	GetKeyValue()
	var m = new(sync.RWMutex)
	for i := 0; i < 10; i++ {
		go update.WriteMsgToConf("ssu.conf","ssu","key"+string('A' + i),"value"+string('A' + i),m )
	}
	time.Sleep(2 * time.Second)
	keyValue ,err := update.FindAllKeyValue("ssu.conf","ssu",m)
	if err != nil {
		fmt.Println("find error:",err)
	}
	fmt.Println("keyvalue:",keyValue)
	fmt.Println("delete first key",keyValue[0])

	value, err1 := update.CompareKeyFromMap(keyValue,"keyH2")
	if err1 != nil {
		fmt.Println(err1)
	}
	fmt.Println("key is keyH, value is ", value)

}
