package main

import (
	"path/filepath"

	"fmt"
	"os"
	"strings"
)



var ssu = "AD6.5(20160809).ssu"
var curr,_= os.Getwd()
var path = filepath.Join(curr,ssu)

func SSUPath(path string) string {
	return strings.Replace(filepath.Base(path),filepath.Ext(path),"",-1)
}

func main() {
	fmt.Println(filepath.Base(path))
	fmt.Println(filepath.Ext(path))
	fmt.Println(filepath.Clean(path))
	fmt.Println(SSUPath(path))

}
