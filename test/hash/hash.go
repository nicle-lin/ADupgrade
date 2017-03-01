package main

import (
	"fmt"
)


/*
@package_arr = Array.new
packhash = {"packet" => now_package, "type" => "1"}
@package_arr<<packhash
*/
	//[{"packet"=> no_package,"type"=>"1"},{"packet"=>next_package,"type"=>2}]
type SSUType struct {
 	packet string
	typ int
}
func test(str string)error{
	fmt.Println(str)
	return nil
}
func main() {
	h := make([]SSUType,1)
	h[0].packet= "now"
	h[0].typ = 0

	var ssu SSUType
	ssu.packet = "next"
	ssu.typ = 1
	h = append(h,ssu)
	fmt.Println(h)
	for k, v := range h{
		fmt.Println("K:",k)
		fmt.Println("v:",v)
	}
	fmt.Printf("TTTT:%T\n",h)
	fmt.Printf("#V#V:%#v\n",h)
	fmt.Printf("+V+V:%+v\n",h)
	var fun func(string) error
	fun = test
	fun("hahahahah")
}

