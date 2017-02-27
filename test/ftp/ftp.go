package main

import "fmt"
import (
	"strings"
)

func main() {
	ssuPath := "ftp://200.200.145.15:21/ad/ad6.5.ssu"
	str := strings.Split(ssuPath,"//")
	host := strings.Split(str[1],"/")[0]
	fmt.Println(host)

}
