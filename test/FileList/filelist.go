package main

import (
	"os"
	"io/ioutil"
	"fmt"
)

func GetFileList(dir string) []os.FileInfo{
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get file list error: %v\n", err)
		return nil
	}
	return files
}

func Get


func main() {

	fmt.Println(GetFileList(os.Args[1]))
	files := GetFileList(os.Args[1])
	for _, v := range files{
		fmt.Println(v.Name())
	}
}
