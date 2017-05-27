package main

import "fmt"

func main() {

	x := [...]int{0,1,2,3,4,5,6,7,8,9}
	fmt.Println("x:",x)

	fmt.Println("x::",x[:])
	fmt.Println("x[2:5]",x[2:5])       //len:5-2=3, cap 10-2=8
	fmt.Println("x[2:5:7]",x[2:5:7])    //len:5-2=3 , cap 7-2=5
	fmt.Println("x[4:]",x[4:])         //len:10-4 = 6, cap 10-4=6

	//第三个参数表示max,第一个表示low, 第二个表示high
	//len = high - low
	//cap = max - low
	fmt.Println("x[:4]",x[:4])                   //len:4  cap 10-0 = 10
	fmt.Println("x[:4:6]",x[:4:6])               //len:4-0 = 4  cap 6-0 = 6
}
