package main

import (
	"os"
	"fmt"
	"github.com/nicle-lin/ADCM/lib/update"
)

func main() {

	if err := update.EncFile(os.Args[1],os.Args[2]); err != nil {
		fmt.Println(err)
	}else {
		fmt.Println("success")
	}

}
