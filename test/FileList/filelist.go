package main

import (
	"os"
	"io/ioutil"
	"fmt"
	"regexp"
	"path/filepath"
)

func GetFileList(dir string) []os.FileInfo{
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Get file list error: %v\n", err)
		return nil
	}
	return files
}

func GetApps(dir string)(apps []string){

	reg := regexp.MustCompile(`app[\d]`)
	files := GetFileList(dir)
	d := filepath.Dir(dir)
	for _, v := range files{
		if reg.FindAllString(v.Name(),-1) != nil{
			apps = append(apps,filepath.Join(d,v.Name()))
		}

	}
	return apps
}


func GetDesApps(DesAppPath string) (desApps []string){
	reg := regexp.MustCompile(`app[\d]_des`)
	files := GetFileList(DesAppPath)
	for _, v := range files{
		//return nil means find the str
		if reg.FindAllString(v.Name(),-1) != nil{
			desApps = append(desApps,v.Name())
		}

	}
	return desApps

}

func main() {

	fmt.Println(GetFileList(os.Args[1]))
	files := GetFileList(os.Args[1])
	for _, v := range files{
		fmt.Println(v.Name())
	}
	fmt.Println(GetApps(os.Args[1]))
	fmt.Println(GetDesApps(os.Args[1]))
	fmt.Println(filepath.Dir(os.Args[1]))
}
