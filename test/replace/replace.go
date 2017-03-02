package main

import (
	"strings"
	"os"
	"fmt"
	"unicode"
)

func main() {
	app := strings.TrimSuffix(os.Args[1],"_des")
	appsh := strings.Replace(app,"app","appsh",2)
	fmt.Println(appsh)

	str := `abc def ghij    klmn
    		123
    		456`
	fmt.Printf("Fields are: %q", strings.FieldsFunc(str, unicode.IsSpace))
}
