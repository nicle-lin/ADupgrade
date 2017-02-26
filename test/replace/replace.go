package main

import (
	"strings"
	"os"
	"fmt"
)

func main() {
	app := strings.TrimSuffix(os.Args[1],"_des")
	appsh := strings.Replace(app,"app","appsh",1)
	fmt.Println(appsh)
}
