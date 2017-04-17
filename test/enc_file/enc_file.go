package main

import (
	"os"
	"fmt"
	"github.com/nicle-lin/ADCM/lib/update"
)

func main() {

<<<<<<< HEAD
	if err := update.EncFile(os.Args[1],os.Args[2]); err != nil {
=======
	fmt.Println(os.Args[1])
	fmt.Println(os.Args[2])
	if err := EncFile(os.Args[1],os.Args[2]); err != nil {
>>>>>>> 33887eee38223abd968bcdb900d717d4a24d76d0
		fmt.Println(err)
	}else {
		fmt.Println("success")
	}

<<<<<<< HEAD
=======
	if err := PutDesApp(os.Args[1]);err != nil {
		fmt.Println(err)
	}
>>>>>>> 33887eee38223abd968bcdb900d717d4a24d76d0
}
